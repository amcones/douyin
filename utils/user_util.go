package utils

import (
	"douyin/models"
	"github.com/hertz-contrib/jwt"
	"time"
)

func CheckUserParameter(username string, password string) (bool, string) {
	if len(username) > 32 || username == "" {
		return false, "用户名不合法"
	}
	if len(password) > 32 || password == "" {
		return false, "密码不合法"
	}
	return true, ""
}

func CreateUserToken(user models.User) (string, time.Time, error) {
	var JwtMiddleware *jwt.HertzJWTMiddleware
	return JwtMiddleware.TokenGenerator(user)
}
