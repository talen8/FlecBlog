package dto

import (
	"flec_blog/internal/model"
	"flec_blog/pkg/utils"
)

// ============ 通用用户响应 ============

// UserResponse 用户信息响应
type UserResponse struct {
	ID             uint            `json:"id"`
	Email          string          `json:"email"`
	EmailHash      string          `json:"email_hash"`       // 邮箱哈希（用于 Gravatar）
	IsVirtualEmail bool            `json:"is_virtual_email"` // 是否为虚拟邮箱（需绑定真实邮箱）
	Nickname       string          `json:"nickname"`
	Avatar         string          `json:"avatar"`
	Badge          string          `json:"badge"`
	Website        string          `json:"website"`
	Role           model.UserRole  `json:"role"`
	HasPassword    bool            `json:"has_password"`
	LinkedOAuths   []string        `json:"linked_oauths"` // ["github", "google", "qq"]
	LastLogin      *utils.JSONTime `json:"last_login,omitempty"`
	CreatedAt      utils.JSONTime  `json:"created_at"`
	IsEnabled      bool            `json:"is_enabled"`
}

// NewUserResponse 从model.User创建UserResponse
func NewUserResponse(user *model.User) *UserResponse {
	linkedOAuths := make([]string, 0)
	if user.GithubID != "" {
		linkedOAuths = append(linkedOAuths, "github")
	}
	if user.GoogleID != "" {
		linkedOAuths = append(linkedOAuths, "google")
	}
	if user.QQID != "" {
		linkedOAuths = append(linkedOAuths, "qq")
	}
	if user.MicrosoftID != "" {
		linkedOAuths = append(linkedOAuths, "microsoft")
	}
	if user.FeishuOpenID != "" {
		linkedOAuths = append(linkedOAuths, "feishu")
	}

	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		EmailHash:      utils.GetEmailHash(user.Email),
		IsVirtualEmail: utils.IsVirtualEmail(user.Email),
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Badge:          user.Badge,
		Website:        user.Website,
		Role:           user.Role,
		HasPassword:    user.HasPassword,
		LinkedOAuths:   linkedOAuths,
		LastLogin:      utils.ToJSONTime(user.LastLogin),
		CreatedAt:      utils.NewJSONTime(user.CreatedAt),
	}
}

// ============ 前台用户请求 ============

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
	Website  string `json:"website,omitempty" binding:"omitempty,url,max=255"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 用户信息更新请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname,omitempty" binding:"omitempty,min=2,max=32"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Avatar   string `json:"avatar,omitempty" binding:"omitempty,max=255"`
	Badge    string `json:"badge,omitempty" binding:"omitempty,max=50"`
	Website  string `json:"website,omitempty" binding:"omitempty,url,max=255"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// SetPasswordRequest 设置密码请求（针对 OAuth 注册用户首次设置密码）
type SetPasswordRequest struct {
	Password        string `json:"password" binding:"required,min=6,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=32"`
}

// DeactivateAccountRequest 注销账号请求
type DeactivateAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

// ============ 前台用户响应 ============

// LoginResponse 用户登录/注册响应
type LoginResponse struct {
	AccessToken string        `json:"access_token"`
	User        *UserResponse `json:"user,omitempty"`
}

// ============ 后台用户管理请求 ============

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page           int    `form:"page,default=1" binding:"min=1"`
	PageSize       int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword        string `form:"keyword"`          // 搜索关键词（邮箱、昵称）
	Role           string `form:"role"`             // 角色筛选
	IsEnabled      *bool  `form:"is_enabled"`       // 状态筛选
	IsDeleted      *bool  `form:"is_deleted"`       // 是否已删除
	LoginMethod    string `form:"login_method"`     // 登录方式筛选（password/github/google/qq/microsoft）
	LastLoginStart string `form:"last_login_start"` // 最后登录开始时间（格式：2006-01-02）
	LastLoginEnd   string `form:"last_login_end"`   // 最后登录结束时间（格式：2006-01-02）
	StartTime      string `form:"start_time"`       // 注册开始时间（格式：2006-01-02）
	EndTime        string `form:"end_time"`         // 注册结束时间（格式：2006-01-02）
}

// AdminCreateUserRequest 管理员创建用户请求
type AdminCreateUserRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6,max=32"`
	Nickname string         `json:"nickname" binding:"required,min=2,max=32"`
	Avatar   string         `json:"avatar,omitempty" binding:"omitempty,max=255"`
	Badge    string         `json:"badge,omitempty" binding:"omitempty,max=50"`
	Website  string         `json:"website,omitempty" binding:"omitempty,url,max=255"`
	Role     model.UserRole `json:"role" binding:"required,oneof=super_admin admin user guest"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Email     string         `json:"email,omitempty" binding:"omitempty,email"`
	Nickname  string         `json:"nickname,omitempty" binding:"omitempty,min=2,max=32"`
	Avatar    string         `json:"avatar,omitempty" binding:"omitempty,max=255"`
	Badge     string         `json:"badge,omitempty" binding:"omitempty,max=50"`
	Website   string         `json:"website,omitempty" binding:"omitempty,url,max=255"`
	Role      model.UserRole `json:"role,omitempty" binding:"omitempty,oneof=super_admin admin user guest"`
	IsEnabled *bool          `json:"is_enabled,omitempty"`                                // 使用指针以区分未设置和false
	Password  string         `json:"password,omitempty" binding:"omitempty,min=6,max=32"` // 可选的密码字段，用于更新密码
}

// ============ 后台用户管理响应 ============

// UserListResponse 用户列表响应
type UserListResponse struct {
	ID           uint            `json:"id"`
	Email        string          `json:"email"`
	Nickname     string          `json:"nickname"`
	Avatar       string          `json:"avatar"`
	Badge        string          `json:"badge"`
	Website      string          `json:"website"`
	Role         model.UserRole  `json:"role"`
	IsEnabled    bool            `json:"is_enabled"`
	LastLogin    *utils.JSONTime `json:"last_login,omitempty"`
	CreatedAt    utils.JSONTime  `json:"created_at"`
	DeletedAt    *utils.JSONTime `json:"deleted_at,omitempty"`
	HasPassword  bool            `json:"has_password"`   // 是否设置了密码
	GithubID     string          `json:"github_id"`      // GitHub ID
	GoogleID     string          `json:"google_id"`      // Google ID
	QQID         string          `json:"qq_id"`          // QQ OpenID
	MicrosoftID  string          `json:"microsoft_id"`   // Microsoft ID
	FeishuOpenID string          `json:"feishu_open_id"` // 飞书 OpenID
}
