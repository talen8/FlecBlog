package service

import (
	"context"
	"errors"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
)

// CategoryService 分类服务
type CategoryService struct {
	repo        *repository.CategoryRepository
	articleRepo *repository.ArticleRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(repo *repository.CategoryRepository, articleRepo *repository.ArticleRepository) *CategoryService {
	return &CategoryService{
		repo:        repo,
		articleRepo: articleRepo,
	}
}

// ============ 前台服务 ============

// ListForWeb 获取前台分类列表
func (s *CategoryService) ListForWeb(ctx context.Context, page, pageSize int) ([]dto.CategoryForWebResponse, int64, error) {
	categories, _, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CategoryForWebResponse, 0)
	for _, category := range categories {
		if count, _ := s.articleRepo.CountByCategory(category.ID, true); count > 0 {
			result = append(result, dto.CategoryForWebResponse{
				ID:          category.ID,
				Name:        category.Name,
				Slug:        category.Slug,
				URL:         fmt.Sprintf("/category/%s", category.Slug),
				Description: category.Description,
				Count:       int(count),
				Sort:        category.Sort,
			})
		}
	}
	return result, int64(len(result)), nil
}

// GetBySlug 根据slug获取分类
func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*dto.CategoryForWebResponse, error) {
	category, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	count, err := s.articleRepo.CountByCategory(category.ID, true)
	if err != nil || count == 0 {
		return nil, errors.New("分类不存在")
	}

	return &dto.CategoryForWebResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		URL:         fmt.Sprintf("/category/%s", category.Slug),
		Description: category.Description,
		Count:       int(count),
		Sort:        category.Sort,
	}, nil
}

// ============ 后台管理服务 ============

// List 获取分类列表
func (s *CategoryService) List(ctx context.Context, page, pageSize int) ([]dto.CategoryResponse, int64, error) {
	categories, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		count, _ := s.articleRepo.CountByCategory(category.ID, false)
		result = append(result, dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			Count:       int(count),
			Sort:        category.Sort,
			CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, total, nil
}

// Create 创建分类
func (s *CategoryService) Create(ctx context.Context, category *model.Category) error {
	// 使用name作为slug
	category.Slug = category.Name

	// 检查slug是否已存在
	if _, err := s.repo.GetBySlug(ctx, category.Slug); err == nil {
		return errors.New("分类已存在")
	}

	return s.repo.Create(ctx, category)
}

// Update 更新分类
func (s *CategoryService) Update(ctx context.Context, id uint, category *model.Category) error {
	category.ID = id
	if category.Name != "" {
		category.Slug = category.Name
		// 检查名称冲突
		if existingCategory, err := s.repo.GetBySlug(ctx, category.Slug); err == nil && existingCategory.ID != id {
			return errors.New("分类名称已存在")
		}
	}

	return s.repo.Update(ctx, category)
}

// Delete 删除分类
func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
