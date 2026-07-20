package tools

import (
	"context"
	"fmt"
	"strings"

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

	defaultArticleReadPageSize = 20
	maxArticleReadPageSize     = 100
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

// ArticleDetailItem 文章详情项
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
	Action  string               `json:"action"`
	Payload ArticleManagePayload `json:"payload"`
}

// ArticleManagePayload article_manage 载荷
type ArticleManagePayload struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Summary  string `json:"summary"`

	// 用于 create/update
	AISummary  string `json:"ai_summary"`
	Cover      string `json:"cover"`
	Location   string `json:"location"`
	IsPublish  *bool  `json:"is_publish"`
	IsTop      *bool  `json:"is_top"`
	IsEssence  *bool  `json:"is_essence"`
	IsOutdated *bool  `json:"is_outdated"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// ArticleManageOutput article_manage 聚合 tool 输出
type ArticleManageOutput struct {
	List          []ArticleItem      `json:"list,omitempty"`
	Total         int64              `json:"total,omitempty"`
	Page          int                `json:"page,omitempty"`
	PageSize      int                `json:"page_size,omitempty"`
	Item          *ArticleDetailItem `json:"item,omitempty"`
	DeleteSuccess *bool              `json:"delete_success,omitempty"`
	ID            *uint              `json:"id,omitempty"`
	Error         string             `json:"error,omitempty"`
}

// ============ 细粒度 Tool 输入/输出类型============

// ArticleListInput article_list 输入
type ArticleListInput struct {
	Page      int   `json:"page,omitempty"`
	PageSize  int   `json:"page_size,omitempty"`
	IsPublish *bool `json:"is_publish,omitempty"`
}

// ArticleSearchInput article_search 输入
type ArticleSearchInput struct {
	Keyword   string `json:"keyword"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	IsPublish *bool  `json:"is_publish,omitempty"`
}

// ArticleGetInput article_get 输入
type ArticleGetInput struct {
	ID uint `json:"id"`
}

// ArticleCollectionOutput 文章列表/搜索输出
type ArticleCollectionOutput struct {
	Items    []ArticleItem `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// ArticleGetOutput 文章详情输出
type ArticleGetOutput struct {
	Item ArticleDetailItem `json:"item"`
}

// ArticleCreateDraftInput article_create_draft 输入
type ArticleCreateDraftInput struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Summary    string `json:"summary,omitempty"`
	Cover      string `json:"cover,omitempty"`
	Location   string `json:"location,omitempty"`
	IsTop      *bool  `json:"is_top,omitempty"`
	IsEssence  *bool  `json:"is_essence,omitempty"`
	IsOutdated *bool  `json:"is_outdated,omitempty"`
	CategoryID *uint  `json:"category_id,omitempty"`
	TagIDs     []uint `json:"tag_ids,omitempty"`
}

// ArticleUpdateDraftInput article_update_draft 输入
type ArticleUpdateDraftInput struct {
	ID         uint    `json:"id"`
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty"`
	Summary    *string `json:"summary,omitempty"`
	AISummary  *string `json:"ai_summary,omitempty"`
	Cover      *string `json:"cover,omitempty"`
	Location   *string `json:"location,omitempty"`
	IsTop      *bool   `json:"is_top,omitempty"`
	IsEssence  *bool   `json:"is_essence,omitempty"`
	IsOutdated *bool   `json:"is_outdated,omitempty"`
	CategoryID *uint   `json:"category_id,omitempty"`
	TagIDs     []uint  `json:"tag_ids,omitempty"`
}

// ArticleDraftOutput 草稿写入输出
type ArticleDraftOutput struct {
	Item ArticleDetailItem `json:"item"`
}

// ArticlePublishInput article_publish 输入
type ArticlePublishInput struct {
	ID uint `json:"id"`
}

// ArticlePublishOutput article_publish 输出
type ArticlePublishOutput struct {
	Item             ArticleDetailItem `json:"item"`
	AlreadyPublished bool              `json:"already_published"`
}

// ArticleUpdatePublishedInput article_update_published 输入
type ArticleUpdatePublishedInput struct {
	ID         uint    `json:"id"`
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty"`
	Summary    *string `json:"summary,omitempty"`
	AISummary  *string `json:"ai_summary,omitempty"`
	Cover      *string `json:"cover,omitempty"`
	Location   *string `json:"location,omitempty"`
	IsTop      *bool   `json:"is_top,omitempty"`
	IsEssence  *bool   `json:"is_essence,omitempty"`
	IsOutdated *bool   `json:"is_outdated,omitempty"`
	CategoryID *uint   `json:"category_id,omitempty"`
	TagIDs     []uint  `json:"tag_ids,omitempty"`
}

