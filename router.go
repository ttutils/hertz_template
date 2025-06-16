package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	handler "hertz_template/biz/handler"
)

func customizedRegister(r *server.Hertz) {
	r.GET("/api/ping", handler.Ping)
	r.GET("/api/metrics", handler.MetricsHandler)
}
