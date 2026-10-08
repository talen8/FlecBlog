package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"

	"gorm.io/gorm"
)

// ThemeService 主题配置服务
type ThemeService struct {
	themeRepo   *repository.ThemeRepository
	fileService *FileService
}

// NewThemeService 创建主题配置服务
func NewThemeService(themeRepo *repository.ThemeRepository, fileService *FileService) *ThemeService {
	return &ThemeService{
		themeRepo:   themeRepo,
		fileService: fileService,
	}
}

// GetThemeConfig 获取后台主题配置
func (s *ThemeService) GetThemeConfig() (*dto.ThemeConfigResponse, error) {
	config, err := s.loadConfig()
	if err != nil {
		return nil, err
	}

	menus, err := s.loadMenus()
	if err != nil {
		return nil, err
	}

	pages, err := s.loadPages()
	if err != nil {
		return nil, err
	}

	return buildThemeConfigResponse(config, menus, pages)
}

// GetThemeConfigForWeb 获取前台主题配置（菜单只返回启用项）
func (s *ThemeService) GetThemeConfigForWeb() (*dto.ThemeConfigResponse, error) {
	config, err := s.loadConfig()
	if err != nil {
		return nil, err
	}

	menus, err := s.loadMenus()
	if err != nil {
		return nil, err
	}

	pages, err := s.loadPages()
	if err != nil {
		return nil, err
	}

	return buildThemeConfigResponse(config, filterEnabledMenuGroups(menus), pages)
}

// UpdateConfig 更新主题配置
func (s *ThemeService) UpdateConfig(req *dto.ConfigUpdateRequest) (json.RawMessage, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(req.Config, &data); err != nil {
		return nil, fmt.Errorf("config 不是合法 JSON: %w", err)
	}

	items := make([]model.ThemeConfig, 0, len(data))
	for key, value := range data {
		if isReservedConfigKey(key) {
			continue
		}
		encoded, err := encodeConfigValue(value)
		if err != nil {
			return nil, err
		}
		items = append(items, model.ThemeConfig{Key: key, Value: encoded})
	}

	oldConfig, err := s.loadConfig()
	if err != nil {
		return nil, err
	}
	if err := s.updateConfigImageUsage(oldConfig, data); err != nil {
		return nil, err
	}

	if err := s.themeRepo.UpsertMany(items); err != nil {
		return nil, err
	}

	return req.Config, nil
}

// UpdateMenus 更新主题菜单
func (s *ThemeService) UpdateMenus(req *dto.MenuUpdateRequest) (map[string][]dto.MenuDataItem, error) {
	oldMenus, err := s.loadMenus()
	if err != nil {
		return nil, err
	}

	nextMenus, err := normalizeMenuGroups(req.Menus, collectMenuIDs(oldMenus), nextMenuID(oldMenus))
	if err != nil {
		return nil, err
	}
	if err := validateUniqueMenuIDs(nextMenus); err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(nextMenus)
	if err != nil {
		return nil, err
	}
	if err := s.themeRepo.Upsert(model.ThemeConfigKeyMenus, string(encoded)); err != nil {
		return nil, err
	}
	if err := s.updateMenuIconUsage(oldMenus, nextMenus); err != nil {
		return nil, err
	}

	return nextMenus, nil
}

// UpdatePages 更新主题页面
func (s *ThemeService) UpdatePages(req *dto.PageUpdateRequest) (map[string]dto.PageDataItem, error) {
	nextPages := make(map[string]dto.PageDataItem, len(req.Pages))
	for path, item := range req.Pages {
		item.Title = strings.TrimSpace(item.Title)
		item.Description = strings.TrimSpace(item.Description)
		if item.Title == "" && item.Description == "" {
			continue
		}
		nextPages[path] = item
	}

	encoded, err := json.Marshal(nextPages)
	if err != nil {
		return nil, err
	}
	if err := s.themeRepo.Upsert(model.ThemeConfigKeyPages, string(encoded)); err != nil {
		return nil, err
	}

	return nextPages, nil
}

// loadPages 读取页面文案配置
func (s *ThemeService) loadPages() (map[string]dto.PageDataItem, error) {
	item, err := s.themeRepo.Get(model.ThemeConfigKeyPages)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]dto.PageDataItem{}, nil
		}
		return nil, err
	}
	if item.Value == "" {
		return map[string]dto.PageDataItem{}, nil
	}

	var pages map[string]dto.PageDataItem
	if err := json.Unmarshal([]byte(item.Value), &pages); err != nil {
		return nil, fmt.Errorf("pages 不是合法 JSON: %w", err)
	}
	if pages == nil {
		pages = map[string]dto.PageDataItem{}
	}
	return pages, nil
}

