package models

import (
	"time"
)

type Comment struct {
	ID         int `json:"id,-"`
	UserID     int `json:"-"`
	VideoID    int `json:"-"`
	User       `gorm:"-" json:"user"`
	Content    string    `gorm:"type:varchar(999) not null" json:"content,-"`
	CreateDate time.Time `gorm:"autoCreateTime" json:"create_date,-"`
}
