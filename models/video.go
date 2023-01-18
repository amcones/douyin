package models

import "time"

type Video struct {
	ID            uint     `json:"id"`
	Author        UserInfo `gorm:"-" json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount uint     `json:"favorite_count"`
	CommentCount  uint     `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FavoriteUsers []*UserInfo `gorm:"many2many:user_favor_videos"`
	Comments      []*Comment  `gorm:"many2many:video_comments"`
}
