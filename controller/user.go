package controller

import (
	"context"
	"douyin/common"
	"douyin/config"
	"douyin/models"
	"douyin/utils"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
	"net/http"
)

type UserResponse struct {
	Response
	User models.User `json:"user"`
}

type UserRegisterResponse struct {
	Response
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

func User(_ context.Context, c *app.RequestContext) {
	id := c.Query("user_id")
	token := c.Query(config.IdentityKey)
	var user models.User
	result := models.Db.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{0, "成功"},
			User:     user,
		})
		return
	}
	user.Avatar = utils.GetSignUrl(user.AvatarKey)
	user.BackgroundImage = utils.GetSignUrl(user.BackgroundImageKey)
	var userObj interface{}
	if len(token) != 0 {
		userObj, _ = c.Get(config.IdentityKey)
		user.IsFollow = user.GetIsFollow(userObj.(models.User).ID)
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{1, "未找到用户"},
			User:     user,
		})
		return
	}
	user.FetchRedisData()
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{0, "成功"},
		User:     user,
	})
}

func UserRegister(_ context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	argValid, reason := utils.CheckUserParameter(username, password)
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
	token, _, _ := common.JwtMiddleware.TokenGenerator(user)
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: Response{0, "ok"},
		Token:    token,
		UserID:   user.ID,
	})
}
