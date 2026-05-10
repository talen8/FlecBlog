/**
 * 访问统计追踪插件
 *
 * 客户端插件，用于收集和上报网站访问统计数据，包括：
 * - 页面浏览量（PV）
 * - 页面停留时长
 * - 自定义事件追踪
 *
 * 数据通过 POST 请求发送到后端 /collect 接口。
 *
 * @module plugins/tracker.client
 */

import { post } from '@/utils/request';

/**
 * 访问统计追踪插件
 *
 * 自动追踪页面浏览和停留时长，并提供自定义事件追踪方法。
 *
 * @example
 * ```ts
 * // 在组件中使用
 * const { $tracker } = useNuxtApp();
 *
 * // 追踪自定义事件
 * $tracker.trackEvent('click_button', { button_id: 'submit' });
 *
 * // 手动发送停留时长
 * $tracker.sendDuration();
 * ```
 */
export default defineNuxtPlugin(() => {
  if (!import.meta.client) return;

  /**
   * 收集并发送统计数据
   *
   * @param data - 统计数据
   * @param data.url - 页面 URL
   * @param data.title - 页面标题
   * @param data.type - 统计类型（pageview/event/duration）
   * @param data.article_id - 文章 ID（可选）
   * @param data.duration - 停留时长（秒）
   * @param data.event_name - 事件名称
   * @param data.event_data - 事件数据
   */
  const collect = (data: {
    url?: string;
    title?: string;
    type?: 'pageview' | 'event' | 'duration';
    article_id?: number;
    duration?: number;
    event_name?: string;
    event_data?: Record<string, unknown>;
  }) => {
    const payload = {
      ...data,
      hostname: window.location.hostname,
      referrer: document.referrer || undefined,
      language: navigator.language,
      screen: `${screen.width}x${screen.height}`,
      timestamp: Date.now(),
    };

    try {
      post('/collect', payload);
    } catch {
      // 静默失败，不影响用户体验
    }
  };

  /** 页面进入时间 */
  let pageEnterTime = Date.now();

  const router = useRouter();

  // 路由变化时追踪页面浏览
  router.afterEach(to => {
    collect({
      url: to.fullPath,
      title: document.title,
      type: 'pageview',
    });

    pageEnterTime = Date.now();
  });

  /**
   * 发送页面停留时长
   *
   * 计算当前页面停留时长并发送到后端。
   */
  const sendDuration = () => {
    const duration = Math.floor((Date.now() - pageEnterTime) / 1000);
    if (duration > 0) {
      collect({
        url: window.location.pathname,
        type: 'duration',
        duration,
      });
    }
  };

  // 页面关闭时发送停留时长
  window.addEventListener('beforeunload', sendDuration);

  /**
   * 追踪自定义事件
   *
   * @param eventName - 事件名称
   * @param eventData - 事件数据（可选）
   *
   * @example
   * ```ts
   * $tracker.trackEvent('click_download', { file: 'document.pdf' });
   * ```
   */
  const trackEvent = (eventName: string, eventData?: Record<string, unknown>) => {
    collect({
      url: window.location.pathname,
      type: 'event',
      event_name: eventName,
      event_data: eventData,
    });
  };

  return {
    provide: {
      tracker: {
        collect,
        trackEvent,
        sendDuration,
      },
    },
  };
});
