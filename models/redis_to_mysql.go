package models

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jasonlvhit/gocron"
	"gorm.io/gorm"
	"strconv"
)

type UserFavorVideo struct {
	VideoID int
	UserID  int
}

var Scheduler = gocron.NewScheduler()

// GetStructSlice 将获得的 redis 数据转化为结构体切片形式
func GetStructSlice(sliceMap map[int][]int) []UserFavorVideo {
	var userFavorVideoSlice []UserFavorVideo
	for key, valArr := range sliceMap {
		for _, val := range valArr {
			userFavorVideoSlice = append(userFavorVideoSlice, UserFavorVideo{
				VideoID: key,
				UserID:  val,
			})
		}
	}
	return userFavorVideoSlice
}

func SaveFavorData() {
	redisConn := GetRedis()

	// 获取 redis key 列表
	keyArr, err := redis.Strings(redisConn.Do("KEYS", "favorite:video:*"))
	if err != nil {
		panic(err)
	}

	// 获取 redis set 数据
	dataMap := make(map[int][]int, len(keyArr))
	countMap := make(map[int]int, len(keyArr))
	for _, key := range keyArr {
		videoId, _ := strconv.Atoi(key[15:len(key)])
		dataMap[videoId], _ = redis.Ints(redisConn.Do("SMEMBERS", key))
		countMap[videoId] = len(dataMap[videoId]) //点赞数
	}

	// 转换数据为切片
	dataSlice := GetStructSlice(dataMap)

	// 持久化的mysql
	err = Db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("DELETE FROM user_favor_videos")
		tx.Create(&dataSlice)
		for key, val := range countMap {
			tx.Table("videos").Where("id = ? ", key).Update("favorite_count", val)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// RunTask 定时执行，间隔时间为1min
func RunTask() chan bool {
	err := Scheduler.Every(60).Seconds().Do(SaveFavorData)
	if err != nil {
		panic(err)
	}
	return Scheduler.Start()
}
