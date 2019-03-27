package models

import "strings"

// 更新token
func UpdateToken(token string, uid int64) error {
	return db.Table(TABLE_USER_AUTH).Where("id = ?", uid).Update("token", token).Error
}

// 判断帐号是否存在  存在为true 反之false
func IsAccountExist(flag, account string) bool {
	return !db.Select("id").Where(flag+" = ?", account).Take(&Auth{}).RecordNotFound()
}

// 判断昵称是否存在  exist true
func IsNameExist(name string) bool {
	return !db.Select("id").Where("nickname = ?", name).Take(&User{}).RecordNotFound()
}

// 判断头像是否存储在OSS  而非第三方网站头像链接
func IsAvatarOss(str string) bool {
	return strings.Contains(str,"avatar")
}
// 过滤昵称中的空格和换行
func ReplaceName(name string) string {
	name = strings.Replace(name, " ", "", -1)
	name = strings.Replace(name, "\r\n", "", -1)
	return name
}

// 根据验证码获取上次发送验证码时间  存在返回 true 反之 false
func GetLastSendTime(flag, account, num string,code *Code) bool {
	if !db.Select("id,created_at").Where(flag+" = ?", account).Where("status = ?", false).
		Where("code = ?", num).Last(code).RecordNotFound() {
		return true
	}
	return false
}

// 根据uid获取用户表信息
func GetUserProfile(id int64,user *User) error{
	if err := db.Select("id,nickname,avatar,gender,birthday,bio").Where("id = ?", id).
		Take(user).Error; err != nil {
		return err
	}
	return nil
}

// 更新用户信息
func UpdateProfile(id int64, info map[string]interface{}) error {
	return db.Table(TABLE_USER).Where("id = ?", id).Updates(info).Error
}

