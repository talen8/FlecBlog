package upload

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"flec_blog/config"
	"flec_blog/pkg/random"
	"flec_blog/pkg/upload/storage"

	"github.com/gin-gonic/gin"
)

// ============================================
// 类型定义和常量
// ============================================

// Type 上传类型
type Type string

// AllowedFileTypes 支持的文件类型
var AllowedFileTypes = []string{
	// 图片类型
	"image/jpeg",
	"image/jpg",
	"image/png",
	"image/gif",
	"image/webp",
	"image/avif",
	"image/svg+xml",
	"image/bmp",
	"image/tiff",
	// 视频类型
	"video/mp4",
	"video/webm",
	"video/quicktime",
	"video/x-msvideo",  // avi
	"video/x-matroska", // mkv
	"video/mpeg",
	"video/3gpp",
	"video/x-flv",
	// 音频类型
	"audio/mpeg", // mp3
	"audio/wav",
	"audio/ogg",
	"audio/aac",
	"audio/flac",
	// 文档类型
	"text/plain",
	"application/pdf",
	"application/msword", // doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document", // docx
	// 压缩文件
	"application/zip",
	"application/x-rar-compressed",
	"application/x-7z-compressed",
	// JSON类型（用于配置文件等）
	"application/json",
}

// FileInfo 文件信息
type FileInfo struct {
	OriginalName string    `json:"original_name"` // 原始文件名
	FileName     string    `json:"file_name"`     // 存储文件名
	FilePath     string    `json:"file_path"`     // 文件路径
	FileSize     int64     `json:"file_size"`     // 文件大小(字节)
	FileType     string    `json:"file_type"`     // 文件MIME类型
	FileHash     string    `json:"file_hash"`     // 文件SHA256哈希值
	UploadType   Type      `json:"upload_type"`   // 上传用途
	StorageType  string    `json:"storage_type"`  // 存储类型
	UserID       uint      `json:"user_id"`       // 上传用户ID
	UploadTime   time.Time `json:"upload_time"`   // 上传时间
	FileURL      string    `json:"file_url"`      // 访问URL
}

// Request 上传请求
type Request struct {
	File               *multipart.FileHeader `json:"file"`                 // 上传文件
	UploadType         Type                  `json:"upload_type"`          // 上传类型
	UserID             uint                  `json:"user_id"`              // 用户ID
	SkipSizeValidation bool                  `json:"skip_size_validation"` // 跳过大小验证（后台管理使用）
}

// Response 上传响应
type Response struct {
	Success  bool      `json:"success"`   // 上传是否成功
	Message  string    `json:"message"`   // 响应消息
	FileInfo *FileInfo `json:"file_info"` // 文件信息
}

// Validator 验证器接口
type Validator interface {
	Validate(file *multipart.FileHeader, config *config.UploadConfig, skipSizeValidation bool) error
}

// ============================================
// 文件验证器
// ============================================

// FileValidator 文件验证器
type FileValidator struct{}

// NewValidator 创建验证器
func NewValidator() Validator {
	return &FileValidator{}
}

// Validate 验证文件
func (v *FileValidator) Validate(file *multipart.FileHeader, cfg *config.UploadConfig, skipSizeValidation bool) error {
	if file == nil {
		return errors.New("文件不能为空")
	}

	if strings.TrimSpace(file.Filename) == "" {
		return errors.New("文件名不能为空")
	}

	// 只有在不跳过大小验证时才检查文件大小
	if !skipSizeValidation && cfg.MaxFileSize > 0 && file.Size > cfg.MaxFileSize*1024*1024 {
		return fmt.Errorf("文件大小超出限制，最大允许: %dMB", cfg.MaxFileSize)
	}

	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		return errors.New("无法确定文件类型")
	}

	allowed := false
	for _, allowedType := range AllowedFileTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return errors.New("文件必须有扩展名")
	}

	dangerousExts := []string{".exe", ".bat", ".cmd", ".com", ".pif", ".scr", ".vbs", ".js", ".jar", ".sh"}
	for _, dangerous := range dangerousExts {
		if ext == dangerous {
			return fmt.Errorf("不允许上传可执行文件: %s", ext)
		}
	}

	return nil
}

