package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/random"
	"flec_blog/pkg/utils"
	"flec_blog/pkg/wechatmp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/gorm"
)

// ArticleService 文章服务
type ArticleService struct {
	articleRepo       *repository.ArticleRepository
	tagRepo           *repository.TagRepository
	categoryRepo      *repository.CategoryRepository
	commentRepo       *repository.CommentRepository
	fileService       *FileService
	subscriberService *SubscriberService
	db                *gorm.DB
	md                goldmark.Markdown
	httpClient        *http.Client
}

// NewArticleService 创建文章服务实例
func NewArticleService(articleRepo *repository.ArticleRepository, tagRepo *repository.TagRepository, categoryRepo *repository.CategoryRepository, commentRepo *repository.CommentRepository, fileService *FileService, db *gorm.DB) *ArticleService {
	// 初始化 goldmark（用于微信导出时渲染 Markdown）
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return &ArticleService{
		articleRepo:  articleRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
		commentRepo:  commentRepo,
		fileService:  fileService,
		db:           db,
		md:           md,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// SetSubscriberService 设置订阅者服务（避免循环依赖）
func (s *ArticleService) SetSubscriberService(subscriberService *SubscriberService) {
	s.subscriberService = subscriberService
}

// ============ 前台服务 ============

// ListForWeb 获取前台文章列表
func (s *ArticleService) ListForWeb(ctx context.Context, req *dto.ListArticlesRequest) ([]dto.ArticleWebResponse, int64, error) {
	articles, total, err := s.articleRepo.ListForWeb(req.Page, req.PageSize, req.Year, req.Month, req.Category, req.Tag)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取文章评论数
	articleSlugs := make([]string, len(articles))
	for i, article := range articles {
		articleSlugs[i] = article.Slug
	}

	commentCounts := make(map[string]int64)
	if len(articleSlugs) > 0 && s.commentRepo != nil {
		commentCounts, err = s.commentRepo.CountByTargetKeys(ctx, "article", articleSlugs)
		if err != nil {
			// 如果获取评论数失败，不影响主流程，只记录错误
			commentCounts = make(map[string]int64)
		}
	}

	// 转换为前台响应格式
	var response []dto.ArticleWebResponse
	for _, article := range articles {
		item := dto.ArticleWebResponse{
			ID:           article.ID,
			Title:        article.Title,
			Summary:      article.Summary,
			Cover:        article.Cover,
			Location:     article.Location,
			IsTop:        article.IsTop,
			IsEssence:    article.IsEssence,
			IsOutdated:   article.IsOutdated,
			URL:          fmt.Sprintf("/posts/%s", article.Slug),
			CommentCount: commentCounts[article.Slug],
			PublishTime:  utils.ToJSONTime(article.PublishTime),
			UpdateTime:   utils.ToJSONTime(article.UpdateTime),
		}

		// 填充分类信息
		if article.Category.ID > 0 {
			item.Category.ID = article.Category.ID
			item.Category.Name = article.Category.Name
			item.Category.URL = fmt.Sprintf("/category/%s", article.Category.Slug)
		}

		// 填充标签信息
		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				ID:   tag.ID,
				Name: tag.Name,
				URL:  fmt.Sprintf("/tag/%s", tag.Slug),
			})
		}

		response = append(response, item)
	}

	return response, total, nil
}

// Search 搜索文章
func (s *ArticleService) Search(ctx context.Context, req *dto.SearchArticlesRequest) ([]dto.ArticleWebResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.articleRepo.Search(req.Keyword, offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取文章评论数
	articleSlugs := make([]string, len(articles))
	for i, article := range articles {
		articleSlugs[i] = article.Slug
	}

	commentCounts := make(map[string]int64)
	if len(articleSlugs) > 0 && s.commentRepo != nil {
		commentCounts, err = s.commentRepo.CountByTargetKeys(ctx, "article", articleSlugs)
		if err != nil {
			// 如果获取评论数失败，不影响主流程，只记录错误
			commentCounts = make(map[string]int64)
		}
	}

	var response []dto.ArticleWebResponse
	for _, article := range articles {
		item := dto.ArticleWebResponse{
			ID:           article.ID,
			Title:        article.Title,
			Summary:      article.Summary,
			Cover:        article.Cover,
			Location:     article.Location,
			IsTop:        article.IsTop,
			IsEssence:    article.IsEssence,
			URL:          fmt.Sprintf("/posts/%s", article.Slug),
			Excerpt:      utils.GenerateExcerpt(article.Content, req.Keyword, 40), // 生成包含关键词的摘录
			CommentCount: commentCounts[article.Slug],
			PublishTime:  utils.ToJSONTime(article.PublishTime),
			UpdateTime:   utils.ToJSONTime(article.UpdateTime),
		}

		if article.Category.ID > 0 {
			item.Category.ID = article.Category.ID
			item.Category.Name = article.Category.Name
			item.Category.URL = fmt.Sprintf("/category/%s", article.Category.Slug)
		}

		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				ID:   tag.ID,
				Name: tag.Name,
				URL:  fmt.Sprintf("/tag/%s", tag.Slug),
			})
		}

		response = append(response, item)
	}

	return response, total, nil
}

