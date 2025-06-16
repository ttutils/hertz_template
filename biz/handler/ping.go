package handler

import (
	"context"

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
