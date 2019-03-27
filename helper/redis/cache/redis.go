package cache

import (
	"campus/helper/logging"
	"campus/helper/setting"
	"github.com/go-redis/redis"
)

var Cache *redis.Client

func Setup(){
	Cache = redis.NewClient(&redis.Options{
		Addr:setting.RedisSetting.Host,
		Password:setting.RedisSetting.Password,
		DB:setting.RedisSetting.CacheDB,
	})

	if _, err := Cache.Ping().Result(); err != nil {
		logging.ErrorMsg("redis-连接失败",err)
		return
	}else {
		logging.Info("redis-连接成功")
	}
}

func Clost()  {
	Cache.Close()
}