// GetBySlug 通过slug获取文章详情
func (s *ArticleService) GetBySlug(ctx context.Context, slug string) (*dto.ArticleDetailResponse, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// 异步增加浏览数
	go func() {
		_ = s.articleRepo.IncrementViewCount(article.ID)
	}()

	// 获取文章评论数
	var commentCount int64
	if s.commentRepo != nil {
		commentCounts, err := s.commentRepo.CountByTargetKeys(ctx, "article", []string{article.Slug})
		if err == nil {
			commentCount = commentCounts[article.Slug]
		}
	}

	response := &dto.ArticleDetailResponse{
		ID:           article.ID,
		Title:        article.Title,
		Slug:         article.Slug,
		Content:      article.Content,
		Summary:      article.Summary,
		AISummary:    article.AISummary,
		Cover:        article.Cover,
		Location:     article.Location,
		IsTop:        article.IsTop,
		IsEssence:    article.IsEssence,
		IsOutdated:   article.IsOutdated,
		ViewCount:    article.ViewCount,
		CommentCount: commentCount,
		URL:          fmt.Sprintf("/posts/%s", article.Slug),
		PublishTime:  utils.ToJSONTime(article.PublishTime),
		UpdateTime:   utils.ToJSONTime(article.UpdateTime),
	}

	// 填充分类信息
	if article.Category.ID > 0 {
		response.Category.ID = article.Category.ID
		response.Category.Name = article.Category.Name
		response.Category.URL = fmt.Sprintf("/category/%s", article.Category.Slug)
	}

	// 填充标签信息
	for _, tag := range article.Tags {
		response.Tags = append(response.Tags, struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			ID:   tag.ID,
			Name: tag.Name,
			URL:  fmt.Sprintf("/tag/%s", tag.Slug),
		})
	}

	// 查询上一篇文章
	if prevArticle, err := s.articleRepo.GetPrevArticle(article.PublishTime); err == nil {
		response.Prev = &struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}{
			Title: prevArticle.Title,
			URL:   fmt.Sprintf("/posts/%s", prevArticle.Slug),
		}
	}

	// 查询下一篇文章
	if nextArticle, err := s.articleRepo.GetNextArticle(article.PublishTime); err == nil {
		response.Next = &struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}{
			Title: nextArticle.Title,
			URL:   fmt.Sprintf("/posts/%s", nextArticle.Slug),
		}
	}

	return response, nil
}

// ============ 后台管理服务 ============

// List 获取文章列表
func (s *ArticleService) List(ctx context.Context, req *dto.ListArticlesRequest) ([]dto.ArticleListResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.articleRepo.List(offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取文章评论数
	articleSlugs := make([]string, len(articles))
	for i, article := range articles {
		articleSlugs[i] = article.Slug
	}

	commentCounts := make(map[string]int64)
	if len(articleSlugs) > 0 && s.commentRepo != nil {
		commentCounts, err = s.commentRepo.CountByTargetKeys(ctx, "article", articleSlugs)
		if err != nil {
			// 如果获取评论数失败，不影响主流程
			commentCounts = make(map[string]int64)
		}
	}

	// 转换为后台列表响应格式
	var response []dto.ArticleListResponse
	for _, article := range articles {
		item := dto.ArticleListResponse{
			ID:           article.ID,
			Title:        article.Title,
			Cover:        article.Cover,
			Location:     article.Location,
			IsPublish:    article.IsPublish,
			IsTop:        article.IsTop,
			IsEssence:    article.IsEssence,
			IsOutdated:   article.IsOutdated,
			ViewCount:    article.ViewCount,
			CommentCount: commentCounts[article.Slug],
			PublishTime:  utils.ToJSONTime(article.PublishTime),
			UpdateTime:   utils.ToJSONTime(article.UpdateTime),
		}

		item.Category.ID = article.Category.ID
		item.Category.Name = article.Category.Name

		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{tag.ID, tag.Name})
		}

		response = append(response, item)
	}

	return response, total, nil
}

