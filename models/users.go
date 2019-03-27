package models

import (
	"campus/helper/encrypt"
	"campus/helper/oauth"
	"campus/helper/oss"
	"campus/helper/setting"
	"campus/helper/snowflake"
	"campus/helper/utils"
	"campus/models/input"
	"campus/response"
	"math/rand"
	"strconv"
	"time"
)

// 向数据库写入注册信息
func Register(r *input.Register,code *Code) error {
	id := uid.SnowflakeId()
	tx := db.Begin()
	if r.Flag == "email" {
		if err := tx.Create(&Auth{ID: id, Email: r.Account, Pwd: utils.PasswordEncrypt(r.Password)}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Create(&Auth{ID: id, Phone: r.Account, Pwd: utils.PasswordEncrypt(r.Password)}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Create(&User{ID: id, Nickname: r.Nickname}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(code).Update("status", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
 //  查看是否注册 已注册 为 true
func QQLoginConfirm(q *input.QQLogin,auth *Auth) bool {
	return !db.Select("id,status,token").Where("qq_id = ?",q.OpenID).Take(auth).RecordNotFound()
}

func QQLoginRegister(q *input.QQLogin, userInfo *oauth.OpenQQRespond) (int64,error) {
	id := uid.SnowflakeId()
	tx := db.Begin()
	if err := tx.Create(&Auth{ID:id,QQID:q.OpenID,Pwd:utils.PasswordEncrypt(strconv.FormatInt(time.Now().Unix(), 10))}).Error; err != nil {
		tx.Rollback()
		return 0,err
	}
	var user User
	if userInfo.Gender == "男" {
		user.Gender = 1
	} else {
		user.Gender = 0
	}
	user.Nickname = userInfo.Nickname
	if IsNameExist(userInfo.Nickname) {
		user.Nickname += strconv.Itoa(rand.Intn(100))
	}
	user.ID = id
	user.Avatar = userInfo.Avatar
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return 0,err
	}
	tx.Commit()
	return id,nil
}

// 获取用户 status
func GetUserStatus(auth *Auth) error{
	return db.Select("id,status,token").Where("id = ?",auth.ID).Take(auth).Error
}
// 修改密码
func UpdatePassword(r *input.Retrieve,codeID int64) error {
	tx := db.Begin()
	if err := tx.Table(TABLE_USER_AUTH).Where(r.Flag+"= ?", r.Account).Update("pwd", utils.PasswordEncrypt(r.Password)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table(TABLE_CODE).Where("id = ?",codeID).Update("status", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 查看用户是否存在  不存在返回false 存在返回true
func IsUserExist(l *input.Login, auth *Auth) bool {
	if db.Select("id,status,token,pwd").Where(l.Flag+" = ?",l.Account).Take(auth).RecordNotFound() {
		return false
	}
	return true
}

func GetUserInfo(id int64,info* response.UserInfo) error {
	var auth Auth
	var user User
	if err := db.Select("email,phone,qq_id,weapp_id").Where("id = ?", id).Take(&auth).Error; err != nil {
		return err
	}
	if err := db.Select("nickname,avatar,gender,school,birthday,bio").Where("id = ?", id).Take(&user).Error; err != nil {
		return err
	}
	user.AvatarToOriginal()
	info.UID = encrypt.Encode(id)
	info.Nickname = user.Nickname
	info.Email = auth.Email
	info.Phone = auth.Phone
	info.QQID = auth.QQID
	info.Avatar = user.Avatar
	info.Gender = user.Gender
	info.School = user.School
	info.Birthday = user.Birthday
	info.Bio = user.Bio
	return nil
}

//  查询 微信小程序是否第一次登录
func WeChatAppLoginConfirm(openID string,auth *Auth) bool {
	return !db.Select("id,status,token").Where("weapp_id = ?", openID).Take(auth).RecordNotFound()
}
//  微信小程序  注册账号
func WeChatAppRegister(openID string,userInfo *input.WeAppLogin) (int64,error) {
	id := uid.SnowflakeId()
	tx := db.Begin()
	if err := tx.Create(&Auth{ID:id,WeappID:openID,Pwd:utils.PasswordEncrypt(strconv.FormatInt(time.Now().Unix(), 10))}).Error; err != nil {
		tx.Rollback()
		return 0,err
	}
	var user User
	// 微信返回数据 0未知 1男 2 女
	if userInfo.Gender == 2 {
		user.Gender = 0
	}
	if userInfo.Gender == 0 {
		user.Gender = 2
	}
	user.Gender = userInfo.Gender
	// 查看昵称是否存在
	user.Nickname = userInfo.Nickname
	if IsNameExist(userInfo.Nickname) {
		user.Nickname += strconv.Itoa(rand.Intn(100))
	}
	user.ID = id
	user.Avatar = userInfo.Avatar
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return 0,err
	}
	tx.Commit()
	return id,nil
}

/**
	为头像添加域名格式
	网站存储的头像只有ID 没有前后域名和格式
	第三方获取的头像有前后缀
 */
func (user *User) AvatarToOriginal() {
	if IsAvatarOss(user.Avatar) {
		user.Avatar = setting.OSSSetting.Path + user.Avatar + oss.IMAGE_ORIGINAL
	}
}
func (user *User) AvatarToB60() {
	if IsAvatarOss(user.Avatar) {
		user.Avatar = setting.OSSSetting.Path + user.Avatar + oss.IMAGE_AVATAR_B_60
	}
}
func (user *User) GenderToString() string {
	if user.Gender == 0 {
		return "女"
	} else if user.Gender == 1 {
		return "男"
	}
	return "保密"
}
