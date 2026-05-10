package repository

import (
	"context"
	"fmt"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// NotificationRepository 通知仓储
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 创建通知仓储实例
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// ============ 查询方法 ============

// List 获取用户特定类型的通知列表
func (r *NotificationRepository) List(ctx context.Context, userID uint, notifTypes []model.NotificationType, offset, limit int) ([]model.UserNotification, int64, error) {
	var userNotifications []model.UserNotification
	var total int64

	// 子查询：获取符合类型条件的通知ID
	subQuery := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Select("id").
		Where("type IN ?", notifTypes)

	// 主查询：根据用户ID和通知类型过滤
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND notification_id IN (?)", userID, subQuery)

	// 获取总数
	if err := query.Model(&model.UserNotification{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表（关联查询通知和发送者信息）
	err := query.
		Preload("Notification.Sender").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&userNotifications).Error

	return userNotifications, total, err
}

// GetUnreadCount 获取用户特定类型的未读通知数量
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uint, notifTypes []model.NotificationType) (int64, error) {
	var count int64

	// 子查询：获取符合类型条件的通知ID
	subQuery := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Select("id").
		Where("type IN ?", notifTypes)

	err := r.db.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("user_id = ? AND is_read = ? AND notification_id IN (?)", userID, false, subQuery).
		Count(&count).Error
	return count, err
}

// ============ 基础 CRUD ============

// Create 创建通知
func (r *NotificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// BatchCreateUserNotifications 批量创建用户通知关联
func (r *NotificationRepository) BatchCreateUserNotifications(ctx context.Context, userNotifications []model.UserNotification) error {
	if len(userNotifications) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&userNotifications).Error
}

// MarkAsRead 标记通知为已读
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// MarkAllAsRead 标记所有通知为已读
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// ============ 辅助方法 ============

// GetAllAdmins 获取所有管理员用户 ID
func (r *NotificationRepository) GetAllAdmins(ctx context.Context) ([]uint, error) {
	var adminIDs []uint
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("role IN ? AND is_enabled = ?", []string{"admin", "super_admin"}, true).
		Pluck("id", &adminIDs).Error
	return adminIDs, err
}

// GetAllSuperAdmins 获取所有超级管理员用户 ID
func (r *NotificationRepository) GetAllSuperAdmins(ctx context.Context) ([]uint, error) {
	var adminIDs []uint
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("role = ? AND is_enabled = ?", model.RoleSuperAdmin, true).
		Pluck("id", &adminIDs).Error
	return adminIDs, err
}

// ExistsVersionUpdateNotification 检查指定版本是否已创建过版本更新通知
func (r *NotificationRepository) ExistsVersionUpdateNotification(ctx context.Context, latestVersion string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("type = ? AND data::jsonb ->> 'alert_type' = ? AND data::jsonb ->> 'latest_version' = ?", model.TypeSystemAlert, model.AlertTypeVersionUpdate, latestVersion).
		Count(&count).Error
	return count > 0, err
}

// GetUserByID 根据ID获取用户信息
func (r *NotificationRepository) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	return &user, err
}

// ExistsAnnouncementNotification 检查指定公告ID是否已创建过通知
func (r *NotificationRepository) ExistsAnnouncementNotification(ctx context.Context, announcementID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("type = ? AND data::jsonb ->> 'alert_type' = ? AND data::jsonb ->> 'announcement_id' = ?",
			model.TypeSystemAlert, model.AlertTypeAnnouncement, fmt.Sprintf("%d", announcementID)).
		Count(&count).Error
	return count > 0, err
}

// GetAllUserIDs 获取所有用户ID
func (r *NotificationRepository) GetAllUserIDs(ctx context.Context) ([]uint, error) {
	var userIDs []uint
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("is_enabled = ?", true).
		Pluck("id", &userIDs).Error
	return userIDs, err
}
