package repository

import (
	"flec_blog/internal/model"
	"time"

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
func (r *FileRepository) List(offset, limit int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	if err := r.db.Model(&model.File{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&files).Error

	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// FileListFilter 文件列表筛选条件
type FileListFilter struct {
	Keyword    string
	FileType   string
	Status     *int
	UploadType string
	MinSize    int64
	MaxSize    int64
	StartTime  string
	EndTime    string
}

// GetByFilter 根据筛选条件获取文件列表
func (r *FileRepository) GetByFilter(filter *FileListFilter, offset, limit int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.db.Model(&model.File{})

	if filter.Keyword != "" {
		query = query.Where("file_name LIKE ? OR original_name LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if filter.FileType != "" {
		query = query.Where("file_type LIKE ?", filter.FileType+"%")
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.UploadType != "" {
		query = query.Where("upload_type = ?", filter.UploadType)
	}
	if filter.MinSize > 0 {
		query = query.Where("file_size >= ?", filter.MinSize)
	}
	if filter.MaxSize > 0 {
		query = query.Where("file_size <= ?", filter.MaxSize)
	}
	if filter.StartTime != "" {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if filter.EndTime != "" {
		query = query.Where("created_at <= ?", filter.EndTime+" 23:59:59")
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

// GetUnusedFiles 获取超过指定天数未使用的文件
func (r *FileRepository) GetUnusedFiles(days int) ([]model.File, error) {
	var files []model.File
	cutoffTime := time.Now().AddDate(0, 0, -days)

	err := r.db.Where("status = ? AND created_at < ?", 0, cutoffTime).
		Find(&files).Error

	return files, err
}

// DeleteByIDs 批量删除文件记录
func (r *FileRepository) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Unscoped().Delete(&model.File{}, ids).Error
}
