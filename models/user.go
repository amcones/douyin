package models

import (
	"douyin/common"
	"errors"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type User struct {
	ID            int     `json:"id"`
	Name          string  `gorm:"type:varchar(32) not null;uniqueIndex" json:"name"`
	Password      string  `gorm:"type:varchar(255) not null;" json:"-"`
	FollowCount   uint    `gorm:"-" json:"follow_count"`
	FollowerCount uint    `gorm:"-" json:"follower_count"`
	IsFollow      bool    `gorm:"-" json:"is_follow"`
	Videos        []Video `gorm:"foreignKey:AuthorID" json:"-"`
	Followers     []*User `gorm:"many2many:user_followers" json:"-"`
}

func (user *User) FetchRedisData() bool {
	conn := GetRedis()
	data, err := redis.Values(conn.Do("HGETALL", common.GetRedisRelationField(user.ID)))
	if err != nil {
		return false
	}
	for i := 0; i < len(data); i += 2 {
		key := string(data[i].([]uint8))
		value := string(data[i+1].([]uint8))
		intValue, _ := strconv.Atoi(value)
		if key == common.RedisFollowerField {
			user.FollowerCount = uint(intValue)
		} else if key == common.RedisFolloweeField {
			user.FollowCount = uint(intValue)
		}
	}
	return true
}

// ValidatePassword 校验密码
func (user *User) ValidatePassword(password string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// GetUserInfoById 通过ID获取UserInfo实例
func GetUserInfoById(id interface{}) User {
	userInfo := User{}
	res := Db.First(&userInfo, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Fatal(res.Error)
	}
	return userInfo
}

// GetUserInfoByName 通过Name获取UserInfo实例
func GetUserInfoByName(name string) (User, bool) {
	userInfo := User{}
	res := Db.Where("name = ?", name).First(&userInfo)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return userInfo, false
	}
	return userInfo, true
}

// CreateUserInfo 通过Name和Password创建UserInfo实例
func CreateUserInfo(name string, password string) User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return User{}
	}
	userInfo := User{
		Name:     name,
		Password: string(hashedPassword),
	}
	Db.Create(&userInfo)
	return userInfo
}
