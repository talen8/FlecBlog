package v1

import (
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// MomentController 动态控制器
type MomentController struct {
	momentService *service.MomentService
}

// NewMomentController 创建动态控制器
func NewMomentController(momentService *service.MomentService) *MomentController {
	return &MomentController{
		momentService: momentService,
	}
}

// ============ 前台接口 ============

// ListForWeb 前台获取动态列表
//
//	@Summary		动态列表
//	@Description	获取所有公开的动态
//	@Tags			动态
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Router			/moments [get]
func (c *MomentController) ListForWeb(ctx *gin.Context) {
	var req dto.ListMomentsForWebRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	moments, total, err := c.momentService.ListForWeb(ctx.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, moments, total, req.Page, req.PageSize)
}

// ============ 后台管理接口 ============

// List 获取动态列表
//
//	@Summary		动态列表
//	@Description	获取所有动态，支持筛选
//	@Tags			动态管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page			query		int		false	"页码"
//	@Param			page_size		query		int		false	"每页数量（不传则返回全部）"
//	@Param			keyword			query		string	false	"搜索关键词（文本内容）"
//	@Param			tags			query		string	false	"标签"
//	@Param			location		query		string	false	"发布地点"
//	@Param			is_publish		query		bool	false	"是否发布"
//	@Param			has_images		query		bool	false	"是否有图片"
//	@Param			has_video		query		bool	false	"是否有视频"
//	@Param			has_music		query		bool	false	"是否有音乐"
//	@Param			has_link		query		bool	false	"是否有链接"
//	@Param			start_time		query		string	false	"发布开始时间"
//	@Param			end_time		query		string	false	"发布结束时间"
//	@Success		200				{object}	response.Response{data=response.PageResult}
//	@Failure		401				{object}	response.Response
//	@Failure		403				{object}	response.Response
//	@Router			/admin/moments [get]
func (c *MomentController) List(ctx *gin.Context) {
	var req dto.ListMomentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	moments, total, err := c.momentService.List(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, moments, total, req.Page, req.PageSize)
}

// Get 获取动态详情
//
//	@Summary		动态详情
//	@Description	获取动态详细信息
//	@Tags			动态管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"动态 ID"
//	@Success		200	{object}	response.Response{data=dto.MomentListResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/moments/{id} [get]
func (c *MomentController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	moment, err := c.momentService.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, moment)
}

// Create 创建动态
//
//	@Summary		创建动态
//	@Description	创建新动态
//	@Tags			动态管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateMomentRequest	true	"动态信息"
//	@Success		201		{object}	response.Response{data=dto.MomentListResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/moments [post]
func (c *MomentController) Create(ctx *gin.Context) {
	var req dto.CreateMomentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	moment, err := c.momentService.Create(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, moment)
}

// Update 更新动态
//
//	@Summary		更新动态
//	@Description	修改动态信息
//	@Tags			动态管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"动态 ID"
//	@Param			request	body		dto.UpdateMomentRequest	true	"动态信息"
//	@Success		200		{object}	response.Response{data=dto.MomentListResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/moments/{id} [put]
func (c *MomentController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateMomentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	moment, err := c.momentService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, moment)
}

// Delete 删除动态
//
//	@Summary		删除动态
//	@Description	硬删除动态，不可恢复
//	@Tags			动态管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"动态 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/moments/{id} [delete]
func (c *MomentController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.momentService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}
