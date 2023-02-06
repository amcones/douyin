package middleware

import (
	"context"
	"douyin/config"
	"douyin/controller"
	"douyin/models"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"log"
	"net/http"
	"time"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
)

func checkUserParameter(username string, password string) (bool, string) {
	if len(username) > 32 || username == "" {
		return false, "用户名不合法"
	}
	if len(password) > 32 || password == "" {
		return false, "密码不合法"
	}
	return true, ""
}

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "douyin",
		Key:           []byte(config.Conf.JWT.SecretKey),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt, form: token",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			userID, _ := c.Get(config.UserIDKey)
			log.Printf("LoginResponse %v %v %v\n", code, token, expire)
			c.JSON(code, controller.UserRegisterResponse{
				Response: controller.Response{
					StatusCode: 0,
					StatusMsg:  "登录成功",
				},
				Token:  token,
				UserID: userID.(int),
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			username := c.Query("username")
			password := c.Query("password")
			argValid, reason := checkUserParameter(username, password)
			if !argValid {
				return nil, errors.New(reason)
			}
			user, exist := models.GetUserInfoByName(username)
			if !exist {
				return nil, errors.New("用户名或密码不匹配")
			}
			passwordValid := user.ValidatePassword(password)
			if !passwordValid {
				return nil, errors.New("用户名或密码不匹配")
			}
			c.Set(config.UserIDKey, user.ID)
			log.Printf("Authenticator\n")
			return user, nil
		},
		IdentityKey: config.IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return models.GetUserInfoById(int(claims[config.UserIDKey].(float64)))
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(models.User); ok {
				return jwt.MapClaims{
					config.UserIDKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK,
				controller.Response{
					StatusCode: int32(code),
					StatusMsg:  message,
				})
		},
	})
	if err != nil {
		panic(err)
	}
}
