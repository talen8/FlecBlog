package repository

import (
	"context"
	"time"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// FriendRepository 友链仓储
type FriendRepository struct {
	db *gorm.DB
}

// NewFriendRepository 创建友链仓储
func NewFriendRepository(db *gorm.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

// ============ 友链类型 ============

// ListTypes 获取所有友链类型
func (r *FriendRepository) ListTypes(ctx context.Context, page, pageSize int) ([]model.FriendType, int64, error) {
	var types []model.FriendType
	var total int64

	query := r.db.WithContext(ctx).Model(&model.FriendType{})

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order("sort DESC, id ASC")

	// 分页处理
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err = query.Find(&types).Error
	if err != nil {
		return nil, 0, err
	}

	return types, total, nil
}

// CountFriendsByTypeID 统计某个类型下的友链数量
func (r *FriendRepository) CountFriendsByTypeID(ctx context.Context, typeID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Friend{}).Where("type = ?", typeID).Count(&count).Error
	return count, err
}

// GetTypeByID 根据 ID 获取友链类型
func (r *FriendRepository) GetTypeByID(ctx context.Context, id uint) (*model.FriendType, error) {
	var friendType model.FriendType
	err := r.db.WithContext(ctx).First(&friendType, id).Error
	if err != nil {
		return nil, err
	}
	return &friendType, nil
}

// CreateType 创建友链类型
func (r *FriendRepository) CreateType(ctx context.Context, friendType *model.FriendType) error {
	return r.db.WithContext(ctx).Create(friendType).Error
}

// UpdateType 更新友链类型
func (r *FriendRepository) UpdateType(ctx context.Context, friendType *model.FriendType) error {
	return r.db.WithContext(ctx).Save(friendType).Error
}

// DeleteType 删除友链类型
func (r *FriendRepository) DeleteType(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.FriendType{}, id).Error
}

// ============ 基础CRUD ============

