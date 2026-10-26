package repository

import (
	"context"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// TagRepository 标签仓储
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// ============ 基础CRUD ============

// List 获取标签列表
func (r *TagRepository) List(ctx context.Context, page, pageSize int) ([]model.Tag, int64, error) {
	var tags []model.Tag
	var total int64

	err := r.db.WithContext(ctx).Model(&model.Tag{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Order("created_at DESC")

	// 分页处理
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err = query.Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// GetBySlug 根据slug获取标签
func (r *TagRepository) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create 创建标签
func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// Update 更新标签
func (r *TagRepository) Update(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete 删除标签
func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	// 删除文章-标签关联
	r.db.WithContext(ctx).Exec("DELETE FROM article_tags WHERE tag_id = ?", id)
	return r.db.WithContext(ctx).Unscoped().Delete(&model.Tag{}, id).Error
}
