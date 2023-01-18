package models

type Comment struct {
	Id         int `gorm:"primarykey;Type:uint not null auto_increment;"`
	UserInfo   `gorm:"embedded;embeddedPrefix:user_"`
	Content    string `gorm:"type:varchar not null"`
	CreateDate string `gorm:"type:varchar(255) not null"`
}
