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
	err = Db.AutoMigrate(&User{}, &Comment{}, &Video{})
	if err != nil {
		panic(err)
	}
}

var (
	redisConfig = config.Conf.Redis
	RedisPool   *redis.Pool
)

// ConnRedis 实例化一个连接池
func ConnRedis() {
	RedisPool = &redis.Pool{
		MaxIdle:     16,  //最初的连接数量
		MaxActive:   0,   //连接池最大连接数量,暂时不确定用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial(redisConfig.Net, redisConfig.Address)
		},
	}
}

func GetRedis() redis.Conn {
	redisConn := RedisPool.Get()
	if _, err := redisConn.Do("AUTH", redisConfig.Password); err != nil {
		panic(err)
	}
	return redisConn
}
