package models

import (
	"douyin/config"
	"fmt"
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
	err = Db.AutoMigrate(&UserInfo{}, &Comment{}, &Video{}, &Favorite{})
	if err != nil {
		panic(err)
	}
}
