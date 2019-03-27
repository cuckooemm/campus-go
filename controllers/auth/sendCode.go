package auth

import (
	"campus/helper/check"
	"campus/helper/email"
	"campus/helper/sns"
	"campus/helper/utils"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"time"
)

func SendCode(c *gin.Context) {
	var entry input.SendCode
	rsp := response.NewResponse(c)
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.ERROR_ACCOUNT_INVALID)
		return
	}
	if err := check.IsEmail(entry.Account); err == nil {
		entry.Flag = "email"
	}
	if err := check.IsPhone(entry.Account); err == nil {
		entry.Flag = "phone"
	}
	if entry.Flag == "" {
		rsp.ErrorMsg(e.ERROR_ACCOUNT_INVALID)
		return
	}
	// operation 1 is register
	if entry.Operation == 1 {
		if models.IsAccountExist(entry.Flag,entry.Account) {
			rsp.ErrorMsg(e.ERROR_ACCOUNT_REGISTERED)
			return
		}
	}
	// 2 is retrieve
	if entry.Operation == 2 {
		if !models.IsAccountExist(entry.Flag,entry.Account) {
			rsp.ErrorMsg(e.ERROR_ACCOUNT_UNREGISTERED)
			return
		}
	}
	if entry.Operation > 2 {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	// 从数据库中获取上次发送验证码时间
	createdAt,ok := models.GetLatelySendTime(&entry)
	if ok {
		if createdAt.Add(time.Minute).Unix() > time.Now().Unix() {
			rsp.ErrorMsg(e.CODE_SENTED)
			return
		}
	}
	// 获取6位随机验证码
	code := utils.Random6Code()
	if entry.Flag == "email" {
		// 向数据库保存此次记录
		ins := &models.Code{
			Email: entry.Account,
			Code:  code,
		}
		if err := ins.Create(); err != nil {
			rsp.Error()
			return
		}
		// 发送邮箱的逻辑
		go func(account string, code string) {
			email.SendEmail(account,code)
		}(entry.Account, code)
	} else {
		// 向数据库保存此次记录
		ins := &models.Code{
			Phone: entry.Account,
			Code:  code,
		}
		if err := ins.Create(); err != nil {
			rsp.Error()
			return
		}
		// 发送短信
		go func(account, code string) {
			sns.SendSns(account,code)
		}(entry.Account, code)
	}
	rsp.Success("验证码马上就要送到您的手中了",nil)
}