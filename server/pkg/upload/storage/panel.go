package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"flec_blog/pkg/panel"
)

// ============================================
// 托管存储实现（经 Panel 转发至托管存储服务）
// ============================================

// SaveMetadata 托管存储保存后的元数据（用于回传给调用方落库）
type SaveMetadata struct {
	URL string // 托管存储完整访问 URL
}

// MetaSaver 支持返回元数据的存储扩展接口
type MetaSaver interface {
	SaveWithMeta(reader io.Reader, fileName string) (SaveMetadata, error)
}

// HostedStorage 托管存储
type HostedStorage struct {
	client *panel.Client
}

// NewHostedStorage 创建托管存储实例
func NewHostedStorage() *HostedStorage {
	return &HostedStorage{client: panel.New()}
}

// Save 实现 Storage 接口 - 保存文件（元数据由 SaveWithMeta 提供）
func (s *HostedStorage) Save(reader io.Reader, path string) error {
	_, err := s.SaveWithMeta(reader, filepath.Base(path))
	return err
}

// SaveWithMeta 上传文件并返回托管存储元数据
func (s *HostedStorage) SaveWithMeta(reader io.Reader, fileName string) (SaveMetadata, error) {
	fileName = filepath.Base(fileName)
	var resp struct {
		URL string `json:"url"`
	}
	if err := s.client.PostMultipart(context.Background(), "/api/storage/upload",
		"file", fileName, reader, &resp); err != nil {
		return SaveMetadata{}, fmt.Errorf("面板存储上传失败: %w", err)
	}
	if resp.URL == "" {
		return SaveMetadata{}, fmt.Errorf("面板存储上传失败: 响应缺少 url")
	}
	return SaveMetadata{URL: resp.URL}, nil
}

// Delete 实现 Storage 接口 - 删除文件（托管存储无需物理删除）
func (s *HostedStorage) Delete(path string) error {
	return nil
}

// GetURL 实现 Storage 接口 - 获取文件访问 URL
func (s *HostedStorage) GetURL(path string, _ string) string {
	// 托管存储的文件 URL 由上传时返回，保存于文件记录中，此处直接返回传入路径
	return path
}

// Exists 实现 Storage 接口 - 检查文件是否存在
func (s *HostedStorage) Exists(path string) bool {
	return false
}

// HealthCheck 检查存储可用性
func (s *HostedStorage) HealthCheck() error {
	// 托管存储连通性由 Panel 保证，面板不可达时上传/删除会立即失败
	return nil
}
