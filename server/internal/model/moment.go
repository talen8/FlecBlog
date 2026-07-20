package model

import (
	"time"
)

// Moment 动态模型
type Moment struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Content     string     `gorm:"type:json;not null" json:"content"` // 内容（JSON）- 包含text、images、location、link、music、video等
	IsPublish   bool       `gorm:"default:true" json:"is_publish"`    // 是否发布
	PublishTime *time.Time `json:"publish_time"`                      // 发布时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