// Get 获取文章详情
func (s *ArticleService) Get(_ context.Context, id uint) (*dto.ArticleAdminDetailResponse, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("获取文章失败: %w", err)
	}

	response := &dto.ArticleAdminDetailResponse{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		Summary:     article.Summary,
		AISummary:   article.AISummary,
		Cover:       article.Cover,
		Location:    article.Location,
		IsPublish:   article.IsPublish,
		IsTop:       article.IsTop,
		IsEssence:   article.IsEssence,
		IsOutdated:  article.IsOutdated,
		PublishTime: utils.ToJSONTime(article.PublishTime),
		UpdateTime:  utils.ToJSONTime(article.UpdateTime),
	}

	// 填充分类信息
	response.Category.ID = article.Category.ID
	response.Category.Name = article.Category.Name

	// 填充标签信息
	for _, tag := range article.Tags {
		response.Tags = append(response.Tags, struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{tag.ID, tag.Name})
	}

	return response, nil
}

// Create 创建文章
func (s *ArticleService) Create(ctx context.Context, req *dto.CreateArticleRequest) (*dto.ArticleAdminDetailResponse, error) {
	// 验证分类是否存在
	if req.CategoryID != nil && *req.CategoryID > 0 {
		_, err := s.categoryRepo.Get(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("分类不存在: %w", err)
		}
	}

	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		Location:   req.Location,
		CategoryID: req.CategoryID,
	}

	// 设置置顶状态
	if req.IsTop != nil {
		article.IsTop = *req.IsTop
	}

	// 设置精选状态
	if req.IsEssence != nil {
		article.IsEssence = *req.IsEssence
	}

	// 设置过时状态
	if req.IsOutdated != nil {
		article.IsOutdated = *req.IsOutdated
	}

	// 设置发布状态
	if req.IsPublish != nil {
		article.IsPublish = *req.IsPublish
	}

	// 如果是发布状态，自动设置发布时间
	if article.IsPublish {
		now := utils.Now().Time
		article.PublishTime = &now
	}

	// 生成唯一slug
	generatedSlug, err := random.UniqueCode(8, s.articleRepo.CheckSlugExists)
	if err != nil {
		return nil, fmt.Errorf("生成 slug 失败: %w", err)
	}
	article.Slug = generatedSlug

	// 创建文章并关联标签
	if err := s.articleRepo.Create(article, req.TagIDs); err != nil {
		return nil, err
	}

	// 标记封面为使用中
	if req.Cover != "" && s.fileService != nil {
		_ = s.fileService.MarkAsUsed(req.Cover)
	}

	// 标记内容中的图片为使用中
	s.markContentImagesAsUsed(req.Content)

	// 如果是发布状态，异步发送订阅推送
	if article.IsPublish && s.subscriberService != nil {
		go func() {
			if err := s.subscriberService.SendArticleNotification(context.Background(), article); err != nil {
				logger.Warn("发送文章推送失败 (文章ID: %d): %v", article.ID, err)
			}
		}()
	}

	return s.Get(ctx, article.ID)
}

// Update 更新文章
func (s *ArticleService) Update(ctx context.Context, id uint, req *dto.UpdateArticleRequest) (*dto.ArticleAdminDetailResponse, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, err
	}

	// 验证新分类是否存在
	if req.CategoryID != nil && *req.CategoryID > 0 {
		if _, err := s.categoryRepo.Get(ctx, *req.CategoryID); err != nil {
			return nil, fmt.Errorf("分类不存在: %w", err)
		}
	}

	// 保存旧值用于后续处理
	oldCover := article.Cover
	oldContent := article.Content
	oldIsPublish := article.IsPublish

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	article.Summary = req.Summary
	article.AISummary = req.AISummary
	article.Cover = req.Cover
	article.Location = req.Location
	article.CategoryID = req.CategoryID
	if req.IsTop != nil {
		article.IsTop = *req.IsTop
	}

	// 处理精选状态
	if req.IsEssence != nil {
		article.IsEssence = *req.IsEssence
	}

	// 处理过时状态
	if req.IsOutdated != nil {
		article.IsOutdated = *req.IsOutdated
	}

	// 处理发布状态
	if req.IsPublish != nil {
		article.IsPublish = *req.IsPublish
	}

	// 先处理请求中的 PublishTime（仅当传入非空时间时才更新）
	if req.PublishTime != nil && !req.PublishTime.IsZero() {
		article.PublishTime = utils.FromJSONTime(req.PublishTime)
	}

	// 如果是发布状态且没有发布时间，自动设置发布时间
	if article.IsPublish && article.PublishTime == nil {
		now := utils.Now().Time
		article.PublishTime = &now
	}
	if req.UpdateTime != nil {
		article.UpdateTime = utils.FromJSONTime(req.UpdateTime)
	}

	if err := s.articleRepo.Update(article, req.TagIDs); err != nil {
		return nil, err
	}

	// 处理封面变化
	if s.fileService != nil && oldCover != req.Cover {
		if oldCover != "" {
			_ = s.fileService.MarkAsUnused(oldCover)
		}
		if req.Cover != "" {
			_ = s.fileService.MarkAsUsed(req.Cover)
		}
	}

	// 处理内容图片变化
	if req.Content != "" {
		s.updateContentFileStatus(oldContent, req.Content)
	}

	// 如果从草稿变为发布状态，异步发送订阅推送
	if !oldIsPublish && article.IsPublish && s.subscriberService != nil {
		go func() {
			if err := s.subscriberService.SendArticleNotification(context.Background(), article); err != nil {
				logger.Warn("发送文章推送失败 (文章ID: %d): %v", article.ID, err)
			}
		}()
	}

	return s.Get(ctx, id)
}

