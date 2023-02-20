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

// doRedisFollowHandle 开启redis的关注事务相关处理
func doRedisFollowHandle(ctx context.Context, fromUser *models.User, toUser *models.User, offset int) error {
	redisConn := models.GetRedis()
	defer redisConn.Close()
	var innerErr error
	innerErr = redisConn.Send("HINCRBY", common.GetRedisRelationField(fromUser.ID), common.RedisFolloweeField, offset)
	if innerErr != nil {
		return innerErr
	}
	innerErr = redisConn.Send("HINCRBY", common.GetRedisRelationField(toUser.ID), common.RedisFollowerField, offset)
	if innerErr != nil {
		return innerErr
	}
	_, innerErr = redisConn.Do("")
	if innerErr != nil {
		return innerErr
	}
	return nil
}

func RelationAction(ctx context.Context, c *app.RequestContext) {
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token获取失败",
		})
		return
	}
	fromUser := userObj.(models.User)

	toId, _ := strconv.Atoi(c.Query("to_user_id"))
	if toId == fromUser.ID {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "不能关注自己",
		})
		return
	}
	var toUser models.User
	models.Db.First(&toUser, toId)

	actionType := c.Query("action_type")

	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "参数错误",
		})
		return
	}

	err := models.Db.Transaction(func(tx *gorm.DB) error {
		var innerErr error
		follower := tx.Model(&toUser).Association("Followers")
		var tmp = models.User{ID: fromUser.ID}
		innerErr = follower.Find(&tmp)
		if innerErr != nil {
			return innerErr
		}
		dbExists := tmp.Name != ""
		if (dbExists && actionType == "1") || (!dbExists && actionType == "2") {
			return errors.New("重复操作")
		}
		offset := 0
		if actionType == "1" {
			offset = 1
			innerErr = follower.Append(&fromUser)
		} else if actionType == "2" {
			offset = -1
			innerErr = follower.Delete(&fromUser)
		}
		innerErr = doRedisFollowHandle(ctx, &fromUser, &toUser, offset)
		if innerErr != nil {
			return innerErr
		}
		return nil
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "创建关注失败 from: %v to: %v %v", fromUser.ID, toUser.ID, err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "关注异常",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "成功",
	})
}
