package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInfo struct {
	ID            int    `json:"id"`
	Name          string `gorm:"type:varchar(32) not null;uniqueIndex" json:"name"`
	Password      string `gorm:"type:varchar(255) not null;"`
	FollowCount   int    `gorm:"default:0" json:"follow_count"`
	FollowerCount int    `gorm:"default:0" json:"follower_count"`
}

// FrontedUserInfo IsFollow字段不应该出现在数据库中，固分离
type FrontedUserInfo struct {
	UserInfo
	IsFollow bool `json:"is_follow"`
}

// ValidatePassword 校验密码
func (userInfo *UserInfo) ValidatePassword(password string) bool {
	hashedStoredPassword := []byte(userInfo.Password)
	passwordToValidate := []byte(password)
	err := bcrypt.CompareHashAndPassword(hashedStoredPassword, passwordToValidate)
	return err != nil
}

// GetUserInfoById 通过ID获取UserInfo实例
func GetUserInfoById(id interface{}) UserInfo {
	userInfo := UserInfo{}
	res := Db.First(&userInfo, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		panic(res.Error)
	}
	return userInfo
}

// GetUserInfoByName 通过Name获取UserInfo实例
func GetUserInfoByName(name string) UserInfo {
	userInfo := UserInfo{}
	res := Db.Where("name = ?", name).First(&userInfo)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		panic(res.Error)
	}
	return userInfo
}

// CreateUserInfo 通过Name和Password创建UserInfo实例
func CreateUserInfo(name string, password string) UserInfo {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	userInfo := UserInfo{
		Name:     name,
		Password: string(hashedPassword),
	}
	Db.Create(&userInfo)
	return userInfo
}
