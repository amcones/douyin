package models

import (
	"douyin/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnDB() {
	DbConf := config.Conf.DB
	var err error
	Db, _ := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=%v&loc=%s",
		DbConf.Username,
		DbConf.Password,
		DbConf.Net,
		DbConf.Addr,
		DbConf.DbName,
		DbConf.Charset,
		DbConf.ParseTime,
		DbConf.Loc,
	)), &gorm.Config{})
	err = Db.AutoMigrate() //TODO:在这里添加需要创建的表结构来自动创建表
	if err != nil {
		panic(err)
	}
}
