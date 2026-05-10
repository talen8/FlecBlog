/**
 * TypeScript 类型定义
 *
 * 本文件包含所有 API 响应、实体模型、请求参数等类型定义。
 * 这些类型与后端 API 结构对应，确保前后端数据类型一致。
 *
 * @module types
 */

/* ============================================
 * 通用响应类型
 * ============================================ */

/**
 * API 统一响应格式
 *
 * 所有 API 接口的统一响应结构。
 *
 * @template T - 响应数据类型
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export interface ApiResponse<T = any> {
  /** 状态码，0 表示成功 */
  code: number;
  /** 响应消息 */
  message: string;
  /** 响应数据 */
  data: T;
}

/**
 * 分页查询参数
 */
export interface PaginationQuery {
  /** 当前页码，从 1 开始 */
  page?: number;
  /** 每页数量 */
  page_size?: number;
}

/**
 * 分页数据响应
 *
 * @template T - 列表项数据类型
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export interface PaginationData<T = any> {
  /** 数据列表 */
  list: T[];
  /** 总数量 */
  total: number;
  /** 当前页码 */
  page: number;
  /** 每页数量 */
  page_size: number;
}

/* ============================================
 * 文章相关类型
 * ============================================ */

/**
 * 文章导航信息
 *
 * 用于文章详情页的上一篇/下一篇导航。
 */
export interface ArticleNav {
  /** 文章标题 */
  title: string;
  /** 文章链接 */
  url: string;
}

/**
 * 文章实体
 */
export interface Article {
  /** 文章 ID */
  id: number;
  /** 文章标题 */
  title: string;
  /** 文章别名（URL 友好） */
  slug?: string;
  /** 文章链接 */
  url: string;
  /** 文章内容（Markdown） */
  content?: string;
  /** 文章摘要 */
  summary: string;
  /** AI 生成的摘要 */
  ai_summary?: string;
  /** 文章摘录 */
  excerpt?: string;
  /** 封面图片 */
  cover?: string;
  /** 是否置顶 */
  is_top: boolean;
  /** 是否精选 */
  is_essence: boolean;
  /** 是否过时 */
  is_outdated?: boolean;
  /** 浏览次数 */
  view_count?: number;
  /** 评论数量 */
  comment_count: number;
  /** 发布时间 */
  publish_time: string;
  /** 更新时间 */
  update_time: string;
  /** 发布地点 */
  location?: string;
  /** 所属分类 */
  category: {
    id: number;
    name: string;
    url: string;
  };
  /** 文章标签 */
  tags: Array<{
    id: number;
    name: string;
    url: string;
  }>;
  /** 上一篇文章 */
  prev?: ArticleNav;
  /** 下一篇文章 */
  next?: ArticleNav;
}

/**
 * 文章查询参数
 */
export interface ArticleQuery {
  /** 页码 */
  page?: number;
  /** 每页数量 */
  page_size?: number;
  /** 年份（归档查询） */
  year?: string;
  /** 月份（归档查询） */
  month?: string;
  /** 分类别名 */
  category?: string;
  /** 标签别名 */
  tag?: string;
}

/* ============================================
 * 分类相关类型
 * ============================================ */

/**
 * 分类实体
 */
export interface Category {
  /** 分类 ID */
  id: number;
  /** 分类名称 */
  name: string;
  /** 分类别名 */
  slug: string;
  /** 分类链接 */
  url: string;
  /** 分类描述 */
  description: string;
  /** 文章数量 */
  count: number;
  /** 排序序号 */
  sort: number;
}

/* ============================================
 * 标签相关类型
 * ============================================ */

/**
 * 标签实体
 */
export interface Tag {
  /** 标签 ID */
  id: number;
  /** 标签名称 */
  name: string;
  /** 标签别名 */
  slug: string;
  /** 标签链接 */
  url: string;
  /** 标签描述 */
  description: string;
  /** 文章数量 */
  count: number;
}

/* ============================================
 * 评论相关类型
 * ============================================ */

/**
 * 评论目标类型
 */
export type CommentTargetType = 'article' | 'page' | 'moment';

/**
 * 评论实体
 */
