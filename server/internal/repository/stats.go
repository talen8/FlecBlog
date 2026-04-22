package repository

import (
	"strings"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// StatsRepository 统计仓储
type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository 创建统计仓储
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// ============ 访问记录 ============

// CreateVisitLog 创建访问记录
func (r *StatsRepository) CreateVisitLog(log *model.Visit) error {
	return r.db.Create(log).Error
}

// GetTodayPV 获取今日浏览量
func (r *StatsRepository) GetTodayPV() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.Visit{}).
		Where("visit_date = ?", today).
		Count(&count).Error
	return count, err
}

// GetTodayUV 获取今日访客量（独立访客数）
func (r *StatsRepository) GetTodayUV() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.Visit{}).
		Where("visit_date = ?", today).
		Distinct("visitor_id").
		Count(&count).Error
	return count, err
}

// GetYesterdayPV 获取昨日浏览量
func (r *StatsRepository) GetYesterdayPV() (int64, error) {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := r.db.Model(&model.Visit{}).
		Where("visit_date = ?", yesterday).
		Count(&count).Error
	return count, err
}

// GetYesterdayUV 获取昨日访客量
func (r *StatsRepository) GetYesterdayUV() (int64, error) {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := r.db.Model(&model.Visit{}).
		Where("visit_date = ?", yesterday).
		Distinct("visitor_id").
		Count(&count).Error
	return count, err
}

// GetMonthPV 获取本月浏览量
func (r *StatsRepository) GetMonthPV() (int64, error) {
	var count int64
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	err := r.db.Model(&model.Visit{}).
		Where("visit_date >= ?", firstDayOfMonth).
		Count(&count).Error
	return count, err
}

// GetTotalPV 获取总浏览量
func (r *StatsRepository) GetTotalPV() (int64, error) {
	var count int64
	err := r.db.Model(&model.Visit{}).Count(&count).Error
	return count, err
}

// GetTotalUV 获取总访客量（历史独立访客数）
func (r *StatsRepository) GetTotalUV() (int64, error) {
	var count int64
	err := r.db.Model(&model.Visit{}).
		Distinct("visitor_id").
		Count(&count).Error
	return count, err
}

// ============ 文章统计 ============

// GetTotalPublishedArticles 获取已发布文章总数
func (r *StatsRepository) GetTotalPublishedArticles() (int64, error) {
	var count int64
	err := r.db.Model(&model.Article{}).
		Where("is_publish = ?", true).
		Count(&count).Error
	return count, err
}

// GetTodayArticles 获取今日发布文章数
func (r *StatsRepository) GetTodayArticles() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.Article{}).
		Where("is_publish = ? AND DATE(publish_time) = ?", true, today).
		Count(&count).Error
	return count, err
}

// ============ 动态统计 ============

// GetTotalMoments 获取公开动态总数
func (r *StatsRepository) GetTotalMoments() (int64, error) {
	var count int64
	err := r.db.Model(&model.Moment{}).
		Where("is_publish = ?", true).
		Count(&count).Error
	return count, err
}

// ============ 友链统计 ============

// GetTotalFriends 获取启用友链总数
func (r *StatsRepository) GetTotalFriends() (int64, error) {
	var count int64
	err := r.db.Model(&model.Friend{}).
		Where("is_invalid = ?", false).
		Count(&count).Error
	return count, err
}

// ============ 评论统计 ============

// GetTotalComments 获取总评论数
func (r *StatsRepository) GetTotalComments() (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// GetTotalVisibleComments 获取可见评论总数
func (r *StatsRepository) GetTotalVisibleComments() (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).
		Where("status = ? AND deleted_at IS NULL", 1).
		Count(&count).Error
	return count, err
}

// GetTodayComments 获取今日评论数
func (r *StatsRepository) GetTodayComments() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.Comment{}).
		Where("DATE(created_at) = ? AND deleted_at IS NULL", today).
		Count(&count).Error
	return count, err
}

// GetYesterdayComments 获取昨日评论数
func (r *StatsRepository) GetYesterdayComments() (int64, error) {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := r.db.Model(&model.Comment{}).
		Where("DATE(created_at) = ? AND deleted_at IS NULL", yesterday).
		Count(&count).Error
	return count, err
}

