package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

var (
	logFile      *os.File
	errorLogFile *os.File
	mu           sync.Mutex
)

// 初始化日志文件
func init() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0750); err != nil {
		log.Printf("创建日志目录失败: %v", err)
		return
	}

	var err error
	logFile, err = os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Printf("打开日志文件失败: %v", err)
	}

	errorLogFile, err = os.OpenFile(filepath.Join(logDir, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Printf("打开错误日志文件失败: %v", err)
	}
}

// Info 记录信息级别日志
func Info(format string, v ...interface{}) {
	writeLog("INFO", ColorCyan, format, v...)
}

// Warn 记录警告级别日志
func Warn(format string, v ...interface{}) {
	writeLog("WARN", ColorYellow, format, v...)
}

// Error 记录错误级别日志（同时写入 app.log 和 error.log）
func Error(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, v...)
	fileMsg := fmt.Sprintf("[%s] [ERROR] %s", timestamp, msg)

	if logFile != nil {
		fmt.Fprintln(logFile, fileMsg)
	}
	if errorLogFile != nil {
		fmt.Fprintln(errorLogFile, fileMsg)
	}

	consoleMsg := fmt.Sprintf("%s[%s]%s [%sERROR%s] %s",
		ColorGray, timestamp, ColorReset,
		ColorRed, ColorReset,
		msg,
	)
	log.Println(consoleMsg)
}

// HTTPRequest 记录 HTTP 请求日志
func HTTPRequest(requestID, method, path string, statusCode int, clientIP string, cost time.Duration, bodySize int, userAgent, referer string) {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level, levelColor := getStatusLevel(statusCode)

	// 写入文件
	fileMsg := fmt.Sprintf("[%s] [%s] [%s] %s %s | Status: %d | IP: %s | Cost: %v | Size: %d | UA: %s",
		timestamp, level, requestID, method, path, statusCode, clientIP, cost, bodySize, userAgent)
	if referer != "" {
		fileMsg += fmt.Sprintf(" | Referer: %s", referer)
	}
	if logFile != nil {
		fmt.Fprintln(logFile, fileMsg)
	}
	// 4xx/5xx 错误同时写入 error.log
	if statusCode >= 400 && errorLogFile != nil {
		fmt.Fprintln(errorLogFile, fileMsg)
	}

	// 控制台输出（带颜色）
	consoleMsg := fmt.Sprintf("%s[%s]%s [%s%s%s] [%s%s%s] %s%-6s%s %s | %sStatus: %d%s | IP: %s | Cost: %s%v%s | Size: %d",
		ColorGray, timestamp, ColorReset,
		levelColor, level, ColorReset,
		ColorCyan, requestID, ColorReset,
		getMethodColor(method), method, ColorReset,
		path,
		levelColor, statusCode, ColorReset,
		clientIP,
		getCostColor(cost), cost, ColorReset,
		bodySize,
	)
	log.Println(consoleMsg)
}

// HTTPError 记录 HTTP 请求中的错误
func HTTPError(requestID string, errMsg string) {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fileMsg := fmt.Sprintf("[%s] [ERROR] [%s] Errors: %s", timestamp, requestID, errMsg)

	if logFile != nil {
		fmt.Fprintln(logFile, fileMsg)
	}
	if errorLogFile != nil {
		fmt.Fprintln(errorLogFile, fileMsg)
	}

	log.Printf("%s[%s]%s [%sERROR%s] [%s] Errors: %s",
		ColorGray, timestamp, ColorReset,
		ColorRed, ColorReset,
		requestID, errMsg)
}

// writeLog 写入日志（同时输出到控制台和文件）
func writeLog(level, color, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, v...)

	if logFile != nil {
		fmt.Fprintf(logFile, "[%s] [%s] %s\n", timestamp, level, msg)
	}

	consoleMsg := fmt.Sprintf("%s[%s]%s [%s%s%s] %s",
		ColorGray, timestamp, ColorReset,
		color, level, ColorReset,
		msg,
	)
	log.Println(consoleMsg)
}

// getStatusLevel 根据状态码返回日志级别和颜色
func getStatusLevel(statusCode int) (string, string) {
	switch {
	case statusCode >= 500:
		return "ERROR", ColorRed
	case statusCode >= 400:
		return "WARN", ColorYellow
	case statusCode >= 300:
		return "INFO", ColorCyan
	default:
		return "INFO", ColorGreen
	}
}

// getMethodColor 根据 HTTP 方法返回颜色
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return ColorBlue
	case "POST":
		return ColorGreen
	case "PUT", "PATCH":
		return ColorYellow
	case "DELETE":
		return ColorRed
	default:
		return ColorWhite
	}
}

// getCostColor 根据耗时返回颜色
func getCostColor(cost time.Duration) string {
	switch {
	case cost > time.Second:
		return ColorRed
	case cost > 500*time.Millisecond:
		return ColorYellow
	default:
		return ColorGreen
	}
}

// Close 关闭日志文件（程序退出时调用）
func Close() {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		_ = logFile.Close()
	}
	if errorLogFile != nil {
		_ = errorLogFile.Close()
	}
}
