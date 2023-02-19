package router

import (
	"douyin/controller"
	"douyin/middleware"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoute(h *server.Hertz) {
	// public directory is used to serve public resources
	h.Static("/public", ".")

	apiRouter := h.Group("/douyin")

	//// basic apis
	apiRouter.GET("/feed/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.Feed}...)
	apiRouter.GET("/user/", controller.User)
	apiRouter.POST("/user/register/", controller.UserRegister)
	apiRouter.POST("/user/login/", middleware.JwtMiddleware.LoginHandler)
	apiRouter.POST("/publish/action/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.Publish}...)
	apiRouter.GET("/publish/list/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.PublishList}...)
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.FavoriteAction}...)
	apiRouter.GET("/favorite/list/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.FavoriteList}...)
	apiRouter.POST("/comment/action/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.CommentAction}...)
	apiRouter.GET("/comment/list/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.CommentList}...)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.RelationAction}...)
	//apiRouter.GET("/relation/follow/list/", controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	//apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.MessageChat}...)
	apiRouter.POST("/message/action/", []app.HandlerFunc{middleware.JwtMiddleware.MiddlewareFunc(), controller.MessageAction}...)
}
