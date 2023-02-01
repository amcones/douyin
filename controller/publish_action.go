package controller

import (
	"bytes"
	"context"
	"douyin/models"
	"douyin/service"
	"douyin/utils"
	"image/jpeg"
	"io"
	"strconv"

	//"douyin/utils"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"log"
	"net/http"
	"os"
	"time"
)

// Publish 上传视频接口
func Publish(_ context.Context, c *app.RequestContext) {
	// 1. 根据token得到author
	token := c.FormValue("token")
	if token == nil {
		c.JSON(http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  "请登录后再操作",
			})
		return
	}
	user, exist := service.SelectToken(string(token))
	if !exist {
		c.JSON(http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  "请登录后再操作",
			})
		return
	}

	// 2. 生成文件名
	file, _ := c.FormFile("data")
	var idx int64
	var video models.Video
	models.Db.Where("author_id=?", user.ID).Find(&video).Count(&idx)
	idx += 1
	playKey := "videos/" + strconv.FormatInt(idx, 10) + ".mp4"
	coverKey := "covers/" + strconv.FormatInt(idx, 10) + ".png"

	// 3. 存入cos
	var r io.Reader
	r, _ = file.Open()
	err := utils.UploadFile(playKey, r)

	// 4. 截取视频封面并上传
	playUrl := utils.GetSignUrl(playKey)
	r, err = readFrameAsJpeg(playUrl)
	if err != nil {
		return
	}
	err = utils.UploadFile(coverKey, r)

	// 5. 保存key到数据库
	video = models.Video{
		Author:        user,
		PlayKey:       playKey,
		CoverKey:      coverKey,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         string(c.FormValue("title")),
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		FavoriteUsers: nil,
		Comments:      nil,
	}
	err = models.Db.Model(&user).Association("Videos").Append(&video)
	if err != nil {
		c.JSON(http.StatusNotImplemented,
			Response{
				StatusCode: http.StatusNotImplemented,
				StatusMsg:  fmt.Sprintf("%s", err),
			})
		log.Println(err)
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: http.StatusOK,
			StatusMsg:  "publish succeeded",
		})
	}
}

// readFrameAsJpeg 从视频中截取1帧并返回
func readFrameAsJpeg(filePath string) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(filePath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf, err
}
