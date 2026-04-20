package middleware

import (
	"bytes"
	"flec_blog/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	panicLogFile *os.File
)

const (
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

// 初始化 panic 日志文件
func init() {
	// 确保日志目录存在
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0750); err != nil {
		logger.Warn("创建日志目录失败: %v", err)
		return
	}

	// 打开 panic 日志文件（追加模式）
	panicLogPath := filepath.Join(logDir, "panic.log")
	var err error
	// #nosec G304 - panic.log 是固定的日志文件名，不是用户输入
	panicLogFile, err = os.OpenFile(panicLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logger.Warn("打开 panic 日志文件失败: %v", err)
		return
	}
}

// Recovery 增强版错误恢复中间件
// 功能：
//   - 捕获 panic 异常，防止服务崩溃
//   - 记录完整的堆栈跟踪信息
//   - 记录详细的请求信息（Headers、Body等）
//   - 同时输出到控制台和日志文件
//   - 关联请求ID便于追踪
//
// 使用: router.Use(middleware.Recovery())
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求ID（如果有）
				requestID, _ := c.Get("request_id")
				requestIDStr := ""
				if requestID != nil {
					requestIDStr = requestID.(string)[:8]
				}

				// 获取堆栈跟踪
				stack := debug.Stack()

				// 获取请求信息
				method := c.Request.Method
				path := c.Request.URL.Path
				query := c.Request.URL.RawQuery
				clientIP := c.ClientIP()
				userAgent := c.Request.UserAgent()

				// 拼接完整路径
				fullPath := path
				if query != "" {
					fullPath = path + "?" + query
				}

				// 读取请求体（如果有）
				body := ""
				if c.Request.Body != nil {
					bodyBytes, _ := io.ReadAll(c.Request.Body)
					// 恢复 body 以便后续可能的读取
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					// 限制 body 长度（避免日志过大）
					if len(bodyBytes) > 0 {
						if len(bodyBytes) > 500 {
							body = string(bodyBytes[:500]) + "... (truncated)"
						} else {
							body = string(bodyBytes)
						}
					}
				}

				// 获取关键 Headers
				contentType := c.GetHeader("Content-Type")
				authorization := c.GetHeader("Authorization")
				if authorization != "" {
					authorization = "Bearer ***" // 隐藏敏感信息
				}

				// 格式化时间
				timestamp := time.Now().Format("2006-01-02 15:04:05")

				// 构建详细的错误日志
				var logBuilder strings.Builder
				logBuilder.WriteString("\n")
				logBuilder.WriteString("╔═══════════════════════════════════ PANIC RECOVERED ═══════════════════════════════════╗\n")
				fmt.Fprintf(&logBuilder, "║ Time: %s                                                              ║\n", timestamp)
				if requestIDStr != "" {
					fmt.Fprintf(&logBuilder, "║ Request ID: %s                                                                     ║\n", requestIDStr)
				}
				logBuilder.WriteString("╠═══════════════════════════════════════════════════════════════════════════════════════╣\n")
				logBuilder.WriteString("║ 🔥 Panic Error:                                                                       ║\n")
				fmt.Fprintf(&logBuilder, "║ %s\n", formatMultiline(fmt.Sprintf("%v", err), 85))
				logBuilder.WriteString("╠═══════════════════════════════════════════════════════════════════════════════════════╣\n")
				logBuilder.WriteString("║ 📨 Request Info:                                                                      ║\n")
				fmt.Fprintf(&logBuilder, "║   Method: %s                                                                          ║\n", padRight(method, 80))
				fmt.Fprintf(&logBuilder, "║   Path: %s\n", formatMultiline(fullPath, 83))
				fmt.Fprintf(&logBuilder, "║   Client IP: %s                                                                   ║\n", padRight(clientIP, 77))
				fmt.Fprintf(&logBuilder, "║   User-Agent: %s\n", formatMultiline(userAgent, 82))
				if contentType != "" {
					fmt.Fprintf(&logBuilder, "║   Content-Type: %s                                                              ║\n", padRight(contentType, 78))
				}
				if authorization != "" {
					fmt.Fprintf(&logBuilder, "║   Authorization: %s                                                             ║\n", padRight(authorization, 77))
				}
				if body != "" {
					logBuilder.WriteString("║   Request Body:                                                                       ║\n")
					fmt.Fprintf(&logBuilder, "║   %s\n", formatMultiline(body, 85))
				}
				logBuilder.WriteString("╠═══════════════════════════════════════════════════════════════════════════════════════╣\n")
				logBuilder.WriteString("║ 📚 Stack Trace (Full):                                                                ║\n")
				logBuilder.WriteString(formatStackTrace(string(stack), false)) // 完整堆栈
				logBuilder.WriteString("╚═══════════════════════════════════════════════════════════════════════════════════════╝\n")

				// 完整日志（写入 panic.log）
				fullLog := logBuilder.String()
				if panicLogFile != nil {
					fmt.Fprint(panicLogFile, fullLog)
					_ = panicLogFile.Sync() // 立即刷新到磁盘，确保不丢失
				}

				// 简化日志（控制台输出）
				consoleLog := buildConsoleLog(timestamp, requestIDStr, err, method, fullPath, clientIP, string(stack))
				fmt.Printf("%s%s%s", ColorRed, consoleLog, ColorReset)

				// 返回 500 错误给客户端
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":       500,
					"message":    "服务器内部错误",
					"request_id": requestIDStr,
				})
			}
		}()
		c.Next()
	}
}

