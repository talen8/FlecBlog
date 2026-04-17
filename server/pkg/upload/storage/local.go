package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"flec_blog/pkg/random"
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
	fullPath := filepath.Join(s.basePath, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		os.Remove(fullPath)
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

// GenerateFilePath 生成文件路径（uploadType作为字符串传入）
func (h *Helper) GenerateFilePath(uploadType string, userID uint, originalName string, pathPattern string) string {
	now := time.Now()

	// 获取文件信息
	ext := filepath.Ext(originalName)
	filename := strings.TrimSuffix(originalName, ext)
	timestamp := now.Format("20060102150405")
	randomStr := random.Code(8)

	result := pathPattern

	// 时间占位符
	result = strings.ReplaceAll(result, "YYYY", now.Format("2006"))
	result = strings.ReplaceAll(result, "MM", now.Format("01"))
	result = strings.ReplaceAll(result, "DD", now.Format("02"))
	result = strings.ReplaceAll(result, "HH", now.Format("15"))
	result = strings.ReplaceAll(result, "mm", now.Format("04"))
	result = strings.ReplaceAll(result, "ss", now.Format("05"))
	result = strings.ReplaceAll(result, "{type}", uploadType)
	result = strings.ReplaceAll(result, "{userid}", fmt.Sprintf("%d", userID))
	result = strings.ReplaceAll(result, "{timestamp}", timestamp)
	result = strings.ReplaceAll(result, "{random}", randomStr)
	result = strings.ReplaceAll(result, "{filename}", filename)
	result = strings.ReplaceAll(result, "{ext}", ext)

	return result
}

// CreateUploadDir 创建上传目录
func (h *Helper) CreateUploadDir(basePath string) error {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("创建上传目录失败: %v", err)
	}
	return nil
}
