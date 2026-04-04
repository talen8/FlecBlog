package repository

import (
	"context"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// CommentRepository 评论仓储
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ============ 查询方法 ============

// GetByTarget 获取目标的顶级评论列表
func (r *CommentRepository) GetByTarget(ctx context.Context, targetType, targetKey string, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 查询顶级评论（包含已删除的记录）
	query := r.db.WithContext(ctx).Unscoped().Model(&model.Comment{}).
		Where("target_type = ? AND target_key = ? AND root_id IS NULL", targetType, targetKey)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		})

	// 只有当 page 和 pageSize 都大于 0 时才应用分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err := query.Find(&comments).Error
	return comments, total, err
}

// GetRepliesByRootIDs 批量获取多个根评论的所有回复
func (r *CommentRepository) GetRepliesByRootIDs(ctx context.Context, rootIDs []uint) ([]model.Comment, error) {
	if len(rootIDs) == 0 {
		return []model.Comment{}, nil
	}

	var replies []model.Comment
	err := r.db.WithContext(ctx).Unscoped().
		Where("root_id IN ?", rootIDs).
		Order("created_at ASC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		Preload("ReplyUser", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		Find(&replies).Error

	return replies, err
}

// GetForWeb 获取评论详情（前台）
func (r *CommentRepository) GetForWeb(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		Preload("ReplyUser", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// List 获取评论列表（后台管理）
func (r *CommentRepository) List(ctx context.Context, offset, limit int, status *int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 包含软删除的记录
	query := r.db.WithContext(ctx).Unscoped().Model(&model.Comment{})

	// 状态过滤
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, email, nickname, avatar, badge, deleted_at")
		}).
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	return comments, total, err
}

// Get 获取评论详情（后台管理）
func (r *CommentRepository) Get(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).Unscoped().
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, email, nickname, avatar, badge, deleted_at")
		}).
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}


// ExistsByContentURL 检查是否有评论内容引用该文件（包含软删除评论，便于恢复后仍可展示）
func (r *CommentRepository) ExistsByContentURL(url string) (bool, error) {
	var count int64
	err := r.db.Unscoped().Model(&model.Comment{}).
		Where("status = ? AND content LIKE ?", 1, "%"+url+"%").
		Count(&count).Error
	return count > 0, err
}

// Create 创建评论
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// Update 更新评论
func (r *CommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// UpdateStatus 更新评论状态
func (r *CommentRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 软删除评论
func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	// 只删除评论本身，子评论保留
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

// Restore 恢复已删除的评论
func (r *CommentRepository) Restore(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().
		Model(&model.Comment{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// CountByTargetKeys 批量获取多个目标的评论数
func (r *CommentRepository) CountByTargetKeys(ctx context.Context, targetType string, targetKeys []string) (map[string]int64, error) {
	if len(targetKeys) == 0 {
		return make(map[string]int64), nil
	}

	type CountResult struct {
		TargetKey string
		Count     int64
	}

	var results []CountResult
	err := r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Select("target_key, COUNT(*) as count").
		Where("target_type = ? AND target_key IN ? AND status = 1 AND deleted_at IS NULL", targetType, targetKeys).
		Group("target_key").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 构建映射表
	countMap := make(map[string]int64, len(targetKeys))
	for _, key := range targetKeys {
		countMap[key] = 0
	}
	for _, result := range results {
		countMap[result.TargetKey] = result.Count
	}

	return countMap, nil
}
