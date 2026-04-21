package v1

import (
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// FriendController 友链控制器
type FriendController struct {
	friendService *service.FriendService
}

// NewFriendController 创建友链控制器
func NewFriendController(friendService *service.FriendService) *FriendController {
	return &FriendController{friendService: friendService}
}

// ============ 前台接口 ============

// ListForWeb 获取友链列表
//
//	@Summary		友链列表
//	@Description	按类型获取友链列表
//	@Tags			友链
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.GroupedFriendsResponse}
//	@Router			/friends [get]
func (c *FriendController) ListForWeb(ctx *gin.Context) {
	result, err := c.friendService.ListForWeb(ctx.Request.Context())
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}

// ApplyFriend 申请友链
//
//	@Summary		申请友链
//	@Description	用户提交友链申请，需登录
//	@Tags			友链
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ApplyFriendRequest	true	"友链申请信息"
//	@Success		201		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/friends/apply [post]
func (c *FriendController) ApplyFriend(ctx *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "请先登录")
		return
	}

	var req dto.ApplyFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.friendService.ApplyFriend(ctx.Request.Context(), &req, userID.(uint)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, nil)
}

// ============ 后台管理接口 - 友链类型 ============

// ListTypes 获取所有友链类型
//
//	@Summary		友链类型列表 [类型]
//	@Description	获取所有友链类型（包括禁用的）
//	@Tags			友链管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/friends/types [get]
func (c *FriendController) ListTypes(ctx *gin.Context) {
	var req dto.ListFriendTypeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	types, total, err := c.friendService.ListTypes(ctx.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, types, total, req.Page, req.PageSize)
}

// GetType 获取友链类型详情
//
//	@Summary		友链类型详情 [类型]
//	@Description	获取友链类型详细信息
//	@Tags			友链管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"类型 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/friends/types/{id} [get]
func (c *FriendController) GetType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	friendType, err := c.friendService.GetType(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, friendType)
}

// CreateType 创建友链类型
//
//	@Summary		创建友链类型 [类型]
//	@Description	创建新的友链类型
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateFriendTypeRequest	true	"类型信息"
//	@Success		201		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/friends/types [post]
func (c *FriendController) CreateType(ctx *gin.Context) {
	var req dto.CreateFriendTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.friendService.CreateType(ctx.Request.Context(), &req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, nil)
}

// UpdateType 更新友链类型
//
//	@Summary		更新友链类型 [类型]
//	@Description	修改友链类型信息
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"类型 ID"
//	@Param			request	body		dto.UpdateFriendTypeRequest	true	"类型信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/friends/types/{id} [put]
func (c *FriendController) UpdateType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateFriendTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.friendService.UpdateType(ctx.Request.Context(), uint(id), &req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// DeleteType 删除友链类型
//
//	@Summary		删除友链类型 [类型]
//	@Description	删除友链类型，关联友链的类型会被设置为 NULL
//	@Tags			友链管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"类型 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/friends/types/{id} [delete]
func (c *FriendController) DeleteType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.friendService.DeleteType(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// ============ 后台管理接口 - 友链 ============

// List 获取友链列表
//
//	@Summary		友链列表
//	@Description	获取所有友链用于管理
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page				query		int		false	"页码"
//	@Param			page_size			query		int		false	"每页数量（不传则返回全部）"
//	@Param			keyword				query		string	false	"关键词搜索（名称、链接、描述）"
//	@Param			type_id				query		int		false	"友链类型ID"
//	@Param			is_invalid			query		bool	false	"是否失效"
//	@Param			is_pending			query		bool	false	"是否待审核"
//	@Param			accessible_status	query		string	false	"可访问性状态: normal=正常, abnormal=异常, ignored=忽略检查"
//	@Param			rss_status			query		string	false	"RSS状态: no_rss=无订阅, normal=正常订阅, warning=三个月未更新, danger=六个月未更新"
//	@Param			has_screenshot		query		bool	false	"是否包含截图"
//	@Success		200					{object}	response.Response
//	@Failure		401					{object}	response.Response
//	@Failure		403					{object}	response.Response
//	@Router			/admin/friends [get]
func (c *FriendController) List(ctx *gin.Context) {
	var req dto.ListFriendRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	friends, total, err := c.friendService.List(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, friends, total, req.Page, req.PageSize)
}

// Get 获取友链信息
//
//	@Summary		友链详情
//	@Description	获取友链详细信息
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"友链 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/friends/{id} [get]
func (c *FriendController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	friend, err := c.friendService.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, friend)
}

// Create 创建友链
//
//	@Summary		创建友链
//	@Description	创建新友链
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateFriendRequest	true	"友链信息"
//	@Success		201		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/friends [post]
func (c *FriendController) Create(ctx *gin.Context) {
	var req dto.CreateFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if _, err := c.friendService.Create(ctx.Request.Context(), &req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, nil)
}

// Update 更新友链
//
//	@Summary		更新友链
//	@Description	修改友链信息
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"友链 ID"
//	@Param			request	body		dto.UpdateFriendRequest	true	"友链信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/friends/{id} [put]
func (c *FriendController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.friendService.Update(ctx.Request.Context(), uint(id), &req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// Delete 删除友链
//
//	@Summary		删除友链
//	@Description	硬删除友链，不可恢复
//	@Tags			友链管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"友链 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/friends/{id} [delete]
func (c *FriendController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.friendService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}
