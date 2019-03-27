package models

import (
	"campus/helper/logging"
	"campus/helper/setting"
	"fmt"
	"github.com/jinzhu/gorm"
)
import _ "github.com/jinzhu/gorm/dialects/postgres"

var db *gorm.DB

func Setup() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.DBName,
		setting.DatabaseSetting.Password))

	if err != nil {
		logging.ErrorMsg("数据库连接失败", err)
		return
	}
	db.LogMode(setting.DatabaseSetting.LogModel)
	db.AutoMigrate(&Auth{}, &Code{}, &User{}, &Dynamic{}, &DynamicComment{}, &DynamicPraise{},&Feedback{})
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	// 移除 数据库插入后的查询功能
	gorm.DefaultCallback.Create().Remove("gorm:force_reload_after_create")
	logging.Info("数据库连接成功")
}
func CloseDB() {
	defer db.Close()
}
