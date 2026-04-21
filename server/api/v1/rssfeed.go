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
