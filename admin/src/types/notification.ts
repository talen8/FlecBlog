import type { PaginationQuery } from './request';

// 通知类型（后台管理系统）
export type NotificationType =
  | 'comment_new' // 新评论
  | 'feedback_new' // 反馈投诉
  | 'system_alert' // 系统告警
  | 'friend_apply'; // 友链申请

// 通知数据类型定义
export interface CommentNotificationData {
  article_title: string;
  article_slug: string;
  comment_id: number;
  comment_content: string;
  parent_comment_id?: number;
}

export interface FeedbackNotificationData {
  ticket_no: string;
  report_url: string;
  report_type: string;
  form_content: Record<string, unknown>;
  status: string;
}

export interface SystemAlertNotificationData {
  alert_type: string;
  message: string;
  severity: string;
}

export interface FriendApplyNotificationData {
  site_name: string;
  site_url: string;
  description?: string;
}

// 通知对象
export interface Notification {
  id: number;
  type: NotificationType;
  type_text: string; // 类型中文文本（后端提供，前端直接显示）

  // 前端显示字段（直接使用）
  title: string;
  content: string;
  link: string;

  // 详细数据（备用）
  data:
    | CommentNotificationData
    | FeedbackNotificationData
    | SystemAlertNotificationData
    | FriendApplyNotificationData
    | Record<string, unknown>;
  target_id?: number;
  is_read: boolean;
  read_at: string | null;
  created_at: string;
  sender: string | null;
}

// 通知列表数据
export interface NotificationListData {
  list: Notification[];
  total: number;
  page: number;
  page_size: number;
  unread_count: number; // 未读数量
}

// 通知列表查询参数
export type NotificationQueryParams = PaginationQuery;
