package model

import "time"

// ThemeConfig KV 配置表
type ThemeConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 主题配置常量
const (
	ThemeConfigKeyMenus = "menus"
)