// ============ 用户统计 ============

// GetTotalUsers 获取总用户数
func (r *StatsRepository) GetTotalUsers() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetTodayUsers 获取今日注册用户数
func (r *StatsRepository) GetTodayUsers() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.User{}).
		Where("DATE(created_at) = ?", today).
		Count(&count).Error
	return count, err
}

// GetYesterdayUsers 获取昨日注册用户数
func (r *StatsRepository) GetYesterdayUsers() (int64, error) {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := r.db.Model(&model.User{}).
		Where("DATE(created_at) = ?", yesterday).
		Count(&count).Error
	return count, err
}

// ============================
// 趋势数据
// ============================

// GetTrendData 获取趋势数据
func (r *StatsRepository) GetTrendData(startDate, endDate time.Time, groupType string) ([]dto.TrendData, error) {
	// 查询数据库中实际有数据的记录
	var dbResults []dto.TrendData

	// 根据类型选择不同的日期格式
	var dateFormat string
	switch groupType {
	case "daily":
		dateFormat = "YYYY-MM-DD"
	case "weekly":
		dateFormat = "IYYY-IW"
	case "monthly":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM-DD"
	}

	err := r.db.Model(&model.Visit{}).
		Select("TO_CHAR(visit_date, ?) as date, COUNT(*) as pv_count, COUNT(DISTINCT visitor_id) as uv_count", dateFormat).
		Where("visit_date >= ? AND visit_date <= ?", startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&dbResults).Error

	if err != nil {
		return nil, err
	}

	// 前端已实现日期填充逻辑，后端只需返回有数据的记录
	// 减少后端计算量和数据传输量
	return dbResults, nil
}

// ============ 分类/标签统计 ============

// GetCategoryStats 获取分类统计数据
func (r *StatsRepository) GetCategoryStats() ([]dto.CategoryStats, error) {
	var results []dto.CategoryStats

	// 实时统计每个分类下已发布文章的数量
	err := r.db.Model(&model.Category{}).
		Select("categories.name, COUNT(articles.id) as count").
		Joins("LEFT JOIN articles ON articles.category_id = categories.id AND articles.is_publish = true").
		Group("categories.id, categories.name").
		Having("COUNT(articles.id) > 0").
		Order("count DESC").
		Scan(&results).Error

	return results, err
}

// GetPublishedCategoryCount 获取有已发布文章的分类数量
func (r *StatsRepository) GetPublishedCategoryCount() (int64, error) {
	var count int64

	err := r.db.Model(&model.Category{}).
		Joins("LEFT JOIN articles ON articles.category_id = categories.id AND articles.is_publish = ?", true).
		Group("categories.id").
		Having("COUNT(articles.id) > 0").
		Count(&count).Error

	return count, err
}

// GetTagStats 获取标签统计数据
func (r *StatsRepository) GetTagStats() ([]dto.TagStats, error) {
	var results []dto.TagStats

	// 实时统计每个标签下已发布文章的数量
	err := r.db.Model(&model.Tag{}).
		Select("tags.name, COUNT(DISTINCT articles.id) as count").
		Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
		Joins("LEFT JOIN articles ON articles.id = article_tags.article_id AND articles.is_publish = true").
		Group("tags.id, tags.name").
		Having("COUNT(DISTINCT articles.id) > 0").
		Order("count DESC").
		Scan(&results).Error

	return results, err
}

// GetPublishedTagCount 获取有已发布文章的标签数量
func (r *StatsRepository) GetPublishedTagCount() (int64, error) {
	var count int64

	err := r.db.Model(&model.Tag{}).
		Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
		Joins("LEFT JOIN articles ON articles.id = article_tags.article_id AND articles.is_publish = ?", true).
		Group("tags.id").
		Having("COUNT(DISTINCT articles.id) > 0").
		Count(&count).Error

	return count, err
}

// ============================
// 文章贡献统计
// ============================

