package router

import (
	"douyin/common"
	"douyin/controller"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoute(h *server.Hertz) {
	// public directory is used to serve public resources
	h.Static("/public", ".")

	apiRouter := h.Group("/douyin")

	//// basic apis
	apiRouter.GET("/feed/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.Feed}...)
	apiRouter.GET("/user/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.User}...)
	apiRouter.POST("/user/register/", controller.UserRegister)
	apiRouter.POST("/user/login/", common.JwtMiddleware.LoginHandler)
	apiRouter.POST("/publish/action/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.Publish}...)
	apiRouter.GET("/publish/list/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.PublishList}...)
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.FavoriteAction}...)
	apiRouter.GET("/favorite/list/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.FavoriteList}...)
	apiRouter.POST("/comment/action/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.CommentAction}...)
	apiRouter.GET("/comment/list/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.CommentList}...)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.RelationAction}...)
	apiRouter.GET("/relation/follow/list/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.FollowList}...)
	apiRouter.GET("/relation/follower/list/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.FollowerList}...)
	apiRouter.GET("/relation/friend/list/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.FriendList}...)
	apiRouter.GET("/message/chat/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.MessageChat}...)
	apiRouter.POST("/message/action/", []app.HandlerFunc{common.JwtMiddleware.MiddlewareFunc(), controller.MessageAction}...)
}
