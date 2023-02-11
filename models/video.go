package models

import (
	"time"
)

type Video struct {
	ID            int    `json:"id"`
	AuthorID      int    `json:"-"`               //外键，与用户表建立has many关系
	Author        User   `gorm:"-" json:"author"` // 作者与视频是has many关系
	PlayKey       string `json:"-"`               //由于使用对象存储，不会有长期有效的url，在需要时通过key获取
	CoverKey      string `json:"-"`
	PlayUrl       string `gorm:"-" json:"play_url"`
	CoverUrl      string `gorm:"-" json:"cover_url"`
	FavoriteCount uint   `gorm:"default:0" json:"favorite_count"`
	CommentCount  uint   `gorm:"default:0" json:"comment_count"`
	IsFavorite    bool   `gorm:"default:0" json:"is_favorite"`
	Title         string `json:"title"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FavoriteUsers []*User    `gorm:"many2many:user_favor_videos"`
	Comments      []*Comment `gorm:"many2many:video_comments"`
}
