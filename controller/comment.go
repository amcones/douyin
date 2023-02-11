package controller

import (
	"context"
	"douyin/config"
	"douyin/models"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"sort"
	"strconv"
)

type CommentActionResponse struct {
	Response
	models.Comment
}

type CommentListResponse struct {
	Response
	comments []models.Comment
}

// AddComment 把评论加入到数据库
func AddComment(comment *models.Comment) {
	models.Db.Create(comment)
}

// DeleteComment 删除某个视频下的某条评论
func DeleteComment(videoID int, commentID int) {
	// 删除评论
	models.Db.Delete(models.Comment{}, commentID)
	// 删除连接表中的id
	models.Db.Table("video_comments").Where("video_id = ? && comment_id = ?", videoID, commentID).Delete(nil)
}

// GetComments 获取某个视频下的所有评论切片
func GetComments(videoID int) []models.Comment {
	comments := make([]models.Comment, 0)
	commentIDs := make([]int, 0)
	models.Db.Table("video_comments").Where("video_id = ?", videoID).Find(commentIDs)
	for _, id := range commentIDs {
		var comment models.Comment
		models.Db.First(&comment, id)
		comments = append(comments, comment)
	}
	return comments
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(_ context.Context, c *app.RequestContext) {
	actionType := c.Query("action_type")
	userObj, _ := c.Get(config.IdentityKey)
	if userObj == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token获取失败"})
		return
	}

	if actionType == "1" {
		text := c.Query("comment_text")
		newComment := &models.Comment{
			User:    userObj.(models.User),
			Content: text,
		}
		AddComment(newComment)
		c.JSON(http.StatusOK, CommentActionResponse{Response{StatusCode: 0}, *newComment})
		return
	}
	if actionType == "2" {
		videoID, err := strconv.Atoi(c.Query("video_id"))
		if err != nil {
			fmt.Println(err)
			return
		}
		commentID, err := strconv.Atoi(c.Query("comment_id"))
		if err != nil {
			fmt.Println(err)
			return
		}
		DeleteComment(videoID, commentID)
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
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
	comments := GetComments(videoID)
	sort.Slice(comments, func(i int, j int) bool {
		return comments[i].CreateDate.After(comments[j].CreateDate)
	})
	c.JSON(http.StatusOK, CommentListResponse{Response{StatusCode: 0}, comments})
}
