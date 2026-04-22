package dto

// ============ 前台统计请求 ============

// CollectRequest 统计数据采集请求
type CollectRequest struct {
	URL       string                 `json:"url"`        // 访问URL
	Hostname  string                 `json:"hostname"`   // 主机名
	Referrer  string                 `json:"referrer"`   // 来源页面
	Language  string                 `json:"language"`   // 浏览器语言
	Screen    string                 `json:"screen"`     // 屏幕分辨率
	Title     string                 `json:"title"`      // 页面标题
	Timestamp int64                  `json:"timestamp"`  // 时间戳
	Type      string                 `json:"type"`       // 类型: pageview/event/duration
	Duration  int64                  `json:"duration"`   // 停留时长（秒）
	EventName string                 `json:"event_name"` // 事件名称
	EventData map[string]interface{} `json:"event_data"` // 事件数据
	ArticleID *uint                  `json:"article_id"` // 文章ID（可选）
}

// ============ 前台统计响应 ============

// SiteStatsResponse 前台网站统计信息
type SiteStatsResponse struct {
	TotalWords         string `json:"total_words"`         // 本站总字数
	TotalVisitors      int64  `json:"total_visitors"`      // 本站访客数
	TotalPageViews     int64  `json:"total_page_views"`    // 本站总浏览量
	OnlineUsers        int64  `json:"online_users"`        // 当前在线人数
	TodayVisitors      int64  `json:"today_visitors"`      // 今日访客
	TodayPageviews     int64  `json:"today_pageviews"`     // 今日访问
	YesterdayVisitors  int64  `json:"yesterday_visitors"`  // 昨日访客
	YesterdayPageviews int64  `json:"yesterday_pageviews"` // 昨日访问
	MonthPageviews     int64  `json:"month_pageviews"`     // 本月访问
	TotalArticles      int64  `json:"total_articles"`      // 已发布文章数
	TotalComments      int64  `json:"total_comments"`      // 公开可见评论数
	TotalFriends       int64  `json:"total_friends"`       // 友链数
	TotalMoments       int64  `json:"total_moments"`       // 动态数
	TotalCategories    int64  `json:"total_categories"`    // 已发布文章分类数
	TotalTags          int64  `json:"total_tags"`          // 已发布文章标签数
}

// ArchiveItem 归档数据项
type ArchiveItem struct {
	Year  string `json:"year"`  // 年份
	Month string `json:"month"` // 月份
	Count int    `json:"count"` // 文章数量
}

// ArchivesResponse 归档数据响应
type ArchivesResponse struct {
	Archives []ArchiveItem `json:"archives"` // 归档列表
}

// ============ 后台统计请求 ============

// GetTrendRequest 获取趋势数据请求
type GetTrendRequest struct {
	StartDate string `form:"start_date" binding:"required"` // 开始日期 YYYY-MM-DD
	EndDate   string `form:"end_date" binding:"required"`   // 结束日期 YYYY-MM-DD
	Type      string `form:"type"`                          // 类型: daily/weekly/monthly
}

// GetContributionRequest 获取文章贡献数据请求
type GetContributionRequest struct {
	Year  *int `form:"year"`  // 年份（可选）
	Month *int `form:"month"` // 月份 1-12（可选，需配合year使用）
}

// ============ 后台统计响应 ============

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	// 基础统计
	TotalArticles int64 `json:"total_articles"` // 已发布文章数
	TotalFriends  int64 `json:"total_friends"`  // 友链数
	TotalMoments  int64 `json:"total_moments"`  // 动态数
	TotalViews    int64 `json:"total_views"`    // 总浏览量(PV)
	TotalVisitors int64 `json:"total_visitors"` // 总访客量(UV)
	TotalComments int64 `json:"total_comments"` // 总评论数
	TotalUsers    int64 `json:"total_users"`    // 总用户数

	// 今日统计
	TodayViews    int64 `json:"today_views"`    // 今日浏览量
	TodayVisitors int64 `json:"today_visitors"` // 今日访客量
	TodayComments int64 `json:"today_comments"` // 今日评论数
	TodayUsers    int64 `json:"today_users"`    // 今日注册用户

	// 趋势对比（较昨日）
	ViewsGrowth    float64 `json:"views_growth"`    // 浏览量增长率
	VisitorsGrowth float64 `json:"visitors_growth"` // 访客量增长率
	CommentsGrowth float64 `json:"comments_growth"` // 评论数增长率
	UsersGrowth    float64 `json:"users_growth"`    // 用户增长率
}

// TrendData 趋势图表数据
type TrendData struct {
	Date    string `json:"date"`     // 日期
	PVCount int    `json:"pv_count"` // 浏览量
	UVCount int    `json:"uv_count"` // 访客量
}

// CategoryStats 分类统计数据
type CategoryStats struct {
	Name  string `json:"name"`  // 分类名称
	Count int    `json:"count"` // 文章数量
}

// TagStats 标签统计数据
type TagStats struct {
	Name  string `json:"name"`  // 标签名称
	Count int    `json:"count"` // 文章数量
}

// ArticleContribution 文章贡献数据
type ArticleContribution struct {
	Date  string `json:"date"`  // 日期 YYYY-MM-DD
	Count int    `json:"count"` // 文章发布数量
}

// ============ 访问日志请求 ============

// GetVisitLogsRequest 获取访问日志请求
type GetVisitLogsRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`              // 页码
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量
	Keyword    string `form:"keyword"`                                     // 搜索关键词（页面URL）
	VisitorID  string `form:"visitor_id"`                                  // 访客ID
	IP         string `form:"ip"`                                          // IP地址
	ExcludeIPs string `form:"exclude_ips"`                                 // 排除的IP地址，多个用逗号分隔
	Location   string `form:"location"`                                    // 地理位置
	Browser    string `form:"browser"`                                     // 浏览器
	OS         string `form:"os"`                                          // 操作系统
	StartTime  string `form:"start_time"`                                  // 开始时间（格式：2006-01-02）
	EndTime    string `form:"end_time"`                                    // 结束时间（格式：2006-01-02）
}

// ============ 访问日志响应 ============

// VisitLogItem 访问日志项
type VisitLogItem struct {
	ID        uint   `json:"id"`         // 记录ID
	VisitorID string `json:"visitor_id"` // 访客唯一标识
	IP        string `json:"ip"`         // 访客IP
	PageURL   string `json:"page_url"`   // 访问页面URL
	UserAgent string `json:"user_agent"` // 浏览器UA
	Location  string `json:"location"`   // 地理位置
	Browser   string `json:"browser"`    // 浏览器
	OS        string `json:"os"`         // 操作系统
	Referer   string `json:"referer"`    // 来源页面
	CreatedAt string `json:"created_at"` // 创建时间
}
