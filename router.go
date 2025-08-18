package main

import (
	handler "hertz_template/biz/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func customizedRegister(r *server.Hertz) {
	r.GET("/api/ping", handler.Ping)
	r.GET("/api/metrics", handler.MetricsHandler)
	r.GET("/api/server_info", handler.ServerInfo)
}
