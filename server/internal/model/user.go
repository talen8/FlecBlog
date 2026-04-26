package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色枚举
type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin" // 超级管理员
	RoleAdmin      UserRole = "admin"       // 管理员
	RoleUser       UserRole = "user"        // 普通用户
	RoleGuest      UserRole = "guest"       // 游客
)

// User 用户模型
type User struct {
	gorm.Model
	Email        string     `gorm:"size:100" json:"email"`
	Password     string     `gorm:"size:100" json:"-"`
	HasPassword  bool       `gorm:"default:false" json:"has_password"` // 是否设置了密码
	Nickname     string     `gorm:"size:50" json:"nickname"`
	Avatar       string     `gorm:"size:255" json:"avatar"`
	Badge        string     `gorm:"size:50" json:"badge"`           // 铭牌标识
	Website      string     `gorm:"size:255" json:"website"`        // 用户网站地址
	IsEnabled    bool       `gorm:"default:true" json:"is_enabled"` // 是否启用
	Role         UserRole   `gorm:"default:'user'" json:"role"`
	LastLogin    *time.Time `json:"last_login"`             // 最后登录时间，未登录为 null
	TokenVersion uint       `gorm:"default:0" json:"-"`     // Token版本号，用于撤销所有token
	GithubID     string     `gorm:"size:50;index" json:"-"` // GitHub 用户ID
	GoogleID     string     `gorm:"size:50;index" json:"-"` // Google 用户ID
	QQID         string     `gorm:"size:50;index" json:"-"` // QQ OpenID
	FeishuOpenID string     `gorm:"size:50;index" json:"-"` // 飞书 OpenID
	MicrosoftID  string     `gorm:"size:50;index" json:"-"` // Microsoft 用户ID
}

// TokenBlacklist Token黑名单模型
type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey"`
	TokenHash string    `gorm:"type:varchar(64);uniqueIndex;not null"` // token的SHA256哈希
	UserID    uint      `gorm:"index;not null"`                        // 关联的用户ID
	ExpiresAt time.Time `gorm:"index;not null"`                        // 过期时间
}
