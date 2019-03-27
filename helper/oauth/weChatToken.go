package oauth

import (
	"campus/helper/logging"
	"campus/helper/redis/cache"
	"campus/helper/setting"
	"github.com/go-redis/redis"
	"github.com/medivhzhan/weapp/token"
)

func AccessToken() (string,bool){
	tk,err := cache.Cache.Get("WeChatAccessToken").Result()
	if err == redis.Nil{
		tok, exp, err := token.AccessToken(setting.WeChatSetting.AppID,setting.WeChatSetting.AppSecret)
		if err != nil {
			logging.ErrorMsg("WeChat AccessToken 获取失败",err)
			return "",false
		}
		cache.Cache.Set("WeChatAccessToken",tok,exp)
		return tok,true
	}else if err != nil {
		return "",false
	}
	return tk,true
}