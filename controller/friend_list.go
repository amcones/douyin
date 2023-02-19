package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type FriendListResponse struct {
	Response
	FriendList []models.User `json:"user_list"`
}

func FriendList(_ context.Context, c *app.RequestContext) {
	var friendList []models.User = nil
	var friendIDList []int64 = nil
	id := c.Query("user_id")
	models.Db.Table("user_followers u1").Select("u1.follower_id").Joins("join user_followers u2 on u1.user_id = u2.follower_id and u1.user_id = u2.follower_id").Where("u1.user_id = ?", id).Find(&friendIDList)
	for _, i := range friendIDList {
		var user models.User
		models.Db.First(&user, i)
		user.Avatar = utils.GetSignUrl(user.AvatarKey)
		user.BackgroundImage = utils.GetSignUrl(user.BackgroundImageKey)
		friendList = append(friendList, user)
	}
	c.JSON(http.StatusOK, FriendListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "关注列表获取成功",
		},
		FriendList: friendList,
	})
}
