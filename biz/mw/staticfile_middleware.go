package mw

import (
	"context"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// StaticFileMiddleware 处理静态文件请求，并支持 fallback 到 index.html（用于 SPA）
func StaticFileMiddleware(staticFS fs.FS) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		filePath := string(c.Path())

		skipPrefixes := []string{"/api"}

		// 跳过 API 路由
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(filePath, prefix) {
				c.Next(ctx)
				return
			}
		}

		if filePath == "" || filePath == "/" {
			filePath = "/index.html"
		}

		relPath := strings.TrimPrefix(filePath, "/")
		fullPath := "static/" + relPath
		indexPath := "static/index.html"

		// 返回指定路径的文件
		if serveFileFromFS(c, staticFS, fullPath) {
			c.Abort()
			return
		}

		// fallback 到 index.html
		if serveFileFromFS(c, staticFS, indexPath) {
			hlog.Debugf("文件 %s 不存在，使用 index.html 代替", fullPath)
			c.Abort()
			return
		}

		hlog.Infof("文件 %s 和 index.html 都不存在，返回 404", fullPath)
		c.String(http.StatusNotFound, "404 not found")
		c.Abort()
	}
}

func serveFileFromFS(c *app.RequestContext, filesystem fs.FS, path string) bool {
	file, err := filesystem.Open(path)
	if err != nil {
		return false
	}
	defer func(file fs.File) {
		if err := file.Close(); err != nil {
			hlog.Debugf("关闭文件 %s 时出错: %v", path, err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return false
	}

	// 设置 content-type
	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Response.Header.Set("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
	return true
}
