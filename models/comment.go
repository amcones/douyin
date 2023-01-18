package models

type Comment struct {
	ID         int
	UserInfo   `gorm:"embedded;embeddedPrefix:user_"`
	Content    string `gorm:"type:varchar(999) not null"`
	CreateDate string `gorm:"type:varchar(255) not null"`
}