export interface Comment {
  /** 评论 ID */
  id: number;
  /** 评论内容 */
  content: string;
  /** 是否已删除 */
  is_deleted: boolean;
  /** 父评论 ID */
  parent_id: number | null;
  /** 创建时间 */
  created_at: string;
  /** 发布地点 */
  location?: string;
  /** 浏览器信息 */
  browser?: string;
  /** 操作系统信息 */
  os?: string;
  /** 评论用户信息 */
  user: {
    role: UserRole;
    badge?: string;
    id: number;
    email_hash: string;
    nickname: string;
    avatar: string;
    website?: string;
  };
  /** 回复目标用户信息 */
  reply_user?: {
    role: UserRole;
    badge?: string;
    id: number;
    email_hash: string;
    nickname: string;
    avatar: string;
    website?: string;
  };
  /** 子评论列表 */
  replies: Comment[];
}

/**
 * 创建评论参数
 */
export interface CreateCommentParams {
  /** 评论目标类型 */
  target_type: CommentTargetType;
  /** 评论目标标识 */
  target_key: string | number;
  /** 评论内容 */
  content: string;
  /** 父评论 ID（回复时使用） */
  parent_id?: number;
  /** 昵称（游客评论时使用） */
  nickname?: string;
  /** 邮箱（游客评论时使用） */
  email?: string;
  /** 网站（游客评论时使用） */
  website?: string;
}

/* ============================================
 * 动态相关类型
 * ============================================ */

/**
 * 动态视频信息
 */
export interface MomentVideo {
  /** 视频链接 */
  url: string;
  /** 视频平台 */
  platform?: 'youtube' | 'bilibili' | 'local';
  /** 视频ID */
  video_id?: string;
}

/**
 * 动态音频信息
 */
export interface MomentAudio {
  /** 音频链接 */
  url: string;
}

/**
 * 动态音乐信息（网易云/QQ音乐）
 */
export interface MomentMusic {
  /** 音乐平台 */
  server: 'netease' | 'tencent';
  /** 资源类型 */
  type: 'song' | 'playlist' | 'album' | 'artist';
  /** 资源ID */
  id: string;
}

/**
 * 动态链接卡片信息
 */
export interface MomentLink {
  /** 链接地址 */
  url: string;
  /** 链接标题 */
  title: string;
  /** 网站图标 */
  favicon?: string;
}

/**
 * 动态内容
 */
export interface MomentContent {
  /** 文字内容 */
  text?: string;
  /** 图片列表 */
  images?: string[];
  /** 发布地点 */
  location?: string;
  /** 标签 */
  tags?: string;
  /** 视频信息 */
  video?: MomentVideo;
  /** 音频信息 */
  audio?: MomentAudio;
  /** 音乐信息 */
  music?: MomentMusic;
  /** 链接卡片信息 */
  link?: MomentLink;
}

/**
 * 动态实体
 */
export interface Moment {
  /** 动态 ID */
  id: number;
  /** 动态内容 */
  content: MomentContent;
  /** 发布时间 */
  publish_time: string;
}

/**
 * 动态列表响应
 */
export interface MomentListResponse {
  /** 动态列表 */
  list: Moment[];
  /** 总数量 */
  total: number;
  /** 当前页码 */
  page: number;
  /** 每页数量 */
  page_size: number;
}

/* ============================================
 * 友链相关类型
 * ============================================ */

/**
 * 友链类型
 */
export interface FriendType {
  /** 类型 ID */
  id: number;
  /** 类型名称 */
  name: string;
  /** 是否可见 */
  is_visible: boolean;
  /** 排序序号 */
  sort: number;
}

/**
 * 友链实体
 */
export interface Friend {
  /** 友链 ID */
  id: number;
  /** 网站名称 */
  name: string;
  /** 网站链接 */
  url: string;
  /** 网站描述 */
  description?: string;
  /** 头像 */
  avatar?: string;
  /** 网站截图 */
  screenshot?: string;
  /** 排序序号 */
  sort: number;
  /** 是否失效 */
  is_invalid: boolean;
  /** 类型 ID */
  type_id?: number;
  /** 类型信息 */
  type?: FriendType;
}

/**
 * 友链分组
 */
export interface FriendGroup {
  /** 类型 ID */
  type_id: number | null;
  /** 类型名称 */
  type_name: string;
  /** 类型排序 */
  type_sort: number;
  /** 友链列表 */
  friends: Friend[];
}

