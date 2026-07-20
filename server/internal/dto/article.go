package dto

import (
	"flec_blog/pkg/utils"
)

// ============ 通用文章响应 ============

// ArticleDetailResponse 文章详情响应（前台专用）
type ArticleDetailResponse struct {
	ID           uint            `json:"id"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Content      string          `json:"content"`
	Summary      string          `json:"summary"`
	AISummary    string          `json:"ai_summary,omitempty"`
	Cover        string          `json:"cover"`
	Location     string          `json:"location"`
	IsTop        bool            `json:"is_top"`
	IsEssence    bool            `json:"is_essence"`
	IsOutdated   bool            `json:"is_outdated"`
	ViewCount    int             `json:"view_count"`
	CommentCount int64           `json:"comment_count"`
	URL          string          `json:"url"`
	PublishTime  *utils.JSONTime `json:"publish_time"`
	UpdateTime   *utils.JSONTime `json:"update_time"`
	Category     struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"category"`
	Tags []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tags"`
	Prev *struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"prev,omitempty"`
	Next *struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"next,omitempty"`
}

// ============ 前台文章请求 ============

// ListArticlesForWebRequest 前台文章列表请求
type ListArticlesForWebRequest struct {
	Page     int    `form:"page,default=1" binding:"min=0"`
	PageSize int    `form:"page_size,default=10" binding:"min=0"`
	Year     string `form:"year"`
	Month    string `form:"month"`
	Category string `form:"category"`
	Tag      string `form:"tag"`
}

// SearchArticlesRequest 文章搜索请求
type SearchArticlesRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=50"`
}

// ============ 前台文章响应 ============

// ArticleWebResponse 前台文章列表响应
type ArticleWebResponse struct {
	ID           uint            `json:"id"`
	Title        string          `json:"title"`
	Summary      string          `json:"summary"`
	Excerpt      string          `json:"excerpt,omitempty"`
	Cover        string          `json:"cover"`
	Location     string          `json:"location"`
	IsTop        bool            `json:"is_top"`
	IsEssence    bool            `json:"is_essence"`
	IsOutdated   bool            `json:"is_outdated"`
	URL          string          `json:"url"`
	CommentCount int64           `json:"comment_count"`
	PublishTime  *utils.JSONTime `json:"publish_time"`
	UpdateTime   *utils.JSONTime `json:"update_time"`
	Category     struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"category"`
	Tags []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tags"`
}

// ============ 后台文章管理请求 ============

// ListArticlesRequest 后台文章列表请求（支持筛选）
type ListArticlesRequest struct {
	Page       int    `form:"page,default=1" binding:"min=0"`
	PageSize   int    `form:"page_size,default=10" binding:"min=0"`
	Keyword    string `form:"keyword"`     // 搜索关键词（标题/内容）
	CategoryID uint   `form:"category_id"` // 分类ID
	TagIDs     []uint `form:"tag_ids"`     // 标签ID列表
	Location   string `form:"location"`    // 发布地点
	IsPublish  *bool  `form:"is_publish"`  // 是否发布
	IsTop      *bool  `form:"is_top"`      // 是否置顶
	IsEssence  *bool  `form:"is_essence"`  // 是否精选
	IsOutdated *bool  `form:"is_outdated"` // 是否过时
	StartTime  string `form:"start_time"`  // 发布开始时间
	EndTime    string `form:"end_time"`    // 发布结束时间
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug"` // 自定义文章Slug
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	Location   string `json:"location"`    // 发布地点
	IsPublish  *bool  `json:"is_publish"`  // 是否发布
	IsTop      *bool  `json:"is_top"`      // 是否置顶
	IsEssence  *bool  `json:"is_essence"`  // 是否精选
	IsOutdated *bool  `json:"is_outdated"` // 是否过时
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title       string          `json:"title"`
	Slug        string          `json:"slug"` // 自定义文章Slug
	Content     string          `json:"content"`
	Summary     string          `json:"summary"`
	AISummary   string          `json:"ai_summary"` // AI 总结
	Cover       string          `json:"cover"`
	Location    string          `json:"location"`    // 发布地点
	IsPublish   *bool           `json:"is_publish"`  // 是否发布
	IsTop       *bool           `json:"is_top"`      // 是否置顶
	IsEssence   *bool           `json:"is_essence"`  // 是否精选
	IsOutdated  *bool           `json:"is_outdated"` // 是否过时
	CategoryID  *uint           `json:"category_id"`
	TagIDs      []uint          `json:"tag_ids"`
	PublishTime *utils.JSONTime `json:"publish_time"`
	UpdateTime  *utils.JSONTime `json:"update_time"`
}

// ============ 后台文章管理响应 ============

// ArticleListResponse 后台文章列表响应
type ArticleListResponse struct {
	ID           uint            `json:"id"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Cover        string          `json:"cover"`
	Location     string          `json:"location"`
	IsPublish    bool            `json:"is_publish"`
	IsTop        bool            `json:"is_top"`
	IsEssence    bool            `json:"is_essence"`
	IsOutdated   bool            `json:"is_outdated"`
	ViewCount    int             `json:"view_count"`
	CommentCount int64           `json:"comment_count"`
	PublishTime  *utils.JSONTime `json:"publish_time"`
	UpdateTime   *utils.JSONTime `json:"update_time"`
	Category     struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Tags []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

// ArticleAdminDetailResponse 后台文章详情响应
type ArticleAdminDetailResponse struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Content     string          `json:"content"`
	Summary     string          `json:"summary"`
	AISummary   string          `json:"ai_summary"`
	Cover       string          `json:"cover"`
	Location    string          `json:"location"`
	IsPublish   bool            `json:"is_publish"`
	IsTop       bool            `json:"is_top"`
	IsEssence   bool            `json:"is_essence"`
	IsOutdated  bool            `json:"is_outdated"`
	PublishTime *utils.JSONTime `json:"publish_time"`
	UpdateTime  *utils.JSONTime `json:"update_time"`
	Category    struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Tags []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

// ============ 数据导入导出 ============

// ImportArticlesResult 导入文章结果
type ImportArticlesResult struct {
	Total           int                  `json:"total"`
	Success         int                  `json:"success"`
	Failed          int                  `json:"failed"`
	CategoriesAdded int                  `json:"categories_added"` // 新增分类数量
	TagsAdded       int                  `json:"tags_added"`       // 新增标签数量
	Errors          []ImportArticleError `json:"errors,omitempty"`
}

// AddError 添加错误信息
func (r *ImportArticlesResult) AddError(filename, title, errMsg string) {
	r.Failed++
	r.Errors = append(r.Errors, ImportArticleError{
		Filename: filename,
		Title:    title,
		Error:    errMsg,
	})
}

// ImportArticleError 导入错误信息
type ImportArticleError struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
	Error    string `json:"error"`
}

// ============ 微信公众号导出 ============

// WeChatExportResult 微信导出结果
type WeChatExportResult struct {
	HTML string `json:"html"` // 公众号 HTML，用于复制粘贴到微信公众平台
}
