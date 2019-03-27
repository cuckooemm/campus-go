package auth

import (
	"campus/helper/check"
	"campus/helper/sensitiveWord"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"time"
)

func Register(c *gin.Context) {
	var entry input.Register
	rsp := response.NewResponse(c)
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
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
	if sensitiveWord.SensitiveWordFilter(entry.Nickname) {
		rsp.ErrorMsg(e.ERROR_NICKNAME_EXIST)
		return
	}
	// 剔除昵称的空格和换行
	entry.Nickname = models.ReplaceName(entry.Nickname)
	if models.IsAccountExist(entry.Flag,entry.Account) {
		rsp.ErrorMsg(e.ERROR_ACCOUNT_REGISTERED)
		return
	}
	if models.IsNameExist(entry.Nickname) {
		rsp.ErrorMsg(e.ERROR_NICKNAME_EXIST)
		return
	}
	var tm models.Code
	if !models.GetLastSendTime(entry.Flag,entry.Account,entry.Code,&tm) {
		rsp.ErrorMsg(e.ERROR_CODE_UNSENT)
		return
	}
	//  10分之有效期
	if tm.CreatedAt.Add(10 * time.Minute).Unix() < time.Now().Unix() {
		rsp.ErrorMsg(e.ERROR_CODE_INVALID)
		return
	}
	if err := models.Register(&entry,&tm);err != nil{
		rsp.ErrorMsg(e.ERROR_REGISTER_FAIL)
		return
	}

	rsp.Success("你有属于自己的帐号啦",nil)
}
