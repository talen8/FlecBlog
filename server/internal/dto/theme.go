package dto

import "encoding/json"

// MenuDataItem 主题菜单项
type MenuDataItem struct {
	ID        int            `json:"id"`
	Title     string         `json:"title" binding:"required,max=100"`
	URL       string         `json:"url" binding:"max=500"`
	Icon      string         `json:"icon" binding:"max=500"`
	Sort      int            `json:"sort"`
	IsEnabled bool           `json:"is_enabled"`
	Children  []MenuDataItem `json:"children"`
}

// ThemeMetaSyncRequest 主题镜像元数据同步请求
type ThemeMetaSyncRequest struct {
	Slug        string          `json:"slug" binding:"required"`
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Author      string          `json:"author"`
	Description string          `json:"description"`
	License     string          `json:"license"`
	Repo        string          `json:"repo"`
	Schema      json.RawMessage `json:"schema" swaggertype:"object"`
}

// ThemeResponse 主题响应
type ThemeResponse struct {
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Author      string          `json:"author"`
	Description string          `json:"description"`
	License     string          `json:"license"`
	Repo        string          `json:"repo"`
	Schema      json.RawMessage `json:"schema" swaggertype:"object"`
	IsActive    bool            `json:"is_active"`
	Config      json.RawMessage `json:"config" swaggertype:"object"`
	Menus       json.RawMessage `json:"menus" swaggertype:"object"`
}

// ThemePublicResponse 前台主题响应
type ThemePublicResponse struct {
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Author      string          `json:"author"`
	Description string          `json:"description"`
	License     string          `json:"license"`
	Repo        string          `json:"repo"`
	IsActive    bool            `json:"is_active"`
	Config      json.RawMessage `json:"config" swaggertype:"object"`
	Menus       json.RawMessage `json:"menus" swaggertype:"object"`
}

// ConfigUpdateRequest 主题配置更新请求
type ConfigUpdateRequest struct {
	Config json.RawMessage `json:"config" binding:"required" swaggertype:"object"`
}

// MenuUpdateRequest 主题菜单更新请求
type MenuUpdateRequest struct {
	Menus map[string][]MenuDataItem `json:"menus" binding:"required" swaggertype:"object"`
}

// ThemeUpdateCheckResponse 主题更新检查响应
type ThemeUpdateCheckResponse struct {
	HasUpdate      bool   `json:"has_update"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	ReleaseURL     string `json:"release_url"`
}
