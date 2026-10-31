package logsetup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// logDir 日志目录
const logDir = "/app/data/logs"

// appLogName 日志文件名
const appLogName = "app.log"

// requestIDKey request_id 在 context 中的键
type requestIDKey struct{}

// logFile 日志文件句柄，由 Setup 创建、Close 关闭
var logFile *os.File

// WithRequestID 把 request_id 写入 context
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID 从 context 中取出 request_id
func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// Setup 创建日志目录、打开日志文件，并把组装好的 logger 注册为 slog 默认 logger
func Setup() error {
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return fmt.Errorf("创建日志目录 %s: %w", logDir, err)
	}

	file, err := os.OpenFile(filepath.Join(logDir, appLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("打开日志文件 %s: %w", appLogName, err)
	}
	logFile = file

	console := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logFileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo})

	slog.SetDefault(slog.New(&fanoutHandler{handlers: []slog.Handler{console, logFileHandler}}))
	return nil
}

// Close 关闭日志文件句柄
func Close() {
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
}
