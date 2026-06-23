/**
 * 我的页面
 * 用户个人中心，包含用户信息、设置等功能
 */

import type { UserInfo, UserRole, OAuthConfig } from '../../types';
import type { IAppOption } from '../../app';
import { APP_CONFIG } from '../../config';
import { getUserProfile, logout as apiLogout } from '../../api/user';
import { getOAuthConfig } from '../../api/setting';
import { getToken, clearAuth, logout } from '../../utils/request';
import { getStorage, setStorage } from '../../utils/storage';

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
    wechat: '微信',
  };
  return names[provider] || provider;
}

Page({
  data: {
    appName: APP_CONFIG.APP_NAME,
    isLoggedIn: false,
    userInfo: null as UserInfo | null,
    roleName: '',
    loginMethods: [] as { name: string; enabled: boolean }[],
    loading: false,
    oauthConfig: null as OAuthConfig | null,
    version: '',
  },

  onLoad() {
    this.setData({ version: app.globalData.version || '1.0.0' });
    this.loadOAuthConfig();
    this.checkLoginStatus();
  },

  onShow() {
    this.checkLoginStatus();
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
   * 加载 OAuth 配置
   * 从 globalData 中获取，如果不存在则从 API 获取
   */
  async loadOAuthConfig() {
    try {
      // 优先使用 globalData 中的配置
      let oauthConfig: Partial<OAuthConfig> | null = app.globalData.oauthConfig;
      if (!oauthConfig || Object.keys(oauthConfig).length === 0) {
        oauthConfig = await getOAuthConfig().catch(() => null) as Partial<OAuthConfig> | null;
      }

      if (oauthConfig) {
        this.setData({ oauthConfig: oauthConfig as OAuthConfig });

        if (this.data.isLoggedIn && this.data.userInfo) {
          this.setData({
            loginMethods: this.buildLoginMethods(this.data.userInfo),
          });
        }
      }
    } catch (error) {
      console.error('加载 OAuth 配置失败:', error);
    }
  },

  /**
   * 构建登录方式列表
   * 只显示在管理端开启了的 OAuth 提供商
   */
  buildLoginMethods(user: UserInfo): { name: string; enabled: boolean }[] {
    const methods: { name: string; enabled: boolean }[] = [];

    if (user.has_password) {
      methods.push({ name: '邮箱', enabled: true });
    }

    const oauthConfig = this.data.oauthConfig;
    if (!oauthConfig) {
      // 如果没有获取到配置，只显示用户已绑定的方式
      const oauthProviders = ['github', 'google', 'qq', 'microsoft', 'oidc', 'wechat'];
      oauthProviders.forEach((provider) => {
        if (user.linked_oauths?.includes(provider)) {
          methods.push({ name: getProviderName(provider), enabled: true });
        }
      });
      return methods;
    }

    // 只显示在管理端开启了的 OAuth 提供商
    const providerConfigMap: Record<string, string> = {
      github: 'github.enabled',
      google: 'google.enabled',
      qq: 'qq.enabled',
      microsoft: 'microsoft.enabled',
      oidc: 'oidc.enabled',
      wechat: 'wechat.enabled',
    };

    Object.entries(providerConfigMap).forEach(([provider, configKey]) => {
      const isEnabled = oauthConfig[configKey as keyof OAuthConfig] === 'true';
      if (isEnabled && user.linked_oauths?.includes(provider)) {
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
      } else {
        // API 请求失败（可能是 token 过期），清除登录状态
        clearAuth();
        this.setData({
          isLoggedIn: false,
          userInfo: null,
          roleName: '',
          loginMethods: [],
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
            await apiLogout().catch(() => {});
          } finally {
            logout();

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
