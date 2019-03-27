package auth

import (
	"campus/helper/oauth"
	"campus/helper/setting"
	"campus/models"
	"campus/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

func UpdateToken(c *gin.Context) {
	rsp := response.NewResponse(c)
	token := c.Request.Header.Get("Authorization")
	// 判断token 是否在黑名单中
	if oauth.TokenBlackList(token) {
		rsp.Authorization()
		return
	}
	tk,err := oauth.ParseToken(token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			if tk != nil {
				if tk.ExpiresAt + int64(setting.AppSetting.JwtRefresh) < time.Now().Unix() {
					rsp.Authorization()
					return
				}
				if newToken, er := oauth.RefreshToken(token); er == nil {
					if err := models.UpdateToken(newToken, tk.UID); err != nil {
						rsp.Authorization()
						return
					}
					oauth.AddTokenBlackList(token)
					rsp.LoginSuccess(newToken)
					return
				}
			}
		default:
			rsp.Authorization()
			return
		}
	}
	if tk.ExpiresAt - time.Now().Unix() < 7200 {
		oauth.AddTokenBlackList(token)
		if token, err = oauth.RefreshToken(token); err != nil {
			rsp.Authorization()
			return
		}
		if err := models.UpdateToken(token, tk.UID); err != nil {
			rsp.Authorization()
			return
		}
	}
	rsp.LoginSuccess(token)
}
