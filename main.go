package main

import (
	"campus/helper/check"
	"campus/helper/encrypt"
	"campus/helper/logging"
	"campus/helper/oss"
	"campus/helper/redis/cache"
	"campus/helper/sensitiveWord"
	"campus/helper/setting"
	"campus/models"
	"campus/routers"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func main() {
	//loading setting
	setting.Setup()
	// hashids
	encrypt.Setup()
	// loading log
	logging.Setup()
	defer logging.Close()

	// loading database
	models.Setup()
	defer models.CloseDB()
	// loading redis
	cache.Setup()
	defer cache.Clost()
	// loading sensitiveWord
	sensitiveWord.Setup()
	// Oss loading
	oss.Setup()

	// config server
	router := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           setting.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	// 升级gin自带的validator  并初始化
	binding.Validator = new(check.DefaultValidator)
	check.ValidatorSetup()
	//err := server.ListenAndServeTLS(setting.ServerSetting.HttpsCert,setting.ServerSetting.HttpsKey)
	err := server.ListenAndServe()
	if err != nil{
		logging.ErrorMsg("服务器启动失败",err)
	}
}