package auth

import (
	"campus/helper/oauth"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
)

func QQLogin(c *gin.Context) {
	var entry input.QQLogin
	rsp := response.NewResponse(c)
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var userInfo oauth.OpenQQRespond
	if err := oauth.GetQQUserInfo(&entry,&userInfo);err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	if userInfo.Ret != 0 {
		rsp.ErrorMsg(e.AUTH_FAIL)
		return
	}
	var auth models.Auth
	// 查询当前用户是否注册
	if models.QQLoginConfirm(&entry, &auth) {
		// 小于1则封禁
		if auth.Status < 0 {
			rsp.ErrorMsg(e.AUTH_FORBIDDEN)
			return
		}
		if auth.Token != "" {
			oauth.AddTokenBlackList(auth.Token)
		}
	}else {
		// 进行注册
		var err error
		auth.ID,err = models.QQLoginRegister(&entry,&userInfo)
		if err != nil{
			rsp.ErrorMsg(e.AUTH_FAIL)
			return
		}
	}

	token, err := oauth.GenerateToken(auth.ID,"qq")
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
