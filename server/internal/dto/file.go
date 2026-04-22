package dto

import (
	"flec_blog/pkg/upload"
	"flec_blog/pkg/utils"
)

// ============ 通用文件请求 ============

// UploadFileRequest 文件上传请求
type UploadFileRequest struct {
	Type string `form:"type" binding:"required"`
}

// ============ 通用文件响应 ============

// FileUploadForWebResponse 文件上传响应
type FileUploadForWebResponse struct {
	OriginalName string `json:"original_name"`
	FileURL      string `json:"file_url"`
}

// ============ 后台文件管理请求 ============

// ListFilesRequest 文件列表请求
type ListFilesRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Keyword    string `form:"keyword"`
	FileType   string `form:"file_type"`
	Status     *int   `form:"status"`
	UploadType string `form:"upload_type"`
	MinSize    int64  `form:"min_size"`
	MaxSize    int64  `form:"max_size"`
	StartTime  string `form:"start_time"`
	EndTime    string `form:"end_time"`
}

// ============ 后台文件管理响应 ============

// FileResponse 文件信息响应
type FileResponse struct {
	ID           uint           `json:"id"`
	OriginalName string         `json:"original_name"`
	FileName     string         `json:"file_name"`
	FileSize     int64          `json:"file_size"`
	FileType     string         `json:"file_type"`
	FileURL      string         `json:"file_url"`
	UploadType   upload.Type    `json:"upload_type"`
	UserID       *uint          `json:"user_id"`
	Status       int            `json:"status"` // 文件状态：0=未使用 1=使用中
	UploadTime   utils.JSONTime `json:"upload_time"`
}
