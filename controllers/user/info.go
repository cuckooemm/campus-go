package user

import (
	"campus/helper/logging"
	"campus/helper/oauth"
	"campus/models"
	"campus/response"
	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	rsp := response.NewResponse(c)
	key,_ := c.Get("claims")
	claims,_ := key.(*oauth.Claims)
	// ID 加密
	var userInfo response.UserInfo
	if err := models.GetUserInfo(claims.UID,&userInfo);err != nil {
		logging.ErrorMsg("数据库连接出错",err)
		rsp.Error()
		return
	}
	rsp.Success("",userInfo)
}