package dto

import "flec_blog/pkg/utils"

// SubmitFeedbackRequest 提交反馈请求
type SubmitFeedbackRequest struct {
	ReportUrl       string   `json:"reportUrl" binding:"required"`
	ReportType      string   `json:"reportType" binding:"required,oneof=copyright inappropriate summary suggestion"`
	Email           string   `json:"email"` // 可选联系邮箱
	Description     string   `json:"description" binding:"required"`
	Reason          string   `json:"reason"`          // 可选原因字段
	AttachmentFiles []string `json:"attachmentFiles"` // 附件文件URL数组
}

// FeedbackContent 存储到数据库的反馈内容结构
type FeedbackContent struct {
	Description     string   `json:"description"`
	Reason          string   `json:"reason,omitempty"`
	AttachmentFiles []string `json:"attachmentFiles,omitempty"`
}

// FeedbackQueryRequest 反馈查询请求
type FeedbackQueryRequest struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword    string `form:"keyword"`     // 搜索关键词（工单号、投诉地址）
	ReportType string `form:"report_type"` // 反馈类型筛选
	Status     string `form:"status"`      // 状态筛选
	StartTime  string `form:"start_time"`  // 反馈开始时间（格式：2006-01-02）
	EndTime    string `form:"end_time"`    // 反馈结束时间（格式：2006-01-02）
}

// FeedbackResponse 反馈响应
type FeedbackResponse struct {
	ID           uint            `json:"id"`
	TicketNo     string          `json:"ticket_no"` // 工单号
	ReportUrl    string          `json:"report_url"`
	ReportType   string          `json:"report_type"`
	FormContent  interface{}     `json:"form_content"` // 表单内容
	Email        string          `json:"email"`        // 联系邮箱
	Status       string          `json:"status"`
	AdminReply   string          `json:"admin_reply"`
	ReplyTime    *utils.JSONTime `json:"reply_time"` // 回复时间
	UserAgent    string          `json:"user_agent"`
	IP           string          `json:"ip"`
	FeedbackTime utils.JSONTime  `json:"feedback_time"` // 反馈时间
}

// UpdateFeedbackRequest 更新反馈请求
type UpdateFeedbackRequest struct {
	Status     string `json:"status" binding:"required,oneof=pending resolved closed"`
	AdminReply string `json:"admin_reply"`
}
