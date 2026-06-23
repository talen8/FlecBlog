package model

import "time"

// ThemeInstance 主题实例
type ThemeInstance struct {
	Slug        string    `gorm:"primaryKey;type:text" json:"slug"`        // 主题唯一标识，与 theme.config.json 的 slug 一致
	Name        string    `gorm:"type:text;default:''" json:"name"`        // 【元数据】显示名
	Version     string    `gorm:"type:text;default:''" json:"version"`     // 【元数据】语义化版本
	Author      string    `gorm:"type:text;default:''" json:"author"`      // 【元数据】作者
	Description string    `gorm:"type:text;default:''" json:"description"` // 【元数据】描述
	License     string    `gorm:"type:text;default:''" json:"license"`     // 【元数据】许可证
	Repo        string    `gorm:"type:text;default:''" json:"repo"`        // 【元数据】仓库地址
	Schema      string    `gorm:"column:schema;type:json" json:"schema"`   // 【元数据】配置表单定义（主题镜像上报）
	IsActive    bool      `gorm:"default:false" json:"is_active"`          // 【状态】是否激活（全表至多一个 TRUE，由 Server 维护）
	Config      string    `gorm:"type:json" json:"config"`                 // 【用户数据】配置取值（站长填写）
	Menus       string    `gorm:"type:json;default:'{}'" json:"menus"`     // 【用户数据】菜单（按 type 分组，见第三节）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
