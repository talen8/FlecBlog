package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"

	"gorm.io/gorm"
)

// ThemeService 主题服务
type ThemeService struct {
	db          *gorm.DB
	themeRepo   *repository.ThemeRepository
	fileService *FileService
}

// NewThemeService 创建主题服务
func NewThemeService(db *gorm.DB, themeRepo *repository.ThemeRepository, fileService *FileService) *ThemeService {
	return &ThemeService{
		db:          db,
		themeRepo:   themeRepo,
		fileService: fileService,
	}
}

// SyncThemeMeta 同步主题镜像元数据，并激活当前主题
func (s *ThemeService) SyncThemeMeta(req *dto.ThemeMetaSyncRequest) error {
	if !json.Valid(req.Schema) {
		return errors.New("schema 不是合法 JSON")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := repository.NewThemeRepository(tx)

		oldTheme, _ := repo.Get(req.Slug)

		if err := repo.DeactivateAll(); err != nil {
			return err
		}

		config := `{}`
		menus := `{}`

		if oldTheme != nil {
			config = migrateConfigBySchema(oldTheme.Config, req.Schema)
			menus = oldTheme.Menus
		}

		return repo.SyncTheme(&model.ThemeInstance{
			Slug:        req.Slug,
			Name:        req.Name,
			Version:     req.Version,
			Author:      req.Author,
			Description: req.Description,
			License:     req.License,
			Repo:        req.Repo,
			Schema:      string(req.Schema),
			IsActive:    true,
			Config:      config,
			Menus:       menus,
		})
	})
}

// GetActiveTheme 获取前台激活主题信息
func (s *ThemeService) GetActiveTheme() (*dto.ThemePublicResponse, error) {
	theme, err := s.themeRepo.GetActive()
	if err != nil {
		return nil, err
	}
	return s.toThemePublicResponse(theme, true)
}

// ListThemes 获取主题列表
func (s *ThemeService) ListThemes() ([]dto.ThemeResponse, error) {
	themes, err := s.themeRepo.List()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ThemeResponse, 0, len(themes))
	for i := range themes {
		resp, err := s.toThemeResponse(&themes[i], false)
		if err != nil {
			return nil, err
		}
		result = append(result, *resp)
	}
	return result, nil
}

// GetTheme 获取后台主题详情
func (s *ThemeService) GetTheme(slug string) (*dto.ThemeResponse, error) {
	theme, err := s.themeRepo.Get(slug)
	if err != nil {
		return nil, err
	}
	return s.toThemeResponse(theme, false)
}

// UpdateConfig 更新主题配置
func (s *ThemeService) UpdateConfig(slug string, req *dto.ConfigUpdateRequest) (json.RawMessage, error) {
	nextConfig := rawJSONOrDefault(req.Config, `{}`)
	if !json.Valid(nextConfig) {
		return nil, errors.New("config 不是合法 JSON")
	}

	theme, err := s.themeRepo.Get(slug)
	if err != nil {
		return nil, err
	}

	if err := s.updateConfigImageUsage(theme, nextConfig); err != nil {
		return nil, err
	}
	if err := s.themeRepo.UpdateConfig(slug, string(nextConfig)); err != nil {
		return nil, err
	}

	return nextConfig, nil
}

// UpdateMenus 整体替换主题菜单
func (s *ThemeService) UpdateMenus(slug string, req *dto.MenuUpdateRequest) (map[string][]dto.MenuDataItem, error) {
	theme, err := s.themeRepo.Get(slug)
	if err != nil {
		return nil, err
	}

	oldMenus, err := parseMenus(theme.Menus)
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
	if err := s.themeRepo.UpdateMenus(slug, string(encoded)); err != nil {
		return nil, err
	}
	if err := s.updateMenuIconUsage(slug, oldMenus, nextMenus); err != nil {
		return nil, err
	}

	return nextMenus, nil
}