/**
 * 友链分组响应
 */
export interface FriendGroupedResponse {
  /** 分组列表 */
  groups: FriendGroup[];
  /** 分组总数 */
  total_groups: number;
  /** 友链总数 */
  total_friends: number;
}

/**
 * 友链查询参数
 */
export interface FriendQueryParams {
  /** 页码 */
  page?: number;
  /** 每页数量 */
  page_size?: number;
}

/**
 * 友链申请参数
 */
export interface FriendApplyRequest {
  /** 网站名称 */
  name: string;
  /** 网站链接 */
  url: string;
  /** 网站描述 */
  description: string;
  /** 头像 */
  avatar: string;
  /** 网站截图 */
  screenshot?: string;
}

/* ============================================
 * 反馈相关类型
 * ============================================ */

/**
 * 反馈类型
 */
export type ReportType = 'copyright' | 'inappropriate' | 'summary' | 'suggestion';

/**
 * 反馈状态
 */
export type FeedbackStatus = 'pending' | 'resolved' | 'closed';

/**
 * 反馈实体
 */
export interface Feedback {
  /** 反馈 ID */
  id: number;
  /** 工单号 */
  ticket_no: string;
  /** 反馈链接 */
  report_url: string;
  /** 反馈类型 */
  report_type: ReportType;
  /** 表单内容 */
  form_content: Record<string, unknown>;
  /** 联系邮箱 */
  email: string;
  /** 处理状态 */
  status: FeedbackStatus;
  /** 管理员回复 */
  admin_reply: string;
  /** 回复时间 */
  reply_time?: string;
  /** 用户代理 */
  user_agent: string;
  /** IP 地址 */
  ip: string;
  /** 反馈时间 */
  feedback_time: string;
}

/**
 * 提交反馈参数
 */
export interface SubmitFeedbackParams {
  /** 反馈链接 */
  reportUrl?: string;
  /** 反馈类型 */
  reportType: ReportType;
  /** 联系邮箱 */
  email?: string;
  /** 问题描述 */
  description: string;
  /** 原因说明 */
  reason?: string;
  /** 附件文件列表 */
  attachmentFiles?: string[];
}

/* ============================================
 * 菜单相关类型
 * ============================================ */

/**
 * 菜单类型
 */
export type MenuType = 'navigation' | 'footer' | 'aggregate';

/**
 * 菜单实体
 */
export interface Menu {
  /** 菜单 ID */
  id: number;
  /** 菜单类型 */
  type: MenuType;
  /** 父菜单 ID */
  parent_id: number | null;
  /** 菜单标题 */
  title: string;
  /** 菜单链接 */
  url: string;
  /** 菜单图标 */
  icon: string;
  /** 排序序号 */
  sort: number;
  /** 子菜单列表 */
  children?: Menu[];
}

/**
 * 菜单列表响应
 */
export interface MenusResponse {
  /** 状态码 */
  code: number;
  /** 响应消息 */
  message: string;
  /** 菜单列表 */
  data: Menu[];
}

/* ============================================
 * 通知相关类型
 * ============================================ */

/**
 * 通知类型
 */
export type NotificationType = 'comment_reply';

/**
 * 评论通知数据
 */
export interface CommentNotificationData {
  /** 文章标题 */
  article_title: string;
  /** 文章别名 */
  article_slug: string;
  /** 评论 ID */
  comment_id: number;
  /** 评论内容 */
  comment_content: string;
  /** 父评论 ID */
  parent_comment_id?: number;
}

/**
 * 通知实体
 */
export interface Notification {
  /** 通知 ID */
  id: number;
  /** 通知类型 */
  type: NotificationType;
  /** 通知类型文本 */
  type_text: string;
  /** 通知标题 */
  title: string;
  /** 通知内容 */
  content: string;
  /** 相关链接 */
  link: string;
  /** 通知数据 */
  data: CommentNotificationData | Record<string, unknown>;
  /** 关联目标 ID */
  target_id?: number;
  /** 是否已读 */
  is_read: boolean;
  /** 已读时间 */
  read_at: string | null;
  /** 创建时间 */
  created_at: string;
  /** 发送者 */
  sender: string | null;
}

/**
 * 通知列表响应
 */