// ArticleUpdatePublishedOutput 已发布文章更新输出
type ArticleUpdatePublishedOutput struct {
	Item ArticleDetailItem `json:"item"`
}

// ArticleUnpublishInput article_unpublish 输入
type ArticleUnpublishInput struct {
	ID uint `json:"id"`
}

// ArticleUnpublishOutput article_unpublish 输出
type ArticleUnpublishOutput struct {
	Item               ArticleDetailItem `json:"item"`
	AlreadyUnpublished bool              `json:"already_unpublished"`
}

// ArticleDeleteInput article_delete 输入
type ArticleDeleteInput struct {
	ID uint `json:"id"`
}

// ArticleDeleteOutput article_delete 输出
type ArticleDeleteOutput struct {
	Deleted bool `json:"deleted"`
	ID      uint `json:"id"`
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

// listArticles 文章列表查询
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
	return nil, ArticleManageOutput{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

// getArticle 文章详情查询
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
		Summary:    payload.Summary,
		Cover:      payload.Cover,
		Location:   payload.Location,
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

// ============ 细粒度 Tool Handler============

// ListArticleTool 分页列出后台文章
func (w *ArticleWrapper) ListArticleTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleListInput,
) (*sdkmcp.CallToolResult, ArticleCollectionOutput, error) {
	page, pageSize := normalizeArticleReadPage(input.Page, input.PageSize)
	return w.queryArticles(ctx, &dto.ListArticlesRequest{
		Page:      page,
		PageSize:  pageSize,
		IsPublish: input.IsPublish,
	})
}

// SearchArticleTool 按关键词搜索后台文章
func (w *ArticleWrapper) SearchArticleTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleSearchInput,
) (*sdkmcp.CallToolResult, ArticleCollectionOutput, error) {
	keyword := strings.TrimSpace(input.Keyword)
	if keyword == "" {
		return nil, ArticleCollectionOutput{}, fmt.Errorf("搜索关键词不能为空")
	}
	page, pageSize := normalizeArticleReadPage(input.Page, input.PageSize)
	return w.queryArticles(ctx, &dto.ListArticlesRequest{
		Page:      page,
		PageSize:  pageSize,
		Keyword:   keyword,
		IsPublish: input.IsPublish,
	})
}

// GetArticleTool 按 ID 获取完整后台文章详情
func (w *ArticleWrapper) GetArticleTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleGetInput,
) (*sdkmcp.CallToolResult, ArticleGetOutput, error) {
	if input.ID == 0 {
		return nil, ArticleGetOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	article, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ArticleGetOutput{}, fmt.Errorf("获取文章失败: %w", err)
	}
	return nil, ArticleGetOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// CreateArticleDraftTool 创建草稿文章
func (w *ArticleWrapper) CreateArticleDraftTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleCreateDraftInput,
) (*sdkmcp.CallToolResult, ArticleDraftOutput, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ArticleDraftOutput{}, fmt.Errorf("文章标题不能为空")
	}
	if strings.TrimSpace(input.Content) == "" {
		return nil, ArticleDraftOutput{}, fmt.Errorf("文章内容不能为空")
	}
	isPublish := false
	article, err := w.articleService.Create(ctx, &dto.CreateArticleRequest{
		Title:      title,
		Content:    input.Content,
		Summary:    input.Summary,
		Cover:      input.Cover,
		Location:   input.Location,
		IsPublish:  &isPublish,
		IsTop:      input.IsTop,
		IsEssence:  input.IsEssence,
		IsOutdated: input.IsOutdated,
		CategoryID: input.CategoryID,
		TagIDs:     input.TagIDs,
	})
	if err != nil {
		return nil, ArticleDraftOutput{}, fmt.Errorf("创建文章草稿失败: %w", err)
	}
	if article.IsPublish {
		return nil, ArticleDraftOutput{}, fmt.Errorf("创建文章草稿失败: 返回结果意外处于已发布状态")
	}
	return nil, ArticleDraftOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// UpdateArticleDraftTool 更新草稿文章
