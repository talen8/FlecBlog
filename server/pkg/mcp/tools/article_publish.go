package tools

import (
	"context"
	"errors"
	"fmt"

	"flec_blog/internal/service"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ArticlePublishInput article_publish 输入。
type ArticlePublishInput struct {
	ID uint `json:"id"`
}

// ArticlePublishOutput article_publish 输出。
type ArticlePublishOutput struct {
	Item             ArticleDetailItem `json:"item"`
	AlreadyPublished bool              `json:"already_published"`
}

// PublishArticleTool 显式发布文章；重复调用保持幂等。
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
		if errors.Is(err, service.ErrArticleUpdateConflict) {
			latest, latestErr := w.articleService.Get(ctx, input.ID)
			if latestErr == nil && latest.IsPublish {
				return nil, ArticlePublishOutput{
					Item:             convertToArticleDetailItem(*latest),
					AlreadyPublished: true,
				}, nil
			}
		}
		return nil, ArticlePublishOutput{}, fmt.Errorf("发布文章失败: %w", err)
	}
	if !article.IsPublish {
		return nil, ArticlePublishOutput{}, fmt.Errorf("发布文章失败: 返回结果仍为草稿状态")
	}

	return nil, ArticlePublishOutput{Item: convertToArticleDetailItem(*article)}, nil
}
