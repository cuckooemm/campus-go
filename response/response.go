package response

import (
	"campus/helper/setting"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type response struct {
	Context *gin.Context
}

func NewResponse(c *gin.Context) *response {
	return &response{Context: c}
}

// login success with token and expiresAt
func (r *response) LoginSuccess(token string) {
	r.Context.JSON(http.StatusOK, gin.H{
		"code":       e.SUCCESS,
		"msg":        "欢迎回来",
		"token":      token,
		"expires_at": time.Now().Add(setting.AppSetting.JwtExpired).Unix() - 60, //token 过期时间 生成token需要时间 并且为了平滑过
		"invalid_at": time.Now().Add(setting.AppSetting.JwtRefresh).Unix() - 60, // token 失效时间 刷新最后期限
	})
}

// success with data 200
func (r *response) Success(message string, data interface{}) {
	r.Context.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  message,
		"data": data,
	})
}

// error
func (r *response) ErrorMsg(code int) {
	r.Context.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}

// error with param 400
func (r *response) ErrorWithMessage(message string) {
	r.Context.JSON(http.StatusOK, gin.H{
		"code": e.INVALID_PARAMS,
		"msg":  message,
		"data": nil,
	})
}

// server error 500
func (r *response) Error() {
	r.Context.JSON(http.StatusOK, gin.H{
		"code": e.SERVER_ERROR,
		"msg":  e.GetMsg(e.SERVER_ERROR),
	})
}

// forbidden 403
func (r *response) ForBidden(message string) {
	r.Context.JSON(http.StatusOK, gin.H{
		"code": http.StatusForbidden,
		"msg":  message,
	})
}

// auth
func (r *response) Authorization() {
	r.Context.JSON(http.StatusOK, gin.H{
		"code": e.AUTH_INVALID,
		"msg":  e.GetMsg(e.AUTH_INVALID),
	})
}