// Delete 删除文章
func (s *ArticleService) Delete(ctx context.Context, id uint) error {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return err
	}

	// 标记封面为未使用
	if s.fileService != nil && article.Cover != "" {
		_ = s.fileService.MarkAsUnused(article.Cover)
	}

	// 标记内容中的图片为未使用
	s.markContentImagesAsUnused(article.Content)

	return s.articleRepo.Delete(id)
}

// ============ 辅助方法 ============

// extractContentImages 从 Markdown/HTML 内容中提取所有图片 URL
func extractContentImages(content string) []string {
	var urls []string
	seen := make(map[string]bool)

	// 提取 Markdown 图片: ![alt](url)
	mdImageRe := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	matches := mdImageRe.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			if url != "" && !seen[url] {
				seen[url] = true
				urls = append(urls, url)
			}
		}
	}

	// 提取 HTML img 标签: <img src="url" />
	htmlImageRe := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches = htmlImageRe.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			if url != "" && !seen[url] {
				seen[url] = true
				urls = append(urls, url)
			}
		}
	}

	return urls
}

// markContentImagesAsUsed 标记内容中的图片为已使用
func (s *ArticleService) markContentImagesAsUsed(content string) {
	if s.fileService == nil {
		return
	}
	for _, url := range extractContentImages(content) {
		_ = s.fileService.MarkAsUsed(url)
	}
}

// markContentImagesAsUnused 标记内容中的图片为未使用
func (s *ArticleService) markContentImagesAsUnused(content string) {
	if s.fileService == nil {
		return
	}
	for _, url := range extractContentImages(content) {
		_ = s.fileService.MarkAsUnused(url)
	}
}

// updateContentFileStatus 对比新旧内容，更新图片文件状态
func (s *ArticleService) updateContentFileStatus(oldContent, newContent string) {
	if s.fileService == nil {
		return
	}

	oldImages := make(map[string]bool)
	for _, url := range extractContentImages(oldContent) {
		oldImages[url] = true
	}

	newImages := make(map[string]bool)
	for _, url := range extractContentImages(newContent) {
		newImages[url] = true
		// 新增的图片标记为使用中
		if !oldImages[url] {
			_ = s.fileService.MarkAsUsed(url)
		}
	}

	// 移除的图片标记为未使用
	for url := range oldImages {
		if !newImages[url] {
			_ = s.fileService.MarkAsUnused(url)
		}
	}
}

// ============ 文章导入方法 ============

// ImportArticles 导入文章数据
func (s *ArticleService) ImportArticles(ctx context.Context, files map[string]string, sourceType string, uploadImages bool, host string) (*dto.ImportArticlesResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("没有找到有效的文章数据")
	}

	result := &dto.ImportArticlesResult{Total: len(files)}
	categoryCache := make(map[string]*model.Category)
	tagCache := make(map[string]*model.Tag)

	for filename, content := range files {
		parsed, err := s.parseAndUploadImages(ctx, filename, content, sourceType, uploadImages, host)
		if err != nil {
			title := "未知标题"
			for _, line := range strings.Split(content, "\n") {
				if line = strings.TrimSpace(line); strings.HasPrefix(line, "title:") {
					title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
					break
				}
			}
			result.AddError(filename, title, err.Error())
			continue
		}

		categoryID, err := s.resolveCategory(ctx, parsed.Category, categoryCache)
		if err != nil {
			result.AddError(filename, parsed.Title, fmt.Sprintf("分类处理失败: %v", err))
			continue
		}

		tagIDs, err := s.resolveTags(ctx, parsed.Tags, tagCache)
		if err != nil {
			result.AddError(filename, parsed.Title, fmt.Sprintf("标签处理失败: %v", err))
			continue
		}

		if err := s.createArticle(parsed, categoryID, tagIDs); err != nil {
			result.AddError(filename, parsed.Title, fmt.Sprintf("保存失败: %v", err))
		} else {
			result.Success++
		}
	}

	result.CategoriesAdded = len(categoryCache)
	result.TagsAdded = len(tagCache)
	return result, nil
}

