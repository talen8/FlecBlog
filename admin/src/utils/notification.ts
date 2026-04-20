/**
 * 系统通知工具类
 */

export interface SystemNotificationOptions {
  title: string;
  body: string;
  icon?: string;
  badge?: string;
  tag?: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data?: any;
  requireInteraction?: boolean;
}

class NotificationManager {
  private permission: NotificationPermission = 'default';

  constructor() {
    if ('Notification' in window) {
      this.permission = Notification.permission;
    }
  }

  /**
   * 请求通知权限
   */
  async requestPermission(): Promise<boolean> {
    if (!('Notification' in window)) {
      console.warn('浏览器不支持通知功能');
      return false;
    }

    if (this.permission === 'granted') {
      return true;
    }

    if (this.permission === 'denied') {
      console.warn('用户已拒绝通知权限');
      return false;
    }

    try {
      this.permission = await Notification.requestPermission();
      return this.permission === 'granted';
    } catch (error) {
      console.error('请求通知权限失败:', error);
      return false;
    }
  }

  /**
   * 显示系统通知
   */
  async show(options: SystemNotificationOptions): Promise<Notification | null> {
    // 检查权限
    if (this.permission !== 'granted') {
      const granted = await this.requestPermission();
      if (!granted) return null;
    }

    try {
      const notification = new Notification(options.title, {
        body: options.body,
        icon: options.icon || '/pwa-192x192.png',
        badge: options.badge || '/pwa-192x192.png',
        tag: options.tag,
        data: options.data,
        requireInteraction: options.requireInteraction || false,
        silent: false,
      });

      // 点击通知时的处理
      notification.onclick = event => {
        event.preventDefault();
        window.focus();

        // 如果有跳转链接，则跳转
        if (options.data?.link) {
          window.location.href = options.data.link;
        }

        notification.close();
      };

      return notification;
    } catch (error) {
      console.error('显示通知失败:', error);
      return null;
    }
  }

  /**
   * 检查是否支持通知
   */
  isSupported(): boolean {
    return 'Notification' in window;
  }

  /**
   * 获取当前权限状态
   */
  getPermission(): NotificationPermission {
    return this.permission;
  }
}

export const notificationManager = new NotificationManager();
