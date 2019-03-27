package models

import (
	"campus/helper/snowflake"
	"campus/models/input"
	"time"
)

func (c *Code) Create() error {
	c.ID = uid.SnowflakeId()
	return db.Create(c).Error
}

// 获取上次发送验证码时间  存在返回时间戳和true 反之 false
func GetLatelySendTime(s *input.SendCode) (*time.Time,bool){
	var tm Code
	if db.Select("id,created_at").Where("status = ?", false).
		Where(s.Flag+" = ?", s.Account).Last(&tm).RecordNotFound() {
		return nil, false
	}
	return &tm.CreatedAt,true
}