package main

import (
	"douyin/models"
	"douyin/router"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	models.ConnDB()

	router.RegisterRoute(h)

	h.Spin()
}
