/**
 * 关于页面
 * 展示博主信息、站点统计
 */

import type { SiteStats } from '../../types';
import type { IAppOption } from '../../app';
import { getSiteStats } from '../../api/stats';

const app = getApp<IAppOption>();

function parseJSON<T>(json: string, defaultValue: T): T {
  try {
    return JSON.parse(json) as T;
  } catch {
    return defaultValue;
  }
}

const PERSONALITY_MAP: Record<string, { name: string; color: string }> = {
  INTJ: { name: '建筑师', color: '#885fb8' },
  INTP: { name: '逻辑学家', color: '#885fb8' },
  ENTJ: { name: '指挥官', color: '#885fb8' },
  ENTP: { name: '辩论家', color: '#885fb8' },
  INFJ: { name: '提倡者', color: '#56a178' },
  INFP: { name: '调停者', color: '#56a178' },
  ENFJ: { name: '主人公', color: '#56a178' },
  ENFP: { name: '竞选者', color: '#56a178' },
  ISTJ: { name: '物流师', color: '#4298b4' },
  ISFJ: { name: '守卫者', color: '#4298b4' },
  ESTJ: { name: '总经理', color: '#4298b4' },
  ESFJ: { name: '执政官', color: '#4298b4' },
  ISTP: { name: '鉴赏家', color: '#e4ae3a' },
  ISFP: { name: '探险家', color: '#e4ae3a' },
  ESTP: { name: '企业家', color: '#e4ae3a' },
  ESFP: { name: '表演者', color: '#e4ae3a' },
};

Page({
  data: {
    loading: true,
    author: '',
    authorAvatar: '',
    describe: '',
    describeTips: '',
    profile: [] as Array<{ label: string; value: string; color: string }>,
    personality: { name: '', code: '', color: '' },
    mottoMain: [] as string[],
    mottoSub: '',
    versions: [] as Array<{ name: string; version: string }>,
    story: '',
    runningDays: 0,
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
      const { siteConfig, blogConfig } = app.globalData;

      const profile = parseJSON<Array<{ label: string; value: string; color: string }>>(
        blogConfig.about_profile || '[]', []
      );
      const mottoMain = parseJSON<string[]>(blogConfig.about_motto_main || '[]', []);
      const versions = parseJSON<Array<{ name: string; version: string }>>(
        blogConfig.about_versions || '[]', []
      );

      const personalityCode = (blogConfig.about_personality || '').substring(0, 4).toUpperCase();
      const personalityInfo = PERSONALITY_MAP[personalityCode];

      const established = blogConfig.established || '';
      let runningDays = 0;
      if (established) {
        const d = new Date(established.includes(' ') ? established.replace(' ', 'T') : established);
        if (!isNaN(d.getTime())) {
          runningDays = Math.floor((Date.now() - d.getTime()) / 86400000);
        }
      }

      this.setData({
        author: siteConfig.author || '',
        authorAvatar: siteConfig.author_avatar || '',
        describe: blogConfig.about_describe || '',
        describeTips: blogConfig.about_describe_tips || '',
        profile,
        personality: {
          name: personalityInfo?.name || '',
          code: personalityCode,
          color: personalityInfo?.color || '',
        },
        mottoMain,
        mottoSub: blogConfig.about_motto_sub || '',
        versions,
        story: blogConfig.about_story || '',
        runningDays,
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

  onCopyLink(e: WechatMiniprogram.TapEvent) {
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
    return { title: '关于', path: '/pages/about/about' };
  },

  onShareTimeline() {
    return { title: '关于' };
  },
});
