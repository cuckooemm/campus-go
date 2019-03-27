package auth

import (
	"campus/helper/logging"
	"campus/helper/oauth"
	"campus/helper/setting"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
)

func WeAppLogin(c *gin.Context) {
	rsp := response.NewResponse(c)
	var entry input.WeAppLogin
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}

	res,err := weapp.Login(setting.WeChatSetting.AppID,setting.WeChatSetting.AppSecret,entry.Code)
	if err != nil{
		rsp.ErrorMsg(e.AUTH_FAIL)
		logging.WarnMsg("微信接口登录失败",err)
		return
	}
	var auth models.Auth
	if models.WeChatAppLoginConfirm(res.OpenID, &auth) {
		if auth.Status < 1 {
			rsp.ErrorMsg(e.AUTH_FORBIDDEN)
			return
		}
		if auth.Token != "" {
			oauth.AddTokenBlackList(auth.Token)
		}
	}else {
		if entry.Nickname == "" {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
		auth.ID,err = models.WeChatAppRegister(res.OpenID,&entry)
		if err != nil {
			rsp.ErrorMsg(e.AUTH_FAIL)
			return
		}
	}

	token, err := oauth.GenerateToken(auth.ID,"weapp")
	if err != nil {
		rsp.Error()
		return
	}
	// 更新token
	if err := models.UpdateToken(token,auth.ID); err != nil {
		rsp.Error()
		return
	}
	rsp.LoginSuccess(token)
}
