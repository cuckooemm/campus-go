package push

import (
	"bytes"
	"campus/helper/logging"
	"campus/helper/setting"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type umengNotificationData struct {
	AppKey    string `json:"appkey"`    // 必填 应用唯一标识
	Timestamp int64  `json:"timestamp"` // 必填 时间戳，10位或者13位均可，时间戳有效期为10分钟
	Type      string `json:"type"`      // 必填 消息发送类型,其值可以为: unicast-单播
	//listcast-列播(要求不超过500个device_token)
	//filecast-文件播 (多个device_token可通过文件形式批量发送）
	//broadcast-广播
	//groupcast-组播 (按照filter条件筛选特定用户群, 具体请参照filter参数)
	//customizedcast(通过开发者自有的alias进行推送),
	//包括以下两种case:
	//- alias: 对单个或者多个alias进行推送
	//- file_id: 将alias存放到文件后，根据file_id来推送
	DeviceTokens string `json:"device_tokens,omitempty"` // 可选 设备唯一表示 当type=unicast时,必填, 表示指定的单个设备 当type=listcast时,必填,要求不超过500个, 多个device_token以英文逗号间隔
	AliasType    string `json:"alias_type,omitempty"`    // 可选 当type=customizedcast时，必填，alias的类型, alias_type可由开发者自定义,开发者在SDK中 调用setAlias(alias, alias_type)时所设置的alias_type
	Alias        string `json:"alias,omitempty"`         // 可选 当type=customizedcast时, 开发者填写自己的alias。 要求不超过50个alias,多个alias以英文逗号间隔。 在SDK中调用setAlias(alias, alias_type)时所设置的alias
	//FileId         string          `json:"file_id,omitempty"`       // 可选 当type=filecast时，file内容为多条device_token, device_token以回车符分隔 当type=customizedcast时，file内容为多条alias， alias以回车符分隔，注意同一个文件内的alias所对应 的alias_type必须和接口参数alias_type一致。 注意，使用文件播前需要先调用文件上传接口获取file_id, 具体请参照"2.4文件上传接口"
	Payload        *androidPayload `json:"payload,omitempty"`
	ProductionMode bool            `json:"production_mode,omitempty"` // 可选 正式/测试模式。测试模式下，广播/组播只会将消息发给测试设备。测试设备需要到web上添加。Android: 测试设备属于正式设备的一个子集。
	Description    string          `json:"description,omitempty"`     // 可选 发送消息描述，建议填写。
	Mipush         bool            `json:"mipush,omitempty"`          // 可选，默认为false。当为true时，表示MIUI、EMUI、Flyme系统设备离线转为系统下发
	MiActivity     string          `json:"mi_activity,omitempty"`     // 可选，mipush值为true时生效，表示走系统通道时打开指定页面acitivity的完整包
}

type Policy struct {
	StartTime  string `json:"start_time,omitempty"`   // 可选 定时发送时间，若不填写表示立即发送。定时发送时间不能小于当前时间 格式: "YYYY-MM-DD HH:mm:ss"。注意, start_time只对任务生效。
	ExpireTime string `json:"expire_time,omitempty"`  // 可选 消息过期时间,其值不可小于发送时间或者 start_time(如果填写了的话),如果不填写此参数，默认为3天后过期。格式同start_time
	MaxSendNum int    `json:"max_send_num,omitempty"` // 可选 发送限速，每秒发送的最大条数。开发者发送的消息如果有请求自己服务器的资源，可以考虑此参数。
	OutBizNo   string `json:"out_biz_no,omitempty"`   // 可选 开发者对消息的唯一标识，服务器会根据这个标识避免重复发送。有些情况下（例如网络异常）开发者可能会重复调用API导致消息多次下发到客户端。如果需要处理这种情况，可以考虑此参数。注意, out_biz_no只对任务生效。
}
type androidPayload struct {
	DisplayType string                 `json:"display_type"`    // 必填，消息类型: notification(通知)、message(消息)
	Body        *AndroidBody           `json:"body,omitempty"`  // 当display_type=notification时，body包含如下参数:
	Extra       map[string]interface{} `json:"extra,omitempty"` // notification 可选  message类型不写  只需填写custom即可，
}

// notification body
type AndroidBody struct {
	Ticker    string      `json:"ticker,omitempty"`   // 必填 通知栏提示文字
	Title     string      `json:"title,omitempty"`    // 必填 通知标题
	Text      string      `json:"text,omitempty"`     // 必填 通知文字描述
	AfterOpen string      `json:"after_open"`         // 必填 值可以为: "go_app": 打开应用 "go_url": 跳转到URL "go_activity": 打开特定的activity "go_custom": 用户自定义内容。
	Url       string      `json:"url,omitempty"`      // 可选 当"after_open"为"go_url"时，必填。 通知栏点击后跳转的URL，要求以http或者https开头
	Activity  string      `json:"activity,omitempty"` // 可选 当"after_open"为"go_activity"时，必填。 通知栏点击后打开的Activity
	Custom    interface{} `json:"custom,omitempty"`   // 可选 display_type=message, 或者 display_type=notification且 "after_open"为"go_custom"时， 该字段必填。用户自定义内容, 可以为字符串或者JSON格式。
}

type result struct {
	Code string `json:"ret,omitempty"`
	Data map[string]string
}

func PushAliasNotification(alias *int64, body *AndroidBody) {
	pushJson := umengNotificationData{AppKey: setting.PushSetting.AppKey, Type: "customizedcast", Timestamp: time.Now().Unix(),
		AliasType: "campus", Alias: strconv.FormatInt(*alias, 10),
		Payload:        &androidPayload{DisplayType: "notification", Body: body},
		ProductionMode: setting.PushSetting.Mode, Description: "发送通知", Mipush: true}
	postBody, err := json.Marshal(pushJson)
	if err != nil {
		logging.ErrorMsg("推送失败", err)
		return
	}
	var result result
	send(&postBody, &result)
	if result.Code == "SUCCESS" {
		return
	}
	logging.Error("推送失败", zap.String("code", result.Data["error_code"]), zap.String("msg", result.Data["error_msg"]))
}

func PushAliasMessage(alias *int64, body *AndroidBody) {
	pushJson := umengNotificationData{AppKey: setting.PushSetting.AppKey, Type: "customizedcast", Timestamp: time.Now().Unix(),
		AliasType: "campus", Alias: strconv.FormatInt(*alias, 10),
		Payload:        &androidPayload{DisplayType: "message", Body: body},
		ProductionMode: setting.PushSetting.Mode, Description: "发送Message", Mipush: true}
	postBody, err := json.Marshal(pushJson)
	if err != nil {
		logging.ErrorMsg("推送失败", err)
		return
	}
	println(string(postBody))
	var result result
	send(&postBody, &result)
	if result.Code == "SUCCESS" {
		return
	}
	logging.Error("推送失败", zap.String("code", result.Data["error_code"]), zap.String("msg", result.Data["error_msg"]))
}
func PushCommentToAlias(alias int64,hint,title string,text string)  {
	pushJson := umengNotificationData{AppKey: setting.PushSetting.AppKey, Type: "customizedcast", Timestamp: time.Now().Unix(),
		AliasType: "campus", Alias: strconv.FormatInt(alias, 10),
		Payload:        &androidPayload{DisplayType: "notification",
		Body: &AndroidBody{Ticker:hint,Title:title,Text:text,AfterOpen:"go_activity",Activity: "cn.campuswall.campus.ui.user.Messageactivity"}},
		ProductionMode: setting.PushSetting.Mode, Description: "发送通知", Mipush: true}
	postBody, err := json.Marshal(pushJson)
	if err != nil {
		logging.ErrorMsg("推送失败", err)
		return
	}
	var result result
	send(&postBody, &result)
	if result.Code == "SUCCESS" {
		return
	}
	logging.Error("推送失败", zap.String("code", result.Data["error_code"]), zap.String("msg", result.Data["error_msg"]))
}
func send(postBody *[]byte, result *result) {
	sign := MD5("POST" + setting.PushSetting.Host + string(*postBody) + setting.PushSetting.AndroidAppMasterSecret)
	url := setting.PushSetting.Host + "?sign=" + sign
	response, err := http.Post(url, "application/json", bytes.NewReader(*postBody))
	defer response.Body.Close()
	if err != nil {
		logging.ErrorMsg("推送失败", err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.ErrorMsg("推送失败", err)
		return
	}
	json.Unmarshal(body, result)
	logging.Info(string(body))
}
func MD5(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
