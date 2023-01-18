package controller

import (
	"douyin/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type CommentActionResponse struct {
	Response
	models.Comment
}

type CommentListResponse struct {
	Response
	comments []models.Comment
}

func SelectToken(token string) (user models.UserInfo, exist bool) {

	return
}

func NewCommentID() (id int) {
	id = 1
	return
}

// AddComment 把评论加入到数据库
func AddComment(comment *models.Comment) error {
	return nil
}

func DeleteComment(videoID int, commentID int) {

}

func GetComments(videoID int) []models.Comment {
	return nil
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := SelectToken(token); exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			newComment := &models.Comment{
				ID:         NewCommentID(),
				UserInfo:   user,
				Content:    text,
				CreateDate: time.Now(),
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
func CommentList(c *gin.Context) {
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
