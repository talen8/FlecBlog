package mcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/logger"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultReadPageSize = 20
	maxReadPageSize     = 100
)

// ============ 类型定义============

// Item 文章列表项
type Item struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	IsPublish bool   `json:"is_publish"`
}

// DetailItem 文章详情项
type DetailItem struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsPublish bool   `json:"is_publish"`
}

// ============ 工具输入/输出类型============

// ListInput article_list 输入
type ListInput struct {
	Page      int   `json:"page,omitempty" jsonschema:"页码，从 1 开始"`
	PageSize  int   `json:"page_size,omitempty" jsonschema:"每页条数，取值 1-100，默认 20"`
	IsPublish *bool `json:"is_publish,omitempty" jsonschema:"按发布状态筛选，不传表示不限"`
}

// SearchInput article_search 输入
type SearchInput struct {
	Keyword   string `json:"keyword" jsonschema:"搜索关键词，匹配文章标题与正文"`
	Page      int    `json:"page,omitempty" jsonschema:"页码，从 1 开始"`
	PageSize  int    `json:"page_size,omitempty" jsonschema:"每页条数，取值 1-100，默认 20"`
	IsPublish *bool  `json:"is_publish,omitempty" jsonschema:"按发布状态筛选，不传表示不限"`
}

// GetInput article_get 输入
type GetInput struct {
	ID uint `json:"id" jsonschema:"文章 ID，来自 article_list 或 article_search 的返回结果"`
}

// CollectionOutput 文章列表/搜索输出
type CollectionOutput struct {
	Items    []Item `json:"items"`
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// CreateInput article_create 输入
type CreateInput struct {
	Title   string `json:"title" jsonschema:"文章标题"`
	Content string `json:"content" jsonschema:"文章正文，Markdown 格式"`
}

// UpdateInput article_update 输入
type UpdateInput struct {
	ID      uint    `json:"id" jsonschema:"文章 ID，来自 article_list 或 article_search 的返回结果"`
	Title   *string `json:"title,omitempty" jsonschema:"文章标题"`
	Content *string `json:"content,omitempty" jsonschema:"文章正文，Markdown 格式"`
}

// DeleteInput article_delete 输入
type DeleteInput struct {
	ID uint `json:"id" jsonschema:"文章 ID，来自 article_list 或 article_search 的返回结果"`
}

// ItemOutput 文章详情输出
type ItemOutput struct {
	Item DetailItem `json:"item"`
}

// DeleteOutput 文章删除输出
type DeleteOutput struct {
	Deleted bool   `json:"deleted"`
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ============ 服务包装器============

// Wrapper 文章服务包装器
type Wrapper struct {
	articleService *service.ArticleService
}

// NewWrapper 创建文章服务包装器
func NewWrapper(articleService *service.ArticleService) *Wrapper {
	return &Wrapper{articleService: articleService}
}

// ============ 细粒度 Tool Handler============

// ListTool 分页列出文章
func (w *Wrapper) ListTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ListInput,
) (*sdkmcp.CallToolResult, CollectionOutput, error) {
	page, pageSize := normalizeReadPage(input.Page, input.PageSize)
	return w.queryArticles(ctx, &dto.ListArticlesRequest{
		Page:      page,
		PageSize:  pageSize,
		IsPublish: input.IsPublish,
	})
}

// SearchTool 按关键词搜索文章
func (w *Wrapper) SearchTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input SearchInput,
) (*sdkmcp.CallToolResult, CollectionOutput, error) {
	keyword := strings.TrimSpace(input.Keyword)
	if keyword == "" {
		return nil, CollectionOutput{}, fmt.Errorf("搜索关键词不能为空")
	}
	page, pageSize := normalizeReadPage(input.Page, input.PageSize)
	return w.queryArticles(ctx, &dto.ListArticlesRequest{
		Page:      page,
		PageSize:  pageSize,
		Keyword:   keyword,
		IsPublish: input.IsPublish,
	})
}

// GetTool 按 ID 获取文章详情
func (w *Wrapper) GetTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input GetInput,
) (*sdkmcp.CallToolResult, ItemOutput, error) {
	if input.ID == 0 {
		return nil, ItemOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	article, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ItemOutput{}, toolError("获取文章", err)
	}
	return nil, ItemOutput{Item: convertToDetailItem(*article)}, nil
}