// parseAndUploadImages 解析文章并上传图片
func (s *ArticleService) parseAndUploadImages(ctx context.Context, filename, content, sourceType string, uploadImages bool, host string) (*ParsedArticle, error) {
	var parsed *ParsedArticle
	var err error

	if sourceType == "markdown" {
		parsed, err = parseMarkdownArticle(filename, content)
	} else {
		parsed, err = parseHexoArticle(content)
	}
	if err != nil {
		return nil, err
	}

	if uploadImages {
		if newContent, err := s.uploadContentImages(ctx, parsed.Content, host); err == nil {
			parsed.Content = newContent
		}
		if parsed.Cover != "" {
			if newCover, err := s.uploadSingleImage(ctx, parsed.Cover, host); err == nil {
				parsed.Cover = newCover
			}
		}
	}
	return parsed, nil
}

// resolveCategory 解析分类ID
func (s *ArticleService) resolveCategory(ctx context.Context, name string, cache map[string]*model.Category) (*uint, error) {
	if name == "" {
		return nil, nil
	}
	if c, ok := cache[name]; ok {
		return &c.ID, nil
	}
	c, err := s.categoryRepo.GetBySlug(ctx, name)
	if err != nil {
		c = &model.Category{Name: name, Slug: name}
		if err := s.categoryRepo.Create(ctx, c); err != nil {
			return nil, err
		}
	}
	cache[name] = c
	return &c.ID, nil
}

// resolveTags 解析标签ID列表
func (s *ArticleService) resolveTags(ctx context.Context, names []string, cache map[string]*model.Tag) ([]uint, error) {
	var ids []uint
	for _, name := range names {
		if t, ok := cache[name]; ok {
			ids = append(ids, t.ID)
			continue
		}
		t, err := s.tagRepo.GetBySlug(ctx, name)
		if err != nil {
			t = &model.Tag{Name: name, Slug: name}
			if err := s.tagRepo.Create(ctx, t); err != nil {
				return nil, err
			}
		}
		cache[name] = t
		ids = append(ids, t.ID)
	}
	return ids, nil
}

// createArticle 创建文章记录
func (s *ArticleService) createArticle(parsed *ParsedArticle, categoryID *uint, tagIDs []uint) error {
	slug := parsed.Slug
	if slug != "" {
		if exists, _ := s.articleRepo.CheckSlugExists(slug); exists {
			slug = ""
		}
	}
	if slug == "" {
		slug, _ = random.UniqueCode(8, s.articleRepo.CheckSlugExists)
	}

	article := &model.Article{
		Title:       parsed.Title,
		Slug:        slug,
		Content:     parsed.Content,
		Summary:     parsed.Summary,
		Cover:       parsed.Cover,
		IsPublish:   false,
		CategoryID:  categoryID,
		PublishTime: parsed.PublishTime,
		UpdateTime:  parsed.UpdateTime,
	}
	return s.articleRepo.Create(article, tagIDs)
}

// uploadContentImages 上传文章内容中的所有图片，返回替换后的内容
func (s *ArticleService) uploadContentImages(ctx context.Context, content string, host string) (string, error) {
	if s.fileService == nil {
		return content, nil
	}

	imageURLs := extractContentImages(content)
	if len(imageURLs) == 0 {
		return content, nil
	}

	uniqueURLs := make(map[string]bool)
	for _, url := range imageURLs {
		uniqueURLs[url] = true
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	replacements := make(map[string]string)
	sem := make(chan struct{}, 10)

	for url := range uniqueURLs {
		if strings.HasPrefix(url, "./") || strings.HasPrefix(url, "../") || strings.HasPrefix(url, "/") {
			continue
		}

		wg.Add(1)
		go func(imgURL string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if newURL, err := s.uploadSingleImage(ctx, imgURL, host); err == nil {
				mu.Lock()
				replacements[imgURL] = newURL
				mu.Unlock()
			}
		}(url)
	}
	wg.Wait()

	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}
	return content, nil
}

