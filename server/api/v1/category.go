package v1

import (
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService *service.CategoryService
}

// NewCategoryController 创建分类控制器
func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// ============ 前台接口 ============

// ListForWeb 获取前台分类列表
//
//	@Summary		分类列表
//	@Description	获取所有分类
//	@Tags			分类
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response
//	@Router			/categories [get]
func (c *CategoryController) ListForWeb(ctx *gin.Context) {
	var req dto.ListCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	categories, total, err := c.categoryService.ListForWeb(ctx.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, categories, total, req.Page, req.PageSize)
}

// GetBySlug 根据 slug 获取分类
//
//	@Summary		分类详情
//	@Description	通过 slug 获取分类信息
//	@Tags			分类
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"分类 slug"
//	@Success		200		{object}	response.Response{data=dto.CategoryForWebResponse}
//	@Failure		404		{object}	response.Response
//	@Router			/categories/{slug} [get]
func (c *CategoryController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.ValidateFailed(ctx, "slug不能为空")
		return
	}

	category, err := c.categoryService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, category)
}

// ============ 后台管理接口 ============

// List 获取分类列表
//
//	@Summary		分类列表
//	@Description	获取所有分类
//	@Tags			分类管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/categories [get]
func (c *CategoryController) List(ctx *gin.Context) {
	var req dto.ListCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	categories, total, err := c.categoryService.List(ctx.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, categories, total, req.Page, req.PageSize)
}

// Create 创建分类
//
//	@Summary		创建分类
//	@Description	创建新分类，自动生成 slug
//	@Tags			分类管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateCategoryRequest	true	"分类信息"
//	@Success		201		{object}	response.Response{data=model.Category}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/categories [post]
func (c *CategoryController) Create(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		Sort:        req.Sort,
	}

	if err := c.categoryService.Create(ctx.Request.Context(), category); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, category)
}

// Update 更新分类
//
//	@Summary		更新分类
//	@Description	修改分类信息
//	@Tags			分类管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"分类 ID"
//	@Param			request	body		dto.UpdateCategoryRequest	true	"分类信息"
//	@Success		200		{object}	response.Response{data=model.Category}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/categories/{id} [put]
func (c *CategoryController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		Sort:        req.Sort,
	}

	if err := c.categoryService.Update(ctx.Request.Context(), uint(id), category); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, category)
}

// Delete 删除分类
//
//	@Summary		删除分类
//	@Description	硬删除分类，不可恢复
//	@Tags			分类管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"分类 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/categories/{id} [delete]
func (c *CategoryController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.categoryService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}
