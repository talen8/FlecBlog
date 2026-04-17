package upload

import (
	"fmt"
	"strings"

	"flec_blog/config"
	"flec_blog/pkg/upload/storage"
)

// StorageType 存储类型常量
const (
	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
	StorageTypeCOS   = "cos"
	StorageTypeOSS   = "oss"
	StorageTypeKodo  = "kodo"
	StorageTypeR2    = "r2"
	StorageTypeMinIO = "minio"
)

// NewStorage 根据配置创建存储实例
func NewStorage(uploadCfg *config.UploadConfig) (storage.Storage, error) {
	storageType := strings.ToLower(uploadCfg.StorageType)

	switch storageType {
	case "", StorageTypeLocal: // 空值默认使用本地存储
		return storage.NewLocalStorage("/app/data/uploads"), nil

	case StorageTypeS3:
		return storage.NewS3UnifiedStorage(*uploadCfg, "s3")

	case StorageTypeCOS:
		return storage.NewS3UnifiedStorage(*uploadCfg, "cos")

	case StorageTypeOSS:
		return storage.NewS3UnifiedStorage(*uploadCfg, "oss")

	case StorageTypeKodo:
		return storage.NewS3UnifiedStorage(*uploadCfg, "kodo")

	case StorageTypeR2:
		return storage.NewS3UnifiedStorage(*uploadCfg, "r2")

	case StorageTypeMinIO:
		return storage.NewS3UnifiedStorage(*uploadCfg, "minio")

	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", storageType)
	}
}

// MustNewStorage 创建存储实例，如果失败则panic（用于启动时初始化）
func MustNewStorage(uploadCfg *config.UploadConfig) storage.Storage {
	s, err := NewStorage(uploadCfg)
	if err != nil {
		panic(fmt.Sprintf("初始化存储失败: %v", err))
	}
	return s
}

// ============================================
// 系统初始化
// ============================================

// InitializeUploadSystem 初始化文件上传系统
func InitializeUploadSystem(globalCfg *config.Config) *Manager {
	uploadStorage := MustNewStorage(&globalCfg.Upload)

	// 使用与存储路径一致的绝对路径创建目录
	_ = storage.NewHelper(uploadStorage).CreateUploadDir("/app/data/uploads")

	return NewManager(uploadStorage, NewValidator(), globalCfg)
}
