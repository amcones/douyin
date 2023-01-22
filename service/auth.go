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
func GetDouyinUserClaims(info models.UserInfo, expire time.Duration) DouyinUserClaims {
	return DouyinUserClaims{
		info.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}
}

var secret = config.Conf.JWT.SecretKey

// CreateToken 通过UserInfo创建一个Token字符串，默认1天过期。
// info: 用户信息实体
func CreateToken(info models.UserInfo) string {
	return CreateTokenWithDuration(info, 24*time.Hour)
}

// CreateTokenWithDuration 通过UserInfo创建一个Token字符串。
// info: 用户信息实体
// expire: token过期时间
func CreateTokenWithDuration(info models.UserInfo, expire time.Duration) string {
	claims := GetDouyinUserClaims(info, expire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return signedString
}

// SelectToken 根据Token获取用户信息，如果不存在，exist为false
func SelectToken(tokenString string) (user models.UserInfo, exist bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); !ok && token.Valid {
		id := claims["id"]
		userInfo := models.GetUserInfoById(id)
		return userInfo, true
	} else {
		fmt.Errorf("unexpected error when solving token: %v", err)
		return models.UserInfo{}, false
	}
}
