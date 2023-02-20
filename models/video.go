package models

import (
	"douyin/common"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	FavoriteCount uint   `gorm:"-" json:"favorite_count"`
	CommentCount  uint   `gorm:"-" json:"comment_count"`
	IsFavorite    bool   `gorm:"-" json:"is_favorite"`
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

func (video *Video) FetchRedisData() {
	redisConn := GetRedis()
	defer redisConn.Close()
	var err error
	count, err := video.getRedisCount(redisConn, common.RedisPrefixFavorVideo)
	if err != nil {
		hlog.Errorf("Video FetchRedisData %v 失败 %v ", common.RedisPrefixFavorVideo, err)
	}
	video.FavoriteCount = count
	count, _ = video.getRedisCount(redisConn, common.RedisPrefixCommentVideo)
	if err != nil {
		hlog.Errorf("Video FetchRedisData %v 失败 %v ", common.RedisPrefixCommentVideo, err)
	}
	video.CommentCount = count
}

func (video *Video) getRedisCount(redisConn redis.Conn, redisFieldName string) (uint, error) {
	isExit, err := redis.Int(redisConn.Do("HEXISTS", redisFieldName, video.ID))
	if isExit == 0 {
		return 0, err
	}
	favorCount, err := redis.Int(redisConn.Do("HGET", redisFieldName, video.ID))
	return uint(favorCount), err
}
