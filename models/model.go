package models

import (
	"time"
)

// 表名
const (
	TABLE_USER            = "users"
	TABLE_USER_AUTH       = "auth"
	TABLE_CODE            = "codes"
	TABLE_DYNAMIC         = "dynamics"
	TABLE_DYNAMIC_COMMENT = "dmc_comments"
	TABLE_DYNAMIC_PRAISE  = "dmc_praise"
	TABLE_FEEDBACK        = "feedback"
)

//登录表
type Auth struct {
	ID        int64  `gorm:"type:bigint;primary_key;unique_index"`
	Email     string `gorm:"type:varchar(50);index"`
	Phone     string `gorm:"type:varchar(20);index"`
	Pwd       string `gorm:"type:varchar(255);not null"`
	QQID      string `gorm:"type:varchar(255);index"`
	WeappID   string `gorm:"type:varchar(255);index"`
	Status    int8   `gorm:"type:smallint;default:1"`
	Token     string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//	用户表
type User struct {
	ID        int64  `gorm:"type:bigint;primary_key;unique_index"`
	Nickname  string `gorm:"type:varchar(20);unique"`
	Avatar    string `gorm:"type:varchar(255);default:'avatar/default'"`
	Gender    int8   `gorm:"type:smallint;default:2"`
	School    int16  `gorm:"type:smallint;index"`
	Birthday  string `gorm:"type:varchar(10);default:'1996-01-01'"`
	Bio       string `gorm:"type:varchar(120);default:'你也太懒了，介绍都没有一个。'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Code struct {
	ID        int64  `gorm:"type:bigint;primary_key"`
	Email     string `gorm:"type:varchar(50);index"`
	Phone     string `gorm:"type:varchar(20);index"`
	Code      string `gorm:"type:varchar(6);not null;index"`
	Status    bool   `gorm:"default:'false'"`
	Send      bool   `gorm:"default:'false'"`
	CreatedAt time.Time
}

// 动态表
type Dynamic struct {
	ID        int64  `gorm:"primary_key;type:bigint;unique_index"`
	User      User   `gorm:"ForeignKey:ID;AssociationForeignKey:UID"`
	UID       int64  `gorm:"type:bigint;not null;index"`
	Content   string `gorm:"type:varchar(512)" form:"content" json:"content" binding:"required,max=512"`
	Images    string `gorm:"type:varchar(512)" form:"images" json:"images" binding:"max=512"`
	School    int16  `gorm:"type:smallint;not null;index" form:"school" json:"school"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type DynamicComment struct {
	ID        int64           `gorm:"primary_key;type:bigint;unique_index" json:"id"`
	User      User            `gorm:"ForeignKey:ID;AssociationForeignKey:UID"`
	UID       int64           `gorm:"type:bigint;not null;index"`
	Dynamic   Dynamic         `gorm:"ForeignKey:ID;AssociationForeignKey:DynamicID"`
	DynamicID int64           `gorm:"type:bigint;not null;index"`
	Receive   User            `gorm:"ForeignKey:ID;AssociationForeignKey:ReceiveID"`
	ReceiveID int64           `gorm:"type:bigint;index"`
	Comment   *DynamicComment `gorm:"ForeignKey:ID;AssociationForeignKey:ParentID"`
	ParentID  int64           `gorm:"type:bigint;index"`
	Content   string          `gorm:"type:varchar(512);not null"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
type DynamicPraise struct {
	ID        int64 `gorm:"primary_key;type:bigint;unique_index"`
	User      User  `gorm:"ForeignKey:ID;AssociationForeignKey:UserID"`
	UID       int64 `gorm:"type:bigint;not null;index"`
	DynamicID int64 `gorm:"type:bigint;index" json:"dynamic_id" form:"dynamic_id"`
	CommentID int64 `gorm:"type:bigint;index" json:"comment_id" form:"comment_id"`
	CreatedAt time.Time
}
type Notice struct {
	ID        int64 `gorm:"primary_key;type:bigint;unique_index"`
}

type Feedback struct {
	ID        int64  `gorm:"type:bigint;primary_key"`
	Status    bool   `gorm:"default:'false'"` // 查看状态
	Content   string `gorm:"type:varchar(512)" form:"content" json:"content" binding:"required,max=512"`
	CreatedAt time.Time
}

func (Auth) TableName() string {
	return TABLE_USER_AUTH
}

func (User) TableName() string {
	return TABLE_USER
}
func (Code) TableName() string {
	return TABLE_CODE
}
func (Dynamic) TableName() string {
	return TABLE_DYNAMIC
}
func (DynamicComment) TableName() string {
	return TABLE_DYNAMIC_COMMENT
}
func (DynamicPraise) TableName() string {
	return TABLE_DYNAMIC_PRAISE
}
func (Feedback) TableName() string {
	return TABLE_FEEDBACK
}
