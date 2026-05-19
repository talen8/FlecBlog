/**
 * 首页
 * 展示文章列表，支持下拉刷新和上拉加载更多
 */

import type { ArticleListItem } from '../../types';
import type { IAppOption } from '../../app';
import { getArticles } from '../../api/article';
import { formatRelativeTime } from '../../utils/format';

const app = getApp<IAppOption>();

Page({
  data: {
    topArticles: [] as ArticleListItem[],
    articles: [] as ArticleListItem[],
    loading: true,
    hasMore: true,
    page: 1,
    error: '',
    isEmpty: false,
  },

  onLoad() {
    wx.setNavigationBarTitle({ title: app.globalData.siteConfig.site_name || 'FlecBlog' });
    this.loadData(true);
  },

  async loadData(reset = false) {
    const { page, loading } = this.data;
    if (!reset && loading) return;

    const currentPage = reset ? 1 : page;
    this.setData({
      loading: !reset,
      error: '',
    });

    try {
      const result = await getArticles({ page: currentPage, page_size: 10 });
      const { list, total } = result;

      const items = list.map((item: ArticleListItem) => ({
        ...item,
        formattedTime: formatRelativeTime(item.publish_time),
        slug: (item.url || '').split('/').pop(),
      }));

      const topItems = reset ? items.filter(item => item.is_top) : [];
      const normalItems = reset
        ? items.filter(item => !item.is_top)
        : items;

      const newArticles = reset ? normalItems : [...this.data.articles, ...normalItems];
      const hasMore = reset
        ? newArticles.length + topItems.length < total
        : newArticles.length + this.data.topArticles.length < total;

      this.setData({
        topArticles: topItems.length ? topItems : this.data.topArticles,
        articles: newArticles,
        page: currentPage + 1,
        hasMore,
        isEmpty: newArticles.length === 0 && topItems.length === 0,
        loading: false,
      });
    } catch (error) {
      const errorMsg = error instanceof Error ? '加载失败' : '请求失败';
      this.setData({ loading: false, error: errorMsg });
      wx.showToast({ title: errorMsg, icon: 'none' });
    }
  },

  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadData(false);
    }
  },

  onRefresh() {
    this.loadData(true);
  },

  onPullDownRefresh() {
    this.loadData(true).finally(() => {
      wx.stopPullDownRefresh();
    });
  },

  onCategoryTap(e: any) {
    const { category } = e.currentTarget.dataset;
    if (category?.slug) {
      wx.showToast({ title: `分类: ${category.name}`, icon: 'none' });
    }
  },

  onTagTap(e: any) {
    const { tag } = e.currentTarget.dataset;
    if (tag?.slug) {
      wx.showToast({ title: `标签: ${tag.name}`, icon: 'none' });
    }
  },

  onArticleTap(e: any) {
    const { article } = e.currentTarget.dataset;
    if (article?.slug) {
      wx.navigateTo({ url: `/pages/article/detail/index?slug=${article.slug}` });
    }
  },

  onShareAppMessage() {
    return { title: app.globalData.siteConfig.site_name || 'FlecBlog', path: '/pages/index/index' };
  },

  onShareTimeline() {
    return { title: app.globalData.siteConfig.site_name || 'FlecBlog' };
  },
});
