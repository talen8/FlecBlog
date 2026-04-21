package tools

import (
	"context"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	taxonomyTargetCategory = "category"
	taxonomyTargetTag      = "tag"

	taxonomyActionList         = "list"
	taxonomyActionCreate       = "create"
	taxonomyActionUpdate       = "update"
	taxonomyActionDelete       = "delete"
	taxonomyActionListArticles = "list_articles"
)

// ============ Taxonomy Manage 输入/输出类型============

// TaxonomyManageInput taxonomy_manage 聚合 tool 输入
type TaxonomyManageInput struct {
	Target  string                `json:"target"` // category|tag
	Action  string                `json:"action"` // list|create|update|delete|list_articles
	Payload TaxonomyManagePayload `json:"payload"`
}

// TaxonomyManagePayload taxonomy_manage 载荷
type TaxonomyManagePayload struct {
	// 用于 list / list_articles
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// 用于 update/delete/list_articles
	ID uint `json:"id"`

	// 用于 create/update
	Name        string `json:"name"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

// TaxonomyManageOutput taxonomy_manage 聚合 tool 输出
type TaxonomyManageOutput struct {
	// list 结果
	List     []TaxonomyItem `json:"list,omitempty"`
	Total    int64          `json:"total,omitempty"`
	Page     int            `json:"page,omitempty"`
	PageSize int            `json:"page_size,omitempty"`

	// list_articles 结果
	Articles *ArticleListResult `json:"articles,omitempty"`

	// create/update 结果
	Item *TaxonomyItem `json:"item,omitempty"`

	// delete 结果
	DeleteSuccess *bool `json:"delete_success,omitempty"`
	ID            *uint `json:"id,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`
}

// TaxonomyItem 分类/标签项
type TaxonomyItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Count       int    `json:"count"`
	Sort        int    `json:"sort"`
}

