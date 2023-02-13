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
	"sync"
)

var mutexMap = make(map[int]*sync.Mutex)

func getOrCreateMutex(userId int) *sync.Mutex {
	value, exists := mutexMap[userId]
	if !exists {
		var mutex sync.Mutex
		mutexMap[userId] = &mutex
		return &mutex
	}
	return value
}

var (
	mutex sync.Mutex
)

func RelationAction(ctx context.Context, c *app.RequestContext) {
	hlog.CtxDebugf(ctx, "[%v] 连接建立，准备抢锁", c.RemoteAddr())
	mutex.Lock()
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token获取失败",
		})
		return
	}
	fromUser := userObj.(models.User)

	//mutex := *getOrCreateMutex(fromUser.ID)
	hlog.CtxDebugf(ctx, "[%v] 拿到锁", c.RemoteAddr())

	toId := c.Query("to_user_id")
	var toUser models.User
	models.Db.First(&toUser, toId)

	actionType := c.Query("action_type")
	hlog.CtxDebugf(ctx, "[%v] action_type: %v ,mutex: %v", c.RemoteAddr(), actionType, &mutex)
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
		tx.Model(&fromUser).Updates(models.User{ID: fromUser.ID, FollowCount: fromUser.FollowCount})
		tx.Model(&toUser).Updates(models.User{ID: toUser.ID, FollowerCount: toUser.FollowerCount})
		hlog.CtxDebugf(ctx, "[%v] UpdateComplete fromUser.FollowCount = %v , toUser.FollowerCount = %v", c.RemoteAddr(), fromUser.FollowCount, toUser.FollowerCount)
		return nil
	})
	defer mutex.Unlock()
	defer hlog.CtxDebugf(ctx, "[%v] 锁释放", c.RemoteAddr())
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
