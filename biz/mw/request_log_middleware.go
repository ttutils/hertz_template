package mw

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().PathOriginal())

		// 记录请求详细信息（仅包含/api的路径）
		if strings.Contains(path, "/api") {
			// 复制请求体内容并重新设置
			bodyBytes := c.Request.BodyBytes()
			body := string(bodyBytes)
			c.Request.SetBody(bodyBytes) // 重置body供后续处理使用

			// 构建详细日志
			var logContent strings.Builder
			logContent.WriteString("\n\033[0;31m请求方式：\033[0;32m")
			logContent.WriteString(string(c.Request.Method()))
			logContent.WriteString("\033[0m\n\033[0;31m请求地址：\033[0;32m")
			logContent.WriteString(c.ClientIP())
			logContent.WriteString("\033[0m\n\033[0;31m请求接口：\033[0;32m")
			logContent.WriteString(path)
			logContent.WriteString("\033[0m\n\033[0;31m请求头：\033[0;32m")

			// 遍历所有请求头
			c.Request.Header.VisitAll(func(key, value []byte) {
				logContent.WriteString(fmt.Sprintf("%s: %s; ", string(key), string(value)))
			})

			logContent.WriteString("\033[0m\n\033[0;31m请求体：\033[0;32m")
			logContent.WriteString(body)
			logContent.WriteString("\033[0m\n\033[0;31m查询参数：\033[0;32m")
			logContent.WriteString(c.QueryArgs().String())
			logContent.WriteString("\033[0m")

			hlog.CtxInfof(ctx, logContent.String())
		}

		c.Next(ctx)

		// 记录响应内容和处理时间
		latency := time.Since(start)
		c.Response.Header.Set("X-Process-Time", fmt.Sprintf("%.5f", latency.Seconds()))

		// 获取响应体内容
		responseBody := string(c.Response.Body())

		// 构建响应日志
		responseLog := fmt.Sprintf("\n\033[0;31m响应状态：\033[0;32m%d\033[0m\n\033[0;31m\033[0;31m响应内容：\033[0;32m%s\033[0m",
			c.Response.StatusCode(),
			responseBody)

		hlog.CtxInfof(ctx, responseLog)
		hlog.CtxInfof(ctx, "\033[0;32m耗时: %.5f s\033[0m", latency.Seconds())
	}
}
