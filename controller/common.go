package controller

import (
	"douyin/models"
	"douyin/utils"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,-"`
}

func FetchVideoList(videoList []models.Video, userObj interface{}) {
	for i := range videoList {
		var user models.User
		models.Db.First(&user, videoList[i].AuthorID)
		videoList[i].Author = user
		// 使用key计算得到预签名url
		videoList[i].PlayUrl = utils.GetSignUrl(videoList[i].PlayKey)
		videoList[i].CoverUrl = utils.GetSignUrl(videoList[i].CoverKey)
		// 获取Redis存储的相关数据
		videoList[i].FetchRedisData()
		// 判断是否已登录，若登录，获取点赞状态
		if userObj != nil {
			videoList[i].IsFavorite = videoList[i].GetIsFavorite(models.Db, userObj.(models.User).ID)
			videoList[i].Author.IsFollow = videoList[i].Author.GetIsFollow(userObj.(models.User).ID)
		}
	}
}
