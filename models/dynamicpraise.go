package models

import (
	"campus/helper/snowflake"
)



func (t *DynamicPraise)Create() error {
	t.ID = uid.SnowflakeId()
	return db.Create(t).Error
}
// 根据用户ID 获取是否点过赞  无 true 有 false
func (t *DynamicPraise)IsUserPraiseD() bool {
	return db.Where("dynamic_id = ?",t.DynamicID).Where("uid = ?",t.UID).Take(t).RecordNotFound()
}

func (t *DynamicPraise)IsUserPraiseC() bool {
	return db.Where("comment_id = ?",t.CommentID).Where("uid = ?",t.UID).Take(t).RecordNotFound()
}