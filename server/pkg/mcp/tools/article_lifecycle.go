package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"flec_blog/internal/service"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ArticleUpdatePublishedInput article_update_published 输入。
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

// ArticleUpdatePublishedOutput 已发布文章更新输出。
type ArticleUpdatePublishedOutput struct {
	Item ArticleDetailItem `json:"item"`
}

// ArticleUnpublishInput article_unpublish 输入。
type ArticleUnpublishInput struct {
	ID uint `json:"id"`
}

// ArticleUnpublishOutput article_unpublish 输出。
type ArticleUnpublishOutput struct {
	Item               ArticleDetailItem `json:"item"`
	AlreadyUnpublished bool              `json:"already_unpublished"`
}

// ArticleDeleteInput article_delete 输入。
type ArticleDeleteInput struct {
	ID uint `json:"id"`
}

// ArticleDeleteOutput article_delete 输出。
type ArticleDeleteOutput struct {
	Deleted bool `json:"deleted"`
	ID      uint `json:"id"`
}

// UpdatePublishedArticleTool 更新现有已发布文章，拒绝修改草稿。
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
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("更新已发布文章失败: %w", err)
	}
	if !article.IsPublish {
		return nil, ArticleUpdatePublishedOutput{}, fmt.Errorf("更新已发布文章失败: 返回结果意外变为草稿状态")
	}

	return nil, ArticleUpdatePublishedOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// UnpublishArticleTool 显式取消发布文章；重复调用保持幂等。
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
		if errors.Is(err, service.ErrArticleUpdateConflict) {
			latest, latestErr := w.articleService.Get(ctx, input.ID)
			if latestErr == nil && !latest.IsPublish {
				return nil, ArticleUnpublishOutput{
					Item:               convertToArticleDetailItem(*latest),
					AlreadyUnpublished: true,
				}, nil
			}
		}
		return nil, ArticleUnpublishOutput{}, fmt.Errorf("取消发布文章失败: %w", err)
	}
	if article.IsPublish {
		return nil, ArticleUnpublishOutput{}, fmt.Errorf("取消发布文章失败: 返回结果仍为已发布状态")
	}

	return nil, ArticleUnpublishOutput{Item: convertToArticleDetailItem(*article)}, nil
}

// DeleteArticleTool 显式删除文章。该操作不可恢复。
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
