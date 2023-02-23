package controller

import (
	"context"
	"douyin/config"
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
	var err error

	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token获取失败",
		})
		return
	}
	user := userObj.(models.User)

	models.Db.Table("user_followers").Select("follower_id").Where("user_id = ?", user.ID).Find(&followerIDList)
	for _, i := range followerIDList {
		var follower models.User
		models.Db.First(&follower, i)
		follower.Avatar = utils.GetSignUrl(user.AvatarKey)
		follower.BackgroundImage = utils.GetSignUrl(user.BackgroundImageKey)
		follower.IsFollow = follower.GetIsFollow(user.ID)
		follower.FetchRedisData()
		if err != nil {
			c.JSON(http.StatusOK, FollowerListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "is_follow获取出错",
				},
				FollowerList: followerList,
			})
		}
		followerList = append(followerList, follower)
	}
	c.JSON(http.StatusOK, FollowerListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "关注列表获取成功",
		},
		FollowerList: followerList,
	})
}
