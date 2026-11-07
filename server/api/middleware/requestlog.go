package middleware

import (
	"log/slog"
	"strings"
	"time"

	"flec_blog/pkg/logsetup"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger HTTP 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
		c.Set("request_id", requestID)
		c.Request = c.Request.WithContext(logsetup.WithRequestID(c.Request.Context(), requestID))

		// 记录开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		// 处理请求
		c.Next()

		// 拼接完整路径
		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		statusCode := c.Writer.Status()
		cost := time.Since(start)

		// 记录 HTTP 请求日志，级别由状态码推导
		slog.Log(c.Request.Context(), httpLevel(statusCode), "HTTP 请求",
			slog.String("method", method),
			slog.String("path", fullPath),
			slog.Int("status", statusCode),
			slog.String("ip", clientIP),
			slog.Int64("cost_ms", cost.Milliseconds()),
			slog.Int("size", c.Writer.Size()),
			slog.String("ua", userAgent),
			slog.String("referer", referer),
		)

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			slog.ErrorContext(c.Request.Context(), "HTTP 请求错误",
				slog.String("errors", c.Errors.String()),
			)
		}
	}
}

// httpLevel 根据 HTTP 状态码推导日志级别
func httpLevel(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
