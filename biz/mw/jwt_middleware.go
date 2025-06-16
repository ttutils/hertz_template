package mw

import (
	"context"
	"hertz_template/utils"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// JWTAuthMiddleware 鉴权中间件
func JWTAuthMiddleware(excludedPaths []string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取请求路径并转换为字符串
		path := string(c.Request.Path())

		skipPrefixes := []string{"/api"}

		// 需要鉴权的路径
		for _, prefix := range skipPrefixes {
			if !strings.HasPrefix(path, prefix) {
				c.Next(ctx)
				return
			}
		}

		// 如果路径是 swagger 文档，则跳过鉴权
		if strings.HasPrefix(path, "/api/swagger") {
			c.Next(ctx)
			return
		}

		// 如果路径在排除列表中，则跳过鉴权
		for _, excludedPath := range excludedPaths {
			if path == excludedPath {
				c.Next(ctx) // 跳过中间件，直接处理请求
				return
			}
		}

		// 获取 Authorization Header
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code": consts.StatusUnauthorized,
				"msg":  "缺少token",
			})
			c.Abort() // 终止后续处理
			return
		}

		// 验证 token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code": consts.StatusUnauthorized,
				"msg":  err.Error(),
			})
			c.Abort() // 终止后续处理
			return
		}

		// 将 claims 保存到上下文
		for k, v := range claims {
			c.Set(k, v)
		}
		c.Set("userid", claims["userid"])
		c.Set("username", claims["username"])

		// 如果验证通过，继续处理请求
		c.Next(ctx)
	}
}
