package service

import (
	"crypto/md5" // #nosec G501 - MD5 用于生成访客标识，不需要加密安全
	"fmt"
	"math"
	"time"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/utils"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
	config    *config.Config
}

func NewStatsService(statsRepo *repository.StatsRepository, conf *config.Config) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		config:    conf,
	}
}

// ============ 通用服务 ============

// RecordVisit 记录访问
func (s *StatsService) RecordVisit(ip, pageURL, userAgent, referer string, articleID *uint) error {
	// 生成访客唯一标识
	visitorID := s.generateVisitorID(ip, userAgent)

	// 解析地理位置和设备信息
	location := utils.GetIPLocation(ip)
	browser, os := utils.ParseUserAgent(userAgent)

	log := &model.Visit{
		VisitorID: visitorID,
		IP:        ip,
		PageURL:   pageURL,
		ArticleID: articleID,
		UserAgent: userAgent,
		Location:  location,
		Browser:   browser,
		OS:        os,
		Referer:   referer,
		VisitDate: time.Now(),
	}

	return s.statsRepo.CreateVisitLog(log)
}

// ============ 前台服务 ============

// GetSiteStats 获取网站统计信息
func (s *StatsService) GetSiteStats() (*dto.SiteStatsResponse, error) {
	stats := &dto.SiteStatsResponse{}

	// 获取总字数
	totalWords, err := s.statsRepo.GetTotalWords()
	if err != nil {
		return nil, err
	}
	stats.TotalWords = formatNumber(totalWords)

	// 获取总访客数
	totalVisitors, err := s.statsRepo.GetTotalUV()
	if err != nil {
		return nil, err
	}
	stats.TotalVisitors = totalVisitors

	// 获取总浏览量
	totalPageViews, err := s.statsRepo.GetTotalPV()
	if err != nil {
		return nil, err
	}
	stats.TotalPageViews = totalPageViews

	// 获取在线用户数
	onlineUsers, err := s.statsRepo.GetOnlineUsers()
	if err != nil {
		return nil, err
	}
	stats.OnlineUsers = onlineUsers

	// 获取今日访客数
	todayVisitors, err := s.statsRepo.GetTodayUV()
	if err != nil {
		return nil, err
	}
	stats.TodayVisitors = todayVisitors

	// 获取今日浏览量
	todayPageviews, err := s.statsRepo.GetTodayPV()
	if err != nil {
		return nil, err
	}
	stats.TodayPageviews = todayPageviews

	// 获取昨日访客数
	yesterdayVisitors, err := s.statsRepo.GetYesterdayUV()
	if err != nil {
		return nil, err
	}
	stats.YesterdayVisitors = yesterdayVisitors

	// 获取昨日浏览量
	yesterdayPageviews, err := s.statsRepo.GetYesterdayPV()
	if err != nil {
		return nil, err
	}
	stats.YesterdayPageviews = yesterdayPageviews

	// 获取本月浏览量
	monthPageviews, err := s.statsRepo.GetMonthPV()
	if err != nil {
		return nil, err
	}
	stats.MonthPageviews = monthPageviews

	totalArticles, err := s.statsRepo.GetTotalPublishedArticles()
	if err != nil {
		return nil, err
	}
	stats.TotalArticles = totalArticles

	totalComments, err := s.statsRepo.GetTotalVisibleComments()
	if err != nil {
		return nil, err
	}
	stats.TotalComments = totalComments

	totalFriends, err := s.statsRepo.GetTotalFriends()
	if err != nil {
		return nil, err
	}
	stats.TotalFriends = totalFriends

	totalMoments, err := s.statsRepo.GetTotalMoments()
	if err != nil {
		return nil, err
	}
	stats.TotalMoments = totalMoments

	totalCategories, err := s.statsRepo.GetPublishedCategoryCount()
	if err != nil {
		return nil, err
	}
	stats.TotalCategories = totalCategories

	totalTags, err := s.statsRepo.GetPublishedTagCount()
	if err != nil {
		return nil, err
	}
	stats.TotalTags = totalTags

	return stats, nil
}

// GetArchives 获取文章归档数据
func (s *StatsService) GetArchives() (*dto.ArchivesResponse, error) {
	archives, err := s.statsRepo.GetArchives()
	if err != nil {
		return nil, err
	}

	return &dto.ArchivesResponse{
		Archives: archives,
	}, nil
}

// ============ 后台管理服务 ============

