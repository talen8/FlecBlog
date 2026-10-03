import type { Notification, GetNotificationsParams } from '../../types/notification';
import { getNotifications, markAsRead, markAllAsRead } from './api/notification';

/**
 * 通知管理（列表 + 未读数 + 已读/全部已读 + 分页）
 * @returns notifications - 通知列表
 * @returns total - 总数
 * @returns currentPage / pageSize - 分页状态
 * @returns unreadCount - 未读数
 * @returns loading - 加载状态
 * @returns fetchNotifications - 获取列表
 * @returns markNotificationAsRead - 单条已读
 * @returns markAllNotificationsAsRead - 全部已读
 * @returns resetPage / clearNotifications - 重置/清空
 */
export function useNotifications() {
  const notifications = useState<Notification[]>('notifications', () => []);
  const total = useState<number>('notifications-total', () => 0);
  const currentPage = useState<number>('notifications-currentPage', () => 1);
  const pageSize = useState<number>('notifications-pageSize', () => 10);
  const unreadCount = useState<number>('notifications-unreadCount', () => 0);
  const loading = useState<boolean>('notifications-loading', () => false);

  const fetchNotifications = async (params?: Partial<GetNotificationsParams>) => {
    loading.value = true;
    try {
      const response = await getNotifications({
        page: params?.page ?? currentPage.value,
        page_size: params?.page_size ?? pageSize.value,
      });
      notifications.value = response.list || [];
      total.value = response.total || 0;
      unreadCount.value = response.unread_count || 0;
      if (params?.page) {
        currentPage.value = params.page;
      }
    } catch (error) {
      console.error('获取通知列表失败:', error);
      notifications.value = [];
      total.value = 0;
      unreadCount.value = 0;
    } finally {
      loading.value = false;
    }
  };

  const markNotificationAsRead = async (id: number) => {
    try {
      await markAsRead(id);
      const notification = notifications.value.find(n => n.id === id);
      if (notification?.is_read === false) {
        notification.is_read = true;
        notification.read_at = new Date().toISOString();
        unreadCount.value = Math.max(0, unreadCount.value - 1);
      }
    } catch (error) {
      console.error('标记通知已读失败:', error);
      throw error;
    }
  };

  const markAllNotificationsAsRead = async () => {
    try {
      await markAllAsRead();
      notifications.value.forEach(n => {
        n.is_read = true;
        n.read_at = new Date().toISOString();
      });
      unreadCount.value = 0;
    } catch (error) {
      console.error('标记所有通知已读失败:', error);
      throw error;
    }
  };

  return {
    notifications,
    total,
    currentPage,
    pageSize,
    unreadCount,
    loading,
    fetchNotifications,
    markNotificationAsRead,
    markAllNotificationsAsRead,
    resetPage: () => (currentPage.value = 1),
    clearNotifications: () => {
      notifications.value = [];
      total.value = 0;
      currentPage.value = 1;
      unreadCount.value = 0;
    },
  };
}
