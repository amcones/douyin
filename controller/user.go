package controller

import (
	"context"
	"douyin/models"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
	"net/http"
)

type UserResponse struct {
	Response
	User models.User
}

func User(_ context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	var user models.User
	result := models.Db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{http.StatusOK, "cannot find user"},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{http.StatusOK, "succeeded"},
			User:     user,
		})
	}
}
