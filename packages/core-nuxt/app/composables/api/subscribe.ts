import { createApi } from './createApi';

const subscribeApi = createApi<unknown>('/subscribe');

/** 订阅本站 */
export const subscribe = (email: string) => subscribeApi.post('', { email });

/** 退订（通过邮件中的 token） */
export const unsubscribe = (token: string) =>
  subscribeApi.get(`/unsubscribe?token=${encodeURIComponent(token)}`);