// loadMenus 读取菜单配置
func (s *ThemeService) loadMenus() (map[string][]dto.MenuDataItem, error) {
	item, err := s.themeRepo.Get(model.ThemeConfigKeyMenus)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string][]dto.MenuDataItem{}, nil
		}
		return nil, err
	}
	return parseMenus(item.Value)
}

// updateConfigImageUsage 根据配置变更同步文件使用状态
func (s *ThemeService) updateConfigImageUsage(oldConfig map[string]interface{}, patch map[string]interface{}) error {
	nextConfig := make(map[string]interface{}, len(oldConfig)+len(patch))
	for key, value := range oldConfig {
		nextConfig[key] = value
	}
	for key, value := range patch {
		nextConfig[key] = value
	}

	oldRefs := collectStrings(oldConfig)
	nextRefs := collectStrings(nextConfig)

	changed := make([]string, 0, len(oldRefs)+len(nextRefs))
	for url := range oldRefs {
		if !nextRefs[url] {
			changed = append(changed, url)
		}
	}
	for url := range nextRefs {
		if !oldRefs[url] {
			changed = append(changed, url)
		}
	}

	for _, url := range changed {
		if nextRefs[url] {
			if err := s.fileService.MarkAsUsed(url); err != nil {
				return err
			}
		} else {
			_ = s.fileService.MarkAsUnused(url)
		}
	}
	return nil
}

// updateMenuIconUsage 根据菜单变更同步菜单图标文件使用状态
func (s *ThemeService) updateMenuIconUsage(oldMenus map[string][]dto.MenuDataItem, nextMenus map[string][]dto.MenuDataItem) error {
	oldIcons := collectMenusIcons(oldMenus)
	nextIcons := collectMenusIcons(nextMenus)
	for icon := range oldIcons {
		if !nextIcons[icon] {
			_ = s.fileService.MarkAsUnused(icon)
		}
	}
	for icon := range nextIcons {
		if !oldIcons[icon] {
			if err := s.fileService.MarkAsUsed(icon); err != nil {
				return err
			}
		}
	}
	return nil
}

// loadConfig 读取全部配置项为对象，菜单键不参与
func (s *ThemeService) loadConfig() (map[string]interface{}, error) {
	items, err := s.themeRepo.List()
	if err != nil {
		return nil, err
	}

	config := make(map[string]interface{}, len(items))
	for i := range items {
		if isReservedConfigKey(items[i].Key) {
			continue
		}
		config[items[i].Key] = decodeConfigValue(items[i].Value)
	}
	return config, nil
}

// isReservedConfigKey 判断配置键是否由专项接口维护，普通配置读写需跳过
func isReservedConfigKey(key string) bool {
	return key == model.ThemeConfigKeyMenus || key == model.ThemeConfigKeyPages
}

// buildThemeConfigResponse 组装配置、菜单与页面文案响应
func buildThemeConfigResponse(
	config map[string]interface{},
	menus map[string][]dto.MenuDataItem,
	pages map[string]dto.PageDataItem,
) (*dto.ThemeConfigResponse, error) {
	configRaw, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	menusRaw, err := json.Marshal(menus)
	if err != nil {
		return nil, err
	}
	pagesRaw, err := json.Marshal(pages)
	if err != nil {
		return nil, err
	}
	return &dto.ThemeConfigResponse{Config: configRaw, Menus: menusRaw, Pages: pagesRaw}, nil
}

// decodeConfigValue 将数据库文本还原为 JSON 值：可解析为 JSON 的按类型还原，否则按纯字符串处理
func decodeConfigValue(raw string) interface{} {
	if raw == "" {
		return ""
	}

	var value interface{}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return raw
	}
	return value
}

// encodeConfigValue 将 JSON 值序列化为数据库文本：字符串直接存原文，其余存 JSON 文本
func encodeConfigValue(value interface{}) (string, error) {
	if text, ok := value.(string); ok {
		return text, nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("配置值序列化失败: %w", err)
	}
	return string(data), nil
}

// collectStrings 深度遍历配置值，收集所有非空字符串叶子
func collectStrings(value interface{}) map[string]bool {
	out := make(map[string]bool)

	var walk func(v interface{})
	walk = func(v interface{}) {
		switch t := v.(type) {
		case string:
			if t != "" {
				out[t] = true
			}
		case map[string]interface{}:
			for _, item := range t {
				walk(item)
			}
		case []interface{}:
			for _, item := range t {
				walk(item)
			}
		}
	}
	walk(value)

	return out
}

// parseMenus 将菜单 JSON 字符串解析为按类型分组的菜单树
func parseMenus(raw string) (map[string][]dto.MenuDataItem, error) {
	if raw == "" {
		return map[string][]dto.MenuDataItem{}, nil
	}

	var menus map[string][]dto.MenuDataItem
	if err := json.Unmarshal([]byte(raw), &menus); err != nil {
		return nil, fmt.Errorf("menus 不是合法 JSON: %w", err)
	}
	if menus == nil {
		menus = map[string][]dto.MenuDataItem{}
	}
	return menus, nil
}