// uploadSingleImage 上传单张图片，返回新的URL
func (s *ArticleService) uploadSingleImage(ctx context.Context, imgURL string, host string) (string, error) {
	if s.fileService == nil || imgURL == "" {
		return imgURL, nil
	}

	if strings.HasPrefix(imgURL, "./") || strings.HasPrefix(imgURL, "../") || strings.HasPrefix(imgURL, "/") {
		return imgURL, nil
	}

	data, ext, err := s.fetchImage(ctx, imgURL)
	if err != nil {
		return imgURL, fmt.Errorf("下载图片失败: %w", err)
	}

	hashBytes := sha256.Sum256(data)
	hashStr := fmt.Sprintf("%x", hashBytes)[:12]
	filename := fmt.Sprintf("import_%s%s", hashStr, ext)

	mimeType := "image/jpeg"
	switch strings.ToLower(ext) {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	case ".avif":
		mimeType = "image/avif"
	case ".svg":
		mimeType = "image/svg+xml"
	case ".bmp":
		mimeType = "image/bmp"
	case ".tiff", ".tif":
		mimeType = "image/tiff"
	}

	reader := bytes.NewReader(data)
	uploadedURL, err := s.fileService.UploadFromReader(reader, filename, mimeType, "文章图片", 0, host)
	if err != nil {
		return imgURL, fmt.Errorf("上传图片失败: %w", err)
	}

	if err := s.fileService.MarkAsUsed(uploadedURL); err != nil {
		logger.Warn("标记文件状态失败: %v", err)
	}

	return uploadedURL, nil
}

// ParsedArticle 解析后的文章数据
type ParsedArticle struct {
	Title       string
	Slug        string
	Content     string
	Summary     string
	Cover       string
	Category    string
	Tags        []string
	PublishTime *time.Time
	UpdateTime  *time.Time
}

// generateSummary 从内容生成摘要
func generateSummary(content string, maxLen int) string {
	content = strings.NewReplacer("#", "", "*", "", "`", "", "\n", " ").Replace(content)
	content = strings.TrimSpace(content)

	runes := []rune(content)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return content
}

// parseHexoArticle 解析Hexo文章格式（Front Matter + Markdown）
func parseHexoArticle(content string) (*ParsedArticle, error) {
	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("无效的Hexo格式：缺少Front Matter")
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("无效的Hexo格式：Front Matter格式错误")
	}

	frontMatter := parts[1]
	markdown := strings.TrimSpace(parts[2])

	parsed := &ParsedArticle{
		Content: markdown,
	}

	// 日期格式列表
	dateFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	// 解析日期的辅助函数
	parseDate := func(dateStr string) *time.Time {
		for _, format := range dateFormats {
			if t, err := time.Parse(format, dateStr); err == nil {
				return &t
			}
		}
		return nil
	}

	lines := strings.Split(frontMatter, "\n")
	var tagLines []string
	inTags := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if inTags {
			if strings.HasPrefix(line, "-") {
				tagValue := strings.TrimSpace(strings.TrimPrefix(line, "-"))
				tagValue = strings.Trim(tagValue, "\"'")
				if tagValue != "" {
					tagLines = append(tagLines, tagValue)
				}
			} else {
				inTags = false
			}
		}

		if strings.Contains(line, ":") && !strings.HasPrefix(line, "-") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"'")
			}

			switch key {
			case "title":
				parsed.Title = value
			case "date":
				parsed.PublishTime = parseDate(value)
			case "updated":
				parsed.UpdateTime = parseDate(value)
			case "categories", "category":
				if value != "" {
					parsed.Category = value
				}
			case "tags":
				if value != "" {
					value = strings.Trim(value, "[]")
					for _, tag := range strings.Split(value, ",") {
						tag = strings.TrimSpace(tag)
						tag = strings.Trim(tag, "\"'")
						if tag != "" {
							parsed.Tags = append(parsed.Tags, tag)
						}
					}
				} else {
					inTags = true
				}
			case "cover", "thumbnail":
				parsed.Cover = value
			case "description", "excerpt":
				parsed.Summary = value
			case "slug", "abbrlink":
				parsed.Slug = value
			}
		}
	}

	if len(tagLines) > 0 {
		parsed.Tags = append(parsed.Tags, tagLines...)
	}

	if parsed.Title == "" {
		return nil, fmt.Errorf("文章缺少标题")
	}

	if parsed.Summary == "" {
		parsed.Summary = generateSummary(parsed.Content, 200)
	}

	return parsed, nil
}

// parseMarkdownArticle 解析Markdown格式文章
func parseMarkdownArticle(filename, content string) (*ParsedArticle, error) {
	parsed := &ParsedArticle{
		Tags:        []string{},
		PublishTime: nil,
		UpdateTime:  nil,
	}

	if filename != "" {
		lowerName := strings.ToLower(filename)
		if strings.HasSuffix(lowerName, ".md") {
			parsed.Title = strings.TrimSpace(filename[:len(filename)-3])
		} else {
			parsed.Title = strings.TrimSpace(filename)
		}
	}

	if parsed.Title == "" {
		parsed.Title = "未命名文章"
	}

	parsed.Summary = generateSummary(content, 200)
	parsed.Content = content

	return parsed, nil
}

