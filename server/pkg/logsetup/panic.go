package logsetup

import (
	"context"
	"log/slog"
	"runtime/debug"
)

// Recover 捕获 panic 并记录 panic 值与完整堆栈，用于没有请求上下文的场景
func Recover(tag string, attrs ...slog.Attr) {
	// recover 只在被 defer 的这个函数里直接调用才有效，不能再往下包一层
	if r := recover(); r != nil {
		logPanic(context.Background(), tag, r, attrs...)
	}
}

// RecoverContext 捕获 panic 并记录 panic 值与完整堆栈，日志带上请求上下文
func RecoverContext(ctx context.Context, tag string, attrs ...slog.Attr) {
	if r := recover(); r != nil {
		logPanic(ctx, tag, r, attrs...)
	}
}

// logPanic 记录 panic 值与完整堆栈
// recover 之后 runtime 不再打印 panic 现场，堆栈必须在这里自己抓
func logPanic(ctx context.Context, tag string, r any, attrs ...slog.Attr) {
	args := make([]any, 0, len(attrs)+3)
	args = append(args,
		slog.String("tag", tag),
		slog.Any("panic", r),
		slog.String("stack", string(debug.Stack())),
	)
	for _, attr := range attrs {
		args = append(args, attr)
	}
	slog.ErrorContext(ctx, "goroutine panic 已恢复", args...)
}
