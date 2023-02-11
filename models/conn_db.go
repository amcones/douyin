package models

import (
	"douyin/config"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnDB() {
	DbConf := config.Conf.DB
	var err error
	Db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=%v&loc=%s",
		DbConf.Username,
		DbConf.Password,
		DbConf.Net,
		DbConf.Addr,
		DbConf.DbName,
		DbConf.Charset,
		DbConf.ParseTime,
		DbConf.Loc,
	)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = Db.AutoMigrate(&User{}, &Comment{}, &Video{}, &Favorite{})
	if err != nil {
		panic(err)
	}
}

var redisConn redis.Conn

func ConnRedis() {
	redisConfig := config.Conf.Redis
	redisConn, _ = redis.Dial(redisConfig.Net, redisConfig.Address)
	if _, err := redisConn.Do("AUTH", redisConfig.Password); err != nil {
		panic(err)
	}
}
