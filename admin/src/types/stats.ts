// Dashboard 统计数据
export interface DashboardStats {
  total_articles: number;
  total_friends: number;
  total_moments: number;
  total_views: number;
  total_visitors: number;
  total_comments: number;
  total_users: number;
  today_views: number;
  today_visitors: number;
  today_comments: number;
  today_users: number;
  views_growth: number;
  visitors_growth: number;
  comments_growth: number;
  users_growth: number;
}

// 趋势数据项
export interface TrendDataItem {
  date: string;
  pv_count: number;
  uv_count: number;
}

// 趋势查询参数
export interface TrendQuery {
  start_date: string;
  end_date: string;
  type: 'daily' | 'monthly';
}

// 标签统计数据
export interface TagStats {
  name: string;
  count: number;
}

// 文章贡献数据
export interface ArticleContribution {
  date: string;
  count: number;
}

// 文章贡献查询参数
export interface ContributionQuery {
  year?: number; // 年份（可选）
  month?: number; // 月份 1-12（可选）
}

// 分类统计数据
export interface CategoryStats {
  name: string;
  count: number;
}

// 访问日志实体
export interface Visit {
  id: number;
  visitor_id: string;
  ip: string;
  page_url: string;
  user_agent: string;
  location: string;
  browser: string;
  os: string;
  referer: string;
  created_at: string;
}

// 访问日志查询参数
export interface VisitListQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  visitor_id?: string;
  ip?: string;
  location?: string;
  browser?: string;
  os?: string;
  start_time?: string;
  end_time?: string;
}

// 访问日志列表数据
export interface VisitListData {
  list: Visit[];
  total: number;
  page: number;
  page_size: number;
}
