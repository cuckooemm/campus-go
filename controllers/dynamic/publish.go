package dynamic

import (
	"campus/helper/oauth"
	"campus/helper/sensitiveWord"
	"campus/models"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	rsp := response.NewResponse(c)
	key,_ := c.Get("claims")
	claims,_ := key.(*oauth.Claims)
	// TODO 用户每日发动态次数限制

	var dynamic models.Dynamic
	if err := c.ShouldBind(&dynamic); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	// 屏蔽当中的铭感字
	if err := sensitiveWord.SensitiveWordReplace(&dynamic.Content);err != nil {
		rsp.Error()
		return
	}
	dynamic.UID = claims.UID
	dynamic.School = 1
	if err := dynamic.Create(); err != nil {
		rsp.ErrorMsg(e.SERVER_ERROR)
		return
	}
	rsp.Success("动态已成功发布",nil)
}