// toThemeResponse 将主题模型转换为接口响应，并按需过滤未启用菜单
func (s *ThemeService) toThemeResponse(theme *model.ThemeInstance, filterEnabledMenus bool) (*dto.ThemeResponse, error) {
	menusRaw := rawJSONOrDefault([]byte(theme.Menus), `{}`)
	if filterEnabledMenus {
		menus, err := parseMenus(theme.Menus)
		if err != nil {
			return nil, err
		}
		filtered := filterEnabledMenuGroups(menus)
		encoded, err := json.Marshal(filtered)
		if err != nil {
			return nil, err
		}
		menusRaw = encoded
	}

	return &dto.ThemeResponse{
		Slug:        theme.Slug,
		Name:        theme.Name,
		Version:     theme.Version,
		Author:      theme.Author,
		Description: theme.Description,
		License:     theme.License,
		Repo:        theme.Repo,
		Schema:      rawJSONOrDefault([]byte(theme.Schema), `{}`),
		IsActive:    theme.IsActive,
		Config:      rawJSONOrDefault([]byte(theme.Config), `{}`),
		Menus:       menusRaw,
	}, nil
}

// toThemePublicResponse 转换为前台主题响应（不含 schema）
func (s *ThemeService) toThemePublicResponse(theme *model.ThemeInstance, filterEnabledMenus bool) (*dto.ThemePublicResponse, error) {
	menusRaw := rawJSONOrDefault([]byte(theme.Menus), `{}`)
	if filterEnabledMenus {
		menus, err := parseMenus(theme.Menus)
		if err != nil {
			return nil, err
		}
		filtered := filterEnabledMenuGroups(menus)
		encoded, err := json.Marshal(filtered)
		if err != nil {
			return nil, err
		}
		menusRaw = encoded
	}

	return &dto.ThemePublicResponse{
		Slug:        theme.Slug,
		Name:        theme.Name,
		Version:     theme.Version,
		Author:      theme.Author,
		Description: theme.Description,
		License:     theme.License,
		Repo:        theme.Repo,
		IsActive:    theme.IsActive,
		Config:      rawJSONOrDefault([]byte(theme.Config), `{}`),
		Menus:       menusRaw,
	}, nil
}

// updateConfigImageUsage 根据配置变更同步主题图片文件使用状态
func (s *ThemeService) updateConfigImageUsage(theme *model.ThemeInstance, nextConfig json.RawMessage) error {
	if s.fileService == nil {
		return nil
	}

	oldURLs := collectConfigImageURLs([]byte(theme.Schema), []byte(theme.Config))
	nextURLs := collectConfigImageURLs([]byte(theme.Schema), nextConfig)
	for url := range oldURLs {
		if !nextURLs[url] {
			_ = s.fileService.MarkAsUnused(url)
		}
	}
	for url := range nextURLs {
		if !oldURLs[url] {
			if err := s.fileService.MarkAsUsedWithType(url, theme.Slug); err != nil {
				return err
			}
		}
	}
	return nil
}

// updateMenuIconUsage 根据菜单变更同步菜单图标文件使用状态
func (s *ThemeService) updateMenuIconUsage(slug string, oldMenus map[string][]dto.MenuDataItem, nextMenus map[string][]dto.MenuDataItem) error {
	if s.fileService == nil {
		return nil
	}

	oldIcons := collectMenusIcons(oldMenus)
	nextIcons := collectMenusIcons(nextMenus)
	for icon := range oldIcons {
		if !nextIcons[icon] {
			_ = s.fileService.MarkAsUnused(icon)
		}
	}
	for icon := range nextIcons {
		if !oldIcons[icon] {
			if err := s.fileService.MarkAsUsedWithType(icon, slug); err != nil {
				return err
			}
		}
	}
	return nil
}

