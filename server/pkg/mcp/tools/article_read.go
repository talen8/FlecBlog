package tools

import (
	"context"
	"fmt"
	"strings"

	"flec_blog/internal/dto"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultArticleReadPageSize = 20
	maxArticleReadPageSize     = 100
)

// ArticleListInput article_list 输入。
type ArticleListInput struct {
	Page      int   `json:"page,omitempty"`
	PageSize  int   `json:"page_size,omitempty"`
	IsPublish *bool `json:"is_publish,omitempty"`
}

// ArticleSearchInput article_search 输入。
type ArticleSearchInput struct {
	Keyword   string `json:"keyword"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	IsPublish *bool  `json:"is_publish,omitempty"`
}

// ArticleGetInput article_get 输入。
type ArticleGetInput struct {
	ID uint `json:"id"`
}

// ArticleCollectionOutput 文章列表/搜索输出。
type ArticleCollectionOutput struct {
	Items    []ArticleItem `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// ArticleGetOutput 文章详情输出。
type ArticleGetOutput struct {
	Item ArticleDetailItem `json:"item"`
}

// ListArticleTool 分页列出后台文章，可选按发布状态筛选。
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

// SearchArticleTool 按标题/正文关键词搜索后台文章，包含草稿与已发布文章。
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

// GetArticleTool 按 ID 获取完整后台文章详情。
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
