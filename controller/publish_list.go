package controller

import (
	"context"
	"douyin/config"
	"douyin/models"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type PublishListResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
}

// PublishList 根据id查询用户所有投稿视频
func PublishList(_ context.Context, c *app.RequestContext) {
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "token获取失败",
			},
			VideoList: []models.Video{},
		})
		return
	}
	id := c.Query("user_id")
	var user models.User
	models.Db.First(&user, id)
	var videoList []models.Video
	models.Db.Where("author_id=?", id).Find(&videoList)
	FetchVideoList(videoList, userObj)
	c.JSON(http.StatusOK, PublishListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "发布列表获取成功",
		},
		VideoList: videoList,
	})
}
