import type { UserRole } from './user';

/**
 * 评论目标类型
 */
export type CommentTargetType = 'article' | 'page' | 'moment';

/**
 * 评论数据结构
 */
export interface Comment {
  id: number;
  content: string;
  is_deleted: boolean;
  parent_id: number | null;
  created_at: string;
  location?: string; // 地理位置
  browser?: string; // 浏览器内核
  os?: string; // 操作系统
  user: {
    role: UserRole;
    badge?: string;
    id: number;
    email_hash: string;
    nickname: string;
    avatar: string;
    website?: string;
  };
  reply_user?: {
    role: UserRole;
    badge?: string;
    id: number;
    email_hash: string;
    nickname: string;
    avatar: string;
    website?: string;
  };
  replies: Comment[];
}

/**
 * 创建评论参数
 */
export interface CreateCommentParams {
  target_type: CommentTargetType;
  target_key: string | number;
  content: string;
  parent_id?: number;

  // 游客信息（可选，未登录时使用）
  nickname?: string;
  email?: string;
  website?: string;
}
