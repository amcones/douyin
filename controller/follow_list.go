package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type FollowListResponse struct {
	Response
	FollowList []models.User `json:"user_list"`
}

func FollowList(_ context.Context, c *app.RequestContext) {
	var followList []models.User = nil
	var followIDList []int64 = nil

	id := c.Query("user_id")
	models.Db.Table("user_followers").Select("user_id").Where("follower_id = ?", id).Find(&followIDList)
	for _, i := range followIDList {
		var user models.User
		models.Db.First(&user, i)
		user.Avatar = utils.GetSignUrl(user.AvatarKey)
		user.BackgroundImage = utils.GetSignUrl(user.BackgroundImageKey)
		user.IsFollow = true
		user.FetchRedisData()
		followList = append(followList, user)
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "关注列表获取成功",
		},
		FollowList: followList,
	})
}