// ============ 微信公众号导出 ============

// ExportToWeChat 将文章渲染为微信公众号 HTML 格式
func (s *ArticleService) ExportToWeChat(_ context.Context, id uint) *dto.WeChatExportResult {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return &dto.WeChatExportResult{}
	}

	// 预处理并渲染 Markdown
	processed := wechatmp.ConvertCustomBlocks(article.Content)
	processed = wechatmp.ConvertLinksToFootnotes(processed)
	processed = wechatmp.PreprocessMarkdown(processed)

	var htmlBuf bytes.Buffer
	if err := s.md.Convert([]byte(processed), &htmlBuf); err != nil {
		return &dto.WeChatExportResult{}
	}

	result, err := wechatmp.ConvertMarkdownToWeChatHTML(htmlBuf.String())
	if err != nil {
		return &dto.WeChatExportResult{}
	}

	return &dto.WeChatExportResult{HTML: result.HTML}
}

// fetchImage 下载图片，返回数据和扩展名
func (s *ArticleService) fetchImage(ctx context.Context, imgURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// 从 URL 或 Content-Type 获取扩展名
	ext := ".jpg"
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		switch ct {
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		case "image/avif":
			ext = ".avif"
		}
	} else if idx := strings.LastIndex(imgURL, "."); idx > 0 {
		if e := imgURL[idx:]; len(e) <= 5 {
			ext = e
		}
	}

	return data, ext, nil
}

// ============ 文章下载导出 ============

// imageDownloadResult 图片下载结果
type imageDownloadResult struct {
	url      string
	data     []byte
	ext      string
	filename string
	err      error
}

// extractFilenameFromURL 从 URL 中提取文件名并清理非法字符
func extractFilenameFromURL(imgURL string) string {
	// 移除查询参数
	if idx := strings.Index(imgURL, "?"); idx > 0 {
		imgURL = imgURL[:idx]
	}
	// 提取路径最后一部分
	var filename string
	if idx := strings.LastIndex(imgURL, "/"); idx >= 0 && idx < len(imgURL)-1 {
		filename = imgURL[idx+1:]
	}
	if filename == "" {
		return ""
	}
	// 清理文件名中的非法字符
	filename = strings.Map(func(r rune) rune {
		if strings.ContainsRune("<>:\"/\\|?*", r) {
			return '_'
		}
		return r
	}, filename)
	return filename
}

