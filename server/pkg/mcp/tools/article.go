package tools

import (
	"context"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	articleActionList   = "list"
	articleActionGet    = "get"
	articleActionCreate = "create"
	articleActionUpdate = "update"
	articleActionDelete = "delete"
)

// ============ MCP 类型定义============

// ArticleItem 文章列表项
type ArticleItem struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Cover        string       `json:"cover"`
	Location     string       `json:"location"`
	IsPublish    bool         `json:"is_publish"`
	IsTop        bool         `json:"is_top"`
	IsEssence    bool         `json:"is_essence"`
	IsOutdated   bool         `json:"is_outdated"`
	ViewCount    int          `json:"view_count"`
	CommentCount int64        `json:"comment_count"`
	PublishTime  *string      `json:"publish_time"`
	UpdateTime   *string      `json:"update_time"`
	Category     CategoryItem `json:"category"`
	Tags         []TagItem    `json:"tags"`
}

// CategoryItem 分类项
type CategoryItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TagItem 标签项
type TagItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ArticleDetailItem 文章详情项（用于聚合输出）
type ArticleDetailItem struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Summary     string       `json:"summary"`
	AISummary   string       `json:"ai_summary"`
	Cover       string       `json:"cover"`
	Location    string       `json:"location"`
	IsPublish   bool         `json:"is_publish"`
	IsTop       bool         `json:"is_top"`
	IsEssence   bool         `json:"is_essence"`
	IsOutdated  bool         `json:"is_outdated"`
	PublishTime *string      `json:"publish_time"`
	UpdateTime  *string      `json:"update_time"`
	Category    CategoryItem `json:"category"`
	Tags        []TagItem    `json:"tags"`
}

// ============ 聚合 Tool 输入/输出类型============

// ArticleManageInput article_manage 聚合 tool 输入
type ArticleManageInput struct {
	Action  string               `json:"action"` // list|get|create|update|delete
	Payload ArticleManagePayload `json:"payload"`
}

// ArticleManagePayload article_manage 载荷
type ArticleManagePayload struct {
	// 用于 list
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// 用于 get/update/delete
	ID uint `json:"id"`

	// 用于 create/update
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Summary    *string `json:"summary"`
	AISummary  *string `json:"ai_summary"`
	Cover      *string `json:"cover"`
	Location   *string `json:"location"`
	IsPublish  *bool   `json:"is_publish"`
	IsTop      *bool   `json:"is_top"`
	IsEssence  *bool   `json:"is_essence"`
	IsOutdated *bool   `json:"is_outdated"`
	CategoryID *uint   `json:"category_id"`
	TagIDs     []uint  `json:"tag_ids"`
}

// ArticleManageOutput article_manage 聚合 tool 输出
type ArticleManageOutput struct {
	// list 结果
	List     []ArticleItem `json:"list,omitempty"`
	Total    int64         `json:"total,omitempty"`
	Page     int           `json:"page,omitempty"`
	PageSize int           `json:"page_size,omitempty"`

	// get/create/update 结果
	Item *ArticleDetailItem `json:"item,omitempty"`

	// delete 结果
	DeleteSuccess *bool `json:"delete_success,omitempty"`
	ID            *uint `json:"id,omitempty"`

	// 错误信息（如果有）
	Error string `json:"error,omitempty"`
}

// ============ 服务包装器============

// ArticleWrapper 文章服务包装器
type ArticleWrapper struct {
	articleService *service.ArticleService
}

// NewArticleWrapper 创建文章服务包装器
func NewArticleWrapper(articleService *service.ArticleService) *ArticleWrapper {
	return &ArticleWrapper{articleService: articleService}
}

// ============ 聚合 Tool Handler============

