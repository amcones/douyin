package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type FollowerListResponse struct {
	Response
	FollowerList []models.User `json:"user_list"`
}

func FollowerList(_ context.Context, c *app.RequestContext) {
	var followerList []models.User = nil
	var followerIDList []int64 = nil
	id := c.Query("user_id")
	models.Db.Table("user_followers").Select("follower_id").Where("user_id = ?", id).Find(&followerIDList)
	for _, i := range followerIDList {
		var user models.User
		models.Db.First(&user, i)
		user.Avatar = utils.GetSignUrl(user.AvatarKey)
		user.BackgroundImage = utils.GetSignUrl(user.BackgroundImageKey)
		followerList = append(followerList, user)
	}
	c.JSON(http.StatusOK, FollowerListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "关注列表获取成功",
		},
		FollowerList: followerList,
	})
}
