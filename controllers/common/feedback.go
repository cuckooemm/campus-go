package common

import (
	"campus/models"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	rsp := response.NewResponse(c)
	var entry models.Feedback
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	if err := entry.Create(); err != nil {
		rsp.Error()
		return
	}
	rsp.Success("已接收到您的反馈",nil)
}