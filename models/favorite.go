package models

import "time"

type Favorite struct {
	ID        int
	UserId    int
	VideoId   int
	Status    int `gorm:"default:1;comment:'1-已赞；2-取消赞'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
