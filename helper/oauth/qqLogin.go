package oauth

import (
	"campus/models/input"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const host = "https://graph.qq.com"
const api_userInfo = "/user/get_user_info"
const appid = "1106710088"
const appkey = "jcABu2ctRgEHvjf9"

/**
	返回的数据
 */
type OpenQQRespond struct {
	Ret      uint64 `json:"ret"`
	Msg      string `json:"msg"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Province string `json:"province"`
	City     string `json:"city"`
	Avatar   string `json:"figureurl_qq_2"`
}

// 通过QQ open_id 获取用户信息
func GetQQUserInfo(q *input.QQLogin, result *OpenQQRespond) error {
	// 添加参数
	url := host + api_userInfo + "?oauth_consumer_key=" + appid + "&access_token=" + q.Token + "&openid=" + q.OpenID + "&format=json"
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, result); err != nil {
		return err
	}
	return nil
}

// 通过微信open_id 获取用户信息
func GetWXUserInfo(openid, token string) {

}