func (w *ArticleWrapper) UpdateArticleDraftTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleUpdateDraftInput,
) (*sdkmcp.CallToolResult, ArticleDraftOutput, error) {
	if input.ID == 0 {
		return nil, ArticleDraftOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return nil, ArticleDraftOutput{}, fmt.Errorf("文章标题不能为空")
	}
	if input.Content != nil && strings.TrimSpace(*input.Content) == "" {
		return nil, ArticleDraftOutput{}, fmt.Errorf("文章内容不能为空")
	}
	current, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ArticleDraftOutput{}, fmt.Errorf("获取当前文章失败: %w", err)
	}
	if current.IsPublish {
		return nil, ArticleDraftOutput{}, fmt.Errorf("article_update_draft 仅允许修改草稿文章")
	}
	payload := ArticleManagePayload{
		ID:         input.ID,
		Summary:    stringValuePtr(input.Summary),
		AISummary:  stringValuePtr(input.AISummary),
		Cover:      stringValuePtr(input.Cover),
		Location:   stringValuePtr(input.Location),
		IsTop:      input.IsTop,
		IsEssence:  input.IsEssence,
		IsOutdated: input.IsOutdated,
		CategoryID: input.CategoryID,
		TagIDs:     input.TagIDs,
	}
	if input.Title != nil {
		payload.Title = *input.Title
	}
	if input.Content != nil {
		payload.Content = *input.Content
	}
	article, err := w.articleService.Update(ctx, input.ID, buildArticleUpdateRequest(current, payload))
	if err != nil {
		return nil, ArticleDraftOutput{}, fmt.Errorf("更新文章草稿失败: %w", err)
	}
	if article.IsPublish {
		return nil, ArticleDraftOutput{}, fmt.Errorf("更新文章草稿失败: 返回结果意外处于已发布状态")
	}
	return nil, ArticleDraftOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// PublishArticleTool 发布文章
func (w *ArticleWrapper) PublishArticleTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticlePublishInput,
) (*sdkmcp.CallToolResult, ArticlePublishOutput, error) {
	if input.ID == 0 {
		return nil, ArticlePublishOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	ctx := context.Background()
	current, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ArticlePublishOutput{}, fmt.Errorf("获取当前文章失败: %w", err)
	}
	if current.IsPublish {
		return nil, ArticlePublishOutput{
			Item:             convertToArticleDetailItem(*current),
			AlreadyPublished: true,
		}, nil
	}
	publish := true
	article, err := w.articleService.Update(ctx, input.ID, buildArticleUpdateRequest(current, ArticleManagePayload{
		ID:        input.ID,
		IsPublish: &publish,
	}))
	if err != nil {
		return nil, ArticlePublishOutput{}, fmt.Errorf("发布文章失败: %w", err)
	}
	if !article.IsPublish {
		return nil, ArticlePublishOutput{}, fmt.Errorf("发布文章失败: 返回结果仍为草稿状态")
	}
	return nil, ArticlePublishOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// UpdatePublishedArticleTool 更新已发布文章
func (w *ArticleWrapper) UpdatePublishedArticleTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleUpdatePublishedInput,
) (*sdkmcp.CallToolResult, ArticleUpdatePublishedOutput, error) {
	if input.ID == 0 {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("文章标题不能为空")
	}
	if input.Content != nil && strings.TrimSpace(*input.Content) == "" {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("文章内容不能为空")
	}
	current, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("获取当前文章失败: %w", err)
	}
	if !current.IsPublish {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("article_update_published 仅允许修改已发布文章")
	}
	payload := ArticleManagePayload{
		ID:         input.ID,
		Summary:    stringValuePtr(input.Summary),
		AISummary:  stringValuePtr(input.AISummary),
		Cover:      stringValuePtr(input.Cover),
		Location:   stringValuePtr(input.Location),
		IsTop:      input.IsTop,
		IsEssence:  input.IsEssence,
		IsOutdated: input.IsOutdated,
		CategoryID: input.CategoryID,
		TagIDs:     input.TagIDs,
	}
	if input.Title != nil {
		payload.Title = *input.Title
	}
	if input.Content != nil {
		payload.Content = *input.Content
	}
	article, err := w.articleService.Update(ctx, input.ID, buildArticleUpdateRequest(current, payload))
	if err != nil {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("更新已发布文章失败: %w", err)
	}
	if !article.IsPublish {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("更新已发布文章失败: 返回结果意外变为草稿状态")
	}
	return nil, ArticleUpdatePublishedOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// UnpublishArticleTool 取消发布文章
