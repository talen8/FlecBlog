package v1

import (
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// RssFeedController RSS订阅控制器
type RssFeedController struct {
	rssFeedService *service.RssFeedService
}

// NewRssFeedController 创建RSS订阅控制器
func NewRssFeedController(rssFeedService *service.RssFeedService) *RssFeedController {
	return &RssFeedController{rssFeedService: rssFeedService}
}

// List 获取RSS文章列表
//
//	@Summary		RSS文章列表
//	@Description	获取RSS订阅文章列表
//	@Tags			RSS订阅管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页数量"
//	@Param			keyword		query		string	false	"搜索关键词"
//	@Param			friend_id	query		int		false	"友链ID筛选"
//	@Param			is_read		query		bool	false	"已读状态筛选"
//	@Param			start_time	query		string	false	"发布开始时间（格式：2006-01-02）"
//	@Param			end_time	query		string	false	"发布结束时间（格式：2006-01-02）"
//	@Success		200			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/rssfeed [get]
func (c *RssFeedController) List(ctx *gin.Context) {
	var req dto.ListRssArticleRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	result, err := c.rssFeedService.List(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}

// FriendCircle 获取友圈文章
//
//	@Summary		友圈
//	@Description	返回友圈文章，可按关键词搜索；无筛选条件时返回最新文章
//	@Tags			友链
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页数量"
//	@Param			keyword		query		string	false	"关键词，匹配文章标题或友链站点名"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Router			/friends/circle [get]
func (c *RssFeedController) FriendCircle(ctx *gin.Context) {
	var req dto.FriendCircleRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	result, err := c.rssFeedService.GetFriendCircle(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}

// FriendCircleRandom 随机获取一篇友圈文章
//
//	@Summary		友圈随机文章
//	@Description	从友圈全部文章中随机返回一篇
//	@Tags			友链
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/friends/circle/random [get]
func (c *RssFeedController) FriendCircleRandom(ctx *gin.Context) {
	result, err := c.rssFeedService.GetFriendCircleRandom(ctx.Request.Context())
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}

// FriendCircleStats 获取友圈统计数据
//
//	@Summary		友圈统计
//	@Description	统计友圈的文章总数、站点数量、今日更新数
//	@Tags			友链
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/friends/circle/stats [get]
func (c *RssFeedController) FriendCircleStats(ctx *gin.Context) {
	result, err := c.rssFeedService.GetFriendCircleStats(ctx.Request.Context())
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}

// MarkRead 标记文章已读
//
//	@Summary		标记文章已读
//	@Description	将指定文章标记为已读，限超级管理员
//	@Tags			RSS订阅管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文章ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/rssfeed/{id}/read [put]
func (c *RssFeedController) MarkRead(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, "无效的文章ID")
		return
	}

	if err := c.rssFeedService.MarkRead(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// SetBlockCircle 设置文章的友圈屏蔽状态
//
//	@Summary		设置文章友圈屏蔽
//	@Description	屏蔽后该文章不再出现在前台友圈，限超级管理员
//	@Tags			RSS订阅管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"文章ID"
//	@Param			body	body		dto.SetArticleBlockCircleRequest	true	"屏蔽状态"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/rssfeed/{id}/block-circle [put]
func (c *RssFeedController) SetBlockCircle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, "无效的文章ID")
		return
	}

	var req dto.SetArticleBlockCircleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.rssFeedService.SetArticleBlockCircle(ctx.Request.Context(), uint(id), req.BlockCircle); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// MarkAllRead 全部标记已读
//
//	@Summary		全部标记已读
//	@Description	将所有未读文章标记为已读，限超级管理员
//	@Tags			RSS订阅管理
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Router			/admin/rssfeed/read-all [put]
func (c *RssFeedController) MarkAllRead(ctx *gin.Context) {
	affected, err := c.rssFeedService.MarkAllRead(ctx.Request.Context())
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"affected": affected})
}