// CreateTool 创建文章
func (w *Wrapper) CreateTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input CreateInput,
) (*sdkmcp.CallToolResult, ItemOutput, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ItemOutput{}, fmt.Errorf("文章标题不能为空")
	}
	if strings.TrimSpace(input.Content) == "" {
		return nil, ItemOutput{}, fmt.Errorf("文章内容不能为空")
	}
	isPublish := false
	article, err := w.articleService.Create(ctx, &dto.CreateArticleRequest{
		Title:     title,
		Content:   input.Content,
		IsPublish: &isPublish,
	})
	if err != nil {
		return nil, ItemOutput{}, toolError("创建文章", err)
	}
	return nil, ItemOutput{Item: convertToDetailItem(*article)}, nil
}

// UpdateTool 更新草稿文章
func (w *Wrapper) UpdateTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input UpdateInput,
) (*sdkmcp.CallToolResult, ItemOutput, error) {
	if input.ID == 0 {
		return nil, ItemOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return nil, ItemOutput{}, fmt.Errorf("文章标题不能为空")
	}
	if input.Content != nil && strings.TrimSpace(*input.Content) == "" {
		return nil, ItemOutput{}, fmt.Errorf("文章内容不能为空")
	}
	if input.Title == nil && input.Content == nil {
		return nil, ItemOutput{}, fmt.Errorf("请至少提供 title 或 content 中的一项")
	}
	current, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, ItemOutput{}, toolError("更新文章", err)
	}
	if current.IsPublish {
		return nil, ItemOutput{}, errors.New("文章已发布，MCP 只允许修改草稿文章，请在管理后台处理")
	}
	article, err := w.articleService.UpdateForMcp(ctx, input.ID, input.Title, input.Content)
	if err != nil {
		return nil, ItemOutput{}, toolError("更新文章", err)
	}
	return nil, ItemOutput{Item: convertToDetailItem(*article)}, nil
}

// DeleteTool 删除草稿文章
func (w *Wrapper) DeleteTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input DeleteInput,
) (*sdkmcp.CallToolResult, DeleteOutput, error) {
	if input.ID == 0 {
		return nil, DeleteOutput{}, fmt.Errorf("文章 ID 不能为空")
	}
	current, err := w.articleService.Get(ctx, input.ID)
	if err != nil {
		return nil, DeleteOutput{}, toolError("删除文章", err)
	}
	if current.IsPublish {
		return nil, DeleteOutput{}, errors.New("文章已发布，MCP 只允许删除草稿文章，请在管理后台处理")
	}
	if err := w.articleService.Delete(ctx, input.ID); err != nil {
		return nil, DeleteOutput{}, toolError("删除文章", err)
	}
	return nil, DeleteOutput{
		Deleted: true,
		ID:      input.ID,
		Title:   current.Title,
		Content: current.Content,
	}, nil
}

// ============ 辅助函数============

// toolError 将服务层错误转换为供调用方理解的提示
func toolError(action string, err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "record not found"), strings.Contains(msg, "未找到"):
		logger.Warn("MCP %s失败: %v", action, err)
		return fmt.Errorf("%s失败：文章不存在，可先用 article_search 确认文章 ID", action)
	default:
		logger.Error("MCP %s失败: %v", action, err)
		return fmt.Errorf("%s失败：请检查参数是否合法后重试", action)
	}
}

// queryArticles 查询文章列表并转换为列表输出
func (w *Wrapper) queryArticles(
	ctx context.Context,
	req *dto.ListArticlesRequest,
) (*sdkmcp.CallToolResult, CollectionOutput, error) {
	articles, total, err := w.articleService.List(ctx, req)
	if err != nil {
		return nil, CollectionOutput{}, toolError("获取文章列表", err)
	}
	items := make([]Item, len(articles))
	for i, article := range articles {
		items[i] = convertToItem(article)
	}
	return nil, CollectionOutput{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// normalizeReadPage 归一化分页参数，页码从 1 开始，每页条数限制在 1-100
func normalizeReadPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultReadPageSize
	}
	if pageSize > maxReadPageSize {
		pageSize = maxReadPageSize
	}
	return page, pageSize
}

// convertToItem 将后台文章列表响应转换为列表项
func convertToItem(item dto.ArticleListResponse) Item {
	return Item{
		ID:        item.ID,
		Title:     item.Title,
		IsPublish: item.IsPublish,
	}
}

// convertToDetailItem 将后台文章详情响应转换为详情项
func convertToDetailItem(item dto.ArticleAdminDetailResponse) DetailItem {
	return DetailItem{
		ID:        item.ID,
		Title:     item.Title,
		Content:   item.Content,
		IsPublish: item.IsPublish,
	}
}
