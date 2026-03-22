package repository

import (
	"context"
	"time"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// RssFeedRepository RSS订阅仓储
type RssFeedRepository struct {
	db *gorm.DB
}

// NewRssFeedRepository 创建RSS订阅仓储
func NewRssFeedRepository(db *gorm.DB) *RssFeedRepository {
	return &RssFeedRepository{db: db}
}

// List 获取RSS文章列表
func (r *RssFeedRepository) List(ctx context.Context, page, pageSize int) ([]model.RssArticle, int64, error) {
	var articles []model.RssArticle
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RssArticle{}).Preload("Friend")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("published_at DESC NULLS LAST, created_at DESC")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByID 根据ID获取RSS文章
func (r *RssFeedRepository) GetByID(ctx context.Context, id uint) (*model.RssArticle, error) {
	var article model.RssArticle
	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateBatch 批量创建RSS文章（忽略重复）
func (r *RssFeedRepository) CreateBatch(ctx context.Context, articles []model.RssArticle) error {
	if len(articles) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(articles, 100).Error
}

// MarkRead 标记文章已读
func (r *RssFeedRepository) MarkRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllRead 全部标记已读
func (r *RssFeedRepository) MarkAllRead(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("is_read = ?", false).Update("is_read", true)
	return result.RowsAffected, result.Error
}

// CountUnread 统计未读数量
func (r *RssFeedRepository) CountUnread(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("is_read = ?", false).Count(&count).Error
	return count, err
}

// ListUnread 获取未读文章列表
func (r *RssFeedRepository) ListUnread(ctx context.Context, limit int) ([]model.RssArticle, error) {
	var articles []model.RssArticle
	query := r.db.WithContext(ctx).Model(&model.RssArticle{}).
		Preload("Friend").
		Where("is_read = ?", false).
		Order("published_at DESC NULLS LAST, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

// DeleteOrphaned 删除孤立文章（友链不存在、失效或RSS地址为空）
func (r *RssFeedRepository) DeleteOrphaned(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("friend_id NOT IN (SELECT id FROM friends) OR friend_id IN (SELECT id FROM friends WHERE is_invalid = true OR rss_url = '' OR rss_url IS NULL)").
		Delete(&model.RssArticle{})
	return result.RowsAffected, result.Error
}

// ExistsByLink 检查链接是否已存在
func (r *RssFeedRepository) ExistsByLink(ctx context.Context, link string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("link = ?", link).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFriendsWithRSS 获取所有配置了RSS且未失效的友链
func (r *RssFeedRepository) GetFriendsWithRSS(ctx context.Context) ([]model.Friend, error) {
	var friends []model.Friend
	err := r.db.WithContext(ctx).Where("rss_url != '' AND rss_url IS NOT NULL AND is_invalid = ?", false).Find(&friends).Error
	return friends, err
}

// UpdateFriendRSSLatime 更新友链RSS最后更新时间
func (r *RssFeedRepository) UpdateFriendRSSLatime(ctx context.Context, friendID uint, latime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Friend{}).Where("id = ?", friendID).Update("rss_latime", latime).Error
}

// GetLatestPublishedTime 获取指定友链最新文章的发布时间
func (r *RssFeedRepository) GetLatestPublishedTime(ctx context.Context, friendID uint) (*time.Time, error) {
	var article model.RssArticle
	err := r.db.WithContext(ctx).
		Where("friend_id = ?", friendID).
		Where("published_at IS NOT NULL").
		Order("published_at DESC").
		Limit(1).
		Find(&article).Error
	if err != nil {
		return nil, err
	}
	if article.ID == 0 {
		return nil, nil
	}
	return article.PublishedAt, nil
}
