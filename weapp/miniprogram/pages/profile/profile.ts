/**
 * 我的页面
 * 用户个人中心，包含用户信息、设置等功能
 */

import type { UserInfo, UserRole } from '../../types';
import type { IAppOption } from '../../app';
import { getUserProfile, logout } from '../../api/user';
import { getToken, clearToken } from '../../utils/request';
import { getStorage, setStorage, removeStorage } from '../../utils/storage';

const app = getApp<IAppOption>();

/**
 * 获取角色名称
 */
function getRoleName(role: UserRole): string {
  const roleMap: Record<UserRole, string> = {
    super_admin: '超级管理员',
    admin: '管理员',
    user: '普通用户',
  };
  return roleMap[role] || role;
}

/**
 * 获取 OAuth 提供商名称
 */
function getProviderName(provider: string): string {
  const names: Record<string, string> = {
    github: 'GitHub',
    google: 'Google',
    qq: 'QQ',
    microsoft: 'Microsoft',
    oidc: 'OIDC',
  };
  return names[provider] || provider;
}

Page({
  data: {
    isLoggedIn: false,
    userInfo: null as UserInfo | null,
    roleName: '',
    loginMethods: [] as { name: string; enabled: boolean }[],
    cacheSize: '0 KB',
    loading: false,
  },

  onLoad() {
    this.checkLoginStatus();
  },

  onShow() {
    if (this.data.isLoggedIn) {
      this.loadUserData();
    }
    this.calculateCacheSize();
  },

  /**
   * 检查登录状态
   */
  checkLoginStatus() {
    const token = getToken();
    const cachedUser = getStorage<UserInfo>('user_info');

    if (token && cachedUser) {
      this.setData({
        isLoggedIn: true,
        userInfo: cachedUser,
        roleName: getRoleName(cachedUser.role),
        loginMethods: this.buildLoginMethods(cachedUser),
      });
      this.loadUserData();
    } else {
      this.setData({
        isLoggedIn: false,
        userInfo: null,
        roleName: '',
        loginMethods: [],
      });
    }
  },

  /**
   * 构建登录方式列表
   */
  buildLoginMethods(user: UserInfo): { name: string; enabled: boolean }[] {
    const methods: { name: string; enabled: boolean }[] = [];

    if (user.has_password) {
      methods.push({ name: '密码', enabled: true });
    }

    const oauthProviders = ['github', 'google', 'qq', 'microsoft', 'oidc'];
    oauthProviders.forEach((provider) => {
      if (user.linked_oauths?.includes(provider)) {
        methods.push({ name: getProviderName(provider), enabled: true });
      }
    });

    return methods;
  },

  /**
   * 加载用户数据
   */
  async loadUserData() {
    if (this.data.loading) return;

    this.setData({ loading: true });

    try {
      const userInfo = await getUserProfile().catch(() => null);

      if (userInfo) {
        setStorage('user_info', userInfo);
        this.setData({
          userInfo,
          roleName: getRoleName(userInfo.role),
          loginMethods: this.buildLoginMethods(userInfo),
        });
      }
    } catch (error) {
      console.error('加载用户数据失败:', error);
    } finally {
      this.setData({ loading: false });
    }
  },

  /**
   * 退出登录
   */
  async handleLogout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      confirmColor: '#0052d9',
      success: async (res) => {
        if (res.confirm) {
          try {
            await logout().catch(() => {});
          } finally {
            clearToken();
            removeStorage('user_info');

            this.setData({
              isLoggedIn: false,
              userInfo: null,
              roleName: '',
              loginMethods: [],
            });

            wx.showToast({ title: '已退出登录', icon: 'success' });
          }
        }
      },
    });
  },

  /**
   * 清除缓存
   */
  handleClearCache() {
    wx.showModal({
      title: '提示',
      content: '确定要清除缓存吗？',
      confirmColor: '#0052d9',
      success: (res) => {
        if (res.confirm) {
          try {
            const keys = wx.getStorageInfoSync().keys;
            const keepKeys = ['access_token', 'user_info', 'site_config', 'blog_config'];

            keys.forEach((key) => {
              if (!keepKeys.includes(key)) {
                wx.removeStorageSync(key);
              }
            });

            this.calculateCacheSize();
            wx.showToast({ title: '清除成功', icon: 'success' });
          } catch (error) {
            console.error('清除缓存失败:', error);
            wx.showToast({ title: '清除失败', icon: 'none' });
          }
        }
      },
    });
  },

  /**
   * 计算缓存大小
   */
  calculateCacheSize() {
    try {
      const info = wx.getStorageInfoSync();
      const sizeKB = info.currentSize;
      let sizeStr: string;

      if (sizeKB < 1024) {
        sizeStr = `${sizeKB} KB`;
      } else {
        sizeStr = `${(sizeKB / 1024).toFixed(2)} MB`;
      }

      this.setData({ cacheSize: sizeStr });
    } catch {
      this.setData({ cacheSize: '0 KB' });
    }
  },

  /**
   * 跳转到关于页面
   */
  goToAbout() {
    wx.switchTab({ url: '/pages/about/about' });
  },

  /**
   * 跳转到登录页面
   */
  goToLogin() {
    wx.navigateTo({ url: '/pages/login/login' });
  },

  onShareAppMessage() {
    const title = this.data.userInfo
      ? `${this.data.userInfo.nickname}的主页`
      : '我的';
    return { title, path: '/pages/profile/profile' };
  },

  onShareTimeline() {
    const title = this.data.userInfo
      ? `${this.data.userInfo.nickname}的主页`
      : '我的';
    return { title };
  },
});
