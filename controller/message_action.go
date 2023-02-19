package controller

import (
	"context"
	"douyin/config"
	"douyin/models"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
)

func MessageAction(_ context.Context, c *app.RequestContext) {
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token获取失败"})
		return
	}

	userId := userObj.(models.User).ID
	friendId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	content := c.Query("content")

	err := models.AddMessage(int64(userId), friendId, content)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请求出错"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "message sending succeeded"})
}