// nextMenuID 返回当前菜单集合中的下一个可分配 ID
func nextMenuID(menus map[string][]dto.MenuDataItem) int {
	maxID := 0
	for _, items := range menus {
		walkMenuItems(items, func(item dto.MenuDataItem) {
			if item.ID > maxID {
				maxID = item.ID
			}
		})
	}
	return maxID + 1
}

// collectMenuIDs 收集当前菜单集合中已经存在的正数 ID
func collectMenuIDs(menus map[string][]dto.MenuDataItem) map[int]bool {
	ids := make(map[int]bool)
	for _, items := range menus {
		walkMenuItems(items, func(item dto.MenuDataItem) {
			if item.ID > 0 {
				ids[item.ID] = true
			}
		})
	}
	return ids
}

// normalizeMenuGroups 规范化提交的菜单集合并校验非新增 ID 来源
func normalizeMenuGroups(menus map[string][]dto.MenuDataItem, existingIDs map[int]bool, nextID int) (map[string][]dto.MenuDataItem, error) {
	if menus == nil {
		return map[string][]dto.MenuDataItem{}, nil
	}

	result := make(map[string][]dto.MenuDataItem, len(menus))
	for menuType, items := range menus {
		nextItems, err := normalizeMenuItemsWithCounter(items, existingIDs, &nextID)
		if err != nil {
			return nil, err
		}
		result[menuType] = nextItems
	}
	return result, nil
}

// normalizeMenuItemsWithCounter 递归规范化菜单项并为新增项分配 ID
func normalizeMenuItemsWithCounter(items []dto.MenuDataItem, existingIDs map[int]bool, nextID *int) ([]dto.MenuDataItem, error) {
	if items == nil {
		return []dto.MenuDataItem{}, nil
	}

	result := make([]dto.MenuDataItem, 0, len(items))
	for _, item := range items {
		if item.ID <= 0 {
			item.ID = *nextID
			(*nextID)++
		} else if !existingIDs[item.ID] {
			return nil, fmt.Errorf("菜单 ID 非法: %d", item.ID)
		}

		children, err := normalizeMenuItemsWithCounter(item.Children, existingIDs, nextID)
		if err != nil {
			return nil, err
		}
		item.Children = children
		result = append(result, item)
	}
	return result, nil
}

// validateUniqueMenuIDs 校验所有菜单项 ID 全局唯一
func validateUniqueMenuIDs(menus map[string][]dto.MenuDataItem) error {
	seen := make(map[int]bool)
	for _, items := range menus {
		var duplicateID int
		walkMenuItems(items, func(item dto.MenuDataItem) {
			if duplicateID != 0 || item.ID == 0 {
				return
			}
			if seen[item.ID] {
				duplicateID = item.ID
				return
			}
			seen[item.ID] = true
		})
		if duplicateID != 0 {
			return fmt.Errorf("菜单 ID 重复: %d", duplicateID)
		}
	}
	return nil
}

// walkMenuItems 深度遍历菜单树并对每个菜单项执行回调
func walkMenuItems(items []dto.MenuDataItem, visit func(dto.MenuDataItem)) {
	for _, item := range items {
		visit(item)
		walkMenuItems(item.Children, visit)
	}
}

// filterEnabledMenuGroups 过滤菜单分组，只保留包含启用项的分组
func filterEnabledMenuGroups(menus map[string][]dto.MenuDataItem) map[string][]dto.MenuDataItem {
	filtered := make(map[string][]dto.MenuDataItem)
	for menuType, items := range menus {
		nextItems := filterEnabledMenuItems(items)
		if len(nextItems) > 0 {
			filtered[menuType] = nextItems
		}
	}
	return filtered
}

// filterEnabledMenuItems 递归过滤菜单树，只保留启用菜单项
func filterEnabledMenuItems(items []dto.MenuDataItem) []dto.MenuDataItem {
	result := make([]dto.MenuDataItem, 0, len(items))
	for _, item := range items {
		if !item.IsEnabled {
			continue
		}
		item.Children = filterEnabledMenuItems(item.Children)
		result = append(result, item)
	}
	return result
}

// collectMenusIcons 收集菜单集合中所有非空图标地址
func collectMenusIcons(menus map[string][]dto.MenuDataItem) map[string]bool {
	icons := make(map[string]bool)
	for _, items := range menus {
		walkMenuItems(items, func(item dto.MenuDataItem) {
			if item.Icon != "" {
				icons[item.Icon] = true
			}
		})
	}
	return icons
}