// formatStackTrace 格式化堆栈跟踪信息
// simplified: true 只显示项目代码，false 显示完整堆栈
func formatStackTrace(stack string, simplified bool) string {
	lines := strings.Split(stack, "\n")
	var result strings.Builder

	maxLines := 30 // 限制堆栈行数
	lineCount := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if lineCount >= maxLines {
			result.WriteString("║   ... (more stack frames)                                                            ║\n")
			break
		}

		if line == "" {
			continue
		}

		trimmedLine := strings.TrimSpace(line)

		// 如果是简化模式，跳过非项目代码
		if simplified {
			// 只保留包含 flec_blog 的行（项目代码）
			if !strings.Contains(trimmedLine, "flec_blog") {
				continue
			}
			// 跳过 middleware 自身的行（recovery.go）
			if strings.Contains(trimmedLine, "middleware/recovery.go") {
				continue
			}
		}

		// 格式化每一行
		if len(trimmedLine) > 85 {
			fmt.Fprintf(&result, "║   %s... ║\n", trimmedLine[:82])
		} else {
			fmt.Fprintf(&result, "║   %-85s ║\n", trimmedLine)
		}
		lineCount++
	}

	return result.String()
}

// buildConsoleLog 构建简化的控制台日志
func buildConsoleLog(timestamp, requestID string, err interface{}, method, path, ip, stack string) string {
	var builder strings.Builder

	// 提取关键堆栈信息（只显示项目代码的第一个出错位置）
	keyStack := extractKeyStack(stack)

	builder.WriteString("\n")
	builder.WriteString("╔════════════════════════ PANIC RECOVERED ═══════════════════════╗\n")
	fmt.Fprintf(&builder, "║ %s | ID: %s\n", timestamp, requestID)
	builder.WriteString("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Fprintf(&builder, "║ 🔥 Error: %v\n", formatErrorLine(fmt.Sprintf("%v", err)))
	fmt.Fprintf(&builder, "║ 📨 %s %s | IP: %s\n", method, formatPathLine(path), ip)
	if keyStack != "" {
		builder.WriteString("╠════════════════════════════════════════════════════════════════╣\n")
		fmt.Fprintf(&builder, "║ 📍 Location: %s\n", keyStack)
	}
	builder.WriteString("╚════════════════════════════════════════════════════════════════╝\n")
	fmt.Fprintf(&builder, "💡 Full details saved to logs/panic.log (Request ID: %s)\n", requestID)
	builder.WriteString("\n")

	return builder.String()
}

// extractKeyStack 提取关键堆栈信息（项目代码中第一个出错的位置）
func extractKeyStack(stack string) string {
	lines := strings.Split(stack, "\n")

	for i := 0; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		// 找到项目代码
		if strings.Contains(line, "flec_blog") && !strings.Contains(line, "middleware/recovery.go") {
			// 下一行通常是文件位置
			if i+1 < len(lines) {
				locationLine := strings.TrimSpace(lines[i+1])
				// 提取文件路径和行号
				if strings.Contains(locationLine, ".go:") {
					// 简化路径：只保留相对路径
					parts := strings.Split(locationLine, "flec_blog/server/")
					if len(parts) > 1 {
						return parts[1]
					}
					return locationLine
				}
			}
		}
	}

	return ""
}

// formatErrorLine 格式化错误信息行
func formatErrorLine(err string) string {
	maxLen := 58
	if len(err) > maxLen {
		return err[:maxLen-3] + "..."
	}
	return fmt.Sprintf("%-*s", maxLen, err) + " ║"
}

// formatPathLine 格式化路径行
func formatPathLine(path string) string {
	maxLen := 35
	if len(path) > maxLen {
		return path[:maxLen-3] + "..."
	}
	return path
}

// formatMultiline 格式化多行文本，确保在框内显示
func formatMultiline(text string, width int) string {
	if len(text) <= width {
		return fmt.Sprintf("%-*s ║", width, text)
	}

	var result strings.Builder
	for i := 0; i < len(text); i += width {
		end := i + width
		if end > len(text) {
			end = len(text)
		}
		line := text[i:end]
		fmt.Fprintf(&result, "%-*s ║\n", width, line)
		if end < len(text) {
			result.WriteString("║   ")
		}
	}

	// 移除最后的换行和前缀
	s := result.String()
	if strings.HasSuffix(s, "\n║   ") {
		s = s[:len(s)-6]
	}
	return s
}

// padRight 右侧填充空格
func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// ClosePanicLogFile 关闭 panic 日志文件（在程序退出时调用）
func ClosePanicLogFile() {
	if panicLogFile != nil {
		_ = panicLogFile.Close()
	}
}