// DownloadZip 下载文章为压缩包
func (s *ArticleService) DownloadZip(ctx context.Context, id uint) ([]byte, string, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	imageMap := make(map[string]string)

	// 收集所有需要下载的图片 URL（封面 + 内容图片）
	var imageURLs []string
	if article.Cover != "" {
		imageURLs = append(imageURLs, article.Cover)
	}
	imageURLs = append(imageURLs, s.extractImageURLs(article.Content)...)

	// 去重
	seen := make(map[string]bool)
	var uniqueURLs []string
	for _, url := range imageURLs {
		if !seen[url] {
			seen[url] = true
			uniqueURLs = append(uniqueURLs, url)
		}
	}

	// 如果没有图片，直接生成 Markdown 文件
	if len(uniqueURLs) == 0 {
		frontMatter := s.buildYAMLFrontMatter(article, imageMap)
		mdContent := frontMatter + "\n" + article.Content
		filename := s.sanitizeFilename(article.Title) + ".md"
		if w, _ := zipWriter.Create(filename); w != nil {
			w.Write([]byte(mdContent))
		}
		zipWriter.Close()
		return buf.Bytes(), s.sanitizeFilename(article.Title) + ".zip", nil
	}

	// 并发下载图片（限制并发数为 10）
	const maxConcurrency = 10
	results := make(chan imageDownloadResult, len(uniqueURLs))
	sem := make(chan struct{}, maxConcurrency)

	// 预先为每个 URL 分配文件名（避免并发竞态）
	filenameMap := make(map[string]string)
	filenameCounter := make(map[string]int)
	for _, url := range uniqueURLs {
		// 从 URL 提取原始文件名
		originalName := extractFilenameFromURL(url)
		if originalName == "" {
			// 从 fetchImage 获取扩展名（这里先使用默认）
			originalName = "image.jpg"
		}

		// 处理文件名冲突
		finalName := "assets/" + originalName
		if count, exists := filenameCounter[originalName]; exists {
			// 文件名冲突，添加序号
			nameWithoutExt := originalName
			ext := ""
			if idx := strings.LastIndex(originalName, "."); idx > 0 {
				nameWithoutExt = originalName[:idx]
				ext = originalName[idx:]
			}
			finalName = fmt.Sprintf("assets/%s_%d%s", nameWithoutExt, count+1, ext)
			filenameCounter[originalName] = count + 1
		} else {
			filenameCounter[originalName] = 1
		}

		// 封面图特殊处理
		if url == article.Cover {
			finalName = "assets/cover.jpg" // 默认扩展名，后续会根据实际类型调整
		}

		filenameMap[url] = finalName
	}

	// 并发下载
	for _, url := range uniqueURLs {
		go func(imgURL string) {
			sem <- struct{}{}
			defer func() { <-sem }()

			result := imageDownloadResult{url: imgURL}
			if data, ext, err := s.fetchImage(ctx, imgURL); err == nil {
				result.data = data
				result.ext = ext

				// 获取预分配的文件名，并根据实际扩展名调整
				filename := filenameMap[imgURL]
				// 替换扩展名
				if idx := strings.LastIndex(filename, "."); idx > 0 {
					filename = filename[:idx] + ext
				}
				result.filename = filename
			} else {
				result.err = err
			}
			results <- result
		}(url)
	}

	// 收集结果并写入 zip
	for range uniqueURLs {
		result := <-results
		if result.err != nil {
			continue
		}
		if w, _ := zipWriter.Create(result.filename); w != nil {
			w.Write(result.data)
			imageMap[result.url] = result.filename
		}
	}

	// 替换图片链接
	content := article.Content
	for url, path := range imageMap {
		content = strings.ReplaceAll(content, url, path)
	}

	// 写入 Markdown 文件
	frontMatter := s.buildYAMLFrontMatter(article, imageMap)
	mdContent := frontMatter + "\n" + content
	filename := s.sanitizeFilename(article.Title) + ".md"
	if w, _ := zipWriter.Create(filename); w != nil {
		w.Write([]byte(mdContent))
	}

	zipWriter.Close()
	return buf.Bytes(), s.sanitizeFilename(article.Title) + ".zip", nil
}

// buildYAMLFrontMatter 构建 YAML Front Matter
func (s *ArticleService) buildYAMLFrontMatter(article *model.Article, imageMap map[string]string) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %q\n", article.Title)
	fmt.Fprintf(&b, "slug: %s\n", article.Slug)

	if article.Summary != "" {
		fmt.Fprintf(&b, "summary: %q\n", article.Summary)
	}
	if article.Cover != "" {
		if path, ok := imageMap[article.Cover]; ok {
			fmt.Fprintf(&b, "cover: %s\n", path)
		} else {
			fmt.Fprintf(&b, "cover: %s\n", article.Cover)
		}
	}
	if article.Location != "" {
		fmt.Fprintf(&b, "location: %q\n", article.Location)
	}

	fmt.Fprintf(&b, "published: %t\n", article.IsPublish)
	fmt.Fprintf(&b, "top: %t\n", article.IsTop)
	fmt.Fprintf(&b, "essence: %t\n", article.IsEssence)
	fmt.Fprintf(&b, "outdated: %t\n", article.IsOutdated)

	if article.Category.ID > 0 {
		fmt.Fprintf(&b, "category: %q\n", article.Category.Name)
	}
	if len(article.Tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range article.Tags {
			fmt.Fprintf(&b, "  - %q\n", tag.Name)
		}
	}
	if article.PublishTime != nil {
		fmt.Fprintf(&b, "date: %s\n", article.PublishTime.Format("2006-01-02 15:04:05"))
	}
	if article.UpdateTime != nil {
		fmt.Fprintf(&b, "updated: %s\n", article.UpdateTime.Format("2006-01-02 15:04:05"))
	}

	b.WriteString("---\n")
	return b.String()
}

// extractImageURLs 提取 Markdown 中的图片 URL
func (s *ArticleService) extractImageURLs(content string) []string {
	re := regexp.MustCompile(`!\[.*?\]\((https?://[^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	urls := make([]string, 0, len(matches))
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			urls = append(urls, m[1])
		}
	}
	return urls
}

// sanitizeFilename 清理文件名
func (s *ArticleService) sanitizeFilename(name string) string {
	result := strings.Map(func(r rune) rune {
		if strings.ContainsRune("<>:\"/\\|?*", r) {
			return '_'
		}
		return r
	}, name)

	if len([]rune(result)) > 100 {
		result = string([]rune(result)[:100])
	}
	return result
}
