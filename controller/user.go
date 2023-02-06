package controller

import (
	"context"
	"douyin/models"
	"douyin/service"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
	"net/http"
)

type UserResponse struct {
	Response
	User models.User
}

type UserRegisterResponse struct {
	Response
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

func User(_ context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	var user models.User
	result := models.Db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{1, "cannot find user"},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{0, "succeeded"},
			User:     user,
		})
	}
}

func checkUserParameter(username string, password string) (bool, string) {
	if len(username) > 32 || username == "" {
		return false, "username not valid"
	}
	if len(password) > 32 || password == "" {
		return false, "password not valid"
	}
	return true, ""
}

func UserLogin(_ context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	argValid, reason := checkUserParameter(username, password)
	if !argValid {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{1, reason},
		})
		return
	}
	user, exist := models.GetUserInfoByName(username)
	if !exist {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{1, "username or password not match"},
		})
		return
	}
	passwordValid := user.ValidatePassword(password)
	if !passwordValid {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{1, "username or password not match"},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: Response{0, "ok"},
		Token:    service.CreateToken(user),
		UserID:   user.ID,
	})
}

func UserRegister(_ context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	argValid, reason := checkUserParameter(username, password)
	if !argValid {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{1, reason},
		})
		return
	}
	_, exist := models.GetUserInfoByName(username)
	if exist {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{1, "用户名已存在"},
		})
		return
	}
	user := models.CreateUserInfo(username, password)
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: Response{0, "ok"},
		Token:    service.CreateToken(user),
		UserID:   user.ID,
	})
}
