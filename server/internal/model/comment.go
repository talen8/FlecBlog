package model

import "gorm.io/gorm"

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content    string `gorm:"type:text;not null" json:"content"`
	TargetType string `gorm:"type:varchar(20);not null;index:idx_target" json:"target_type"` // article/page/moment
	TargetKey  string `gorm:"type:varchar(50);not null;index:idx_target" json:"target_key"`  // 文章slug或页面key
	UserID     uint   `gorm:"not null" json:"user_id"`
	ParentID   *uint  `json:"parent_id"`                       // 直接父评论ID
	RootID     *uint  `gorm:"column:root_id" json:"root_id"`   // 根评论ID（用于扁平化）
	ReplyTo    *uint  `gorm:"column:reply_to" json:"reply_to"` // 回复的目标用户ID
	Status     int    `gorm:"default:1" json:"status"`         // 0:隐藏 1:显示
	IsPinned   bool   `gorm:"default:false" json:"is_pinned"`  // 是否置顶

	// 用户环境信息
	IP       string `gorm:"type:varchar(45)" json:"ip"`        // IP地址（支持IPv6）
	Location string `gorm:"type:varchar(100)" json:"location"` // 地理位置
	Browser  string `gorm:"type:varchar(50)" json:"browser"`   // 浏览器内核
	OS       string `gorm:"type:varchar(50)" json:"os"`        // 操作系统

	// 关联关系
	User      User  `gorm:"foreignKey:UserID" json:"user"`
	ReplyUser *User `gorm:"foreignKey:ReplyTo" json:"reply_user"` // 回复的目标用户
}
