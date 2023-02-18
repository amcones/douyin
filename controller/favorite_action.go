package controller

import (
	"context"
	"douyin/common"
	"douyin/config"
	"douyin/models"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gomodule/redigo/redis"
	"net/http"
)

// GetVideoFavorKey 获取 redis key
func GetVideoFavorKey(videoId string) string {
	return common.RedisPrefixFavorVideo + common.RedisKeySplit + videoId
}

// FindVideoFavorStatus 判断点赞状态
func FindVideoFavorStatus(key string, userId int) bool {
	redisConn := models.GetRedis()
	isMember, err := redis.Bool(redisConn.Do("SISMEMBER", key, userId))
	if err != nil {
		panic(err)
	}
	return isMember
}

// FindVideoFavorCount 获取点赞数
func FindVideoFavorCount(key string, videoId int) int {
	redisConn := models.GetRedis()
	count, err := redis.Int(redisConn.Do("SCARD", key))
	// redis 查找失败，从mysql获取
	if err != nil {
		models.Db.Model(&models.Video{}).Select("favorite_count").Where("id = ?", videoId).Find(count)
		fmt.Printf("ERROR: " + err.Error())
	}
	return count
}

func FavoriteAction(_ context.Context, c *app.RequestContext) {
	userObj, _ := c.Get(config.IdentityKey)

	if userObj == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token获取失败"})
		return
	}

	userId := userObj.(models.User).ID
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	videoFavorKey := GetVideoFavorKey(videoId)
	redisConn := models.GetRedis()

	//更新redis
	if actionType == "1" {
		if _, err := redisConn.Do("SADD", videoFavorKey, userId); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database operate failed"})
			panic(err)
		}
	} else if actionType == "2" {
		if _, err := redisConn.Do("SREM", videoFavorKey, userId); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database operate failed"})
			panic(err)
		}
	} else {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "action_type is valid",
			},
		})
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}
