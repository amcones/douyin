package controller

import (
	"bytes"
	"context"
	"douyin/config"
	"douyin/models"
	"douyin/utils"
	"image/jpeg"
	"io"
	"mime/multipart"
	"path/filepath"

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

var tmpFolder = filepath.Join(os.TempDir(), "douyin")

func publishFail(reason string, c *app.RequestContext) {
	c.JSON(http.StatusOK, PublishListResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  reason,
		},
		VideoList: []models.Video{},
	})
}

func saveTmpPlay(path string, fileHeader *multipart.FileHeader) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}

// Publish 上传视频接口
func Publish(_ context.Context, c *app.RequestContext) {
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
	user := userObj.(models.User)
	log.Printf("用户准备上传视频 ID: %v\n", user.ID)

	// 2. 生成文件名  新：根据文件md5生成文件名，以散列值的前2位作为前缀，其他的作为后缀
	file, _ := c.FormFile("data")
	var video models.Video

	// 3. 存入cos
	var r io.Reader
	r, _ = file.Open()
	var buf bytes.Buffer
	io.TeeReader(r, &buf)
	videoMD5, err := utils.FileMD5(&buf)
	if err != nil {
		publishFail("视频散列失败", c)
		return
	}
	r, _ = file.Open()
	playKey := "videos/" + utils.GetStoragePath(videoMD5) + ".mp4"
	err = utils.UploadFile(playKey, r)
	if err != nil {
		publishFail("视频上传失败", c)
		return
	}

	// 4. 截取视频封面并上传
	// 4.1 保存视频到本地临时文件
	os.MkdirAll(tmpFolder, os.ModePerm)
	tmpPath := filepath.Join(tmpFolder, videoMD5)
	err = saveTmpPlay(tmpPath, file)
	if err != nil {
		log.Printf("保存文件失败 %s %v\n", tmpPath, err)
		publishFail("临时文件保存失败", c)
		return
	}
	// 4.2 本地读取封面
	r, err = readFrameAsJpeg(tmpPath)
	if err != nil {
		publishFail("封面读取失败", c)
		return
	}
	// 4.3 计算封面hash
	io.TeeReader(r, &buf)
	coverMD5, err := utils.FileMD5(&buf)
	if err != nil {
		publishFail("封面散列失败", c)
		return
	}
	// 4.4 保存封面
	err = os.Remove(tmpPath)
	if err != nil {
		publishFail("临时文件删除失败", c)
		return
	}
	// 4.5 上传封面
	coverKey := "covers/" + utils.GetStoragePath(coverMD5) + ".png"
	err = utils.UploadFile(coverKey, r)
	if err != nil {
		publishFail("封面上传失败", c)
		return
	}
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
		c.JSON(http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%s", err),
			})
		log.Println(err)
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "发布成功",
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
