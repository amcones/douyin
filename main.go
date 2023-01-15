package main

import (
	"douyin/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.InitRouter(r)

	_ = r.Run()
}
