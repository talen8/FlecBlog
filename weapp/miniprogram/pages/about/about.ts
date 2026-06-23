/**
 * 关于页面
 * 展示博主信息、站点统计
 */

import type { SiteStats } from '../../types';
import type { IAppOption } from '../../app';
import { getSiteStats } from '../../api/stats';

const app = getApp<IAppOption>();

Page({
  data: {
    loading: true,
    author: '',
    authorAvatar: '',
    stats: null as SiteStats | null,
  },

  onLoad() {
    this.loadAboutData();
  },

  onPullDownRefresh() {
    this.loadAboutData().finally(() => {
      wx.stopPullDownRefresh();
    });
  },

  async loadAboutData() {
    this.setData({ loading: true });

    try {
      const { siteConfig } = app.globalData;

      this.setData({
        author: siteConfig.author || '',
        authorAvatar: siteConfig.author_avatar || '',
        loading: false,
      });

      this.loadStats();
    } catch (error) {
      console.error('加载关于页面数据失败:', error);
      this.setData({ loading: false });
    }
  },

  async loadStats() {
    try {
      const stats = await getSiteStats();
      this.setData({ stats });
    } catch (error) {
      console.error('加载统计数据失败:', error);
    }
  },

  onShareAppMessage() {
    return { title: '关于', path: '/pages/about/about' };
  },

  onShareTimeline() {
    return { title: '关于' };
  },
});
