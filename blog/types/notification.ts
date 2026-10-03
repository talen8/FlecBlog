/**
 * 通知类型枚举（前台用户）
 */
export type NotificationType = 'comment_reply'; // 评论回复

/**
 * 评论通知数据类型
 */
export interface CommentNotificationData {
  article_title: string;
  article_slug: string;
  comment_id: number;
  comment_content: string;
  parent_comment_id?: number;
}

/**
 * 通知对象（前台用户）
 */
export interface Notification {
  id: number;
  type: NotificationType;
  type_text: string; // 类型中文文本（后端提供，前端直接显示）

  // 前端显示字段（直接使用）
  title: string;
  content: string;
  link: string;

  // 详细数据（备用）
  data: CommentNotificationData | Record<string, unknown>;
  target_id?: number;
  is_read: boolean;
  read_at: string | null;
  created_at: string;
  sender: string | null;
}

/**
 * 通知列表响应
 */
export interface NotificationListResponse {
  list: Notification[];
  total: number;
  page: number;
  page_size: number;
  unread_count: number; // 未读数量
}

/**
 * 获取通知列表参数
 */
export interface GetNotificationsParams {
  page: number;
  page_size: number;
}
