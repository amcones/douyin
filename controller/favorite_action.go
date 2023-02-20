package controller

import (
	"context"
	"douyin/common"
	"douyin/config"
	"douyin/models"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// doRedisWatch 开启redis的点赞事务相关处理
func doRedisFavorHandle(videoId int, userId int, authorId int, offset int) error {
	redisConn := models.GetRedis()
	defer redisConn.Close()
	var innerErr error
	innerErr = redisConn.Send("HINCRBY", common.RedisPrefixFavorVideo, videoId, offset)
	if innerErr != nil {
		return innerErr
	}
	innerErr = redisConn.Send("HINCRBY", common.GetRedisUserField(userId), common.RedisFavoriteField, offset)
	if innerErr != nil {
		return innerErr
	}
	innerErr = redisConn.Send("HINCRBY", common.GetRedisUserField(authorId), common.RedisFavoritedField, offset)
	if innerErr != nil {
		return innerErr
	}
	_, innerErr = redisConn.Do("")
	if innerErr != nil {
		return innerErr
	}
	return nil
}

func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token获取失败",
		})
		return
	}

	user := userObj.(models.User)
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType := c.Query("action_type")

	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "参数错误",
		})
		return
	}

	err := models.Db.Transaction(func(tx *gorm.DB) error {
		var video models.Video
		var innerErr error
		tx.First(&video, videoId)
		dbExists := video.GetIsFavorite(tx, user.ID)
		if (dbExists && actionType == "1") || (!dbExists && actionType == "2") {
			return errors.New("重复操作")
		}
		offset := 0
		if actionType == "1" {
			offset = 1
			innerErr = tx.Model(&video).Association("FavoriteUsers").Append(&user)
		} else if actionType == "2" {
			offset = -1
			innerErr = tx.Model(&video).Association("FavoriteUsers").Delete(&user)
		}
		innerErr = doRedisFavorHandle(video.ID, user.ID, video.AuthorID, offset)
		if innerErr != nil {
			return innerErr
		}
		return nil
	})

	if err != nil {
		hlog.CtxErrorf(ctx, "点赞操作失败 videoId: %v userId: %v %v", videoId, user.ID, err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "点赞异常",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "成功",
	})
}