export interface NotificationListResponse {
  /** 通知列表 */
  list: Notification[];
  /** 总数量 */
  total: number;
  /** 当前页码 */
  page: number;
  /** 每页数量 */
  page_size: number;
  /** 未读数量 */
  unread_count: number;
}

/**
 * 获取通知参数
 */
export interface GetNotificationsParams {
  /** 页码 */
  page: number;
  /** 每页数量 */
  page_size: number;
}

/* ============================================
 * 统计相关类型
 * ============================================ */

/**
 * 网站统计数据
 */
export interface SiteStats {
  /** 总字数 */
  total_words: string;
  /** 总访客数 */
  total_visitors: number;
  /** 总浏览量 */
  total_page_views: number;
  /** 在线用户数 */
  online_users: number;
  /** 文章总数 */
  total_articles: number;
  /** 评论总数 */
  total_comments: number;
  /** 友链总数 */
  total_friends: number;
  /** 动态总数 */
  total_moments: number;
  /** 分类总数 */
  total_categories: number;
  /** 标签总数 */
  total_tags: number;
  /** 今日访客数 */
  today_visitors: number;
  /** 今日浏览量 */
  today_pageviews: number;
  /** 昨日访客数 */
  yesterday_visitors: number;
  /** 昨日浏览量 */
  yesterday_pageviews: number;
  /** 本月浏览量 */
  month_pageviews: number;
}

/* ============================================
 * 配置相关类型
 * ============================================ */

/**
 * 配置分组类型
 */
export type SettingGroupType = 'basic' | 'blog' | 'ai' | 'oauth' | 'upload';

/* ============================================
 * 用户相关类型
 * ============================================ */

/**
 * 用户角色
 */
export type UserRole = 'super_admin' | 'admin' | 'user';

/**
 * 用户信息
 */
export interface UserInfo {
  /** 用户 ID */
  id: number;
  /** 邮箱 */
  email: string;
  /** 邮箱哈希值（用于 Gravatar） */
  email_hash: string;
  /** 是否虚拟邮箱 */
  is_virtual_email: boolean;
  /** 头像 */
  avatar?: string;
  /** 徽章 */
  badge?: string;
  /** 昵称 */
  nickname: string;
  /** 个人网站 */
  website?: string;
  /** 最后登录时间 */
  last_login?: string;
  /** 注册时间 */
  created_at: string;
  /** 用户角色 */
  role: UserRole;
  /** 是否设置密码 */
  has_password: boolean;
  /** 已绑定的 OAuth 提供商列表 */
  linked_oauths: string[];
}

/**
 * 更新用户信息参数
 */
export interface UpdateProfileParams {
  /** 昵称 */
  nickname?: string;
  /** 邮箱 */
  email?: string;
  /** 头像 */
  avatar?: string;
  /** 徽章 */
  badge?: string;
  /** 个人网站 */
  website?: string;
}

/**
 * 登录参数
 */
export interface LoginParams {
  /** 邮箱 */
  email: string;
  /** 密码 */
  password: string;
}

/**
 * 登录响应
 */
export interface LoginResponse {
  /** 访问令牌 */
  access_token: string;
  /** 用户信息 */
  user: UserInfo;
}

/**
 * 注册参数
 */
export interface RegisterParams {
  /** 邮箱 */
  email: string;
  /** 昵称 */
  nickname: string;
  /** 密码 */
  password: string;
  /** 个人网站 */
  website?: string;
}

/**
 * 注册响应
 */
export interface RegisterResponse {
  /** 访问令牌 */
  access_token: string;
  /** 用户信息 */
  user: UserInfo;
}

/**
 * 忘记密码参数
 */
export interface ForgotPasswordParams {
  /** 邮箱 */
  email: string;
}

/**
 * 重置密码参数
 */
export interface ResetPasswordParams {
  /** 邮箱 */
  email: string;
  /** 验证码 */
  code: string;
  /** 新密码 */
  password: string;
}

/**
 * 修改密码参数
 */
export interface ChangePasswordParams {
  /** 原密码 */
  old_password: string;
  /** 新密码 */
  new_password: string;
}

/**
 * 注销账号参数
 */
export interface DeactivateAccountParams {
  /** 密码 */
  password: string;
}

/**
 * 刷新令牌响应
 */
export interface RefreshTokenResponse {
  /** 访问令牌 */
  access_token: string;
}
