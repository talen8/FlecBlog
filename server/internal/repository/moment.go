package repository

import (
	"context"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// MomentRepository 动态仓储
type MomentRepository struct {
	db *gorm.DB
}

// NewMomentRepository 创建动态仓储
func NewMomentRepository(db *gorm.DB) *MomentRepository {
	return &MomentRepository{db: db}
}

// ============ 基础CRUD ============

// List 获取动态列表
func (r *MomentRepository) List(ctx context.Context, page, pageSize int, isPublish *bool) ([]model.Moment, int64, error) {
	var moments []model.Moment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Moment{})

	// 根据发布状态过滤
	if isPublish != nil {
		query = query.Where("is_publish = ?", *isPublish)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 排序：优先按发布时间，没有发布时间则按创建时间倒序
	query = query.Order("COALESCE(publish_time, created_at) DESC")

	// 分页处理
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err = query.Find(&moments).Error
	if err != nil {
		return nil, 0, err
	}

	return moments, total, nil
}

// Get 获取动态详情
func (r *MomentRepository) Get(ctx context.Context, id uint) (*model.Moment, error) {
	var moment model.Moment
	err := r.db.WithContext(ctx).First(&moment, id).Error
	if err != nil {
		return nil, err
	}
	return &moment, nil
}

// Create 创建动态
func (r *MomentRepository) Create(ctx context.Context, moment *model.Moment) error {
	return r.db.WithContext(ctx).Create(moment).Error
}

// Update 更新动态
func (r *MomentRepository) Update(ctx context.Context, moment *model.Moment) error {
	return r.db.WithContext(ctx).Save(moment).Error
}

// ExistsByContentURL 检查是否有动态内容引用该文件
func (r *MomentRepository) ExistsByContentURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Moment{}).Where("content LIKE ?", "%"+url+"%").Count(&count).Error
	return count > 0, err
}

// Delete 删除动态
func (r *MomentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Moment{}, id).Error
}
