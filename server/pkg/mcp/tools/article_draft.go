package tools

import (
	"context"
	"fmt"
	"strings"

	"flec_blog/internal/dto"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ArticleCreateDraftInput article_create_draft 输入。
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

// ArticleUpdateDraftInput article_update_draft 输入。
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

// ArticleDraftOutput 草稿写入输出。
type ArticleDraftOutput struct {
	Item ArticleDetailItem `json:"item"`
}

// CreateArticleDraftTool 创建草稿文章，永不直接发布。
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

// UpdateArticleDraftTool 更新现有草稿，拒绝修改已发布文章。
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
		Summary:    input.Summary,
		AISummary:  input.AISummary,
		Cover:      input.Cover,
		Location:   input.Location,
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
