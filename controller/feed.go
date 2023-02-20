package controller

import (
	"context"
	"douyin/config"
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

func GetVideoInfo(videoList []models.Video, token string, c *app.RequestContext) error {
	var err error
	var user models.User
	var userObj interface{}
	redisConn := models.GetRedis()
	if len(token) != 0 {
		userObj, _ = c.Get(config.IdentityKey)
	}

	for i := range videoList {
		models.Db.First(&user, videoList[i].AuthorID)
		videoList[i].Author = user
		// 使用key计算得到预签名url
		videoList[i].PlayUrl = utils.GetSignUrl(videoList[i].PlayKey)
		videoList[i].CoverUrl = utils.GetSignUrl(videoList[i].CoverKey)
		// 获取点赞数
		if videoList[i].FavoriteCount, err = videoList[i].GetFavoriteCount(redisConn); err != nil {
			return err
		}
		// 判断是否已登录，若登录，获取点赞状态
		if userObj != nil {
			videoList[i].IsFavorite = videoList[i].GetIsFavorite(models.Db, userObj.(models.User).ID)
		}
	}
	return nil
}

func Feed(_ context.Context, c *app.RequestContext) {
	var videoList []models.Video = nil

	// 按照投稿时间降序，一次最多30条
	models.Db.Order("updated_at desc").Limit(30).Find(&videoList)
	nextTime := videoList[len(videoList)-1].UpdatedAt.Unix()

	//获取token，完善视频信息
	token := c.Query(config.IdentityKey)
	err := GetVideoInfo(videoList, token, c)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "操作出错",
			},
			VideoList: nil,
		})
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