// ============================================
// 上传管理器
// ============================================

// Manager 上传管理器
type Manager struct {
	storage   storage.Storage
	validator Validator
	config    *config.Config // 全局配置对象引用（支持热重载）
	mu        sync.RWMutex   // 保护存储实例的并发访问
}

// NewManager 创建上传管理器
func NewManager(s storage.Storage, validator Validator, cfg *config.Config) *Manager {
	return &Manager{storage: s, validator: validator, config: cfg}
}

// GetMaxFileSize 获取最大文件大小限制（MB）
func (m *Manager) GetMaxFileSize() int64 {
	return m.config.Upload.MaxFileSize
}

// GetStorageType 获取存储类型
func (m *Manager) GetStorageType() string {
	return m.config.Upload.StorageType
}

// HealthCheck 检查存储服务可用性
func (m *Manager) HealthCheck() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.storage.HealthCheck()
}

// ReloadStorage 重新加载存储实例（用于热重载）
func (m *Manager) ReloadStorage() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	newStorage, err := NewStorage(&m.config.Upload)
	if err != nil {
		return err
	}

	m.storage = newStorage
	return nil
}

// HandleUpload 处理文件上传
func (m *Manager) HandleUpload(req *Request, host string) (*Response, error) {
	// 1. 验证文件
	if err := m.validator.Validate(req.File, &m.config.Upload, req.SkipSizeValidation); err != nil {
		return m.createErrorResponse(fmt.Sprintf("文件验证失败: %v", err)), nil
	}

	// 2. 计算文件hash
	fileHash, err := m.CalculateFileHash(req.File)
	if err != nil {
		return m.createErrorResponse(fmt.Sprintf("计算文件哈希失败: %v", err)), err
	}

	// 3. 生成文件路径
	filePath := m.generateFilePath(req.UploadType, req.UserID, req.File.Filename)

	// 4. 打开上传文件
	file, err := req.File.Open()
	if err != nil {
		return m.createErrorResponse(fmt.Sprintf("打开文件失败: %v", err)), err
	}
	defer func() {
		_ = file.Close()
	}()

	// 5. 保存文件
	if err := m.storage.Save(file, filePath); err != nil {
		return m.createErrorResponse(fmt.Sprintf("文件保存失败: %v", err)), err
	}

	// 6. 生成文件信息（使用动态 host）
	fileInfo := m.createFileInfo(req, filePath, fileHash, host)

	return &Response{
		Success:  true,
		Message:  "文件上传成功",
		FileInfo: fileInfo,
	}, nil
}

// HandleUploadFromReader 从 Reader 上传文件
func (m *Manager) HandleUploadFromReader(reader io.Reader) ([]byte, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", fmt.Errorf("读取数据失败: %w", err)
	}

	hash := sha256.New()
	hash.Write(data)
	fileHash := hex.EncodeToString(hash.Sum(nil))

	return data, fileHash, nil
}

// SaveFileData 保存文件数据
func (m *Manager) SaveFileData(data []byte, fileHash string, originalName string, fileType string, uploadType Type, userID uint, host string) (*FileInfo, error) {
	// 1. 生成文件路径
	filePath := m.generateFilePath(uploadType, userID, originalName)

	// 2. 保存文件
	if err := m.storage.Save(bytes.NewReader(data), filePath); err != nil {
		return nil, fmt.Errorf("文件保存失败: %w", err)
	}

	// 3. 生成文件信息（使用动态 host）
	fileInfo := &FileInfo{
		OriginalName: originalName,
		FileName:     filepath.Base(filePath),
		FilePath:     filePath,
		FileSize:     int64(len(data)),
		FileType:     fileType,
		FileHash:     fileHash,
		UploadType:   uploadType,
		UserID:       userID,
		UploadTime:   time.Now(),
		FileURL:      m.storage.GetURL(filePath, host),
	}

	return fileInfo, nil
}

// DeleteFile 删除文件
func (m *Manager) DeleteFile(filePath string) error {
	return m.storage.Delete(filePath)
}

