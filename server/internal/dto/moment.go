package dto

import "flec_blog/pkg/utils"

// ============ 前台动态请求 ============

// ListMomentsForWebRequest 前台动态列表请求
type ListMomentsForWebRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=1000"`
}

// FetchLinkMetadataRequest 获取链接元数据请求
type FetchLinkMetadataRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// LinkMetadataResponse 链接元数据响应
type LinkMetadataResponse struct {
	Title   string `json:"title"`
	Favicon string `json:"favicon"`
}

// ParseVideoRequest 解析视频请求
type ParseVideoRequest struct {
	URL string `json:"url" binding:"required"`
}

// ParseVideoResponse 解析视频响应
type ParseVideoResponse struct {
	Platform string `json:"platform"` // 'bilibili' | 'youtube'
	VideoID  string `json:"video_id"` // 视频ID
}

// ============ 通用动态响应 ============

// MomentContent 动态内容结构
type MomentContent struct {
	Text     string         `json:"text,omitempty"`     // 文本
	Images   []string       `json:"images,omitempty"`   // 图片列表
	Location string         `json:"location,omitempty"` // 位置信息
	Tags     string         `json:"tags,omitempty"`     // 标签
	Link     *MomentLink    `json:"link,omitempty"`     // 外链
	Music    *MomentMusic   `json:"music,omitempty"`    // 音乐（基于MetingJS）
	Video    *MomentVideo   `json:"video,omitempty"`    // 视频（本地或在线）
	Book     map[string]any `json:"book,omitempty"`     // 书籍
	Movie    map[string]any `json:"movie,omitempty"`    // 电影
}

// MomentLink 外链结构
type MomentLink struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Favicon string `json:"favicon"`
}

// MomentMusic 音乐结构（基于 MetingJS）
type MomentMusic struct {
	Server string `json:"server"` // 音乐平台：netease, tencent, kugou, xiami, baidu
	Type   string `json:"type"`   // 类型：song, playlist, album, search, artist
	ID     string `json:"id"`     // 音乐ID
}

// MomentVideo 视频结构
type MomentVideo struct {
	URL      string `json:"url"`                // 视频URL（本地视频或在线视频链接）
	Platform string `json:"platform,omitempty"` // 平台：bilibili, youtube（本地视频为空）
	VideoID  string `json:"video_id,omitempty"` // 视频ID（在线视频的ID，本地视频为空）
}

// MomentForWebResponse 前台动态响应
type MomentForWebResponse struct {
	ID          uint            `json:"id"`
	Content     MomentContent   `json:"content"`
	IsPublish   bool            `json:"is_publish"`
	PublishTime *utils.JSONTime `json:"publish_time"`
}

// ============ 后台动态管理请求 ============

// ListMomentsRequest 后台动态列表请求（支持筛选）
type ListMomentsRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=1000"`
	Keyword   string `form:"keyword"`    // 搜索关键词（文本内容）
	Tags      string `form:"tags"`       // 标签
	Location  string `form:"location"`   // 发布地点
	IsPublish *bool  `form:"is_publish"` // 是否发布
	HasImages *bool  `form:"has_images"` // 是否有图片
	HasVideo  *bool  `form:"has_video"`  // 是否有视频
	HasMusic  *bool  `form:"has_music"`  // 是否有音乐
	HasLink   *bool  `form:"has_link"`   // 是否有链接
	StartTime string `form:"start_time"` // 发布开始时间
	EndTime   string `form:"end_time"`   // 发布结束时间
}

// CreateMomentRequest 创建动态请求
type CreateMomentRequest struct {
	Content     MomentContent   `json:"content" binding:"required"`
	IsPublish   bool            `json:"is_publish"`   // 是否发布，默认true
	PublishTime *utils.JSONTime `json:"publish_time"` // 发布时间（可选，不填则使用创建时间）
}

// UpdateMomentRequest 更新动态请求
type UpdateMomentRequest struct {
	Content     MomentContent   `json:"content" binding:"required"`
	IsPublish   bool            `json:"is_publish"`   // 是否发布
	PublishTime *utils.JSONTime `json:"publish_time"` // 发布时间（可编辑）
}

// ============ 后台动态管理响应 ============

// MomentListResponse 后台动态列表响应
type MomentListResponse struct {
	ID          uint            `json:"id"`
	Content     MomentContent   `json:"content"`
	IsPublish   bool            `json:"is_publish"`
	PublishTime *utils.JSONTime `json:"publish_time"`
}
