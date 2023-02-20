package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type FavoriteListResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
}

func FavoriteList(_ context.Context, c *app.RequestContext) {
	var videoList []models.Video = nil
	var redisConn = models.GetRedis()
	var err error

	var videoIdList []int

	id := c.Query("user_id")

	models.Db.Table("user_favor_videos").Select("video_id").Where("user_id = ?", id).Find(&videoIdList)

	if len(videoIdList) > 0 {
		models.Db.Where(videoIdList).Find(&videoList)
	}

	var user models.User

	for i := range videoList {
		models.Db.First(&user, videoList[i].AuthorID)
		videoList[i].Author = user
		videoList[i].PlayUrl = utils.GetSignUrl(videoList[i].PlayKey)
		videoList[i].CoverUrl = utils.GetSignUrl(videoList[i].CoverKey)
		if videoList[i].FavoriteCount, err = videoList[i].GetFavoriteCount(redisConn); err != nil {
			c.JSON(http.StatusOK, FavoriteListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "喜欢列表请求出错",
				},
				VideoList: nil,
			})
		}
	}

	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "favorite list get succeeded",
		},
		VideoList: videoList,
	})
}
