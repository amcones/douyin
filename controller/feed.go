package controller

import (
	"context"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type FeedResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
	NextTime  int64          `json:"next_time,-"`
}

func Feed(_ context.Context, c *app.RequestContext) {
	var videoList []models.Video

	// 按照投稿时间降序，一次最多30条
	models.Db.Order("updated_at desc").Limit(30).Find(&videoList)
	var user models.User
	var nextTime int64
	// 使用key计算得到预签名url
	for i := range videoList {
		models.Db.First(&user, videoList[i].AuthorID)
		videoList[i].Author = user
		videoList[i].PlayUrl = utils.GetSignUrl(videoList[i].PlayKey)
		videoList[i].CoverUrl = utils.GetSignUrl(videoList[i].CoverKey)
		// 视频列表按时间倒序，最后一个视频是时间最早的
		nextTime = videoList[i].UpdatedAt.Unix()
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "feed getting succeeded",
		},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
