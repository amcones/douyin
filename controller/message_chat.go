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
	MessageList []models.Message `json:"message_list"`
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

	friendId, _ := strconv.Atoi(c.Query("to_user_id"))
	preMsgTime, _ := strconv.Atoi(c.Query("pre_msg_time"))
	// 获取聊天记录
	messageList = models.GetMessagesById(int64(userObj.(models.User).ID), int64(friendId), int64(preMsgTime))

	c.JSON(http.StatusOK, MessageChatResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "message getting succeeded",
		},
		MessageList: messageList,
	})
}
