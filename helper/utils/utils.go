package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

/**
	生成6位随机验证码
 */
func Random6Code() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(999999))
}

/**
	密码加密
 */
func PasswordEncrypt(password string) string {
	psc, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(psc)
}

// 获取零点时间
func GetZeroTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.AddDate(0, 0, 1)
}