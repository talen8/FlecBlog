package logsetup

import (
	"fmt"
	"os"
	"sync"
)

// 日志文件轮转参数
const (
	// maxLogSize 单个日志文件的大小上限，超过后轮转
	maxLogSize = 10 << 20 // 10MB
	// maxLogBackups 保留的历史文件个数
	maxLogBackups = 5
)

// rotatingWriter 带轮转的日志文件 writer
type rotatingWriter struct {
	mu      sync.Mutex
	path    string
	file    *os.File
	size    int64
	maxSize int64
	backups int
}

// newRotatingWriter 打开日志文件，返回可写的轮转 writer
func newRotatingWriter(path string, maxSize int64, backups int) (*rotatingWriter, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600) //nolint:gosec // G304: 路径非外部输入
	if err != nil {
		return nil, err
	}

	// 取现有大小，避免重启后把已有内容当成空文件
	var size int64
	if info, err := file.Stat(); err == nil {
		size = info.Size()
	}

	return &rotatingWriter{
		path:    path,
		file:    file,
		size:    size,
		maxSize: maxSize,
		backups: backups,
	}, nil
}

// Write 写入日志，超过大小上限时轮转
func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return 0, fmt.Errorf("日志文件不可写: %s", w.path)
	}

	n, err := w.file.Write(p)
	if err != nil {
		return n, err
	}
	w.size += int64(n)

	if w.size >= w.maxSize {
		w.rotate()
	}
	return n, nil
}

// rotate 轮转日志文件
func (w *rotatingWriter) rotate() {
	// 删掉最老的备份
	_ = os.Remove(w.backupPath(w.backups))

	// 依次后移：.N-1 → .N，……，.1 → .2
	for i := w.backups - 1; i >= 1; i-- {
		src := w.backupPath(i)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		_ = os.Rename(src, w.backupPath(i+1))
	}

	// 当前文件变成 .1
	if err := os.Rename(w.path, w.backupPath(1)); err != nil {
		return
	}

	// 旧句柄还指向已改名的 app.log.1：关掉它，再打开新的
	_ = w.file.Close()
	w.file = nil
	w.reopen()
}

// reopen 重新打开日志文件
func (w *rotatingWriter) reopen() {
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		// 打不开就让后续写入返回错误，不让服务崩溃
		w.file = nil
		w.size = 0
		return
	}
	w.file = file
	w.size = 0
}

// backupPath 第 n 个历史文件的路径
func (w *rotatingWriter) backupPath(n int) string {
	return fmt.Sprintf("%s.%d", w.path, n)
}

// Close 关闭日志文件
func (w *rotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}
