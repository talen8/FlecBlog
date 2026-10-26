package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ============================================
// 存储接口定义
// ============================================

// Storage 存储接口
type Storage interface {
	Save(reader io.Reader, path string) error
	Delete(path string) error
	Exists(path string) bool
	GetURL(path string, host string) string
	HealthCheck() error
}

// ============================================
// 本地存储实现
// ============================================

// LocalStorage 本地存储
type LocalStorage struct {
	basePath string // 基础存储路径
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(basePath string) Storage {
	// 确保使用绝对路径
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		absPath = basePath
	}

	return &LocalStorage{
		basePath: absPath,
	}
}

// Save 保存文件
func (s *LocalStorage) Save(reader io.Reader, path string) error {
	// 清理路径，防止路径遍历攻击
	path = filepath.Clean(path)
	// 确保路径不以 .. 开头，防止访问上级目录
	if strings.HasPrefix(path, "..") || strings.Contains(path, "../") || strings.Contains(path, "..") {
		return fmt.Errorf("非法的文件路径")
	}

	fullPath := filepath.Join(s.basePath, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// #nosec G304 - 文件路径已在函数开头进行安全检查
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, reader)
	if err != nil {
		_ = os.Remove(fullPath)
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.basePath, path)

	if !s.Exists(path) {
		return nil
	}

	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	s.removeEmptyDirs(filepath.Dir(fullPath))
	return nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(path string) bool {
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetURL 获取文件 URL
func (s *LocalStorage) GetURL(path string, host string) string {
	urlPath := strings.ReplaceAll(path, "\\", "/")

	if host == "" {
		// 没有 host 时返回相对路径
		return "/uploads/" + urlPath
	}

	// 组装完整 URL
	serverHost := strings.TrimSuffix(host, "/")
	return serverHost + "/uploads/" + urlPath
}

// HealthCheck 检查存储可用性
func (s *LocalStorage) HealthCheck() error {
	info, err := os.Stat(s.basePath)
	if err != nil {
		return fmt.Errorf("存储目录不存在: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("存储路径不是目录")
	}
	return nil
}

// removeEmptyDirs 递归删除空目录
func (s *LocalStorage) removeEmptyDirs(dir string) {
	if dir == s.basePath || dir == "" {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) > 0 {
		return
	}

	if err := os.Remove(dir); err == nil {
		s.removeEmptyDirs(filepath.Dir(dir))
	}
}

// Helper 存储辅助工具
type Helper struct {
	storage Storage
}

// NewHelper 创建辅助工具
func NewHelper(storage Storage) *Helper {
	return &Helper{storage: storage}
}

// CreateUploadDir 创建上传目录
func (h *Helper) CreateUploadDir(basePath string) error {
	if err := os.MkdirAll(basePath, 0750); err != nil {
		return fmt.Errorf("创建上传目录失败: %v", err)
	}
	return nil
}
