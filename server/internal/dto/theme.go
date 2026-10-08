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

// PageDataItem 主题页面项
type PageDataItem struct {
	Title       string `json:"title" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
}

// ThemeConfigResponse 主题配置响应
type ThemeConfigResponse struct {
	Config json.RawMessage `json:"config" swaggertype:"object"`
	Menus  json.RawMessage `json:"menus" swaggertype:"object"`
	Pages  json.RawMessage `json:"pages" swaggertype:"object"`
}

// ConfigUpdateRequest 主题配置更新请求
type ConfigUpdateRequest struct {
	Config json.RawMessage `json:"config" binding:"required" swaggertype:"object"`
}

// MenuUpdateRequest 主题菜单更新请求
type MenuUpdateRequest struct {
	Menus map[string][]MenuDataItem `json:"menus" binding:"required" swaggertype:"object"`
}

// PageUpdateRequest 主题页面更新请求
type PageUpdateRequest struct {
	Pages map[string]PageDataItem `json:"pages" binding:"required" swaggertype:"object"`
}
