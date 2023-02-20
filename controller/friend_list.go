package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
)

type FriendListResponse struct {
	Response
	FriendList []FriendUser `json:"user_list"`
}

type FriendUser struct {
	models.User
	Message string `json:"message"`
	MsgType int64  `json:"msgType"`
}

func FriendList(_ context.Context, c *app.RequestContext) {
	var friendList []FriendUser = nil
	var friendIDList []int64 = nil

	id, _ := strconv.Atoi(c.Query("user_id"))

	models.Db.Table("user_followers u1").Select("u1.follower_id").Joins("join user_followers u2 on u1.user_id = u2.follower_id and u2.user_id = u1.follower_id").Where("u1.user_id = ?", id).Find(&friendIDList)
	for _, i := range friendIDList {
		var user models.User
		var message models.Message
		models.Db.First(&user, i)
		user.Avatar = utils.GetSignUrl(user.AvatarKey)
		user.BackgroundImage = utils.GetSignUrl(user.BackgroundImageKey)
		message.GetLatestMessagesById(int64(id), i)
		friendUser := FriendUser{
			User:    user,
			Message: message.Content,
		}
		// 0-当前用户接收的消息；1-当前用户发送的消息
		if message.ToUserID == int64(user.ID) {
			friendUser.MsgType = 1
		} else {
			friendUser.MsgType = 0
		}
		friendList = append(friendList, friendUser)
	}

	c.JSON(http.StatusOK, FriendListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "好友列表获取成功",
		},
		FriendList: friendList,
	})
}
