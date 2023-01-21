package service

import (
	"douyin/config"
	"douyin/models"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type DouyinUserClaims struct {
	UserID int
	jwt.RegisteredClaims
}

// GetDouyinUserClaims 获取jwt Claims实例
func GetDouyinUserClaims(u models.User, expire time.Duration) DouyinUserClaims {
	return DouyinUserClaims{
		u.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}
}

var secret = config.CommonConf.JWT.SecretKey

// CreateToken 通过UserInfo创建一个Token字符串，默认1天过期。
// info: 用户信息实体
func CreateToken(u models.User) string {
	return CreateTokenWithDuration(u, 24*time.Hour)
}

// CreateTokenWithDuration 通过UserInfo创建一个Token字符串。
// info: 用户信息实体
// expire: token过期时间
func CreateTokenWithDuration(u models.User, expire time.Duration) string {
	claims := GetDouyinUserClaims(u, expire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return signedString
}

// SelectToken 根据Token获取用户信息，如果不存在，exist为false
func SelectToken(tokenString string) (u models.User, exist bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); !ok && token.Valid {
		id := claims["id"]
		userInfo := u.GetUserInfoById(id)
		return userInfo, true
	} else {
		log.Fatal(fmt.Errorf("unexpected error when solving token: %v", err))
		return models.User{}, false
	}
}
