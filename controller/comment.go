package controller

import (
	"context"
	"douyin/common"
	"douyin/config"
	"douyin/models"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
)

type CommentActionResponse struct {
	Response
	Comment any `json:"comment,-"`
}

type CommentListResponse struct {
	Response
	Comments []models.Comment `json:"comment_list,-"`
}

func AddCommentRedis(videoID int, offset int) error {
	redisConn := models.GetRedis()
	defer redisConn.Close()
	_, err := redisConn.Do("HINCRBY", common.RedisPrefixCommentVideo, videoID, offset)
	return err
}

// AddComment 把评论加入到数据库
func AddComment(comment *models.Comment) bool {
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		err := AddCommentRedis(comment.VideoID, 1)
		if err != nil {
			return err
		}
		tx.Create(comment)
		return nil
	})
	if err != nil {
		hlog.Errorf("评论插入数据库失败 %v", err)
		return false
	}
	return true
}

// DeleteComment 删除某个视频下的某条评论
func DeleteComment(commentID int, videoID int) bool {
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		err := AddCommentRedis(videoID, -1)
		if err != nil {
			return err
		}
		tx.Delete(models.Comment{}, commentID)
		return nil
	})
	if err != nil {
		hlog.Errorf("评论删除数据库失败 %v", err)
		return false
	}
	return true
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(_ context.Context, c *app.RequestContext) {
	actionType := c.Query("action_type")
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token获取失败"})
		return
	}
	videoID, _ := strconv.Atoi(c.Query("video_id"))
	if actionType == "1" {
		text := c.Query("comment_text")
		newComment := &models.Comment{
			User:    userObj.(models.User),
			Content: text,
			UserID:  userObj.(models.User).ID,
			VideoID: videoID,
		}
		result := AddComment(newComment)
		if result {
			c.JSON(http.StatusOK, CommentActionResponse{Response{StatusCode: 0}, *newComment})
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{Response{StatusCode: 1}, nil})
		}
		return
	}
	if actionType == "2" {
		commentID, err := strconv.Atoi(c.Query("comment_id"))
		if err != nil {
			fmt.Println(err)
			return
		}
		result := DeleteComment(commentID, videoID)
		if result {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0, StatusMsg: "删除成功"},
				Comment:  nil,
			})
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "删除失败"},
				Comment:  nil,
			})
		}
	}
}

// CommentList all videos have same demo comment list
func CommentList(_ context.Context, c *app.RequestContext) {
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token获取失败"})
		return
	}
	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	var comments []models.Comment
	models.Db.Where("video_id = ?", videoID).Find(&comments)
	if len(comments) == 0 {
		comments = nil
	}
	for i := range comments {
		models.Db.First(&comments[i].User, comments[i].UserID)
	}
	sort.Slice(comments, func(i int, j int) bool {
		return comments[i].CreateDate.After(comments[j].CreateDate)
	})
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get comment list succeeded",
		},
		Comments: comments,
	})
}
