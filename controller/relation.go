package controller

import (
	"context"
	"douyin/config"
	"douyin/models"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

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
		hlog.CtxDebugf(ctx, "[%v] fromUser.FollowCount = %v , toUser.FollowerCount = %v", c.RemoteAddr(), fromUser.FollowCount, toUser.FollowerCount)
		if actionType == "1" {
			innerErr = follower.Append(&fromUser)
			fromUser.FollowCount += 1
			toUser.FollowerCount += 1
		} else if actionType == "2" {
			innerErr = follower.Delete(&fromUser)
			fromUser.FollowCount -= 1
			toUser.FollowerCount -= 1
		}
		if innerErr != nil {
			return innerErr
		}
		tx.Model(&fromUser).Update("FollowCount", fromUser.FollowCount)
		tx.Model(&toUser).Update("FollowerCount", toUser.FollowerCount)
		hlog.CtxDebugf(ctx, "[%v] UpdateComplete fromUser.FollowCount = %v , toUser.FollowerCount = %v", c.RemoteAddr(), fromUser.FollowCount, toUser.FollowerCount)
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
