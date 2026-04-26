package repository

import (
	"flec_blog/internal/model"
	"flec_blog/pkg/random"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ============ 基础CRUD ============

// Create 创建新用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Get 获取用户
func (r *UserRepository) Get(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 软删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 获取用户信息
		var user model.User
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}

		// 生成短随机后缀
		suffix := random.Code(4)

		// 更新邮箱，避免唯一索引冲突
		user.Email = user.Email + "_" + suffix
		user.Avatar = ""
		user.IsEnabled = false
		user.Password = ""

		// 保存更新
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// 软删除
		return tx.Delete(&model.User{}, id).Error
	})
}

// ============ 查询方法 ============

// GetByEmail 通过邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// GetGuestByEmail 通过邮箱获取游客用户
func (r *UserRepository) GetGuestByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND role = ?", email, model.RoleGuest).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *UserRepository) List(
	offset, limit int,
	keyword, role string,
	isEnabled, isDeleted *bool,
	loginMethod, lastLoginStart, lastLoginEnd string,
	startTime, endTime string,
) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 是否已删除
	var query *gorm.DB
	if isDeleted != nil {
		if *isDeleted {
			query = r.db.Unscoped().Model(&model.User{}).Where("deleted_at IS NOT NULL")
		} else {
			query = r.db.Model(&model.User{}).Where("deleted_at IS NULL")
		}
	} else {
		query = r.db.Unscoped().Model(&model.User{})
	}

	// 关键词搜索（邮箱、昵称）
	if keyword != "" {
		searchKeyword := "%" + keyword + "%"
		query = query.Where("email ILIKE ? OR nickname ILIKE ?", searchKeyword, searchKeyword)
	}

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 状态筛选
	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	// 登录方式筛选
	if loginMethod != "" {
		switch loginMethod {
		case "password":
			query = query.Where("has_password = ?", true)
		case "github":
			query = query.Where("github_id IS NOT NULL AND github_id != ''")
		case "google":
			query = query.Where("google_id IS NOT NULL AND google_id != ''")
		case "qq":
			query = query.Where("qq_id IS NOT NULL AND qq_id != ''")
		case "microsoft":
			query = query.Where("microsoft_id IS NOT NULL AND microsoft_id != ''")
		}
	}

	// 最后登录时间范围筛选
	if lastLoginStart != "" {
		query = query.Where("last_login >= ?", lastLoginStart)
	}
	if lastLoginEnd != "" {
		query = query.Where("last_login <= ?", lastLoginEnd+" 23:59:59")
	}

	// 注册时间范围筛选
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime+" 23:59:59")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	err = query.
		Select("id, email, nickname, avatar, badge, website, is_enabled, role, last_login, created_at, updated_at, deleted_at, has_password, github_id, google_id, qq_id, feishu_open_id").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ExistsByAvatar 检查是否有用户头像引用该文件
func (r *UserRepository) ExistsByAvatar(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("avatar = ?", url).Count(&count).Error
	return count > 0, err
}

// CountSuperAdmins 统计超级管理员数量
func (r *UserRepository) CountSuperAdmins() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("role = ?", model.RoleSuperAdmin).Count(&count).Error
	return count, err
}

// UpdateAvatar 更新用户头像
func (r *UserRepository) UpdateAvatar(userID uint, avatarURL string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}

// UpdatePasswordAndIncrementVersion 更新密码并递增TokenVersion（一次数据库操作）
// hasPassword > 0 时同时设置 HasPassword=true（用于 OAuth 用户首次设置密码）
func (r *UserRepository) UpdatePasswordAndIncrementVersion(id uint, hashedPassword string, hasPassword ...bool) error {
	updates := map[string]interface{}{
		"password":      hashedPassword,
		"token_version": gorm.Expr("token_version + 1"),
	}
	if len(hasPassword) > 0 && hasPassword[0] {
		updates["has_password"] = true
	}
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// ============ Token黑名单 ============

// AddTokenToBlacklist 添加token到黑名单
func (r *UserRepository) AddTokenToBlacklist(tokenHash string, userID uint, expiresAt time.Time) error {
	blacklist := &model.TokenBlacklist{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(blacklist).Error
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (r *UserRepository) IsTokenBlacklisted(tokenHash string) bool {
	var count int64
	r.db.Model(&model.TokenBlacklist{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		Count(&count)
	return count > 0
}

// CleanupExpiredTokens 清理过期的黑名单记录
func (r *UserRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.TokenBlacklist{}).Error
}

// IncrementTokenVersion 递增用户Token版本号（用于撤销所有token）
func (r *UserRepository) IncrementTokenVersion(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("token_version", gorm.Expr("token_version + 1")).Error
}

// ============ OAuth 相关 ============

// GetByOAuthID 通过 OAuth ID 获取用户
func (r *UserRepository) GetByOAuthID(provider, providerID string) (*model.User, error) {
	var user model.User
	var query string

	// 根据提供商选择查询字段
	switch provider {
	case "github":
		query = "github_id = ?"
	case "google":
		query = "google_id = ?"
	case "qq":
		query = "qq_id = ?"
	case "microsoft":
		query = "microsoft_id = ?"
	case "feishu":
		query = "feishu_open_id = ?"
	default:
		return nil, gorm.ErrRecordNotFound
	}

	// 执行查询
	err := r.db.Where(query, providerID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateOAuthBinding 更新 OAuth 绑定
func (r *UserRepository) UpdateOAuthBinding(userID uint, provider, providerID string) error {
	// 根据提供商选择更新字段
	switch provider {
	case "github":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("github_id", providerID).Error
	case "google":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("google_id", providerID).Error
	case "qq":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("qq_id", providerID).Error
	case "microsoft":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("microsoft_id", providerID).Error
	case "feishu":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("feishu_open_id", providerID).Error
	default:
		return gorm.ErrInvalidData
	}
}

// ClearOAuthBinding 清除 OAuth 绑定
func (r *UserRepository) ClearOAuthBinding(userID uint, provider string) error {
	// 根据提供商选择清除字段
	switch provider {
	case "github":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("github_id", "").Error
	case "google":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("google_id", "").Error
	case "qq":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("qq_id", "").Error
	case "microsoft":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("microsoft_id", "").Error
	default:
		return gorm.ErrInvalidData
	}
}