// ManageArticle 文章管理聚合入口
func (w *ArticleWrapper) ManageArticle(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleManageInput,
) (*sdkmcp.CallToolResult, ArticleManageOutput, error) {
	switch input.Action {
	case articleActionList:
		return w.listArticles(input.Payload)
	case articleActionGet:
		return w.getArticle(input.Payload)
	case articleActionCreate:
		return w.createArticle(input.Payload)
	case articleActionUpdate:
		return w.updateArticle(input.Payload)
	case articleActionDelete:
		return w.deleteArticle(input.Payload)
	default:
		return nil, ArticleManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

// listArticles 获取文章列表
func (w *ArticleWrapper) listArticles(payload ArticleManagePayload) (*sdkmcp.CallToolResult, ArticleManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	req := &dto.ListArticlesRequest{Page: page, PageSize: pageSize}
	articles, total, err := w.articleService.List(context.Background(), req)
	if err != nil {
		return nil, ArticleManageOutput{Error: fmt.Sprintf("获取文章列表失败: %v", err)}, nil
	}

	list := make([]ArticleItem, len(articles))
	for i, article := range articles {
		list[i] = convertToArticleItem(article)
	}

	return nil, ArticleManageOutput{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// getArticle 获取文章详情
func (w *ArticleWrapper) getArticle(payload ArticleManagePayload) (*sdkmcp.CallToolResult, ArticleManageOutput, error) {
	if payload.ID == 0 {
		return nil, ArticleManageOutput{Error: "文章 ID 不能为空"}, nil
	}

	article, err := w.articleService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, ArticleManageOutput{Error: fmt.Sprintf("获取文章失败: %v", err)}, nil
	}

	item := convertToArticleDetailItem(*article)
	return nil, ArticleManageOutput{Item: &item}, nil
}

// createArticle 创建文章
func (w *ArticleWrapper) createArticle(payload ArticleManagePayload) (*sdkmcp.CallToolResult, ArticleManageOutput, error) {
	if payload.Title == "" {
		return nil, ArticleManageOutput{Error: "文章标题不能为空"}, nil
	}
	if payload.Content == "" {
		return nil, ArticleManageOutput{Error: "文章内容不能为空"}, nil
	}

	isPublish := false
	req := &dto.CreateArticleRequest{
		Title:      payload.Title,
		Content:    payload.Content,
		Summary:    stringValue(payload.Summary),
		Cover:      stringValue(payload.Cover),
		Location:   stringValue(payload.Location),
		IsPublish:  &isPublish,
		IsTop:      payload.IsTop,
		IsEssence:  payload.IsEssence,
		IsOutdated: payload.IsOutdated,
		CategoryID: payload.CategoryID,
		TagIDs:     payload.TagIDs,
	}

	article, err := w.articleService.Create(context.Background(), req)
	if err != nil {
		return nil, ArticleManageOutput{Error: fmt.Sprintf("创建文章失败: %v", err)}, nil
	}

	item := convertToArticleDetailItem(*article)
	return nil, ArticleManageOutput{Item: &item}, nil
}

// updateArticle 更新文章
func (w *ArticleWrapper) updateArticle(payload ArticleManagePayload) (*sdkmcp.CallToolResult, ArticleManageOutput, error) {
	if payload.ID == 0 {
		return nil, ArticleManageOutput{Error: "文章 ID 不能为空"}, nil
	}

	ctx := context.Background()
	current, err := w.articleService.Get(ctx, payload.ID)
	if err != nil {
		return nil, ArticleManageOutput{Error: fmt.Sprintf("获取当前文章失败: %v", err)}, nil
	}

	req := buildArticleUpdateRequest(current, payload)
	article, err := w.articleService.Update(ctx, payload.ID, req)
	if err != nil {
		return nil, ArticleManageOutput{Error: fmt.Sprintf("更新文章失败: %v", err)}, nil
	}

	item := convertToArticleDetailItem(*article)
	return nil, ArticleManageOutput{Item: &item}, nil
}

// deleteArticle 删除文章
func (w *ArticleWrapper) deleteArticle(payload ArticleManagePayload) (*sdkmcp.CallToolResult, ArticleManageOutput, error) {
	if payload.ID == 0 {
		return nil, ArticleManageOutput{Error: "文章 ID 不能为空"}, nil
	}

	err := w.articleService.Delete(context.Background(), payload.ID)
	if err != nil {
		return nil, ArticleManageOutput{Error: fmt.Sprintf("删除文章失败: %v", err)}, nil
	}

	success := true
	return nil, ArticleManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

// ArticleManageInputSchema 返回 article_manage 的自定义输入 schema
func ArticleManageInputSchema() *jsonschema.Schema {
	listPayload := BuildPayloadSchema(map[string]*jsonschema.Schema{
		"page":      {Type: "integer"},
		"page_size": PageSizeSchema(),
	})
	idPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id": {Type: "integer"},
		},
		"id",
	)
	createPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"title":       {Type: "string"},
			"content":     {Type: "string"},
			"summary":     {Type: "string"},
			"cover":       {Type: "string"},
			"location":    {Type: "string"},
			"is_top":      {Type: "boolean"},
			"is_essence":  {Type: "boolean"},
			"is_outdated": {Type: "boolean"},
			"category_id": {Type: "integer"},
			"tag_ids": {
				Type:  "array",
				Items: &jsonschema.Schema{Type: "integer"},
			},
		},
		"title",
		"content",
	)
	updatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":          {Type: "integer"},
			"title":       {Type: "string"},
			"content":     {Type: "string"},
			"summary":     {Type: "string"},
			"ai_summary":  {Type: "string"},
			"cover":       {Type: "string"},
			"location":    {Type: "string"},
			"is_publish":  {Type: "boolean", Description: "发布状态。没有明确需求请保持为false草稿状态，谨慎操作"},
			"is_top":      {Type: "boolean"},
			"is_essence":  {Type: "boolean"},
			"is_outdated": {Type: "boolean"},
			"category_id": {Type: "integer"},
			"tag_ids": {
				Type:  "array",
				Items: &jsonschema.Schema{Type: "integer"},
			},
		},
		"id",
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					articleActionList,
					articleActionGet,
					articleActionCreate,
					articleActionUpdate,
					articleActionDelete,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(articleActionList, "获取文章列表", listPayload),
			BuildActionSchema(articleActionGet, "获取文章详情", idPayload),
			BuildActionSchema(articleActionCreate, "创建文章（默认草稿状态）", createPayload),
			BuildActionSchema(articleActionUpdate, "更新文章内容或状态", updatePayload),
			BuildActionSchema(articleActionDelete, "删除文章。风险操作，谨慎使用，不可恢复", idPayload),
		},
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func buildArticleUpdateRequest(current *dto.ArticleAdminDetailResponse, payload ArticleManagePayload) *dto.UpdateArticleRequest {
	summary := current.Summary
	if payload.Summary != nil {
		summary = *payload.Summary
	}

	aiSummary := current.AISummary
	if payload.AISummary != nil {
		aiSummary = *payload.AISummary
	}

	cover := current.Cover
	if payload.Cover != nil {
		cover = *payload.Cover
	}

	location := current.Location
	if payload.Location != nil {
		location = *payload.Location
	}

	categoryID := articleCategoryIDPtr(current.Category.ID)
	if payload.CategoryID != nil {
		if *payload.CategoryID == 0 {
			categoryID = nil
		} else {
			categoryID = articleCategoryIDPtr(*payload.CategoryID)
		}
	}

	var expectedUpdatedAt = current.Revision
	var expectedUpdatedAtPtr = &expectedUpdatedAt
	if expectedUpdatedAt.IsZero() {
		expectedUpdatedAtPtr = nil
	}

	return &dto.UpdateArticleRequest{
		Title:             payload.Title,
		Content:           payload.Content,
		Summary:           summary,
		AISummary:         aiSummary,
		Cover:             cover,
		Location:          location,
		IsPublish:         payload.IsPublish,
		IsTop:             payload.IsTop,
		IsEssence:         payload.IsEssence,
		IsOutdated:        payload.IsOutdated,
		CategoryID:        categoryID,
		TagIDs:            payload.TagIDs,
		ExpectedUpdatedAt: expectedUpdatedAtPtr,
	}
}