// GetArticleContribution 获取文章贡献数据
func (r *StatsRepository) GetArticleContribution(year *int, month *int) ([]dto.ArticleContribution, error) {
	var results []dto.ArticleContribution

	query := r.db.Model(&model.Article{}).
		Select("TO_CHAR(DATE(publish_time), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("is_publish = ?", true)

	// 根据参数构建查询条件
	switch {
	case year != nil && month != nil:
		// 查询指定年月的数据
		startDate := time.Date(*year, time.Month(*month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second) // 月末最后一秒
		query = query.Where("publish_time >= ? AND publish_time <= ?", startDate, endDate)
	case year != nil:
		// 查询指定年份全年的数据
		startDate := time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(*year, 12, 31, 23, 59, 59, 0, time.UTC)
		query = query.Where("publish_time >= ? AND publish_time <= ?", startDate, endDate)
	default:
		// 默认查询过去一年的数据
		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		query = query.Where("publish_time >= ?", oneYearAgo)
	}

	// 执行查询
	err := query.
		Group("DATE(publish_time)").
		Order("date ASC").
		Scan(&results).Error

	return results, err
}

// ============ 网站数据 ============

// GetTotalWords 获取所有已发布文章的总字数
func (r *StatsRepository) GetTotalWords() (int, error) {
	var articles []model.Article
	err := r.db.Model(&model.Article{}).
		Select("content").
		Where("is_publish = ?", true).
		Find(&articles).Error

	if err != nil {
		return 0, err
	}

	totalWords := 0
	for _, article := range articles {
		// 统计字符数（包括中文和英文）
		totalWords += len([]rune(article.Content))
	}

	return totalWords, nil
}

// GetOnlineUsers 获取当前在线用户数（最近5分钟有访问记录）
func (r *StatsRepository) GetOnlineUsers() (int64, error) {
	var count int64
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	err := r.db.Model(&model.Visit{}).
		Where("created_at >= ?", fiveMinutesAgo).
		Distinct("visitor_id").
		Count(&count).Error

	return count, err
}

// ============================
// 归档统计
// ============================

// GetArchives 获取文章归档数据（按年月分组）
func (r *StatsRepository) GetArchives() ([]dto.ArchiveItem, error) {
	var results []dto.ArchiveItem

	// 查询已发布文章按年月分组统计
	err := r.db.Model(&model.Article{}).
		Select("TO_CHAR(publish_time, 'YYYY') as year, TO_CHAR(publish_time, 'MM') as month, COUNT(*) as count").
		Where("is_publish = ? AND publish_time IS NOT NULL", true).
		Group("year, month").
		Order("year DESC, month DESC").
		Scan(&results).Error

	return results, err
}

// ============================
// 访问日志查询
// ============================

// GetVisitLogs 获取访问日志列表
func (r *StatsRepository) GetVisitLogs(req *dto.GetVisitLogsRequest) ([]model.Visit, int64, error) {
	var visits []model.Visit
	var total int64

	query := r.db.Model(&model.Visit{})

	// 关键词搜索（页面URL）
	if req.Keyword != "" {
		query = query.Where("page_url ILIKE ?", "%"+req.Keyword+"%")
	}

	// 访客ID筛选
	if req.VisitorID != "" {
		query = query.Where("visitor_id = ?", req.VisitorID)
	}

	// IP地址筛选
	if req.IP != "" {
		query = query.Where("ip ILIKE ?", "%"+req.IP+"%")
	}

	// 排除IP地址筛选（支持多个IP，用逗号分隔）
	if req.ExcludeIPs != "" {
		// 解析逗号分隔的IP列表
		ips := strings.Split(req.ExcludeIPs, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				query = query.Where("ip NOT ILIKE ?", "%"+ip+"%")
			}
		}
	}

	// 地理位置筛选
	if req.Location != "" {
		query = query.Where("location ILIKE ?", "%"+req.Location+"%")
	}

	// 浏览器筛选（模糊匹配）
	if req.Browser != "" {
		query = query.Where("browser ILIKE ?", "%"+req.Browser+"%")
	}

	// 操作系统筛选（模糊匹配）
	if req.OS != "" {
		query = query.Where("os ILIKE ?", "%"+req.OS+"%")
	}

	// 时间范围筛选
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime+" 23:59:59")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&visits).Error

	return visits, total, err
}
