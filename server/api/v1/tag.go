package v1

import (
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// TagController 标签控制器
type TagController struct {
	tagService *service.TagService
}

// NewTagController 创建标签控制器
func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

// ============ 前台接口 ============

// ListForWeb 前台获取标签列表
//
//	@Summary		标签列表
//	@Description	获取所有标签
//	@Tags			标签
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response
//	@Router			/tags [get]
func (c *TagController) ListForWeb(ctx *gin.Context) {
	var req dto.ListTagRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	tags, total, err := c.tagService.ListForWeb(ctx.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, tags, total, req.Page, req.PageSize)
}

// GetBySlug 根据 slug 获取标签信息
//
//	@Summary		标签详情
//	@Description	通过 slug 获取标签信息
//	@Tags			标签
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"标签 slug"
//	@Success		200		{object}	response.Response{data=dto.TagForWebResponse}
//	@Failure		404		{object}	response.Response
//	@Router			/tags/{slug} [get]
func (c *TagController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.ValidateFailed(ctx, "slug不能为空")
		return
	}

	tag, err := c.tagService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, tag)
}

// ============ 后台管理接口 ============

// List 获取标签列表
//
//	@Summary		标签列表
//	@Description	获取所有标签
//	@Tags			标签管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/tags [get]
func (c *TagController) List(ctx *gin.Context) {
	var req dto.ListTagRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	tags, total, err := c.tagService.List(ctx.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, tags, total, req.Page, req.PageSize)
}

// Create 创建标签
//
//	@Summary		创建标签
//	@Description	创建新标签，自动生成 slug
//	@Tags			标签管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateTagRequest	true	"标签信息"
//	@Success		201		{object}	response.Response{data=model.Tag}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/tags [post]
func (c *TagController) Create(ctx *gin.Context) {
	var req dto.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	tag := &model.Tag{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := c.tagService.Create(ctx.Request.Context(), tag); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, tag)
}

// Update 更新标签
//
//	@Summary		更新标签
//	@Description	修改标签信息
//	@Tags			标签管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"标签 ID"
//	@Param			request	body		dto.UpdateTagRequest	true	"标签信息"
//	@Success		200		{object}	response.Response{data=model.Tag}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/tags/{id} [put]
func (c *TagController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	tag := &model.Tag{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := c.tagService.Update(ctx.Request.Context(), uint(id), tag); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, tag)
}

// Delete 删除标签
//
//	@Summary		删除标签
//	@Description	硬删除标签，不可恢复
//	@Tags			标签管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"标签 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/tags/{id} [delete]
func (c *TagController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.tagService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}
