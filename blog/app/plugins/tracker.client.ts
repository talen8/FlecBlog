/**
 * 访问统计追踪插件
 * 自动追踪 PV/UV、页面停留时长、支持手动事件追踪
 *
 * 使用 .client.ts 后缀确保只在客户端运行
 * 使用 requestIdleCallback 延迟初始化，避免阻塞首屏渲染
 */

// 类型声明
export interface TrackerPlugin {
  trackPageView: (path?: string, articleId?: number) => void;
  trackEvent: (name: string, data?: Record<string, unknown>) => void;
  setArticleId: (id?: number) => void;
}

declare module '#app' {
  interface NuxtApp {
    $tracker: TrackerPlugin;
  }
}

export default defineNuxtPlugin({
  parallel: true, // 并行加载，不阻塞其他插件
  setup() {
    const router = useRouter();
    const endpoint = `${useRuntimeConfig().public.apiUrl}/collect`;

    let pageStartTime = Date.now();
    let lastPageUrl = location.pathname + location.search;
    let currentArticleId: number | undefined;

    const getBaseData = (url?: string, articleId?: number) => ({
      url: url || location.pathname + location.search,
      hostname: location.hostname,
      referrer: document.referrer,
      language: navigator.language,
      screen: `${screen.width}x${screen.height}`,
      title: document.title,
      timestamp: Date.now(),
      ...(articleId !== undefined && { article_id: articleId }),
    });

    const send = (
      type: string,
      extra: Record<string, unknown> = {},
      url?: string,
      articleId?: number
    ) => {
      const payload = {
        ...getBaseData(url, articleId),
        type,
        ...extra,
      };
      const blob = new Blob([JSON.stringify(payload)], {
        type: 'application/json',
      });
      if (navigator.sendBeacon) {
        navigator.sendBeacon(endpoint, blob);
      } else {
        fetch(endpoint, { method: 'POST', body: blob, keepalive: true }).catch(() => {});
      }
    };

    const sendDuration = (url?: string, articleId?: number) => {
      const sec = Math.floor((Date.now() - pageStartTime) / 1000);
      if (sec > 0) send('duration', { duration: sec }, url, articleId);
    };

    // 页面隐藏/卸载时发送停留时长
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        sendDuration(undefined, currentArticleId);
      } else {
        pageStartTime = Date.now();
      }
    });
    window.addEventListener('beforeunload', () => sendDuration(undefined, currentArticleId));

    // 路由变化时统计
    router.afterEach(to => {
      setTimeout(() => {
        sendDuration(lastPageUrl, currentArticleId);
        pageStartTime = Date.now();
        lastPageUrl = to.path;
        currentArticleId = undefined;
        send('pageview', {}, to.path);
      }, 100);
    });

    return {
      provide: {
        tracker: {
          trackPageView: (path?: string, articleId?: number) =>
            send('pageview', {}, path, articleId),
          trackEvent: (name: string, data?: Record<string, unknown>) =>
            name && send('event', { event_name: name, event_data: data }),
          setArticleId: (id?: number) => {
            currentArticleId = id;
          },
        },
      },
    };
  },
});
