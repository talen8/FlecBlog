/**
 * 登录页面
 * 支持邮箱密码登录
 */

import type { UserInfo } from '../../types';
import type { IAppOption } from '../../app';
import { login } from '../../api/user';
import { setToken } from '../../utils/request';
import { setStorage } from '../../utils/storage';

const app = getApp<IAppOption>();

Page({
  data: {
    email: '',
    password: '',
    loading: false,
    emailError: '',
    passwordError: '',
  },

  /**
   * 输入邮箱
   */
  onEmailInput(e: WechatMiniprogram.Input) {
    this.setData({
      email: e.detail.value,
      emailError: '',
    });
  },

  /**
   * 输入密码
   */
  onPasswordInput(e: WechatMiniprogram.Input) {
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
      const { access_token, user } = await login({
        email: email.trim(),
        password,
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
   * 忘记密码
   */
  handleForgotPassword() {
    wx.showModal({
      title: '提示',
      content: '请在后台管理系统重置密码',
      showCancel: false,
      confirmColor: '#0052d9',
    });
  },

  onShareAppMessage() {
    return { title: '登录', path: '/pages/login/login' };
  },
});
