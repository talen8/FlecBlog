package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"

	"flec_blog/pkg/logsetup"

	"github.com/gin-gonic/gin"
)

// maxPanicBodySize panic 日志中保留的请求体最大字节数，避免日志过大
const maxPanicBodySize = 500

// Recovery panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := requestIDOf(c)

				// 读取请求体
				body := ""
				if c.Request.Body != nil {
					bodyBytes, _ := io.ReadAll(c.Request.Body)
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					if len(bodyBytes) > 0 {
						if len(bodyBytes) > maxPanicBodySize {
							body = string(bodyBytes[:maxPanicBodySize]) + "... (truncated)"
						} else {
							body = string(bodyBytes)
						}
					}
				}

				// 记录完整现场：panic 值、请求信息、完整堆栈
				slog.ErrorContext(c.Request.Context(), "request panic 已恢复",
					slog.Any("panic", err),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("ip", c.ClientIP()),
					slog.String("ua", c.Request.UserAgent()),
					slog.String("content_type", c.GetHeader("Content-Type")),
					slog.String("body", body),
					slog.String("stack", string(debug.Stack())),
				)

				// 返回 500 错误给客户端
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":       500,
					"message":    "服务器内部错误",
					"request_id": requestID,
				})
			}
		}()
		c.Next()
	}
}

// requestIDOf 取出本次请求的 ID
func requestIDOf(c *gin.Context) string {
	if id := logsetup.RequestID(c.Request.Context()); id != "" {
		return id
	}
	if value, ok := c.Get("request_id"); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}
