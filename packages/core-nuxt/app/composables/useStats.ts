import type { SiteStats } from '../../types/stats';
import { getSiteStats, getArchiveStats } from './api/stats';

const defaultSiteStats: SiteStats = {
  total_words: '0',
  total_visitors: 0,
  total_page_views: 0,
  online_users: 0,
  total_articles: 0,
  total_comments: 0,
  total_friends: 0,
  total_moments: 0,
  total_categories: 0,
  total_tags: 0,
  today_visitors: 0,
  today_pageviews: 0,
  yesterday_visitors: 0,
  yesterday_pageviews: 0,
  month_pageviews: 0,
};

/**
 * 获取归档统计（支持 SSR）
 * @returns data - 归档数据，refresh - 刷新方法
 */
export function useArchives() {
  return useAsyncData('archives', async () => {
    const { archives } = await getArchiveStats();
    return archives;
  });
}

/**
 * 获取站点统计信息（文章数/访客/评论等，支持 SSR）
 * @returns siteStats - 站点统计数据（含默认值兜底）
 */
export function useStats() {
  const { data } = useAsyncData('stats-fetch', () => getSiteStats());

  return {
    siteStats: computed(() => data.value ?? defaultSiteStats),
  };
}