// ArticleListResult 文章列表结果
type ArticleListResult struct {
	List     []ArticleItem `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// ============ 服务包装器============

// TaxonomyWrapper 分类/标签服务包装器
type TaxonomyWrapper struct {
	categoryService *service.CategoryService
	tagService      *service.TagService
	articleService  *service.ArticleService
}

// NewTaxonomyWrapper 创建分类/标签服务包装器
func NewTaxonomyWrapper(
	categoryService *service.CategoryService,
	tagService *service.TagService,
	articleService *service.ArticleService,
) *TaxonomyWrapper {
	return &TaxonomyWrapper{
		categoryService: categoryService,
		tagService:      tagService,
		articleService:  articleService,
	}
}

// ============ Tool Handler============

// ManageTaxonomy 分类/标签管理聚合入口
func (w *TaxonomyWrapper) ManageTaxonomy(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input TaxonomyManageInput,
) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	switch input.Target {
	case taxonomyTargetCategory:
		return w.handleCategory(input.Action, input.Payload)
	case taxonomyTargetTag:
		return w.handleTag(input.Action, input.Payload)
	default:
		return nil, TaxonomyManageOutput{}, fmt.Errorf("不支持的目标类型: %s", input.Target)
	}
}

// handleCategory 处理分类操作
func (w *TaxonomyWrapper) handleCategory(action string, payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	switch action {
	case taxonomyActionList:
		return w.listCategories(payload)
	case taxonomyActionCreate:
		return w.createCategory(payload)
	case taxonomyActionUpdate:
		return w.updateCategory(payload)
	case taxonomyActionDelete:
		return w.deleteCategory(payload)
	case taxonomyActionListArticles:
		return w.listCategoryArticles(payload)
	default:
		return nil, TaxonomyManageOutput{}, fmt.Errorf("不支持的操作: %s", action)
	}
}

// handleTag 处理标签操作
func (w *TaxonomyWrapper) handleTag(action string, payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	switch action {
	case taxonomyActionList:
		return w.listTags(payload)
	case taxonomyActionCreate:
		return w.createTag(payload)
	case taxonomyActionUpdate:
		return w.updateTag(payload)
	case taxonomyActionDelete:
		return w.deleteTag(payload)
	case taxonomyActionListArticles:
		return w.listTagArticles(payload)
	default:
		return nil, TaxonomyManageOutput{}, fmt.Errorf("不支持的操作: %s", action)
	}
}

// ============ Category Handlers============

func (w *TaxonomyWrapper) listCategories(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	categories, total, err := w.categoryService.List(context.Background(), page, pageSize)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取分类列表失败: %v", err)}, nil
	}

	list := make([]TaxonomyItem, len(categories))
	for i, category := range categories {
		list[i] = convertCategoryToItem(category)
	}

	return nil, TaxonomyManageOutput{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (w *TaxonomyWrapper) createCategory(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.Name == "" {
		return nil, TaxonomyManageOutput{Error: "分类名称不能为空"}, nil
	}

	category := &model.Category{
		Name:        payload.Name,
		Description: payload.Description,
		Sort:        payload.Sort,
	}

	if err := w.categoryService.Create(context.Background(), category); err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("创建分类失败: %v", err)}, nil
	}

	created, err := w.categoryService.Get(context.Background(), category.ID)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取分类失败: %v", err)}, nil
	}

	item := convertCategoryToItem(*created)
	return nil, TaxonomyManageOutput{Item: &item}, nil
}

func (w *TaxonomyWrapper) updateCategory(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.ID == 0 {
		return nil, TaxonomyManageOutput{Error: "分类 ID 不能为空"}, nil
	}
	if payload.Name == "" {
		return nil, TaxonomyManageOutput{Error: "分类名称不能为空"}, nil
	}

	category := &model.Category{
		Name:        payload.Name,
		Description: payload.Description,
		Sort:        payload.Sort,
	}

	if err := w.categoryService.Update(context.Background(), payload.ID, category); err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("更新分类失败: %v", err)}, nil
	}

	updated, err := w.categoryService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取分类失败: %v", err)}, nil
	}

	item := convertCategoryToItem(*updated)
	return nil, TaxonomyManageOutput{Item: &item}, nil
}

func (w *TaxonomyWrapper) deleteCategory(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.ID == 0 {
		return nil, TaxonomyManageOutput{Error: "分类 ID 不能为空"}, nil
	}

	if err := w.categoryService.Delete(context.Background(), payload.ID); err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("删除分类失败: %v", err)}, nil
	}

	success := true
	return nil, TaxonomyManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

func (w *TaxonomyWrapper) listCategoryArticles(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.ID == 0 {
		return nil, TaxonomyManageOutput{Error: "分类 ID 不能为空"}, nil
	}

	category, err := w.categoryService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取分类失败: %v", err)}, nil
	}

	page, pageSize := NormalizePage(payload.Page, payload.PageSize)
	req := &dto.ListArticlesForWebRequest{Page: page, PageSize: pageSize, Category: category.Slug}
	articles, total, err := w.articleService.ListForWeb(context.Background(), req)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取分类文章列表失败: %v", err)}, nil
	}

	list := make([]ArticleItem, len(articles))
	for i, article := range articles {
		list[i] = convertWebArticleToItem(article)
	}

	return nil, TaxonomyManageOutput{
		Articles: &ArticleListResult{List: list, Total: total, Page: page, PageSize: pageSize},
	}, nil
}

// ============ Tag Handlers============

func (w *TaxonomyWrapper) listTags(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	tags, total, err := w.tagService.List(context.Background(), page, pageSize)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取标签列表失败: %v", err)}, nil
	}

	list := make([]TaxonomyItem, len(tags))
	for i, tag := range tags {
		list[i] = convertTagToItem(tag)
	}

	return nil, TaxonomyManageOutput{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (w *TaxonomyWrapper) createTag(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.Name == "" {
		return nil, TaxonomyManageOutput{Error: "标签名称不能为空"}, nil
	}

	tag := &model.Tag{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := w.tagService.Create(context.Background(), tag); err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("创建标签失败: %v", err)}, nil
	}

	created, err := w.tagService.Get(context.Background(), tag.ID)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取标签失败: %v", err)}, nil
	}

	item := convertTagToItem(*created)
	return nil, TaxonomyManageOutput{Item: &item}, nil
}

func (w *TaxonomyWrapper) updateTag(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.ID == 0 {
		return nil, TaxonomyManageOutput{Error: "标签 ID 不能为空"}, nil
	}
	if payload.Name == "" {
		return nil, TaxonomyManageOutput{Error: "标签名称不能为空"}, nil
	}

	tag := &model.Tag{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := w.tagService.Update(context.Background(), payload.ID, tag); err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("更新标签失败: %v", err)}, nil
	}

	updated, err := w.tagService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取标签失败: %v", err)}, nil
	}

	item := convertTagToItem(*updated)
	return nil, TaxonomyManageOutput{Item: &item}, nil
}

func (w *TaxonomyWrapper) deleteTag(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.ID == 0 {
		return nil, TaxonomyManageOutput{Error: "标签 ID 不能为空"}, nil
	}

	if err := w.tagService.Delete(context.Background(), payload.ID); err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("删除标签失败: %v", err)}, nil
	}

	success := true
	return nil, TaxonomyManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

func (w *TaxonomyWrapper) listTagArticles(payload TaxonomyManagePayload) (*sdkmcp.CallToolResult, TaxonomyManageOutput, error) {
	if payload.ID == 0 {
		return nil, TaxonomyManageOutput{Error: "标签 ID 不能为空"}, nil
	}

	tag, err := w.tagService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取标签失败: %v", err)}, nil
	}

	page, pageSize := NormalizePage(payload.Page, payload.PageSize)
	req := &dto.ListArticlesForWebRequest{Page: page, PageSize: pageSize, Tag: tag.Slug}
	articles, total, err := w.articleService.ListForWeb(context.Background(), req)
	if err != nil {
		return nil, TaxonomyManageOutput{Error: fmt.Sprintf("获取标签文章列表失败: %v", err)}, nil
	}

	list := make([]ArticleItem, len(articles))
	for i, article := range articles {
		list[i] = convertWebArticleToItem(article)
	}

	return nil, TaxonomyManageOutput{
		Articles: &ArticleListResult{List: list, Total: total, Page: page, PageSize: pageSize},
	}, nil
}

// TaxonomyManageInputSchema 返回 taxonomy_manage 的自定义输入 schema
func TaxonomyManageInputSchema() *jsonschema.Schema {
	listPayload := BuildPayloadSchema(map[string]*jsonschema.Schema{
		"page":      {Type: "integer"},
		"page_size": PageSizeSchema(),
	})
	idPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id": {Type: "integer"},
		},
		"id",
	)
	listArticlesPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":        {Type: "integer"},
			"page":      {Type: "integer"},
			"page_size": PageSizeSchema(),
		},
		"id",
	)
	categoryCreatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"name":        {Type: "string"},
			"description": {Type: "string"},
			"sort":        {Type: "integer"},
		},
		"name",
	)
	categoryUpdatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":          {Type: "integer"},
			"name":        {Type: "string"},
			"description": {Type: "string"},
			"sort":        {Type: "integer"},
		},
		"id",
		"name",
	)
	tagCreatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"name":        {Type: "string"},
			"description": {Type: "string"},
		},
		"name",
	)
	tagUpdatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":          {Type: "integer"},
			"name":        {Type: "string"},
			"description": {Type: "string"},
		},
		"id",
		"name",
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"target": {
				Type: "string",
				Enum: []any{taxonomyTargetCategory, taxonomyTargetTag},
			},
			"action": {
				Type: "string",
				Enum: []any{
					taxonomyActionList,
					taxonomyActionCreate,
					taxonomyActionUpdate,
					taxonomyActionDelete,
					taxonomyActionListArticles,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"target", "action", "payload"},
		OneOf: []*jsonschema.Schema{
			buildTargetActionSchema(taxonomyTargetCategory, taxonomyActionList, "获取分类列表", listPayload),
			buildTargetActionSchema(taxonomyTargetCategory, taxonomyActionCreate, "创建分类", categoryCreatePayload),
			buildTargetActionSchema(taxonomyTargetCategory, taxonomyActionUpdate, "更新分类信息", categoryUpdatePayload),
			buildTargetActionSchema(taxonomyTargetCategory, taxonomyActionDelete, "删除分类。风险操作，谨慎使用，不可恢复", idPayload),
			buildTargetActionSchema(taxonomyTargetCategory, taxonomyActionListArticles, "获取分类下的文章列表", listArticlesPayload),
			buildTargetActionSchema(taxonomyTargetTag, taxonomyActionList, "获取标签列表", listPayload),
			buildTargetActionSchema(taxonomyTargetTag, taxonomyActionCreate, "创建标签", tagCreatePayload),
			buildTargetActionSchema(taxonomyTargetTag, taxonomyActionUpdate, "更新标签信息", tagUpdatePayload),
			buildTargetActionSchema(taxonomyTargetTag, taxonomyActionDelete, "删除标签。风险操作，谨慎使用，不可恢复", idPayload),
			buildTargetActionSchema(taxonomyTargetTag, taxonomyActionListArticles, "获取标签下的文章列表", listArticlesPayload),
		},
	}
}

func buildTargetActionSchema(target, action, description string, payload *jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"target": {
				Type: "string",
				Enum: []any{target},
			},
			"action": {
				Type:        "string",
				Enum:        []any{action},
				Description: description,
			},
			"payload": payload,
		},
		Required: []string{"target", "action", "payload"},
	}
}

// ============ 辅助函数============

func convertCategoryToItem(category dto.CategoryResponse) TaxonomyItem {
	return TaxonomyItem{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Count:       category.Count,
		Sort:        category.Sort,
	}
}

func convertTagToItem(tag dto.TagResponse) TaxonomyItem {
	return TaxonomyItem{
		ID:          tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		Count:       tag.Count,
	}
}

func convertWebArticleToItem(item dto.ArticleWebResponse) ArticleItem {
	result := ArticleItem{
		ID:           item.ID,
		Title:        item.Title,
		Cover:        item.Cover,
		Location:     item.Location,
		IsPublish:    true,
		IsTop:        item.IsTop,
		IsEssence:    item.IsEssence,
		IsOutdated:   item.IsOutdated,
		CommentCount: item.CommentCount,
		Category:     CategoryItem{ID: item.Category.ID, Name: item.Category.Name},
		Tags:         make([]TagItem, len(item.Tags)),
	}
	result.PublishTime = ToTimeStringPtr(item.PublishTime)
	result.UpdateTime = ToTimeStringPtr(item.UpdateTime)
	for i, tag := range item.Tags {
		result.Tags[i] = TagItem{ID: tag.ID, Name: tag.Name}
	}
	return result
}
