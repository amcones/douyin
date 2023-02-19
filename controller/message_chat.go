package controller

import (
	"context"
	"douyin/config"
	"douyin/models"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
)

type MessageChatResponse struct {
	Response
	MessageList []models.Message
}

func MessageChat(_ context.Context, c *app.RequestContext) {
	var messageList []models.Message

	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, MessageChatResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "token获取失败",
			},
			MessageList: []models.Message{},
		})
		return
	}
	friendId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	messageList = models.GetMessagesById(int64(userObj.(models.User).ID), friendId)

	c.JSON(http.StatusOK, MessageChatResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "message getting succeeded",
		},
		MessageList: messageList,
	})
}