func articleCategoryIDPtr(id uint) *uint {
	if id == 0 {
		return nil
	}
	value := id
	return &value
}

func convertArticleCategory(category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}) CategoryItem {
	return CategoryItem{ID: category.ID, Name: category.Name}
}

func convertArticleTags(tags []struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}) []TagItem {
	result := make([]TagItem, len(tags))
	for i, tag := range tags {
		result[i] = TagItem{ID: tag.ID, Name: tag.Name}
	}
	return result
}

func convertArticleTimes(publishTime, updateTime interface{ String() string }) (*string, *string) {
	return ToTimeStringPtr(publishTime), ToTimeStringPtr(updateTime)
}

func convertToArticleItem(item dto.ArticleListResponse) ArticleItem {
	publishTime, updateTime := convertArticleTimes(item.PublishTime, item.UpdateTime)
	return ArticleItem{
		ID:           item.ID,
		Title:        item.Title,
		Cover:        item.Cover,
		Location:     item.Location,
		IsPublish:    item.IsPublish,
		IsTop:        item.IsTop,
		IsEssence:    item.IsEssence,
		IsOutdated:   item.IsOutdated,
		ViewCount:    item.ViewCount,
		CommentCount: item.CommentCount,
		PublishTime:  publishTime,
		UpdateTime:   updateTime,
		Category:     convertArticleCategory(item.Category),
		Tags:         convertArticleTags(item.Tags),
	}
}

func convertToArticleDetailItem(item dto.ArticleAdminDetailResponse) ArticleDetailItem {
	publishTime, updateTime := convertArticleTimes(item.PublishTime, item.UpdateTime)
	return ArticleDetailItem{
		ID:          item.ID,
		Title:       item.Title,
		Content:     item.Content,
		Summary:     item.Summary,
		AISummary:   item.AISummary,
		Cover:       item.Cover,
		Location:    item.Location,
		IsPublish:   item.IsPublish,
		IsTop:       item.IsTop,
		IsEssence:   item.IsEssence,
		IsOutdated:  item.IsOutdated,
		PublishTime: publishTime,
		UpdateTime:  updateTime,
		Category:    convertArticleCategory(item.Category),
		Tags:        convertArticleTags(item.Tags),
	}
}
