package models

import (
	"errors"
	"gorm.io/gorm"
)

type UserInfo struct {
	ID            int    `json:"id"`
	Name          string `gorm:"type:varchar(32) not null;uniqueIndex" json:"name"`
	Password      string `gorm:"type:varchar(255) not null;"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func GetUserInfo(id interface{}) (UserInfo, error) {
	userInfo := UserInfo{}
	res := Db.First(&userInfo, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return userInfo, res.Error
	}
	return userInfo, nil
}