// DeleteFileByStorageType 根据存储类型删除文件
func (m *Manager) DeleteFileByStorageType(filePath string, storageType string) error {
	// 如果是当前使用的存储类型，直接使用现有实例
	if storageType == m.config.Upload.StorageType {
		return m.storage.Delete(filePath)
	}

	// 如果是其他存储类型，创建对应的临时存储实例
	var targetStorage storage.Storage
	var err error

	switch storageType {
	case "local":
		targetStorage = storage.NewLocalStorage("/app/data/uploads")
	case "s3":
		targetStorage, err = storage.NewS3UnifiedStorage(m.config.Upload, "s3")
	case "cos":
		targetStorage, err = storage.NewS3UnifiedStorage(m.config.Upload, "cos")
	case "oss":
		targetStorage, err = storage.NewS3UnifiedStorage(m.config.Upload, "oss")
	case "kodo":
		targetStorage, err = storage.NewS3UnifiedStorage(m.config.Upload, "kodo")
	case "r2":
		targetStorage, err = storage.NewS3UnifiedStorage(m.config.Upload, "r2")
	case "minio":
		targetStorage, err = storage.NewS3UnifiedStorage(m.config.Upload, "minio")
	default:
		return fmt.Errorf("不支持的存储类型: %s", storageType)
	}

	if err != nil {
		return fmt.Errorf("创建存储实例失败: %w", err)
	}

	return targetStorage.Delete(filePath)
}

// generateFilePath 生成文件路径（内部方法）
func (m *Manager) generateFilePath(uploadType Type, userID uint, originalName string) string {
	return GenerateFilePath(uploadType, userID, originalName, m.config.Upload.PathPattern)
}

// GenerateFilePath 生成文件路径
func GenerateFilePath(uploadType Type, userID uint, originalName string, pattern string) string {
	now := time.Now()
	ext := filepath.Ext(originalName)
	filename := strings.TrimSuffix(originalName, ext)
	timestamp := now.Format("20060102150405")
	randomStr := random.Code(8)

	// 使用配置的路径模式
	if pattern == "" {
		pattern = "{timestamp}_{random}{ext}"
	}

	// 创建替换映射表
	replacements := map[string]string{
		"YYYY":        now.Format("2006"),
		"MM":          now.Format("01"),
		"DD":          now.Format("02"),
		"HH":          now.Format("15"),
		"mm":          now.Format("04"),
		"ss":          now.Format("05"),
		"{type}":      string(uploadType),
		"{userid}":    fmt.Sprintf("%d", userID),
		"{timestamp}": timestamp,
		"{random}":    randomStr,
		"{filename}":  filename,
		"{ext}":       ext,
	}

	// 执行替换
	result := pattern
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// createErrorResponse 创建错误响应
func (m *Manager) createErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Message: message,
	}
}

// CalculateFileHash 计算文件哈希
func (m *Manager) CalculateFileHash(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// createFileInfo 创建文件信息
func (m *Manager) createFileInfo(req *Request, filePath string, fileHash string, host string) *FileInfo {
	return &FileInfo{
		OriginalName: req.File.Filename,
		FileName:     filepath.Base(filePath),
		FilePath:     filePath,
		FileSize:     req.File.Size,
		FileType:     req.File.Header.Get("Content-Type"),
		FileHash:     fileHash,
		UploadType:   req.UploadType,
		StorageType:  m.config.Upload.StorageType,
		UserID:       req.UserID,
		UploadTime:   time.Now(),
		FileURL:      m.storage.GetURL(filePath, host),
	}
}

// ============================================
// 辅助函数
// ============================================

// ExtractHostFromContext 从 gin.Context 中提取完整的 host 地址（包含 scheme）
func ExtractHostFromContext(c *gin.Context) string {
	// 1. 检查是否有反向代理的 header
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// 检查 X-Forwarded-Proto header（反向代理场景）
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	// 2. 获取 host
	host := c.Request.Host

	// 检查 X-Forwarded-Host header（反向代理场景）
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
