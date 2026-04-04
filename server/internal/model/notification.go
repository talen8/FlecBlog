package model

import "time"

// NotificationType 通知类型
type NotificationType string

const (
	TypeCommentReply   NotificationType = "comment_reply"   // 评论回复
	TypeCommentNew     NotificationType = "comment_new"     // 评论通知
	TypeFeedbackNew    NotificationType = "feedback_new"    // 问题反馈
	TypeFriendApply    NotificationType = "friend_apply"    // 友链申请
	TypeFriendAbnormal NotificationType = "friend_abnormal" // 异常友链
	TypeSystemAlert    NotificationType = "system_alert"    // 系统通知

	AlertTypeVersionUpdate = "version_update" // 版本更新提醒
)

// Notification 通知模型
type Notification struct {
	ID   uint             `gorm:"primarykey" json:"id"`
	Type NotificationType `gorm:"type:varchar(50);not null;index" json:"type"`

	// 前端显示字段（直接使用）
	Title   string `gorm:"type:varchar(200);not null" json:"title"` // 通知标题
	Content string `gorm:"type:text;not null" json:"content"`       // 通知内容
	Link    string `gorm:"type:varchar(500);not null" json:"link"`  // 跳转链接

	// 详细数据（邮件/Webhook用）
	Data string `gorm:"type:json;not null" json:"data"` // 结构化数据（JSON字符串）

	SenderID *uint `gorm:"index" json:"sender_id,omitempty"` // 发送者ID(可选)
	TargetID *uint `gorm:"index" json:"target_id,omitempty"` // 目标对象ID（文章、评论）

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

// UserNotification 用户通知关联模型
type UserNotification struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	NotificationID uint       `gorm:"not null;index:idx_notification" json:"notification_id"`
	UserID         uint       `gorm:"not null;index:idx_user_read" json:"user_id"`
	IsRead         bool       `gorm:"default:false;index:idx_user_read" json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`

	// 关联
	Notification Notification `gorm:"foreignKey:NotificationID" json:"notification"`
}
