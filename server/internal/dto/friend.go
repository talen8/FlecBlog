package dto

import "flec_blog/pkg/utils"

// ============ 友链类型 ============

// ListFriendTypeRequest 友链类型列表请求
type ListFriendTypeRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=1000"`
}

// FriendTypeResponse 友链类型响应（前台用）
type FriendTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Sort int    `json:"sort"`
}

// CreateFriendTypeRequest 创建友链类型请求
type CreateFriendTypeRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Sort      int    `json:"sort" binding:"omitempty"`
	IsVisible *bool  `json:"is_visible"`
}

// UpdateFriendTypeRequest 更新友链类型请求
type UpdateFriendTypeRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Sort      int    `json:"sort" binding:"omitempty"`
	IsVisible *bool  `json:"is_visible"`
}

// FriendTypeListResponse 后台友链类型列表响应
type FriendTypeListResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Sort      int    `json:"sort"`
	IsVisible bool   `json:"is_visible"`
	Count     int    `json:"count"` // 该类型下的友链数量
}

// ============ 分组友链响应 ============

// FriendGroupResponse 友链分组响应
type FriendGroupResponse struct {
	TypeID   *uint                   `json:"type_id"`   // 类型ID，未分类为 null
	TypeName string                  `json:"type_name"` // 类型名称
	TypeSort int                     `json:"type_sort"` // 类型排序值
	Friends  []FriendInGroupResponse `json:"friends"`   // 该类型下的友链（已排序）
}

// FriendInGroupResponse 分组中的友链响应（精简版）
type FriendInGroupResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Screenshot  string `json:"screenshot"`
	Sort        int    `json:"sort"`
	IsInvalid   bool   `json:"is_invalid"` // 是否失效
}

// GroupedFriendsResponse 分组友链列表响应
type GroupedFriendsResponse struct {
	Groups       []FriendGroupResponse `json:"groups"`        // 友链分组列表
	TotalGroups  int                   `json:"total_groups"`  // 分组总数
	TotalFriends int                   `json:"total_friends"` // 友链总数
}

// ============ 前台友链请求 ============

// ListFriendForWebRequest 友链列表请求
type ListFriendForWebRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=1000"`
}

// ============ 前台友链响应 ============

// FriendForWebResponse 前台友链响应
type FriendForWebResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Screenshot  string `json:"screenshot"`
	Sort        int    `json:"sort"`
	IsInvalid   bool   `json:"is_invalid"` // 是否失效
	TypeID      *uint  `json:"type_id"`    // 类型 ID
}

// ============ 后台友链管理请求 ============

// ListFriendRequest 友链列表请求
type ListFriendRequest struct {
	Page             int    `form:"page" binding:"omitempty,min=1"`
	PageSize         int    `form:"page_size" binding:"omitempty,min=1,max=1000"`
	Keyword          string `form:"keyword" binding:"omitempty,max=100"`                                 // 关键词搜索（名称、链接、描述）
	TypeID           *uint  `form:"type_id" binding:"omitempty"`                                         // 友链类型ID
	IsInvalid        *bool  `form:"is_invalid" binding:"omitempty"`                                      // 是否失效
	IsPending        *bool  `form:"is_pending" binding:"omitempty"`                                      // 是否待审核
	AccessibleStatus string `form:"accessible_status" binding:"omitempty,oneof=normal abnormal ignored"` // 可访问性状态: normal=正常, abnormal=异常, ignored=忽略检查
	RSSStatus        string `form:"rss_status" binding:"omitempty,oneof=no_rss normal warning danger"`   // RSS状态: no_rss=无订阅, normal=正常订阅, warning=三个月未更新, danger=六个月未更新
	HasScreenshot    *bool  `form:"has_screenshot" binding:"omitempty"`                                  // 是否包含截图
}

// CreateFriendRequest 创建友链请求
type CreateFriendRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	URL         string `json:"url" binding:"required,url,max=255"`
	Description string `json:"description" binding:"max=500"`
	Avatar      string `json:"avatar" binding:"omitempty,url,max=255"`
	Screenshot  string `json:"screenshot" binding:"omitempty,url,max=255"`
	Sort        int    `json:"sort" binding:"omitempty,min=1,max=10"`
	TypeID      *uint  `json:"type_id" binding:"omitempty"`
	RSSUrl      string `json:"rss_url" binding:"omitempty,url,max=500"` // RSS订阅地址
}

// UpdateFriendRequest 更新友链请求
type UpdateFriendRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	URL         string `json:"url" binding:"required,url,max=255"`
	Description string `json:"description" binding:"max=500"`
	Avatar      string `json:"avatar" binding:"omitempty,url,max=255"`
	Screenshot  string `json:"screenshot" binding:"omitempty,url,max=255"`
	Sort        int    `json:"sort" binding:"omitempty,min=1,max=10"`
	IsInvalid   *bool  `json:"is_invalid"` // 是否失效
	IsPending   *bool  `json:"is_pending"` // 是否为待审核申请
	TypeID      *uint  `json:"type_id" binding:"omitempty"`
	RSSUrl      string `json:"rss_url" binding:"omitempty,url,max=500"` // RSS订阅地址
	Accessible  *int   `json:"accessible"`                              // 可访问性状态
}

// ============ 前台友链申请 ============

// ApplyFriendRequest 申请友链请求
type ApplyFriendRequest struct {
	Name        string `json:"name" binding:"required,max=50"`             // 网站名称
	URL         string `json:"url" binding:"required,url,max=255"`         // 网站链接
	Description string `json:"description" binding:"required,max=500"`     // 网站描述
	Avatar      string `json:"avatar" binding:"required,url,max=255"`      // 网站头像/logo
	Screenshot  string `json:"screenshot" binding:"omitempty,url,max=255"` // 网站截图（可选）
}

// ============ 后台友链管理响应 ============

// FriendListResponse 后台友链列表响应
type FriendListResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	URL         string          `json:"url"`
	Description string          `json:"description"`
	Avatar      string          `json:"avatar"`
	Screenshot  string          `json:"screenshot"`
	Sort        int             `json:"sort"`
	IsInvalid   bool            `json:"is_invalid"`
	IsPending   bool            `json:"is_pending"`
	TypeID      *uint           `json:"type_id"`
	TypeName    string          `json:"type_name,omitempty"`
	RSSUrl      string          `json:"rss_url"`
	RSSLatime   *utils.JSONTime `json:"rss_latime,omitempty"`
	Accessible  int             `json:"accessible"`
}
