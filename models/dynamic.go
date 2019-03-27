package models

import "campus/helper/snowflake"

func (d *Dynamic) Create() error {
	d.ID = uid.SnowflakeId()
	println(d.DeletedAt)
	d.DeletedAt = nil
	return db.Create(d).Error
}
func (d *DynamicComment) Create() error {
	d.ID = uid.SnowflakeId()
	return db.Create(d).Error
}