// GetDashboardStats 获取仪表盘统计数据
func (s *StatsService) GetDashboardStats() (*dto.DashboardStats, error) {
	stats := &dto.DashboardStats{}

	// 基础统计
	totalArticles, err := s.statsRepo.GetTotalPublishedArticles()
	if err != nil {
		return nil, err
	}
	stats.TotalArticles = totalArticles

	totalFriends, err := s.statsRepo.GetTotalFriends()
	if err != nil {
		return nil, err
	}
	stats.TotalFriends = totalFriends

	totalMoments, err := s.statsRepo.GetTotalMoments()
	if err != nil {
		return nil, err
	}
	stats.TotalMoments = totalMoments

	totalViews, err := s.statsRepo.GetTotalPV()
	if err != nil {
		return nil, err
	}
	stats.TotalViews = totalViews

	totalVisitors, err := s.statsRepo.GetTotalUV()
	if err != nil {
		return nil, err
	}
	stats.TotalVisitors = totalVisitors

	totalComments, err := s.statsRepo.GetTotalComments()
	if err != nil {
		return nil, err
	}
	stats.TotalComments = totalComments

	totalUsers, err := s.statsRepo.GetTotalUsers()
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = totalUsers

	// 今日统计
	todayViews, err := s.statsRepo.GetTodayPV()
	if err != nil {
		return nil, err
	}
	stats.TodayViews = todayViews

	todayVisitors, err := s.statsRepo.GetTodayUV()
	if err != nil {
		return nil, err
	}
	stats.TodayVisitors = todayVisitors

	todayComments, err := s.statsRepo.GetTodayComments()
	if err != nil {
		return nil, err
	}
	stats.TodayComments = todayComments

	todayUsers, err := s.statsRepo.GetTodayUsers()
	if err != nil {
		return nil, err
	}
	stats.TodayUsers = todayUsers

	// 计算增长率
	yesterdayViews, _ := s.statsRepo.GetYesterdayPV()
	stats.ViewsGrowth = calculateGrowthRate(todayViews, yesterdayViews)

	yesterdayVisitors, _ := s.statsRepo.GetYesterdayUV()
	stats.VisitorsGrowth = calculateGrowthRate(todayVisitors, yesterdayVisitors)

	yesterdayComments, _ := s.statsRepo.GetYesterdayComments()
	stats.CommentsGrowth = calculateGrowthRate(todayComments, yesterdayComments)

	yesterdayUsers, _ := s.statsRepo.GetYesterdayUsers()
	stats.UsersGrowth = calculateGrowthRate(todayUsers, yesterdayUsers)

	return stats, nil
}

// GetTrendData 获取趋势数据
func (s *StatsService) GetTrendData(startDate, endDate string, trendType string) ([]dto.TrendData, error) {
	// 解析日期
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	// 验证日期范围
	if end.Before(start) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// 限制查询范围
	maxDays := 365
	if end.Sub(start).Hours() > float64(maxDays*24) {
		return nil, fmt.Errorf("date range cannot exceed %d days", maxDays)
	}

	// 验证趋势类型
	if trendType != "daily" && trendType != "weekly" && trendType != "monthly" {
		trendType = "daily"
	}

	return s.statsRepo.GetTrendData(start, end, trendType)
}

// GetCategoryStats 获取分类统计
func (s *StatsService) GetCategoryStats() ([]dto.CategoryStats, error) {
	return s.statsRepo.GetCategoryStats()
}

// GetTagStats 获取标签统计
func (s *StatsService) GetTagStats() ([]dto.TagStats, error) {
	return s.statsRepo.GetTagStats()
}

// GetArticleContribution 获取文章贡献统计
func (s *StatsService) GetArticleContribution(year *int, month *int) ([]dto.ArticleContribution, error) {
	// 验证月份范围
	if month != nil && (*month < 1 || *month > 12) {
		return nil, fmt.Errorf("month must be between 1 and 12")
	}

	// 如果只传了month没传year，返回错误
	if month != nil && year == nil {
		return nil, fmt.Errorf("month parameter requires year parameter")
	}

	return s.statsRepo.GetArticleContribution(year, month)
}

// GetVisitLogs 获取访问日志列表
func (s *StatsService) GetVisitLogs(req *dto.GetVisitLogsRequest) ([]dto.VisitLogItem, int64, int, int, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 查询访问日志
	visits, total, err := s.statsRepo.GetVisitLogs(req)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// 转换为响应结构
	items := make([]dto.VisitLogItem, 0, len(visits))
	for _, visit := range visits {
		items = append(items, dto.VisitLogItem{
			ID:        visit.ID,
			VisitorID: visit.VisitorID,
			IP:        visit.IP,
			PageURL:   visit.PageURL,
			UserAgent: visit.UserAgent,
			Location:  visit.Location,
			Browser:   visit.Browser,
			OS:        visit.OS,
			Referer:   visit.Referer,
			CreatedAt: visit.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return items, total, req.Page, req.PageSize, nil
}

// ============ 辅助方法 ============

// generateVisitorID 生成访客唯一标识
func (s *StatsService) generateVisitorID(ip, userAgent string) string {
	data := fmt.Sprintf("%s:%s", ip, userAgent)
	// #nosec G401 - MD5 用于生成访客标识，不需要加密安全
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// calculateGrowthRate 计算增长率
func calculateGrowthRate(today, yesterday int64) float64 {
	if yesterday == 0 {
		if today > 0 {
			return 100.0
		}
		return 0.0
	}
	rate := float64(today-yesterday) / float64(yesterday) * 100
	return math.Round(rate*100) / 100
}

// formatNumber 格式化数字显示
func formatNumber(num int) string {
	if num < 1000 {
		return fmt.Sprintf("%d", num)
	}
	if num < 10000 {
		return fmt.Sprintf("%.1fk", float64(num)/1000)
	}
	if num < 1000000 {
		return fmt.Sprintf("%.0fk", float64(num)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(num)/1000000)
}
