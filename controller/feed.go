package controller

import (
	"context"
	"douyin/models"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
	NextTime  int64          `json:"next_time"`
}

func Feed(_ context.Context, c *app.RequestContext) {
	var videoList []models.Video

	models.Db.Find(&videoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