func (w *ArticleWrapper) UnpublishArticleTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleUnpublishInput,
) (*sdkmcp.CallToolResult, ArticleUnpublishOutput, error) {
	if input.ID == 0 {
		return nil, ArticleUnpublishOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	ctx := context.Background()
	current, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ArticleUnpublishOutput{}, fmt.Errorf("获取当前文章失败: %w", err)
	}
	if !current.IsPublish {
		return nil, ArticleUnpublishOutput{
			Item:               convertToArticleDetailItem(*current),
			AlreadyUnpublished: true,
		}, nil
	}
	publish := false
	article, err := w.articleService.Update(ctx, input.ID, buildArticleUpdateRequest(current, ArticleManagePayload{
		ID:        input.ID,
		IsPublish: &publish,
	}))
	if err != nil {
		return nil, ArticleUnpublishOutput{}, fmt.Errorf("取消发布文章失败: %w", err)
	}
	if article.IsPublish {
		return nil, ArticleUnpublishOutput{}, fmt.Errorf("取消发布文章失败: 返回结果仍为已发布状态")
	}
	return nil, ArticleUnpublishOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// DeleteArticleTool 删除文章
func (w *ArticleWrapper) DeleteArticleTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleDeleteInput,
) (*sdkmcp.CallToolResult, ArticleDeleteOutput, error) {
	if input.ID == 0 {
		return nil, ArticleDeleteOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	if err := w.articleService.Delete(ctx, input.ID); err != nil {
		return nil, ArticleDeleteOutput{}, fmt.Errorf("删除文章失败: %w", err)
	}
	return nil, ArticleDeleteOutput{Deleted: true, ID: input.ID}, nil
}

// ============ 辅助函数============

func (w *ArticleWrapper) queryArticles(
	ctx context.Context,
	req *dto.ListArticlesRequest,
) (*sdkmcp.CallToolResult, ArticleCollectionOutput, error) {
	articles, total, err := w.articleService.List(ctx, req)
	if err != nil {
		return nil, ArticleCollectionOutput{}, fmt.Errorf("获取文章列表失败: %w", err)
	}
	items := make([]ArticleItem, len(articles))
	for i, article := range articles {
		items[i] = convertToArticleItem(article)
	}
	return nil, ArticleCollectionOutput{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func normalizeArticleReadPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultArticleReadPageSize
	}
	if pageSize > maxArticleReadPageSize {
		pageSize = maxArticleReadPageSize
	}
	return page, pageSize
}

func stringValuePtr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func buildArticleUpdateRequest(current *dto.ArticleAdminDetailResponse, payload ArticleManagePayload) *dto.UpdateArticleRequest {
	summary := current.Summary
	if payload.Summary != "" {
		summary = payload.Summary
	}
	aiSummary := current.AISummary
	if payload.AISummary != "" {
		aiSummary = payload.AISummary
	}
	cover := current.Cover
	if payload.Cover != "" {
		cover = payload.Cover
	}
	location := current.Location
	if payload.Location != "" {
		location = payload.Location
	}
	categoryID := articleCategoryIDPtr(current.Category.ID)
	if payload.CategoryID != nil {
		if *payload.CategoryID == 0 {
			categoryID = nil
		} else {
			categoryID = articleCategoryIDPtr(*payload.CategoryID)
		}
	}
	return &dto.UpdateArticleRequest{
		Title:      payload.Title,
		Content:    payload.Content,
		Summary:    summary,
		AISummary:  aiSummary,
		Cover:      cover,
		Location:   location,
		IsPublish:  payload.IsPublish,
		IsTop:      payload.IsTop,
		IsEssence:  payload.IsEssence,
		IsOutdated: payload.IsOutdated,
		CategoryID: categoryID,
		TagIDs:     payload.TagIDs,
	}
}

func articleCategoryIDPtr(id uint) *uint {
	if id == 0 {
		return nil
	}
	value := id
	return &value
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
			"is_publish":  {Type: "boolean"},
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
			BuildActionSchema(articleActionCreate, "创建文章", createPayload),
			BuildActionSchema(articleActionUpdate, "更新文章", updatePayload),
			BuildActionSchema(articleActionDelete, "删除文章", idPayload),
		},
	}
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
