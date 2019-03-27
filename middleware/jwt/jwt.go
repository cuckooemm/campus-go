package jwt

import (
	"campus/helper/oauth"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		rsp := response.NewResponse(c)
		// query token in header
		token := c.Request.Header.Get("Authorization")
		// 判断token 是否在黑名单中
		if len(token) < 10 {
			rsp.ErrorMsg(e.AUTH_LOGIN)
			c.Abort()
			return
		}
		if oauth.TokenBlackList(token) {
			rsp.Authorization()
			c.Abort()
			return
		}
		claims ,err := oauth.ParseToken(token)
		if err != nil {
			rsp.Authorization()
			c.Abort()
			return
		}
		c.Set("claims",claims)
		c.Next()
	}
}

