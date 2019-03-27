package models

import "github.com/jinzhu/gorm"

// 获取动态列表
func GetDynamicList(offset, page int, dynamic *[]Dynamic) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Select("id,uid,content,images,created_at").Order("id desc").Offset(offset).Limit(page).Find(dynamic).Error
}
func GetDynamicListOfID(page int, dynamic *[]Dynamic) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Select("id,uid,content,images,created_at").Order("id desc").Limit(page).Find(dynamic).Error
}
// 根据laseID 获取更多动态
func GetMoreDynamicListOfID(page int,lastID int64, dynamic *[]Dynamic) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Select("id,uid,content,images,created_at").Order("id desc").
		Where("id < ?",lastID).Limit(page).Find(dynamic).Error
}

// 获取动态评论列表
func GetDynamicCommentList(id int64, offset, page int, comment *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Select("id,uid,content,created_at").Where("dynamic_id = ?", id).Where("parent_id = ?", 0).
		Offset(offset).Limit(page + 1).Find(comment).Error
}

func GetDynamicCommentListOfID(dynamicID int64,page int,comment *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Select("id,uid,content,created_at").Where("dynamic_id = ?", dynamicID).
		Where("parent_id = ?", 0).Limit(page).Find(comment).Error
}
func GetMoreDynamicCommentListOfID(dynamicID,lastID int64,page int,comment *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Select("id,uid,content,created_at").Where("dynamic_id = ?", dynamicID).
		Where("parent_id = ?", 0).Where("id > ?",lastID).Limit(page).Find(comment).Error
}
// 获取评论回复列表
func GetDynamicCommentReplyList(id int64,offset, page int,reply *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Preload("Receive",func(db * gorm.DB) *gorm.DB{
		return db.Select("id,nickname")
	}).Select("id,uid,receive_id,content,created_at").Order("id desc").Where("parent_id = ?",id).
		Offset(offset).Limit(page + 1).Find(reply).Error
}
func GetDynamicCommentReplyListOfID(id int64,page int, reply *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Preload("Receive",func(db * gorm.DB) *gorm.DB {
		return db.Select("id,nickname")
	}).Select("id,uid,receive_id,content,created_at").Order("id desc").
		Where("parent_id = ?",id).Find(reply).Error
}
func GetMoreDynamicCommentReplyListOfID(id, lastID int64,page int, reply *[]DynamicComment) error {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,avatar,gender")
	}).Preload("Receive",func(db * gorm.DB) *gorm.DB {
		return db.Select("id,nickname")
	}).Select("id,uid,receive_id,content,created_at").Order("id desc").
		Where("parent_id = ?",id).Where("id < ?",lastID).Limit(page).Find(reply).Error
}

// 获取动态发布者uid
func GetDynamicUid(id int64, dynamic *Dynamic) error {
	return db.Select("uid").Where("id = ?",id).Take(dynamic).Error
}

// 获取动态评论的信息
func GetDynamicCommentInfo(id int64,reply *DynamicComment) error {
	return db.Select("uid,dynamic_id,parent_id").Where("id = ?",id).Take(reply).Error
}

