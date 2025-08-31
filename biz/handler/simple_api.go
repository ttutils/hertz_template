package handler

import (
	"context"
	"hertz_template/utils/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Ping 测试网络接口
// @Tags 测试
// @Summary 测试网络接口
// @Description 测试网络接口
// @Accept application/json
// @Produce application/json
// @Router /api/ping [get]
func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{
		"msg": "pong",
	})
}

// ServerInfo 服务信息
// @Tags 测试
// @Summary 服务信息
// @Description 服务信息
// @Accept application/json
// @Produce application/json
// @Router /api/server_info [get]
func ServerInfo(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{
		"code": 200,
		"msg":  "获取成功",
		"data": utils.H{
			"name":    config.Cfg.Server.Name,
			"version": config.Cfg.Server.Version,
		},
	})
}

// GetDemo 获取demo状态
// @Tags 测试
// @Summary 获取demo
// @Description 获取demo
// @Accept application/json
// @Produce application/json
// @Router /api/is_demo [get]
func GetDemo(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{
		"code": 200,
		"msg":  "获取成功",
		"data": utils.H{
			"is_demo": config.Cfg.Server.IsDemo,
		},
	})
}
