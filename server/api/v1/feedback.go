package v1

import (
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/errcode"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// FeedbackHandler 反馈处理器
type FeedbackHandler struct {
	service *service.FeedbackService
}

// NewFeedbackHandler 创建反馈处理器实例
func NewFeedbackHandler(service *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		service: service,
	}
}

// Submit 提交反馈投诉
//
//	@Summary	提交反馈投诉
//	@Description	用户提交反馈投诉
//	@Tags		反馈
//	@Accept		json
//	@Produce	json
//	@Param		feedback	body		dto.SubmitFeedbackRequest	true	"反馈信息"
//	@Success	200			{object}	response.Response{data=dto.FeedbackResponse}
//	@Router		/api/v1/feedback [post]
func (h *FeedbackHandler) Submit(c *gin.Context) {
	var req dto.SubmitFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 获取IP和UserAgent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	feedback, err := h.service.Submit(c.Request.Context(), &req, ip, userAgent)
	if err != nil {
		response.Error(c, errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.Success(c, feedback)
}

// GetByTicketNo 查询反馈投诉进度
//
//	@Summary	查询反馈投诉进度
//	@Description	根据工单号查询反馈投诉的处理进度和结果
//	@Tags		反馈
//	@Accept		json
//	@Produce	json
//	@Param		ticket_no	path		string	true	"工单号"
//	@Success	200			{object}	response.Response{data=dto.FeedbackResponse}
//	@Router		/api/v1/feedback/ticket/{ticket_no} [get]
func (h *FeedbackHandler) GetByTicketNo(c *gin.Context) {
	ticketNo := c.Param("ticket_no")
	if ticketNo == "" {
		response.Error(c, errcode.InvalidParams.WithDetails("工单号不能为空"))
		return
	}

	feedback, err := h.service.GetByTicketNo(c.Request.Context(), ticketNo)
	if err != nil {
		response.Error(c, errcode.NotFound.WithDetails("未找到该工单"))
		return
	}

	response.Success(c, feedback)
}

// ============ 后台管理接口 ============

// List 获取反馈列表
//
//	@Summary	获取反馈列表
//	@Tags		反馈管理
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		page		query		int		true	"页码"
//	@Param		page_size	query		int		true	"每页数量"
//	@Param		keyword		query		string	false	"搜索关键词（工单号、投诉地址）"
//	@Param		report_type	query		string	false	"反馈类型筛选（copyright/inappropriate/summary/suggestion）"
//	@Param		status		query		string	false	"状态筛选（pending/resolved/closed）"
//	@Param		start_time	query		string	false	"反馈开始时间（格式：2006-01-02）"
//	@Param		end_time	query		string	false	"反馈结束时间（格式：2006-01-02）"
//	@Success	200			{object}	response.Response{data=response.PageResult{list=[]dto.FeedbackResponse}}
//	@Router		/api/v1/admin/feedback [get]
func (h *FeedbackHandler) List(c *gin.Context) {
	var req dto.FeedbackQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	feedbacks, total, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.Success(c, response.PageResult{
		List:     feedbacks,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// Get 获取反馈详情
//
//	@Summary	获取反馈详情
//	@Tags		反馈管理
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path		int	true	"反馈ID"
//	@Success	200	{object}	response.Response{data=dto.FeedbackResponse}
//	@Router		/api/v1/admin/feedback/{id} [get]
func (h *FeedbackHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails("无效的反馈ID"))
		return
	}

	feedback, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.Success(c, feedback)
}

// Update 更新反馈
//
//	@Summary	更新反馈
//	@Tags		反馈管理
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id			path		int							true	"反馈ID"
//	@Param		feedback	body		dto.UpdateFeedbackRequest	true	"更新信息"
//	@Success	200			{object}	response.Response
//	@Router		/api/v1/admin/feedback/{id} [put]
func (h *FeedbackHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails("无效的反馈ID"))
		return
	}

	var req dto.UpdateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		response.Error(c, errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.Success(c, nil)
}

// Delete 删除反馈
//
//	@Summary	删除反馈
//	@Tags		反馈管理
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path		int	true	"反馈ID"
//	@Success	200	{object}	response.Response
//	@Router		/api/v1/admin/feedback/{id} [delete]
func (h *FeedbackHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails("无效的反馈ID"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.Success(c, nil)
}
