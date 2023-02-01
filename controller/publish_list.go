package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type PublishListResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
}

// PublishList 根据id查询用户所有投稿视频
func PublishList(_ context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	var user models.User
	models.Db.First(&user, id)
	var videoList []models.Video
	models.Db.Where("author_id=?", id).Find(&videoList)
	for i := range videoList {
		models.Db.First(&user, videoList[i].AuthorID)
		videoList[i].Author = user
		videoList[i].PlayUrl = utils.GetSignUrl(videoList[i].PlayKey)
		videoList[i].CoverUrl = utils.GetSignUrl(videoList[i].CoverKey)
	}
	c.JSON(http.StatusOK, PublishListResponse{
		Response: Response{
			StatusCode: http.StatusOK,
			StatusMsg:  "publish list succeeded",
		},
		VideoList: videoList,
	})
}
