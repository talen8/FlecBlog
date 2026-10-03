import { subscribe, unsubscribe } from './api/subscribe';

/**
 * 邮箱订阅管理
 * @returns subscribe - 订阅本站，unsubscribe - 退订（通过邮件中的 token）
 */
export function useSubscribe() {
  return { subscribe, unsubscribe };
}
