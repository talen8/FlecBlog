/**
 * FlecBlog 微信小程序应用入口
 * 应用生命周期管理和全局数据
 */

import type { SiteBasicConfig, BlogConfig } from './types';
import { DEFAULT_SITE_CONFIG, DEFAULT_BLOG_CONFIG, STORAGE_KEYS } from './config';
import { getStorage, setStorage } from './utils/storage';
import { initInterceptors } from './utils/request';

// ==================== Towxml 类型定义 ====================

/**
 * towxml 解析选项
 */
interface TowxmlOption {
  theme?: 'light' | 'dark';
  events?: {
    tap?: (e: { target: { dataset: { href?: string } } }) => void;
  };
}

/**
 * towxml 解析结果
 */
interface TowxmlResult {
  theme: 'light' | 'dark';
  _e: Record<string, unknown>;
}

/**
 * towxml 解析函数类型
 */
type TowxmlParser = (str: string, type: 'markdown' | 'html', option?: TowxmlOption) => TowxmlResult;

// ==================== 全局 App 类型扩展 ====================

export interface IAppOption {
  /** 全局数据 */
  globalData: {
    /** 站点基础配置 */
    siteConfig: Partial<SiteBasicConfig>;
    /** 博客配置 */
    blogConfig: Partial<BlogConfig>;
    /** 系统信息 */
    systemInfo?: ReturnType<typeof wx.getSystemInfoSync>;
    /** 小程序版本 */
    version?: string;
  };
  /** towxml 解析器实例 */
  towxml: TowxmlParser;
  /** 更新站点配置 */
  updateSiteConfig: (config: Partial<SiteBasicConfig>) => void;
  /** 更新博客配置 */
  updateBlogConfig: (config: Partial<BlogConfig>) => void;
  /** 检查更新 */
  checkUpdate: () => void;
  /** 获取系统信息 */
  getSystemInfo: () => void;
  /** 上报错误日志 */
  reportError: (error: string) => void;
}

// ==================== App 实例 ====================

App<IAppOption>({
  /** towxml 解析器实例 */
  towxml: require('./towxml/index'),

  /** 全局数据对象 */
  globalData: {
    siteConfig: getStorage<Partial<SiteBasicConfig>>(STORAGE_KEYS.SITE_CONFIG) || DEFAULT_SITE_CONFIG,
    blogConfig: getStorage<Partial<BlogConfig>>(STORAGE_KEYS.BLOG_CONFIG) || DEFAULT_BLOG_CONFIG,
    systemInfo: undefined,
    version: wx.getAccountInfoSync().miniProgram.version,
  },

  /**
   * 应用初始化完成时触发（全局只触发一次）
   */
  onLaunch() {
    console.log('FlecBlog 小程序启动');

    // 初始化请求拦截器
    initInterceptors();

    // 获取系统信息
    this.getSystemInfo();

    // 检查更新
    this.checkUpdate();
  },

  /**
   * 应用启动或从后台进入前台时触发
   */
  onShow() {
    console.log('FlecBlog 小程序显示');

    // 每次从后台回来都检查更新
    this.checkUpdate();
  },

  /**
   * 应用从前台进入后台时触发
   */
  onHide() {
    console.log('FlecBlog 小程序隐藏');
  },

  /**
   * 应用发生脚本错误时触发
   * @param error - 错误信息
   */
  onError(error: string) {
    // 记录错误日志（可接入日志服务）
    console.error('FlecBlog 脚本错误:', error);

    // 错误上报（可根据需要接入服务）
    this.reportError(error);
  },

  /**
   * 页面不存在时触发
   * @param options - 页面信息
   */
  onPageNotFound(options: { path: string; query: Record<string, string>; isEntryPage: boolean }) {
    console.warn('页面不存在:', options);

    // 重定向到首页
    wx.redirectTo({
      url: '/pages/index/index',
      fail: () => {
        // 如果重定向失败，返回小程序首页
        wx.switchTab({
          url: '/pages/index/index',
        });
      },
    });
  },

  /**
   * 获取系统信息
   */
  getSystemInfo() {
    try {
      const systemInfo = wx.getSystemInfoSync();
      this.globalData.systemInfo = systemInfo;
      console.log('系统信息:', systemInfo);
    } catch (error) {
      console.error('获取系统信息失败:', error);
    }
  },

  /**
   * 检查小程序更新
   * 支持静默更新和用户确认后更新
   */
  checkUpdate() {
    // 基础库 2.8.3 开始支持
    if (typeof wx.getUpdateManager !== 'function') {
      console.warn('当前微信版本不支持自动更新，请升级到最新版本');
      return;
    }

    const updateManager = wx.getUpdateManager();

    // 1. 检查是否有新版本
    updateManager.onCheckForUpdate((res) => {
      if (res.hasUpdate) {
        console.log('有新版本可用，版本号:', this.globalData.version);
      }
    });

    // 2. 新版本下载完成
    updateManager.onUpdateReady(() => {
      wx.showModal({
        title: '更新提示',
        content: '新版本已准备好，是否重启应用以获取最新体验？',
        confirmColor: '#0052d9',
        success: (res) => {
          if (res.confirm) {
            // 强制当前页面栈使用新版本
            updateManager.applyUpdate();
          }
        },
      });
    });

    // 3. 新版本下载失败
    updateManager.onUpdateFailed(() => {
      wx.showModal({
        title: '更新失败',
        content: '新版本下载失败，请检查网络设置后重试',
        confirmColor: '#0052d9',
        showCancel: false,
      });
    });
  },

  /**
   * 上报错误日志
   * @param error - 错误信息
   */
  reportError(error: string) {
    // TODO: 可接入日志服务，如 Sentry、神策等
    // 这里仅做本地记录
    try {
      const errorLog = {
        type: 'app_error',
        message: error,
        version: this.globalData.version,
        time: Date.now(),
        userAgent: this.globalData.systemInfo?.platform || 'unknown',
      };

      // 保存到本地（生产环境可上报到服务器）
      const logs = getStorage<typeof errorLog[]>('error_logs') || [];
      logs.unshift(errorLog);

      // 只保留最近 10 条
      if (logs.length > 10) {
        logs.splice(10);
      }

      setStorage('error_logs', logs);
    } catch {
      // 忽略存储错误
    }
  },

  /**
   * 更新站点配置
   * @param config - 新的站点配置
   */
  updateSiteConfig(config: Partial<SiteBasicConfig>) {
    this.globalData.siteConfig = {
      ...this.globalData.siteConfig,
      ...config,
    };
    setStorage(STORAGE_KEYS.SITE_CONFIG, this.globalData.siteConfig);
  },

  /**
   * 更新博客配置
   * @param config - 新的博客配置
   */
  updateBlogConfig(config: Partial<BlogConfig>) {
    this.globalData.blogConfig = {
      ...this.globalData.blogConfig,
      ...config,
    };
    setStorage(STORAGE_KEYS.BLOG_CONFIG, this.globalData.blogConfig);
  },
});
