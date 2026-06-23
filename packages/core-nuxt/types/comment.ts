import type { UserRole } from './user';

/**
 * 评论目标类型
 */
export type CommentTargetType = 'article' | 'page';

/**
 * 评论数据结构
 */
export interface Comment {
  id: number;
  content: string;
  is_deleted: boolean;
  is_pinned: boolean; // 是否置顶
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

/**
 * 游客信息
 */
export interface GuestInfo {
  nickname?: string;
  email?: string;
  website?: string;
}

/**
 * 表情项
 */
export interface EmojiItem {
  key: string;
  val: string;
}

/**
 * 表情分组
 */
export interface EmojiGroup {
  name: string;
  type: 'emoji' | 'image' | 'emoticon';
  items: EmojiItem[];
}

/**
 * 获取评论列表参数
 */
export interface GetCommentsParams {
  target_type: CommentTargetType;
  target_key: string | number;
  page?: number;
  page_size?: number;
}

/**
 * 扁平化评论项（flattenComments 输出 / CommentList 入参）
 */
export interface FlatComment {
  comment: Comment;
  depth: number;
}
