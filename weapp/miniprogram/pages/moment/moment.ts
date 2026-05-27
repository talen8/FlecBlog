/**
 * 动态页面
 * 展示动态列表，支持下拉刷新和上拉加载更多
 */

import type { MomentListItem } from '../../types';
import { getMoments } from '../../api/moment';
import { formatRelativeTime } from '../../utils/format';

interface MomentDisplay extends MomentListItem {
  formattedTime: string;
}

Page({
  data: {
    moments: [] as MomentDisplay[],
    loading: true,
    hasMore: true,
    page: 1,
    error: '',
    isEmpty: false,
  },

  onLoad() {
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
      const result = await getMoments({ page: currentPage, page_size: 10 });
      const { list, total } = result;

      const items: MomentDisplay[] = list.map((item: MomentListItem) => ({
        ...item,
        formattedTime: formatRelativeTime(item.publish_time),
      }));

      const newMoments = reset ? items : [...this.data.moments, ...items];
      const hasMore = newMoments.length < total;

      this.setData({
        moments: newMoments,
        page: currentPage + 1,
        hasMore,
        isEmpty: newMoments.length === 0,
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

  onPullDownRefresh() {
    this.loadData(true).finally(() => {
      wx.stopPullDownRefresh();
    });
  },

  onRefresh() {
    this.loadData(true);
  },

  onImageTap(e: any) {
    const { images, current } = e.currentTarget.dataset;
    if (images && images.length > 0) {
      wx.previewImage({
        current: current || images[0],
        urls: images,
      });
    }
  },

  onCopyUrl(e: any) {
    const { url } = e.currentTarget.dataset;
    if (url) {
      wx.setClipboardData({
        data: url,
        success: () => {
          wx.showToast({ title: '链接已复制', icon: 'success' });
        },
      });
    }
  },

  onShareAppMessage() {
    return { title: '动态' };
  },

  onShareTimeline() {
    return { title: '动态' };
  },
});
