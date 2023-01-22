package models

import (
	"time"
)

type Video struct {
	ID            uint   `json:"id"`
	Author        User   `gorm:"-" json:"author"`
	PlayKey       string `json:"-"` //由于使用对象存储，不会有长期有效的url，在需要时通过key获取
	CoverKey      string `json:"-"`
	PlayUrl       string `gorm:"-" json:"play_url"`
	CoverUrl      string `gorm:"-" json:"cover_url"`
	FavoriteCount uint   `json:"favorite_count"`
	CommentCount  uint   `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FavoriteUsers []*User    `gorm:"many2many:user_favor_videos"`
	Comments      []*Comment `gorm:"many2many:video_comments"`
}
