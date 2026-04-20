/**
 * 网站统计信息
 */
export interface SiteStats {
  total_words: string; // 本站总字数
  total_visitors: number; // 本站访客数
  total_page_views: number; // 本站总浏览量
  online_users: number; // 当前在线人数
  total_articles: number; // 已发布文章数
  total_comments: number; // 公开可见评论数
  total_friends: number; // 友链数
  total_moments: number; // 动态数
  total_categories: number; // 已发布文章分类数
  total_tags: number; // 已发布文章标签数

  // 详细访问统计
  today_visitors: number; // 今日访客数（UV）
  today_pageviews: number; // 今日访问量（PV）
  yesterday_visitors: number; // 昨日访客数（UV）
  yesterday_pageviews: number; // 昨日访问量（PV）
  month_pageviews: number; // 本月访问量（PV）
}

/**
 * 归档统计项
 */
export interface ArchiveItem {
  year: string;
  month: string;
  count: number;
}

/**
 * 归档统计数据
 */
export interface ArchiveStats {
  archives: ArchiveItem[];
}
