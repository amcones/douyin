package main

import (
	"douyin/middleware"
	"douyin/models"
	"douyin/router"
	"flag"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	bindAddr := flag.String("addr", ":8080", "HTTP port to bind")
	pprof := flag.Bool("pprof", false, "Enable pprof on port localhost:6060")
	flag.Parse()

	// 解析命令行参数，读取服务器地址等配置信息

	h := server.Default(
		server.WithHostPorts(*bindAddr),
		server.WithMaxRequestBodySize(20<<40), // 提高request body的容量到20MB
	)
	models.ConnDB()
	models.ConnRedis()

	middleware.InitJwt()

	router.RegisterRoute(h)
	// 根据参数决定是否启动pprof的http服务器
	if *pprof {
		go func() {
			hlog.Debug(http.ListenAndServe("localhost:6060", nil))
		}()
		hlog.Warn("pprof 在 localhost:6060 启动，请不要在正式环境使用")
	}
	h.Use(middleware.CorsMw())
	h.Spin()
}
