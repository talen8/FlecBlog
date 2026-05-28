/**
 * 微信小程序应用入口
 */

import type { SiteBasicConfig, BlogConfig, OAuthConfig } from './types';
import { APP_CONFIG } from './config';
import { getStorage, setStorage } from './utils/storage';
import { initInterceptors } from './utils/request';
import { getSiteBasicConfig, getBlogConfig, getOAuthConfig } from './api/setting';

// ==================== 全局 App 类型 ====================

export interface IAppOption {
  globalData: {
    siteConfig: Partial<SiteBasicConfig>;
    blogConfig: Partial<BlogConfig>;
    oauthConfig: Partial<OAuthConfig>;
    version?: string;
  };
  towxml: (str: string, type: 'markdown' | 'html', option?: {
    theme?: 'light' | 'dark';
    events?: { tap?: (e: { target: { dataset: { href?: string } } }) => void };
  }) => { theme: 'light' | 'dark'; _e: Record<string, unknown> };
  updateSiteConfig: (config: Partial<SiteBasicConfig>) => void;
  updateBlogConfig: (config: Partial<BlogConfig>) => void;
  updateOAuthConfig: (config: Partial<OAuthConfig>) => void;
  loadConfig: () => Promise<void>;
  checkUpdate: () => void;
  reportError: (error: string) => void;
}

// ==================== App 实例 ====================

App<IAppOption>({
  towxml: require('./towxml/index'),

  globalData: {
    siteConfig: getStorage<Partial<SiteBasicConfig>>('site_config') || {},
    blogConfig: getStorage<Partial<BlogConfig>>('blog_config') || {},
    oauthConfig: getStorage<Partial<OAuthConfig>>('oauth_config') || {},
    version: wx.getAccountInfoSync().miniProgram.version,
  },

  onLaunch(_options: WechatMiniprogram.App.LaunchShowOption) {
    initInterceptors();
    this.loadConfig();
    this.checkUpdate();
  },

  onShow() {
    this.checkUpdate();
  },

  onHide() {},

  onError(error: string) {
    console.error(`${APP_CONFIG.APP_NAME} 脚本错误:`, error);
    this.reportError(error);
  },

  onPageNotFound() {
    wx.redirectTo({
      url: '/pages/index/index',
      fail: () => {
        wx.switchTab({ url: '/pages/index/index' });
      },
    });
  },

  checkUpdate() {
    if (typeof wx.getUpdateManager !== 'function') return;

    const updateManager = wx.getUpdateManager();

    updateManager.onCheckForUpdate(() => {});

    updateManager.onUpdateReady(() => {
      wx.showModal({
        title: '更新提示',
        content: '新版本已准备好，是否重启应用以获取最新体验？',
        confirmColor: '#0052d9',
        success: (res) => {
          if (res.confirm) updateManager.applyUpdate();
        },
      });
    });

    updateManager.onUpdateFailed(() => {
      wx.showModal({
        title: '更新失败',
        content: '新版本下载失败，请检查网络设置后重试',
        confirmColor: '#0052d9',
        showCancel: false,
      });
    });
  },

  reportError(error: string) {
    try {
      const logs = getStorage<{ type: string; message: string; time: number }[]>('error_logs') || [];
      logs.unshift({ type: 'app_error', message: error, time: Date.now() });
      if (logs.length > 20) logs.splice(20);
      setStorage('error_logs', logs);
    } catch {}
  },

  updateSiteConfig(config) {
    this.globalData.siteConfig = { ...this.globalData.siteConfig, ...config };
    setStorage('site_config', this.globalData.siteConfig);
  },

  updateBlogConfig(config) {
    this.globalData.blogConfig = { ...this.globalData.blogConfig, ...config };
    setStorage('blog_config', this.globalData.blogConfig);
  },

  updateOAuthConfig(config) {
    this.globalData.oauthConfig = { ...this.globalData.oauthConfig, ...config };
    setStorage('oauth_config', this.globalData.oauthConfig);
  },

  async loadConfig() {
    try {
      const [basicSettings, blogSettings, oauthSettings] = await Promise.all([
        getSiteBasicConfig().catch(() => null),
        getBlogConfig().catch(() => null),
        getOAuthConfig().catch(() => null),
      ]);

      if (basicSettings) {
        this.updateSiteConfig({
          author: basicSettings['basic.author'] || '',
          author_avatar: basicSettings['basic.author_avatar'] || '',
          blog_url: basicSettings['basic.blog_url'] || '',
        });
      }

      if (blogSettings) {
        this.updateBlogConfig({
          title: blogSettings['blog.title'] || '',
          established: blogSettings['blog.established'] || '',
          about_describe: blogSettings['blog.about_describe'] || '',
          about_describe_tips: blogSettings['blog.about_describe_tips'] || '',
          about_profile: blogSettings['blog.about_profile'] || '[]',
          about_personality: blogSettings['blog.about_personality'] || '',
          about_motto_main: blogSettings['blog.about_motto_main'] || '[]',
          about_motto_sub: blogSettings['blog.about_motto_sub'] || '',
          about_versions: blogSettings['blog.about_versions'] || '[]',
          about_unions: blogSettings['blog.about_unions'] || '[]',
          about_story: blogSettings['blog.about_story'] || '',
        });
      }

      if (oauthSettings) {
        this.updateOAuthConfig(oauthSettings as Partial<OAuthConfig>);
      }
    } catch (error) {
      console.error('加载配置失败:', error);
    }
  },
});
