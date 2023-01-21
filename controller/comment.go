package controller

import (
	"douyin/models"
	"douyin/service"
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

// NewCommentID 获取新的ID（逐渐增加）
func NewCommentID() (id int) {
	id = 1
	return
}

// AddComment 把评论加入到数据库
func AddComment(comment *models.Comment) error {
	return nil
}

// DeleteComment 删除某个视频下的某条评论
func DeleteComment(videoID int, commentID int) {

}

// GetComments 获取某个视频下的所有评论切片
func GetComments(videoID int) []models.Comment {
	return nil
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *app.RequestContext) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := service.SelectToken(token); exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			newComment := &models.Comment{
				ID:      NewCommentID(),
				User:    user,
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
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *app.RequestContext) {
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
