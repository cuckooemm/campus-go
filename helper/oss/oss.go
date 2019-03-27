package oss

import (
	"campus/helper/logging"
	"campus/helper/setting"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"hash"
	"io"
	"time"
)

var campusBucket *oss.Bucket

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}
type OssConfig struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}
func Setup() {
	client, err := oss.New(setting.OSSSetting.Endpoint, setting.OSSSetting.AccesskeyID,
		setting.OSSSetting.AccesskeySecret, oss.Timeout(10, 120))
	if err != nil {
		logging.ErrorMsg("OSS配置失败", err)
		return
	}
	if campusBucket, err = client.Bucket(setting.OSSSetting.BucketName); err != nil {
		logging.ErrorMsg("OSS Bucket绑定失败", err)
		return
	}
	logging.Info("OSS 配置成功")

}

func Detele(object []string) {
	delRes, err := campusBucket.DeleteObjects(object)
	if err != nil {
		logging.Warn("OSS Object 删除失败", zap.Strings("object", delRes.DeletedObjects), zap.String("err", err.Error()))
		return
	}
	logging.Info("OSS Object删除成功", zap.Strings("object", delRes.DeletedObjects))
}

func Find(object string) bool {
	isExist, err := campusBucket.IsObjectExist(object)
	if err != nil {
		logging.Warn("OSS 发生错误 ", zap.String("err", err.Error()))
		return false
	}
	return isExist
}
// TODO 待完善
var upload_dir string = ""

// 用户上传图片需要的签名数据
func GetPolicyToken(token *PolicyToken) error {
	expireEnd := time.Now().Add(setting.OSSSetting.Expiration).Unix()
	var config OssConfig
	config.Expiration = getGmtIso8601(expireEnd)
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	result, _ := json.Marshal(config)
	resByte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(setting.OSSSetting.AccesskeySecret))
	io.WriteString(h, resByte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = setting.OSSSetting.CallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"

	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		logging.WarnMsg("oss 签名出错", err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)
	token.AccessKeyId = setting.OSSSetting.AccesskeyID
	token.Host = setting.OSSSetting.Path
	token.Expire = expireEnd
	token.Signature = string(signedStr)
	token.Directory = upload_dir
	token.Policy = string(resByte)
	token.Callback = string(callbackBase64)
	return nil
}

func getGmtIso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}
