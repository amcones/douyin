package controller

import (
	"context"
	"douyin/config"
	"douyin/models"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
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

	//获取token
	token := c.Query(config.IdentityKey)

	// 使用key计算得到预签名url
	for i := range videoList {
		models.Db.First(&user, videoList[i].AuthorID)
		redisVideoFavorKey := GetVideoFavorKey(strconv.Itoa(videoList[i].ID))
		videoList[i].Author = user
		videoList[i].PlayUrl = utils.GetSignUrl(videoList[i].PlayKey)
		videoList[i].CoverUrl = utils.GetSignUrl(videoList[i].CoverKey)
		videoList[i].FavoriteCount = uint(FindVideoFavorCount(redisVideoFavorKey, videoList[i].ID))
		//判断是否已登录
		if len(token) != 0 {
			userObj, _ := c.Get(config.IdentityKey)
			videoList[i].IsFavorite = FindVideoFavorStatus(redisVideoFavorKey, userObj.(models.User).ID)
		}
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
