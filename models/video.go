package models

import (
	"douyin/common"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
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
	FavoriteUsers []*User    `gorm:"many2many:user_favor_videos"` //用户与点赞视频是many to many
	Comments      []*Comment `json:"-"`                           // 视频与评论是has many关系
}

func (video *Video) GetIsFavorite(db *gorm.DB, userId int) bool {
	var count int64
	db.Table("user_favor_videos").Where("video_id = ? AND user_id = ?", video.ID, userId).Count(&count)
	return count != 0
}

func (video *Video) GetFavoriteCount(redisConn redis.Conn) (uint, error) {
	isExit, err := redis.Int(redisConn.Do("HEXISTS", common.RedisPrefixFavorVideo, video.ID))
	if isExit == 0 {
		return 0, err
	}
	favorCount, err := redis.Int(redisConn.Do("HGET", common.RedisPrefixFavorVideo, video.ID))
	return uint(favorCount), err
}
