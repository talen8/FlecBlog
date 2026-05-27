/**
 * 微信扫码登录授权确认页
 * 用户扫码进入小程序后，在此页面确认授权登录
 */

import { confirmWechatLogin } from '../../api/user';

Page({
  data: {
    scene: '',
    loading: false,
    confirmed: false,
    error: '',
  },

  onLoad(options: Record<string, string>) {
    if (options.scene) {
      this.setData({ scene: options.scene });
    }
  },

  /**
   * 确认授权登录
   */
  async handleConfirm() {
    if (this.data.loading || this.data.confirmed) return;
    if (!this.data.scene) {
      this.setData({ error: '缺少场景参数' });
      return;
    }

    this.setData({ loading: true, error: '' });
    try {
      // 调用 wx.login 获取微信身份
      const { code } = await new Promise<WechatMiniprogram.LoginSuccessCallbackResult>((resolve, reject) => {
        wx.login({ success: resolve, fail: reject });
      });

      await confirmWechatLogin(this.data.scene, code);
      this.setData({ confirmed: true });
      wx.showToast({ title: '授权成功', icon: 'success' });
    } catch (error) {
      this.setData({
        error: (error as Error).message || '授权失败，场景可能已过期',
      });
    } finally {
      this.setData({ loading: false });
    }
  },

  onShareAppMessage() {
    return { title: '扫码登录', path: '/pages/wechat-auth/wechat-auth' };
  },
});
