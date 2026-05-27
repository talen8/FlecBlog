/**
 * 微信小程序类型定义
 */

// ==================== 通用类型 ====================

export interface ApiResponse<T> {
  code: number;
  message: string;
  data?: T;
}

export interface PageQuery {
  page?: number;
  page_size?: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// ==================== 文章相关类型 ====================

export interface ArticleCategory {
  id: number;
  name: string;
  url: string;
}

export interface ArticleTag {
  id: number;
  name: string;
  url: string;
}

export interface ArticleNavigation {
  title: string;
  url: string;
}

export interface ArticleListItem {
  id: number;
  title: string;
  summary: string;
  excerpt?: string;
  cover: string;
  location: string;
  is_top: boolean;
  is_essence: boolean;
  is_outdated: boolean;
  url: string;
  comment_count: number;
  publish_time: string;
  update_time: string;
  category: ArticleCategory;
  tags: ArticleTag[];
}

export interface ArticleDetail extends ArticleListItem {
  slug: string;
  content: string;
  ai_summary?: string;
  view_count: number;
  prev?: ArticleNavigation;
  next?: ArticleNavigation;
}

export interface ArticleQuery extends PageQuery {
  year?: string;
  month?: string;
  category?: string;
  tag?: string;
}

// ==================== 信息相关类型 ====================

export interface SiteBasicConfig {
  author: string;
  author_avatar: string;
  blog_url: string;
}

export interface BlogConfig {
  title: string;
  established: string;
  about_describe: string;
  about_describe_tips: string;
  about_profile: string;
  about_personality: string;
  about_motto_main: string;
  about_motto_sub: string;
  about_versions: string;
  about_unions: string;
  about_story: string;
}

// ==================== 动态相关类型 ====================

export interface MomentVideo {
  platform: string;
  url: string;
  video_id: string;
}

export interface MomentLink {
  favicon: string;
  title: string;
  url: string;
}

export interface MomentMusic {
  id: string;
  server: string;
  type: string;
}

export interface MomentContent {
  text: string;
  tags?: string;
  images?: string[];
  video?: MomentVideo;
  link?: MomentLink;
  music?: MomentMusic;
}

export interface MomentListItem {
  id: number;
  is_publish: boolean;
  publish_time: string;
  content: MomentContent;
}

// ==================== 配置相关类型 ====================

export interface OAuthConfig {
  'oauth.github.enabled': string;
  'oauth.google.enabled': string;
  'oauth.qq.enabled': string;
  'oauth.microsoft.enabled': string;
  'oauth.oidc.enabled': string;
  'oauth.wechat.enabled': string;
}

export interface SiteStats {
  total_words: string;
  total_visitors: number;
  total_page_views: number;
  online_users: number;
  today_visitors: number;
  today_pageviews: number;
  yesterday_visitors: number;
  yesterday_pageviews: number;
  month_pageviews: number;
  total_articles: number;
  total_comments: number;
  total_friends: number;
  total_moments: number;
  total_categories: number;
  total_tags: number;
}

// ==================== 用户相关类型 ====================

export type UserRole = 'super_admin' | 'admin' | 'user';

export interface UserInfo {
  id: number;
  email: string;
  email_hash: string;
  is_virtual_email: boolean;
  avatar?: string;
  badge?: string;
  nickname: string;
  website?: string;
  last_login?: string;
  created_at: string;
  role: UserRole;
  has_password: boolean;
  linked_oauths: string[];
}

// ==================== 请求相关类型 ====================

export type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

export interface RequestConfig {
  url: string;
  method?: RequestMethod;
  data?: Record<string, any>;
  params?: Record<string, any>;
  header?: Record<string, string>;
  loading?: boolean;
  loadingText?: string;
  needAuth?: boolean;
}

export interface ApiError {
  code: number;
  message: string;
  url?: string;
  statusCode?: number;
}

export type RequestInterceptor = (config: RequestConfig) => RequestConfig | Promise<RequestConfig>;

export type ErrorInterceptor = (error: ApiError, config: RequestConfig) => void | Promise<void>;
