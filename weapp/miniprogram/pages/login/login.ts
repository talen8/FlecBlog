/**
 * 登录页面
 * 支持邮箱密码登录
 */

import { login, wechatLogin } from '../../api/user';
import { APP_CONFIG } from '../../config';
import { setToken } from '../../utils/request';
import { setStorage } from '../../utils/storage';

Page({
  data: {
    appName: APP_CONFIG.APP_NAME,
    email: '',
    password: '',
    loading: false,
    wechatLoading: false,
    emailError: '',
    passwordError: '',
  },

  /**
   * 输入邮箱
   */
  onEmailInput(e: WechatMiniprogram.InputEvent) {
    this.setData({
      email: e.detail.value,
      emailError: '',
    });
  },

  /**
   * 输入密码
   */
  onPasswordInput(e: WechatMiniprogram.InputEvent) {
    this.setData({
      password: e.detail.value,
      passwordError: '',
    });
  },

  /**
   * 验证表单
   */
  validateForm(): boolean {
    const { email, password } = this.data;
    let valid = true;

    if (!email.trim()) {
      this.setData({ emailError: '请输入邮箱' });
      valid = false;
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      this.setData({ emailError: '邮箱格式不正确' });
      valid = false;
    }

    if (!password) {
      this.setData({ passwordError: '请输入密码' });
      valid = false;
    } else if (password.length < 6) {
      this.setData({ passwordError: '密码至少6位' });
      valid = false;
    }

    return valid;
  },

  /**
   * 提交登录
   */
  async handleSubmit() {
    if (!this.validateForm()) return;
    if (this.data.loading) return;

    const { email, password } = this.data;

    this.setData({ loading: true });

    try {
      // 获取微信 code 用于绑定微信身份
      let wechatCode = '';
      try {
        const loginResult = await new Promise<WechatMiniprogram.LoginSuccessCallbackResult>((resolve, reject) => {
          wx.login({ success: resolve, fail: reject });
        });
        wechatCode = loginResult.code;
      } catch {
        // wx.login 失败不影响邮箱登录
      }

      const { access_token, user } = await login({
        email: email.trim(),
        password,
        wechat_code: wechatCode || undefined,
      });

      setToken(access_token);
      setStorage('user_info', user);

      wx.showToast({ title: '登录成功', icon: 'success' });

      setTimeout(() => {
        wx.switchTab({ url: '/pages/profile/profile' });
      }, 1000);
    } catch (error) {
      const err = error as Error;
      wx.showToast({
        title: err.message || '登录失败',
        icon: 'none',
      });
    } finally {
      this.setData({ loading: false });
    }
  },

  /**
   * 微信登录
   */
  async handleWechatLogin() {
    if (this.data.wechatLoading) return;
    this.setData({ wechatLoading: true });
    try {
      const { code } = await new Promise<WechatMiniprogram.LoginSuccessCallbackResult>((resolve, reject) => {
        wx.login({ success: resolve, fail: reject });
      });
      const { access_token, user } = await wechatLogin(code);
      setToken(access_token);
      setStorage('user_info', user);
      wx.showToast({ title: '登录成功', icon: 'success' });
      setTimeout(() => {
        wx.switchTab({ url: '/pages/profile/profile' });
      }, 1000);
    } catch (error) {
      wx.showToast({
        title: (error as Error).message || '登录失败',
        icon: 'none',
      });
    } finally {
      this.setData({ wechatLoading: false });
    }
  },

  onShareAppMessage() {
    return { title: '登录', path: '/pages/login/login' };
  },
});
