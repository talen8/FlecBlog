package dto

import "flec_blog/pkg/utils"

// ============ 通用评论请求 ============

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content    string `json:"content" binding:"required,min=1,max=500"`
	TargetType string `json:"target_type" binding:"required,oneof=article page"`
	TargetKey  string `json:"target_key" binding:"required"`
	ParentID   *uint  `json:"parent_id"`

	// 游客信息（可选，未登录时必填）
	Nickname string `json:"nickname,omitempty" binding:"omitempty,min=2,max=32"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Website  string `json:"website,omitempty" binding:"omitempty,url,max=255"`

	// 客户端信息（由 handler 层填充）
	IP        string `json:"-"` // 不从请求体获取
	UserAgent string `json:"-"` // 不从请求体获取
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// ============ 通用评论响应 ============

// CommentResponse 评论列表响应（前台用）
type CommentResponse struct {
	ID        uint           `json:"id"`
	Content   string         `json:"content"`
	IsDeleted bool           `json:"is_deleted"`
	ParentID  *uint          `json:"parent_id"`
	CreatedAt utils.JSONTime `json:"created_at"`
	Location  string         `json:"location,omitempty"` // 地理位置
	Browser   string         `json:"browser,omitempty"`  // 浏览器内核
	OS        string         `json:"os,omitempty"`       // 操作系统
	User      struct {
		ID        uint   `json:"id"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Badge     string `json:"badge"`      // 铭牌
		EmailHash string `json:"email_hash"` // 邮箱哈希（用于 Gravatar）
		Website   string `json:"website"`    // 网站
		Role      string `json:"role"`       // 角色
	} `json:"user"`
	ReplyUser *struct {
		ID        uint   `json:"id"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Badge     string `json:"badge"`      // 铭牌
		EmailHash string `json:"email_hash"` // 邮箱哈希（用于 Gravatar）
		Website   string `json:"website"`    // 网站
		Role      string `json:"role"`       // 角色
	} `json:"reply_user,omitempty"`
	Replies []CommentResponse `json:"replies"`
}

// ============ 前台评论请求 ============

// CommentQueryForWebRequest 前台评论查询请求
type CommentQueryForWebRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=10" binding:"min=0"`
	TargetType string `form:"target_type" binding:"required,oneof=article page"`
	TargetKey  string `form:"target_key" binding:"required"`
}

// ============ 后台评论请求 ============

// CommentQueryRequest 后台评论查询请求
type CommentQueryRequest struct {
	Page     int  `form:"page,default=1" binding:"min=1"`
	PageSize int  `form:"page_size,default=10" binding:"min=1,max=100"`
	Status   *int `form:"status"`
}

// ============ 后台评论响应 ============

// CommentListResponse 后台评论列表响应
type CommentListResponse struct {
	ID        uint            `json:"id"`
	Content   string          `json:"content"`
	Status    int             `json:"status"`
	ParentID  *uint           `json:"parent_id"`
	CreatedAt utils.JSONTime  `json:"created_at"`
	DeletedAt *utils.JSONTime `json:"deleted_at,omitempty"`
	Target    struct {
		Type  string `json:"type"`
		Key   string `json:"key"`
		Title string `json:"title"`
	} `json:"target"`
	User struct {
		ID       uint   `json:"id"`
		Email    string `json:"email"` // 后台显示邮箱
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Badge    string `json:"badge"` // 铭牌
	} `json:"user"`
}

// ============ 数据导入导出 ============

// ImportCommentsRequest 导入评论请求
type ImportCommentsRequest struct {
	SourceType string `json:"source_type" form:"source_type" binding:"required,oneof=artalk disqus wordpress"`
}

// ImportCommentsResult 导入评论结果
type ImportCommentsResult struct {
	Total       int                  `json:"total"`
	Success     int                  `json:"success"`
	Failed      int                  `json:"failed"`
	UserCreated int                  `json:"user_created"`
	Errors      []ImportCommentError `json:"errors,omitempty"`
}

// ImportCommentError 导入错误信息
type ImportCommentError struct {
	Index   int    `json:"index"`
	Content string `json:"content"`
	Error   string `json:"error"`
}

// ArtalkCommentData Artalk评论数据格式（完整版，兼容真实导出）
type ArtalkCommentData struct {
	ID            string `json:"id"`
	RID           string `json:"rid"`
	Content       string `json:"content"`
	UA            string `json:"ua"`
	IP            string `json:"ip"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	IsCollapsed   string `json:"is_collapsed"`
	IsPending     string `json:"is_pending"`
	IsPinned      string `json:"is_pinned"` // 置顶标记
	VoteUp        string `json:"vote_up"`
	VoteDown      string `json:"vote_down"`
	Nick          string `json:"nick"`
	Email         string `json:"email"`
	Link          string `json:"link"`
	Password      string `json:"password,omitempty"` // 旧版本字段，可选
	BadgeName     string `json:"badge_name"`
	BadgeColor    string `json:"badge_color"`
	PageKey       string `json:"page_key"`
	PageTitle     string `json:"page_title"`
	PageAdminOnly string `json:"page_admin_only"`
	SiteName      string `json:"site_name"`
	SiteURLs      string `json:"site_urls"`
}
