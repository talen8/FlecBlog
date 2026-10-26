package repository

import (
	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// FileRepository 文件仓储
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓储
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// ============ 基础CRUD ============

// Create 创建文件记录
func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

// Get 获取文件信息
func (r *FileRepository) Get(id uint) (*model.File, error) {
	var file model.File
	err := r.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// Delete 删除文件记录
func (r *FileRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&model.File{}, id).Error
}

// ============ 查询方法 ============

// List 获取文件列表
func (r *FileRepository) List(
	offset, limit int,
	keyword, fileType string,
	status *int,
	uploadType string,
	minSize, maxSize int64,
	startTime, endTime string,
) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.db.Model(&model.File{})

	if keyword != "" {
		query = query.Where("file_name LIKE ? OR original_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if fileType != "" {
		query = query.Where("file_type LIKE ?", fileType+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if uploadType != "" {
		query = query.Where("upload_type = ?", uploadType)
	}
	if minSize > 0 {
		query = query.Where("file_size >= ?", minSize)
	}
	if maxSize > 0 {
		query = query.Where("file_size <= ?", maxSize)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&files).Error

	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// GetByStatus 根据状态获取文件列表
func (r *FileRepository) GetByStatus(status int) ([]model.File, error) {
	var files []model.File
	err := r.db.Where("status = ?", status).Order("created_at ASC").Find(&files).Error
	return files, err
}

// ============ 辅助方法 ============

// UpdateStatus 更新文件使用状态
func (r *FileRepository) UpdateStatus(url string, status int) error {
	return r.db.Model(&model.File{}).
		Where("file_url = ?", url).
		Update("status", status).Error
}

// UpdateFileStatusByUrls 批量更新文件状态
func (r *FileRepository) UpdateFileStatusByUrls(urls []string, status int) error {
	if len(urls) == 0 {
		return nil
	}

	return r.db.Model(&model.File{}).
		Where("file_url IN ?", urls).
		Update("status", status).Error
}

// ============ 维护方法 ============

// DeleteByIDs 批量删除文件记录
func (r *FileRepository) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Unscoped().Delete(&model.File{}, ids).Error
}