// rawJSONOrDefault 在 JSON 为空或 null 时返回默认 JSON
func rawJSONOrDefault(raw []byte, fallback string) json.RawMessage {
	if len(raw) == 0 || string(raw) == "null" {
		return json.RawMessage(fallback)
	}
	return json.RawMessage(raw)
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
		if item.ID == 0 {
			item.ID = *nextID
			(*nextID)++
		} else if item.ID < 0 || !existingIDs[item.ID] {
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

// validateUniqueMenuIDs 校验所有菜单项 ID 在主题内全局唯一
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
		if item.Children == nil {
			item.Children = []dto.MenuDataItem{}
		}
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

// CheckThemeUpdate 检查主题版本更新
func (s *ThemeService) CheckThemeUpdate(ctx context.Context, slug string) (*dto.ThemeUpdateCheckResponse, error) {
	theme, err := s.themeRepo.Get(slug)
	if err != nil {
		return nil, fmt.Errorf("获取主题失败: %w", err)
	}

	if theme.Repo == "" {
		return nil, errors.New("主题未设置仓库地址")
	}

	owner, repoName, ok := parseRepoURL(theme.Repo)
	if !ok {
		return nil, errors.New("仅支持 GitHub 仓库地址")
	}

	resp := &dto.ThemeUpdateCheckResponse{
		CurrentVersion: strings.TrimPrefix(theme.Version, "v"),
	}

	latestVersion, releaseURL, err := fetchLatestRelease(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	resp.LatestVersion = strings.TrimPrefix(latestVersion, "v")
	resp.ReleaseURL = releaseURL

	if resp.CurrentVersion == "" {
		return resp, nil
	}

	cmp, err := compareVersion(latestVersion, resp.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("比较版本失败: %w", err)
	}

	resp.HasUpdate = cmp > 0
	return resp, nil
}

// parseRepoURL 解析 GitHub 仓库地址，返回 owner 和 repo
func parseRepoURL(repo string) (string, string, bool) {
	repo = strings.TrimSuffix(repo, ".git")
	repo = strings.TrimSuffix(repo, "/")

	if !strings.HasPrefix(repo, "https://github.com/") {
		return "", "", false
	}

	parts := strings.SplitN(strings.TrimPrefix(repo, "https://github.com/"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return parts[0], parts[1], true
}

// fetchLatestRelease 从 GitHub API 获取最新 release
func fetchLatestRelease(ctx context.Context, owner, repo string) (tagName, htmlURL string, err error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "FlecBlog-ThemeUpdateChecker")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("请求 GitHub API 失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return "", "", errors.New("仓库不存在或没有 release")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", "", fmt.Errorf("GitHub API 返回错误: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}

	if release.TagName == "" {
		return "", "", errors.New("release 缺少 tag_name")
	}

	return release.TagName, release.HTMLURL, nil
}

// collectConfigImageURLs 按主题 schema 从配置中收集图片地址
func collectConfigImageURLs(schemaRaw []byte, configRaw []byte) map[string]bool {
	urls := make(map[string]bool)

	var configData interface{}
	if err := json.Unmarshal(rawJSONOrDefault(configRaw, `{}`), &configData); err != nil {
		return urls
	}

	var schemaData interface{}
	if err := json.Unmarshal(rawJSONOrDefault(schemaRaw, `{}`), &schemaData); err != nil {
		return urls
	}
	collectSchemaImageURLs(schemaData, configData, urls)
	return urls
}

// collectSchemaImageURLs 递归按 schema 字段定义收集配置中的图片地址
func collectSchemaImageURLs(schemaData interface{}, configData interface{}, urls map[string]bool) {
	schemaObj, ok := schemaData.(map[string]interface{})
	if !ok {
		return
	}

	if isImageField(schemaObj) {
		if value, ok := configData.(string); ok && value != "" {
			urls[value] = true
		}
	}

	if properties, ok := schemaObj["properties"].(map[string]interface{}); ok {
		configObj, _ := configData.(map[string]interface{})
		for key, childSchema := range properties {
			collectSchemaImageURLs(childSchema, configObj[key], urls)
		}
	}

	if itemFields, ok := schemaObj["x-item-fields"].([]interface{}); ok {
		if configItems, ok := configData.([]interface{}); ok {
			for _, item := range configItems {
				itemObj, _ := item.(map[string]interface{})
				for _, fieldSchema := range itemFields {
					fieldObj, ok := fieldSchema.(map[string]interface{})
					if !ok {
						continue
					}
					key, _ := fieldObj["key"].(string)
					collectSchemaImageURLs(fieldSchema, itemObj[key], urls)
				}
			}
		}
	}

	if itemsSchema, ok := schemaObj["items"]; ok {
		if configItems, ok := configData.([]interface{}); ok {
			for _, item := range configItems {
				collectSchemaImageURLs(itemsSchema, item, urls)
			}
		}
	}

	configObj, _ := configData.(map[string]interface{})
	for key, childSchema := range schemaObj {
		if strings.HasPrefix(key, "$") || isSchemaKeyword(key) {
			continue
		}
		childObj, ok := childSchema.(map[string]interface{})
		if !ok {
			continue
		}
		if isThemeFieldSchema(childObj) {
			collectSchemaImageURLs(childSchema, configObj[key], urls)
		} else {
			collectSchemaImageURLs(childSchema, configData, urls)
		}
	}
}

// isImageField 判断 schema 字段是否表示图片或上传类型
func isImageField(schemaObj map[string]interface{}) bool {
	format, _ := schemaObj["format"].(string)
	fieldType, _ := schemaObj["type"].(string)
	return format == "image" || format == "upload" || fieldType == "upload"
}

// isThemeFieldSchema 判断对象是否像一个主题字段 schema
func isThemeFieldSchema(schemaObj map[string]interface{}) bool {
	_, hasType := schemaObj["type"]
	_, hasFormat := schemaObj["format"]
	_, hasItems := schemaObj["items"]
	_, hasProperties := schemaObj["properties"]
	_, hasItemFields := schemaObj["x-item-fields"]
	_, hasEnum := schemaObj["enum"]
	_, hasOptions := schemaObj["options"]
	return hasType || hasFormat || hasItems || hasProperties || hasItemFields || hasEnum || hasOptions
}

// isSchemaKeyword 判断 key 是否为 schema 元信息关键字
func isSchemaKeyword(key string) bool {
	switch key {
	case "type", "title", "label", "description", "default", "enum", "enumNames",
		"options", "format", "placeholder", "widget", "group", "properties", "items",
		"minimum", "maximum", "min", "max", "rows", "uploadType", "x-group",
		"x-component", "x-upload-type", "x-item-fields", "width", "height":
		return true
	default:
		return false
	}
}

// migrateConfigBySchema 按新 schema 中声明的 $from 迁移旧配置
func migrateConfigBySchema(oldConfigStr string, schemaRaw json.RawMessage) string {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(oldConfigStr), &config); err != nil {
		return oldConfigStr
	}

	var schema interface{}
	if err := json.Unmarshal(schemaRaw, &schema); err != nil {
		return oldConfigStr
	}

	collectAliases(schema, func(from, to string) {
		if _, exists := config[to]; exists {
			return
		}
		if val, ok := config[from]; ok {
			config[to] = val
			delete(config, from)
		}
	})

	coerceConfigTypes(config, schema)

	data, _ := json.Marshal(config)
	return string(data)
}

// collectAliases 递归遍历 schema，对每个含 $from 的字段回调 (旧名, 新名)
func collectAliases(schema interface{}, fn func(from, to string)) {
	obj, ok := schema.(map[string]interface{})
	if !ok {
		return
	}

	for key, val := range obj {
		child, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		if from, ok := child["$from"]; ok {
			emitAlias(from, key, fn)
			continue
		}

		collectAliases(val, fn)
	}
}

// coerceConfigTypes 按 schema 声明的 type 隐式转换 config 值类型
func coerceConfigTypes(config map[string]interface{}, schema interface{}) {
	obj, ok := schema.(map[string]interface{})
	if !ok {
		return
	}
	for key, val := range obj {
		child, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		if targetType, ok := child["type"]; ok {
			if v, exists := config[key]; exists {
				config[key] = coerceValue(v, fmt.Sprintf("%v", targetType))
			}
			continue
		}
		coerceConfigTypes(config, val)
	}
}

func coerceValue(v interface{}, targetType string) interface{} {
	switch targetType {
	case "number", "integer":
		switch val := v.(type) {
		case string:
			if n, err := strconv.ParseFloat(val, 64); err == nil {
				if targetType == "integer" {
					return int(n)
				}
				return n
			}
		case bool:
			if val {
				return 1
			}
			return 0
		}
	case "string":
		switch val := v.(type) {
		case float64:
			if val == float64(int64(val)) {
				return fmt.Sprintf("%.0f", val)
			}
			return fmt.Sprintf("%v", val)
		case bool:
			return fmt.Sprintf("%v", val)
		case int:
			return fmt.Sprintf("%d", val)
		}
	case "boolean":
		switch val := v.(type) {
		case string:
			switch val {
			case "true", "1":
				return true
			case "false", "0":
				return false
			}
		case float64:
			return val != 0
		case int:
			return val != 0
		}
	}
	return v
}

func emitAlias(raw interface{}, to string, fn func(from, to string)) {
	switch v := raw.(type) {
	case string:
		fn(v, to)
	case []interface{}:
		for _, alias := range v {
			if s, ok := alias.(string); ok {
				fn(s, to)
			}
		}
	}
}
