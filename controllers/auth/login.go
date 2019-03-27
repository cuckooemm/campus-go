package auth

import (
	"campus/helper/check"
	"campus/helper/logging"
	"campus/helper/oauth"
	"campus/helper/redis/cache"
	"campus/helper/utils"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {

	var entry input.Login
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

	// 查看是否登录错误过多
	loginCount, _ := cache.Cache.Get(cache.PREFIX_LOGIN_ERR + entry.Account).Int()
	if loginCount > 5 {
		logging.Warn("帐号登录错次过多",zap.String("account",entry.Account))
		rsp.ErrorMsg(e.MAX_LOGIN_FAIL)
		return
	}
	var auth models.Auth
	if !models.IsUserExist(&entry,&auth) {
		rsp.ErrorMsg(e.ERROR_ACCOUNT_UNREGISTERED)
		return
	}
	// 对比数据库中的密码
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Pwd), []byte(entry.Password)); err != nil {
		// 账号密码错则增加一次登陆错误次数
		pipe := cache.Cache.Pipeline()
		pipe.IncrBy(cache.PREFIX_LOGIN_ERR + entry.Account, 1)
		pipe.ExpireAt(cache.PREFIX_LOGIN_ERR + entry.Account, utils.GetZeroTime())
		pipe.Exec()
		pipe.Close()
		rsp.ErrorMsg(e.ERROR_LOGIN_MISMATCH)
		return
	}
	// 小于0则封禁
	if auth.Status < 0 {
		rsp.ErrorMsg(e.AUTH_FORBIDDEN)
		return
	}
	token, err := oauth.GenerateToken(auth.ID, entry.Account)
	if err != nil {
		rsp.Error()
		return
	}
	// 查询上次登录token 是否存在
	if auth.Token != "" {
		oauth.AddTokenBlackList(auth.Token)
	}
	// 更新token
	if err := models.UpdateToken(token,auth.ID); err != nil {
		rsp.Error()
		return
	}
	rsp.LoginSuccess(token)
}
