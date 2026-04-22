package v1

import (
	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatsHandler 统计控制器
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler 创建统计控制器
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// ============ 前台接口 ============

// Collect 收集前端追踪数据
//
//	@Summary		数据收集
//	@Description	前端埋点数据收集，记录页面访问等
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.CollectRequest	true	"访问数据"
//	@Success		204
//	@Router			/collect [post]
func (h *StatsHandler) Collect(c *gin.Context) {
	var req dto.CollectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(204)
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	if req.Type == "pageview" || req.Type == "" {
		go func() {
			_ = h.statsService.RecordVisit(ip, req.URL, userAgent, req.Referrer, req.ArticleID)
		}()
	}

	c.Status(204)
}

// GetSiteStats 获取前台网站统计信息
//
//	@Summary		网站统计
//	@Description	获取博客前台公开统计数据，包含规模总览与访问数据
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.SiteStatsResponse}
//	@Router			/stats/site [get]
func (h *StatsHandler) GetSiteStats(c *gin.Context) {
	stats, err := h.statsService.GetSiteStats()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetArchives 获取文章归档数据
//
//	@Summary		文章归档
//	@Description	按年月归档的文章统计，用于归档页面
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.ArchivesResponse}
//	@Router			/stats/archives [get]
func (h *StatsHandler) GetArchives(c *gin.Context) {
	archives, err := h.statsService.GetArchives()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, archives)
}

// ============ 后台管理接口 ============

// GetDashboard 获取仪表盘统计数据
//
//	@Summary		仪表盘统计
//	@Description	获取基础统计、今日数据和趋势对比
//	@Tags			统计管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.DashboardStats}
//	@Router			/admin/stats/dashboard [get]
func (h *StatsHandler) GetDashboard(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetTrend 获取访问趋势数据
//
//	@Summary		访问趋势
//	@Description	指定时间段的访问趋势，支持按天/周/月聚合
//	@Tags			统计管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			start_date	query		string	true	"开始日期 YYYY-MM-DD"			example(2025-09-01)
//	@Param			end_date	query		string	true	"结束日期 YYYY-MM-DD"			example(2025-10-05)
//	@Param			type		query		string	false	"统计类型 daily/weekly/monthly"	default(daily)	Enums(daily, weekly, monthly)
//	@Success		200			{object}	response.Response{data=[]dto.TrendData}
//	@Router			/admin/stats/trend [get]
func (h *StatsHandler) GetTrend(c *gin.Context) {
	var req dto.GetTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(c, err.Error())
		return
	}

	if req.Type == "" {
		req.Type = "daily"
	}

	trends, err := h.statsService.GetTrendData(req.StartDate, req.EndDate, req.Type)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, trends)
}

// GetCategoryStats 获取分类统计数据
//
//	@Summary		分类统计
//	@Description	每个分类的文章数量统计
//	@Tags			统计管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=[]dto.CategoryStats}
//	@Router			/admin/stats/category [get]
func (h *StatsHandler) GetCategoryStats(c *gin.Context) {
	stats, err := h.statsService.GetCategoryStats()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetTagStats 获取标签统计数据
//
//	@Summary		标签统计
//	@Description	每个标签的文章数量统计
//	@Tags			统计管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=[]dto.TagStats}
//	@Router			/admin/stats/tag [get]
func (h *StatsHandler) GetTagStats(c *gin.Context) {
	stats, err := h.statsService.GetTagStats()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetArticleContribution 获取文章贡献数据
//
//	@Summary		文章贡献图
//	@Description	获取文章发布数据，支持按年份或月份查询
//	@Tags			统计管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			year	query		int	false	"年份（可选）"	example(2025)
//	@Param			month	query		int	false	"月份 1-12（可选）"	example(11)
//	@Success		200		{object}	response.Response{data=[]dto.ArticleContribution}
//	@Router			/admin/stats/contribution [get]
func (h *StatsHandler) GetArticleContribution(c *gin.Context) {
	var req dto.GetContributionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(c, err.Error())
		return
	}

	stats, err := h.statsService.GetArticleContribution(req.Year, req.Month)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetVisitLogs 获取访问日志列表
//
//	@Summary		访问日志
//	@Description	获取访问日志列表，支持分页和多种筛选条件
//	@Tags			统计管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int		false	"页码"	default(1)	minimum(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)	minimum(1)	maximum(100)
//	@Param			keyword		query		string	false	"搜索关键词（页面URL）"
//	@Param			visitor_id	query		string	false	"访客ID"
//	@Param			ip			query		string	false	"IP地址"
//	@Param			location	query		string	false	"地理位置"
//	@Param			browser		query		string	false	"浏览器"
//	@Param			os			query		string	false	"操作系统"
//	@Param			start_time	query		string	false	"开始时间（格式：2006-01-02）"
//	@Param			end_time	query		string	false	"结束时间（格式：2006-01-02）"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Router			/admin/stats/visits [get]
func (h *StatsHandler) GetVisitLogs(c *gin.Context) {
	var req dto.GetVisitLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(c, err.Error())
		return
	}

	list, total, page, pageSize, err := h.statsService.GetVisitLogs(&req)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	response.PageSuccess(c, list, total, page, pageSize)
}
