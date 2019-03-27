package models

import "github.com/jinzhu/gorm"

// 获取用户自己的动态信息
func GetUserDynamicList(uid int64,offset, page int, dynamic *[]Dynamic) error {
	return db.Select("id,content,images,created_at").Where("uid = ?",uid).Order("id desc").
		Offset(offset).Limit(page + 1).Find(dynamic).Error
}
func GetUserDynamicListOfID(uid int64,page int,dynamic *[]Dynamic) error {
	return db.Select("id,content,images,created_at").Where("uid = ?",uid).Order("id desc").
		Limit(page).Find(dynamic).Error
}

func GetMoreUserDynamicListOfID(uid, lastID int64, page int, dynamic *[]Dynamic) error {
	return db.Select("id,content,images,created_at").Where("uid = ?",uid).Order("id desc").
		Where("id < ?",lastID).Limit(page).Find(dynamic).Error
}
// 获取动态所属用户详细信息
func GetDynamicUID(ID int64,dynamic *Dynamic) error {
	return db.Select("id,uid").Where("id = ?",ID).Take(dynamic).Error
}

// 删除动态
func DeleteDynamic(dynamic *Dynamic) error {
	return db.Delete(dynamic).Error
}
// 从根据用户ID获取评论表中的回复
func GetMessageToMe(uid int64,offset, page int,message *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Preload("Dynamic", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,content")
	}).Preload("Comment", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,content")
	}).Select("id,uid,dynamic_id,parent_id,content,created_at").Where("receive_id = ?",uid).
		Order("id desc").Offset(offset).Limit(page + 1).Find(message).Error
}

func GetUserMessage(uid int64, page int, message *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Preload("Dynamic", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,content")
	}).Preload("Comment", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,content")
	}).Select("id,uid,dynamic_id,parent_id,content,created_at").Where("receive_id = ?",uid).
		Order("id desc").Limit(page).Find(message).Error
}
func GetMoreUserMessage(uid,lastID int64, page int, message *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Preload("Dynamic", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,content")
	}).Preload("Comment", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,content")
	}).Select("id,uid,dynamic_id,parent_id,content,created_at").Where("receive_id = ?",uid).
		Order("id desc").Where("id < ?",lastID).Limit(page).Find(message).Error
}