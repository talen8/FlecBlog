package dto

import "flec_blog/pkg/utils"

// ListRssArticleRequest RSS文章列表请求
type ListRssArticleRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword   string `form:"keyword"`    // 搜索关键词
	FriendID  uint   `form:"friend_id"`  // 友链ID筛选
	IsRead    *bool  `form:"is_read"`    // 已读状态筛选
	StartTime string `form:"start_time"` // 发布开始时间（格式：2006-01-02）
	EndTime   string `form:"end_time"`   // 发布结束时间（格式：2006-01-02）
}

// RssArticleResponse RSS文章响应
type RssArticleResponse struct {
	ID          uint            `json:"id"`
	FriendID    uint            `json:"friend_id"`
	FriendName  string          `json:"friend_name"`
	FriendURL   string          `json:"friend_url"`
	Title       string          `json:"title"`
	Link        string          `json:"link"`
	PublishedAt *utils.JSONTime `json:"published_at,omitempty"`
	IsRead      bool            `json:"is_read"`
	CreatedAt   *utils.JSONTime `json:"created_at"`
}

// RssArticleListResponse RSS文章列表响应
type RssArticleListResponse struct {
	List        []RssArticleResponse `json:"list"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	UnreadCount int64                `json:"unread_count"`
}

// MarkArticleReadRequest 标记文章已读请求
type MarkArticleReadRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// MarkAllReadRequest 全部标记已读请求
type MarkAllReadRequest struct{}
