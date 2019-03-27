package oauth

import (
	"campus/helper/logging"
	"campus/helper/redis/cache"
	"campus/helper/setting"
	"campus/response/e"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	UID     int64  `json:"uid"`
	Account string `json:"account"`
	jwt.StandardClaims
}

func GenerateToken(id int64, account string) (string, error) {
	claims := Claims{
		id,
		account,
		jwt.StandardClaims{
			ExpiresAt: tokenExpired(),
			Issuer:    "CampusWall",
		},
	}
	return generate(&claims)
}
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
		if claims, ok := tokenClaims.Claims.(*Claims); ok {
			return claims, err
		}
	}
	return nil, err
}

func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = tokenExpired()
		return generate(claims)
	}
	return "", errors.New(e.GetMsg(e.AUTH_INVALID))
}

func generate(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// token expired time
func tokenExpired() int64 {
	nowTime := time.Now()
	return nowTime.Add(setting.AppSetting.JwtExpired).Unix()
}

// 验证token 是否在黑名单
func TokenBlackList(token string) bool {
	_, err := cache.Cache.Get(cache.PREFIX_TOKEN_BLACKLIST + token).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		logging.ErrorMsg("redis连接出现异常", err)
		return false
	}
	return true
}

// 把token 加入黑名单
func AddTokenBlackList(token string) {
	// 检查token 是否出于有效期
	tk,err := ParseToken(token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			if tk != nil {
				if tk.ExpiresAt + int64(setting.AppSetting.JwtRefresh) < time.Now().Unix() {
					return
				}
			}
		default:
			return
		}
	}
	if err := cache.Cache.Set(cache.PREFIX_TOKEN_BLACKLIST + token, 0,
		setting.AppSetting.JwtRefresh).Err(); err != nil {
		logging.ErrorMsg("redis连接出现异常", err)
	}
	logging.Info(token + "已加入黑名单")
}
