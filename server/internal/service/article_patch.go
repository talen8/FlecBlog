package service

import (
	"context"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/utils"
)

// Patch 局部更新文章。
// nil 字段保留原值；非 nil 字段按请求显式更新。
func (s *ArticleService) Patch(ctx context.Context, id uint, req *dto.PatchArticleRequest) (*dto.ArticleAdminDetailResponse, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil && *req.Title == "" {
		return nil, fmt.Errorf("文章标题不能为空")
	}
	if req.Content != nil && *req.Content == "" {
		return nil, fmt.Errorf("文章内容不能为空")
	}

	// 验证新分类是否存在。0 表示清除分类。
	if req.CategoryID != nil && *req.CategoryID > 0 {
		if _, err := s.categoryRepo.Get(ctx, *req.CategoryID); err != nil {
			return nil, fmt.Errorf("分类不存在: %w", err)
		}
	}

	oldCover := article.Cover
	oldContent := article.Content
	oldIsPublish := article.IsPublish

	applyArticlePatch(article, req)

	// 如果是发布状态且没有发布时间，自动设置发布时间。
	if article.IsPublish && article.PublishTime == nil {
		now := utils.Now().Time
		article.PublishTime = &now
	}

	if err := s.articleRepo.Update(article, req.TagIDs); err != nil {
		return nil, err
	}

	// 只有显式更新封面时才处理文件状态。
	if req.Cover != nil && s.fileService != nil && oldCover != article.Cover {
		if oldCover != "" {
			_ = s.fileService.MarkAsUnused(oldCover)
		}
		if article.Cover != "" {
			_ = s.fileService.MarkAsUsed(article.Cover)
		}
	}

	// 只有显式更新正文时才处理正文媒体文件状态。
	if req.Content != nil && oldContent != article.Content {
		s.updateContentFileStatus(oldContent, article.Content)
	}

	// 如果从草稿变为发布状态，异步发送订阅推送。
	if !oldIsPublish && article.IsPublish && s.subscriberService != nil {
		go func(ctx context.Context, articleID uint) {
			if err := s.subscriberService.SendArticleNotification(ctx, article); err != nil {
				logger.Warn("发送文章推送失败 (文章ID: %d): %v", articleID, err)
			}
		}(ctx, article.ID)
	}

	return s.Get(ctx, id)
}

func applyArticlePatch(article *model.Article, req *dto.PatchArticleRequest) {
	if req.Title != nil {
		article.Title = *req.Title
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.Summary != nil {
		article.Summary = *req.Summary
	}
	if req.AISummary != nil {
		article.AISummary = *req.AISummary
	}
	if req.Cover != nil {
		article.Cover = *req.Cover
	}
	if req.Location != nil {
		article.Location = *req.Location
	}
	if req.CategoryID != nil {
		if *req.CategoryID == 0 {
			article.CategoryID = nil
		} else {
			categoryID := *req.CategoryID
			article.CategoryID = &categoryID
		}
	}
	if req.IsTop != nil {
		article.IsTop = *req.IsTop
	}
	if req.IsEssence != nil {
		article.IsEssence = *req.IsEssence
	}
	if req.IsOutdated != nil {
		article.IsOutdated = *req.IsOutdated
	}
	if req.IsPublish != nil {
		article.IsPublish = *req.IsPublish
	}
}
