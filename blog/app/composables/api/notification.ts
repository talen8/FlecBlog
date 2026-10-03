import type { NotificationListResponse, GetNotificationsParams } from '../../../types/notification';
import { createApi } from './createApi';

const notificationApi = createApi<NotificationListResponse>('/notifications');

/** 获取通知列表 */
export const getNotifications = async (params: GetNotificationsParams) => {
  return notificationApi.get<NotificationListResponse>('', params);
};

/** 标记单条通知已读 */
export const markAsRead = async (id: number) => {
  await notificationApi.put(`/${id}/read`);
};

/** 标记全部通知已读 */
export const markAllAsRead = async () => {
  await notificationApi.put('/read-all');
};
