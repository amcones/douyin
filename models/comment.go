package models

type Comment struct {
	ID         int `json:"id,-"`
	UserInfo   `gorm:"embedded;embeddedPrefix:user_" json:"userinfo"`
	Content    string `gorm:"type:varchar not null" json:"content,-"`
	CreateDate string `gorm:"type:varchar(255) not null" json:"create-date,-"`
}
