package models

import "campus/helper/snowflake"

func (f *Feedback)Create() error {
	f.ID = uid.SnowflakeId()
	return db.Create(f).Error
}