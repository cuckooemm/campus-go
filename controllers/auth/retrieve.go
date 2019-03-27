package auth

import (
	"campus/helper/check"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"time"
)

func Retrieve(c *gin.Context) {
	var entry input.Retrieve
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
		rsp.ErrorMsg(e.ERROR_ACCOUNT_UNREGISTERED)
		return
	}
	if !models.IsAccountExist(entry.Flag, entry.Account) {
		rsp.ErrorMsg(e.ERROR_ACCOUNT_UNREGISTERED)
		return
	}

	var tm models.Code
	if !models.GetLastSendTime(entry.Flag, entry.Account, entry.Code, &tm) {
		rsp.ErrorMsg(e.ERROR_CODE_UNSENT)
		return
	}
	//  10分之有效期
	if tm.CreatedAt.Add(10 * time.Minute).Unix() < time.Now().Unix() {
		rsp.ErrorMsg(e.ERROR_CODE_INVALID)
		return
	}
	if err := models.UpdatePassword(&entry, tm.ID); err != nil {
		rsp.ErrorMsg(e.AUTH_UPDATE_PASD_FAIL)
		return
	}
	rsp.Success("密码重置成功", nil)
}
