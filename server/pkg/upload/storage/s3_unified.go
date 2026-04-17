package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"flec_blog/config"
)

// S3UnifiedStorage 基于MinIO SDK的统一S3存储
type S3UnifiedStorage struct {
	cfg         config.UploadConfig
	storageType string
	client      *minio.Client
	bucketName  string
	baseURL     string
	mu          sync.Mutex
}

// NewS3UnifiedStorage 创建统一S3存储实例
func NewS3UnifiedStorage(cfg config.UploadConfig, storageType string) (*S3UnifiedStorage, error) {
	storage := &S3UnifiedStorage{
		cfg:         cfg,
		storageType: storageType,
		bucketName:  cfg.Bucket,
	}
	return storage, nil
}

// ensureClient 确保客户端已初始化
func (s *S3UnifiedStorage) ensureClient() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return nil
	}

	endpoint := s.cfg.Endpoint
	useSSL := s.cfg.UseSSL

	if endpoint == "" {
		switch s.storageType {
		case "s3":
			if s.cfg.Region == "" {
				return fmt.Errorf("AWS S3 需要配置 region")
			}
			endpoint = fmt.Sprintf("s3.%s.amazonaws.com", s.cfg.Region)
			useSSL = true
		case "cos":
			if s.cfg.Region == "" {
				return fmt.Errorf("腾讯云 COS 需要配置 region")
			}
			endpoint = fmt.Sprintf("cos.%s.myqcloud.com", s.cfg.Region)
			useSSL = true
		case "oss":
			if s.cfg.Region == "" {
				return fmt.Errorf("阿里云 OSS 需要配置 region")
			}
			endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", s.cfg.Region)
			useSSL = true
		case "kodo":
			if s.cfg.Region == "" {
				return fmt.Errorf("七牛云 Kodo 需要配置 region")
			}
			endpoint = fmt.Sprintf("s3.%s.qiniucs.com", s.cfg.Region)
			useSSL = true
		default:
			return fmt.Errorf("存储类型 %s 需要配置 endpoint", s.storageType)
		}
	}

	region := s.cfg.Region
	if s.storageType == "minio" && region == "" {
		region = "us-east-1"
	}

	bucketLookup := minio.BucketLookupAuto
	if s.storageType == "cos" || s.storageType == "oss" || s.storageType == "kodo" {
		bucketLookup = minio.BucketLookupDNS
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(s.cfg.AccessKey, s.cfg.SecretKey, ""),
		Secure:       useSSL,
		Region:       region,
		BucketLookup: bucketLookup,
	})
	if err != nil {
		return fmt.Errorf("创建存储实例失败: %w", err)
	}

	s.client = client
	s.baseURL = buildBaseURL(endpoint, s.cfg.Bucket, useSSL, s.cfg.Domain)
	return nil
}

// buildObjectKey 构建完整的对象键
func (s *S3UnifiedStorage) buildObjectKey(filePath string) string {
	return strings.TrimPrefix(filePath, "/")
}

// buildBaseURL 构建基础访问URL
func buildBaseURL(endpoint, bucket string, useSSL bool, customDomain string) string {
	if customDomain != "" {
		return strings.TrimSuffix(customDomain, "/")
	}

	scheme := "http"
	if useSSL {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s.%s", scheme, bucket, endpoint)
}

// Save 实现 Storage 接口 - 保存文件
func (s *S3UnifiedStorage) Save(reader io.Reader, filePath string) error {
	if err := s.ensureClient(); err != nil {
		return err
	}
	objectKey := s.buildObjectKey(filePath)
	_, err := s.client.PutObject(context.Background(), s.bucketName, objectKey, reader, -1, minio.PutObjectOptions{
		ContentType: getContentType(filePath),
	})
	return err
}

// Delete 实现 Storage 接口 - 删除文件
func (s *S3UnifiedStorage) Delete(filePath string) error {
	if err := s.ensureClient(); err != nil {
		return err
	}
	objectKey := s.buildObjectKey(filePath)
	return s.client.RemoveObject(context.Background(), s.bucketName, objectKey, minio.RemoveObjectOptions{})
}

// GetURL 实现 Storage 接口 - 获取文件访问 URL
func (s *S3UnifiedStorage) GetURL(filePath string, _ string) string {
	objectKey := s.buildObjectKey(filePath)

	if s.baseURL != "" {
		return fmt.Sprintf("%s/%s", s.baseURL, objectKey)
	}

	endpointURL := s.client.EndpointURL()
	return fmt.Sprintf("%s://%s/%s", endpointURL.Scheme, endpointURL.Host, objectKey)
}

// GetPreSignedURL 获取预签名URL（用于临时访问）
func (s *S3UnifiedStorage) GetPreSignedURL(filePath string, expires time.Duration) (string, error) {
	objectKey := s.buildObjectKey(filePath)
	presignedURL, err := s.client.PresignedGetObject(context.Background(), s.bucketName, objectKey, expires, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// BatchDelete 批量删除文件
func (s *S3UnifiedStorage) BatchDelete(filePaths []string) error {
	if len(filePaths) == 0 {
		return nil
	}

	ctx := context.Background()
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, filePath := range filePaths {
			objectKey := s.buildObjectKey(filePath)
			objectsCh <- minio.ObjectInfo{Key: objectKey}
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	errorCh := s.client.RemoveObjects(ctx, s.bucketName, objectsCh, opts)

	var errors []string
	for err := range errorCh {
		if err.Err != nil {
			errors = append(errors, fmt.Sprintf("删除 %s 失败: %v", err.ObjectName, err.Err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量删除部分失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Exists 实现 Storage 接口 - 检查文件是否存在
func (s *S3UnifiedStorage) Exists(filePath string) bool {
	objectKey := s.buildObjectKey(filePath)
	_, err := s.client.StatObject(context.Background(), s.bucketName, objectKey, minio.StatObjectOptions{})
	return err == nil
}

// HealthCheck 检查存储可用性
func (s *S3UnifiedStorage) HealthCheck() error {
	if err := s.ensureClient(); err != nil {
		return fmt.Errorf("初始化存储客户端失败: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("存储桶不可访问: %w", err)
	}
	return nil
}

// GetObjectInfo 获取对象信息
func (s *S3UnifiedStorage) GetObjectInfo(filePath string) (*ObjectInfo, error) {
	objectKey := s.buildObjectKey(filePath)
	objInfo, err := s.client.StatObject(context.Background(), s.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &ObjectInfo{
		Key:          objInfo.Key,
		Size:         objInfo.Size,
		LastModified: objInfo.LastModified,
		ContentType:  objInfo.ContentType,
		ETag:         objInfo.ETag,
	}, nil
}

// ObjectInfo 对象信息
type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ContentType  string    `json:"content_type"`
	ETag         string    `json:"etag"`
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".avif": "image/avif",
		".svg":  "image/svg+xml",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".json": "application/json",
		".xml":  "application/xml",
		".css":  "text/css",
		".js":   "application/javascript",
		".html": "text/html",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}
