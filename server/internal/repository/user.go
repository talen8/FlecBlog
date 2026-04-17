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

// List 获取用户列表（后台管理）
func (r *UserRepository) List(offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 包含软删除的用户
	err := r.db.Unscoped().Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Unscoped().
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

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
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

// RevokeAllUserTokens 撤销某用户的所有token
func (r *UserRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Delete(&model.TokenBlacklist{}).Error
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
