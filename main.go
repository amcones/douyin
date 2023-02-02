package main

import (
	"douyin/models"
	"douyin/router"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default(
		server.WithHostPorts("0.0.0.0:8080"),
		server.WithMaxRequestBodySize(20<<40), // 提高request body的容量到20MB
	)

	models.ConnDB()

	models.ConnRedis()

	router.RegisterRoute(h)

	h.Spin()
}