// List 获取友链列表
func (r *FriendRepository) List(ctx context.Context, page, pageSize int, keyword string, typeID *uint, isInvalid, isPending *bool, accessibleStatus string, rssStatus string, hasScreenshot *bool) ([]model.Friend, int64, error) {
	var friends []model.Friend
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Friend{})

	// 关键词搜索（名称、链接、描述）
	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("friends.name ILIKE ? OR friends.url ILIKE ? OR friends.description ILIKE ?", pattern, pattern, pattern)
	}

	// 按类型筛选
	if typeID != nil {
		query = query.Where("friends.type = ?", *typeID)
	}

	// 按失效状态筛选
	if isInvalid != nil {
		query = query.Where("friends.is_invalid = ?", *isInvalid)
	}

	// 按待审核状态筛选
	if isPending != nil {
		query = query.Where("friends.is_pending = ?", *isPending)
	}

	// 按可访问性状态筛选
	if accessibleStatus != "" {
		switch accessibleStatus {
		case "normal":
			query = query.Where("friends.accessible = ?", 0)
		case "abnormal":
			query = query.Where("friends.accessible > ?", 0)
		case "ignored":
			query = query.Where("friends.accessible = ?", -1)
		}
	}

	// 按RSS状态筛选
	if rssStatus != "" {
		switch rssStatus {
		case "no_rss":
			// 无RSS订阅
			query = query.Where("friends.rss_url IS NULL OR friends.rss_url = ''")
		case "normal":
			// 正常订阅（3个月内更新）
			threeMonthsAgo := time.Now().AddDate(0, -3, 0)
			query = query.Where("friends.rss_url IS NOT NULL AND friends.rss_url != '' AND friends.rss_latime >= ?", threeMonthsAgo)
		case "warning":
			// 三个月未更新
			threeMonthsAgo := time.Now().AddDate(0, -3, 0)
			sixMonthsAgo := time.Now().AddDate(0, -6, 0)
			query = query.Where("friends.rss_url IS NOT NULL AND friends.rss_url != '' AND friends.rss_latime < ? AND friends.rss_latime >= ?", threeMonthsAgo, sixMonthsAgo)
		case "danger":
			// 六个月未更新
			sixMonthsAgo := time.Now().AddDate(0, -6, 0)
			query = query.Where("friends.rss_url IS NOT NULL AND friends.rss_url != '' AND friends.rss_latime < ?", sixMonthsAgo)
		}
	}

	// 按是否包含截图筛选
	if hasScreenshot != nil {
		if *hasScreenshot {
			query = query.Where("friends.screenshot IS NOT NULL AND friends.screenshot != ''")
		} else {
			query = query.Where("friends.screenshot IS NULL OR friends.screenshot = ''")
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 预加载类型信息（用于后台显示类型名称）
	query = query.Preload("Type")

	// 后台排序：待审核优先 > 类型排序(大到小) > 友链排序(大到小) > 添加时间(旧到新) > 失效的在最后
	// 需要 LEFT JOIN 友链类型表来按类型排序值排序
	query = query.Select("friends.*").
		Joins("LEFT JOIN friend_types ON friends.type = friend_types.id").
		Order("friends.is_pending DESC, friends.is_invalid ASC, friend_types.sort DESC NULLS LAST, friends.sort DESC, friends.created_at ASC")

	// 分页处理
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err = query.Find(&friends).Error
	if err != nil {
		return nil, 0, err
	}

	return friends, total, nil
}

// Get 获取友链
func (r *FriendRepository) Get(ctx context.Context, id uint) (*model.Friend, error) {
	var friend model.Friend
	err := r.db.WithContext(ctx).First(&friend, id).Error
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

// Create 创建友链
func (r *FriendRepository) Create(ctx context.Context, friend *model.Friend) error {
	return r.db.WithContext(ctx).Create(friend).Error
}

// Update 更新友链
func (r *FriendRepository) Update(ctx context.Context, friend *model.Friend) error {
	return r.db.WithContext(ctx).Save(friend).Error
}

// Delete 删除友链
func (r *FriendRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&model.Friend{}, id).Error
}

// ============ 前台友链查询 ============

// ListForWeb 获取前台友链数据
func (r *FriendRepository) ListForWeb(ctx context.Context) ([]model.FriendType, []model.Friend, error) {
	var types []model.FriendType
	var friends []model.Friend

	// 查询展示的类型
	if err := r.db.WithContext(ctx).
		Where("is_visible = ?", true).
		Order("sort DESC, id ASC").
		Find(&types).Error; err != nil {
		return nil, nil, err
	}

	// 查询友链（排除待审核的申请）
	if err := r.db.WithContext(ctx).
		Preload("Type").
		Joins("LEFT JOIN friend_types ON friends.type = friend_types.id").
		Where("friends.is_pending = ? AND (friends.type IS NULL OR friend_types.is_visible = ?)", false, true).
		Order("friends.sort DESC, friends.id ASC").
		Find(&friends).Error; err != nil {
		return nil, nil, err
	}

	return types, friends, nil
}

// ============ 友链检测 ============

// GetAllForCheck 获取所有友链用于检测（排除待审核、失效和忽略检查的）
func (r *FriendRepository) GetAllForCheck(ctx context.Context) ([]model.Friend, error) {
	var friends []model.Friend
	err := r.db.WithContext(ctx).
		Where("is_pending = ? AND is_invalid = ? AND accessible != ?", false, false, -1).
		Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// ExistsByAvatarOrScreenshot 检查是否有友链引用该文件
func (r *FriendRepository) ExistsByAvatarOrScreenshot(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Friend{}).
		Where("avatar = ? OR screenshot = ?", url, url).
		Count(&count).Error
	return count > 0, err
}

// UpdateCheckStatus 更新友链检测状态
func (r *FriendRepository) UpdateCheckStatus(ctx context.Context, id uint, accessible int) error {
	return r.db.WithContext(ctx).
		Model(&model.Friend{}).
		Where("id = ?", id).
		Update("accessible", accessible).Error
}
