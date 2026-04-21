package repository

import (
	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// ArticleRepository 文章仓储
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// ============ 查询方法 ============

// ListForWeb 获取文章列表（前台）
func (r *ArticleRepository) ListForWeb(page, pageSize int, year, month, categorySlug, tagSlug string) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{}).Where("is_publish = ?", true)

	// 按年份筛选
	if year != "" {
		query = query.Where("EXTRACT(YEAR FROM publish_time) = ?", year)
	}

	// 按月份筛选
	if month != "" && year != "" {
		query = query.Where("EXTRACT(MONTH FROM publish_time) = ?", month)
	}

	// 按分类筛选
	if categorySlug != "" {
		query = query.Joins("JOIN categories ON categories.id = articles.category_id").
			Where("categories.slug = ?", categorySlug)
	}

	// 按标签筛选
	if tagSlug != "" {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.slug = ?", tagSlug).
			Distinct()
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 只有当 page 和 pageSize 都大于 0 时才应用分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Order("is_top DESC, publish_time DESC, created_at DESC").
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// Search 搜索文章
func (r *ArticleRepository) Search(keyword string, offset, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	searchCondition := "is_publish = ? AND (title ILIKE ? OR content ILIKE ?)"
	searchKeyword := "%" + keyword + "%"

	if err := r.db.Model(&model.Article{}).
		Where(searchCondition, true, searchKeyword, searchKeyword).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where(searchCondition, true, searchKeyword, searchKeyword).
		Order("is_top DESC, created_at DESC").
		Preload("Category").
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetBySlug 通过slug获取文章
func (r *ArticleRepository) GetBySlug(slug string) (*model.Article, error) {
	var article model.Article
	if err := r.db.Preload("Category").
		Preload("Tags").
		Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetPrevArticle 获取上一篇文章（发布时间更早的文章）
func (r *ArticleRepository) GetPrevArticle(publishTime interface{}) (*model.Article, error) {
	var article model.Article
	if err := r.db.Where("is_publish = ? AND publish_time < ?", true, publishTime).
		Order("publish_time DESC").
		First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetNextArticle 获取下一篇文章（发布时间更晚的文章）
func (r *ArticleRepository) GetNextArticle(publishTime interface{}) (*model.Article, error) {
	var article model.Article
	if err := r.db.Where("is_publish = ? AND publish_time > ?", true, publishTime).
		Order("publish_time ASC").
		First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetByTag 根据标签查询文章列表
func (r *ArticleRepository) GetByTag(tagID uint, offset, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	subQuery := r.db.Table("article_tags").
		Select("article_id").
		Where("tag_id = ?", tagID)

	if err := r.db.Model(&model.Article{}).
		Where("id IN (?) AND is_publish = ?", subQuery, true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("id IN (?) AND is_publish = ?", subQuery, true).
		Order("is_top DESC, created_at DESC").
		Preload("Category").
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByCategory 根据分类查询文章列表
func (r *ArticleRepository) GetByCategory(categoryID uint, offset, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	if err := r.db.Model(&model.Article{}).
		Where("category_id = ? AND is_publish = ?", categoryID, true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("category_id = ? AND is_publish = ?", categoryID, true).
		Order("is_top DESC, created_at DESC").
		Preload("Category").
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// CountByCategory 统计分类下的文章数量
func (r *ArticleRepository) CountByCategory(categoryID uint, onlyPublished bool) (int64, error) {
	var count int64
	query := r.db.Model(&model.Article{}).Where("category_id = ?", categoryID)
	if onlyPublished {
		query = query.Where("is_publish = ?", true)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByTag 统计标签下的文章数量
func (r *ArticleRepository) CountByTag(tagID uint, onlyPublished bool) (int64, error) {
	var count int64
	query := r.db.Model(&model.Article{}).
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Where("article_tags.tag_id = ?", tagID)
	if onlyPublished {
		query = query.Where("articles.is_publish = ?", true)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ============ 基础CRUD ============

// List 获取文章列表
func (r *ArticleRepository) List(
	offset, limit int,
	keyword, location string,
	categoryID uint,
	tagIDs []uint,
	isPublish, isTop, isEssence, isOutdated *bool,
	startTime, endTime string,
) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{})

	// 关键词搜索（标题或内容）
	if keyword != "" {
		searchKeyword := "%" + keyword + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchKeyword, searchKeyword)
	}

	// 按发布地点筛选
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	// 按分类筛选
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 按标签筛选
	if len(tagIDs) > 0 {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id IN ?", tagIDs).
			Distinct()
	}

	// 按发布状态筛选
	if isPublish != nil {
		query = query.Where("is_publish = ?", *isPublish)
	}

	// 按置顶状态筛选
	if isTop != nil {
		query = query.Where("is_top = ?", *isTop)
	}

	// 按精选状态筛选
	if isEssence != nil {
		query = query.Where("is_essence = ?", *isEssence)
	}

	// 按过时状态筛选
	if isOutdated != nil {
		query = query.Where("is_outdated = ?", *isOutdated)
	}

	// 按发布时间范围筛选
	if startTime != "" {
		query = query.Where("publish_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("publish_time <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("is_publish ASC, publish_time DESC NULLS LAST, created_at DESC").
		Preload("Category").
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// Get 通过ID获取文章
func (r *ArticleRepository) Get(id uint) (*model.Article, error) {
	var article model.Article
	if err := r.db.Preload("Category").
		Preload("Tags").
		First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// Create 创建文章
func (r *ArticleRepository) Create(article *model.Article, tagIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// 关联标签
		if len(tagIDs) > 0 {
			var tags []model.Tag
			if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})
}

// Update 更新文章
func (r *ArticleRepository) Update(article *model.Article, tagIDs []uint) error {
	// 使用 Updates 而不是 Save，确保只更新不创建
	if err := r.db.Model(&model.Article{}).Where("id = ?", article.ID).
		Select("*").Omit("Category", "Tags", "id", "created_at").
		Updates(article).Error; err != nil {
		return err
	}

	// 更新标签关联
	if tagIDs != nil {
		var tags []model.Tag
		if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
		if err := r.db.Model(article).Association("Tags").Replace(tags); err != nil {
			return err
		}
	}

	return nil
}

// Delete 删除文章
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&model.Article{}, id).Error
}

// ExistsByCover 检查是否有文章使用该封面
func (r *ArticleRepository) ExistsByCover(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Where("cover = ?", url).Count(&count).Error
	return count > 0, err
}

// ExistsByContentURL 检查是否有文章正文引用该文件
func (r *ArticleRepository) ExistsByContentURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Where("content LIKE ?", "%"+url+"%").Count(&count).Error
	return count > 0, err
}

// ============ 辅助方法 ============

// CheckSlugExists 检查slug是否已存在
func (r *ArticleRepository) CheckSlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IncrementViewCount 增加文章浏览数
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
