package service

import (
	"context"
	"errors"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
)

// TagService 标签服务
type TagService struct {
	repo        *repository.TagRepository
	articleRepo *repository.ArticleRepository
}

// NewTagService 创建标签服务实例
func NewTagService(repo *repository.TagRepository, articleRepo *repository.ArticleRepository) *TagService {
	return &TagService{
		repo:        repo,
		articleRepo: articleRepo,
	}
}

// ============ 前台服务 ============

// ListForWeb 获取前台标签列表
func (s *TagService) ListForWeb(ctx context.Context, page, pageSize int) ([]dto.TagForWebResponse, int64, error) {
	tags, _, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.TagForWebResponse, 0)
	for _, tag := range tags {
		if count, _ := s.articleRepo.CountByTag(tag.ID, true); count > 0 {
			result = append(result, dto.TagForWebResponse{
				ID:          tag.ID,
				Name:        tag.Name,
				Slug:        tag.Slug,
				URL:         fmt.Sprintf("/tag/%s", tag.Slug),
				Description: tag.Description,
				Count:       int(count),
			})
		}
	}
	return result, int64(len(result)), nil
}

// GetBySlug 根据slug获取标签
func (s *TagService) GetBySlug(ctx context.Context, slug string) (*dto.TagForWebResponse, error) {
	tag, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	count, err := s.articleRepo.CountByTag(tag.ID, true)
	if err != nil || count == 0 {
		return nil, errors.New("标签不存在")
	}

	return &dto.TagForWebResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		URL:         fmt.Sprintf("/tag/%s", tag.Slug),
		Description: tag.Description,
		Count:       int(count),
	}, nil
}

// ============ 后台管理服务 ============

// List 获取标签列表
func (s *TagService) List(ctx context.Context, page, pageSize int) ([]dto.TagResponse, int64, error) {
	tags, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		count, _ := s.articleRepo.CountByTag(tag.ID, false)
		result = append(result, dto.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			Count:       int(count),
			CreatedAt:   tag.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, total, nil
}

// Create 创建标签
func (s *TagService) Create(ctx context.Context, tag *model.Tag) error {
	// 使用name作为slug
	tag.Slug = tag.Name

	// 检查slug是否已存在
	if _, err := s.repo.GetBySlug(ctx, tag.Slug); err == nil {
		return errors.New("标签已存在")
	}

	return s.repo.Create(ctx, tag)
}

// Update 更新标签
func (s *TagService) Update(ctx context.Context, id uint, tag *model.Tag) error {
	tag.ID = id
	if tag.Name != "" {
		tag.Slug = tag.Name
		// 检查名称冲突
		if existingTag, err := s.repo.GetBySlug(ctx, tag.Slug); err == nil && existingTag.ID != id {
			return errors.New("标签名称已存在")
		}
	}

	return s.repo.Update(ctx, tag)
}

// Delete 删除标签
func (s *TagService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
