/**
 * FlecBlog 微信小程序类型定义
 * 与后端 API 数据结构保持一致
 */

// ==================== 通用类型 ====================

/**
 * API 统一响应格式
 */
export interface ApiResponse<T> {
  /** 状态码，0 表示成功 */
  code: number;
  /** 提示信息 */
  message: string;
  /** 响应数据 */
  data?: T;
}

/**
 * 分页查询参数
 */
export interface PageQuery {
  /** 页码，从 1 开始 */
  page?: number;
  /** 每页数量 */
  page_size?: number;
}

/**
 * 分页结果
 */
export interface PageResult<T> {
  /** 数据列表 */
  list: T[];
  /** 总条数 */
  total: number;
  /** 当前页码 */
  page: number;
  /** 每页数量 */
  page_size: number;
}

// ==================== 文章相关类型 ====================

/**
 * 文章分类信息
 */
export interface ArticleCategory {
  /** 分类 ID */
  id: number;
  /** 分类名称 */
  name: string;
  /** 分类链接 */
  url: string;
}

/**
 * 文章标签信息
 */
export interface ArticleTag {
  /** 标签 ID */
  id: number;
  /** 标签名称 */
  name: string;
  /** 标签链接 */
  url: string;
}

/**
 * 前后篇文章导航
 */
export interface ArticleNavigation {
  /** 文章标题 */
  title: string;
  /** 文章链接 */
  url: string;
}

/**
 * 文章列表项（前台展示）
 */
export interface ArticleListItem {
  /** 文章 ID */
  id: number;
  /** 文章标题 */
  title: string;
  /** 文章摘要 */
  summary: string;
  /** 文章摘录 */
  excerpt?: string;
  /** 封面图 */
  cover: string;
  /** 发布地点 */
  location: string;
  /** 是否置顶 */
  is_top: boolean;
  /** 是否精华 */
  is_essence: boolean;
  /** 是否过时 */
  is_outdated: boolean;
  /** 文章链接 */
  url: string;
  /** 评论数量 */
  comment_count: number;
  /** 发布时间 */
  publish_time: string;
  /** 更新时间 */
  update_time: string;
  /** 所属分类 */
  category: ArticleCategory;
  /** 标签列表 */
  tags: ArticleTag[];
}

/**
 * 文章详情
 */
export interface ArticleDetail extends ArticleListItem {
  /** 文章别名 */
  slug: string;
  /** 文章内容（Markdown） */
  content: string;
  /** AI 摘要 */
  ai_summary?: string;
  /** 浏览次数 */
  view_count: number;
  /** 上一篇文章 */
  prev?: ArticleNavigation;
  /** 下一篇文章 */
  next?: ArticleNavigation;
}

/**
 * 文章查询参数
 */
export interface ArticleQuery extends PageQuery {
  /** 按年份筛选 */
  year?: string;
  /** 按月份筛选 */
  month?: string;
  /** 按分类筛选（slug） */
  category?: string;
  /** 按标签筛选（slug） */
  tag?: string;
}

// ==================== 分类相关类型 ====================

/**
 * 分类信息
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
  /** 排序权重 */
  sort: number;
}

// ==================== 标签相关类型 ====================

/**
 * 标签信息
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

// ==================== 站点配置类型 ====================

/**
 * 站点基础配置
 */
export interface SiteBasicConfig {
  /** 站点名称 */
  site_name: string;
  /** 站点描述 */
  site_description: string;
  /** 站点 Logo */
  site_logo: string;
  /** 站点 ICP 备案号 */
  site_icp: string;
  /** 站点版权信息 */
  site_copyright: string;
}

/**
 * 博客配置
 */
export interface BlogConfig {
  /** 博客标题 */
  blog_title: string;
  /** 博客副标题 */
  blog_subtitle: string;
  /** 每页文章数量 */
  page_size: number;
  /** 是否开启评论 */
  enable_comment: boolean;
}

// ==================== 动态相关类型 ====================

/**
 * 动态视频信息
 */
export interface MomentVideo {
  /** 视频平台 */
  platform: string;
  /** 视频链接 */
  url: string;
  /** 视频ID */
  video_id: string;
}

/**
 * 动态链接信息
 */
export interface MomentLink {
  /** 网站图标 */
  favicon: string;
  /** 链接标题 */
  title: string;
  /** 链接地址 */
  url: string;
}

/**
 * 动态音乐信息
 */
export interface MomentMusic {
  /** 音乐ID */
  id: string;
  /** 音乐平台 */
  server: string;
  /** 音乐类型 */
  type: string;
}

/**
 * 动态内容
 */
export interface MomentContent {
  /** 动态文本 */
  text: string;
  /** 标签 */
  tags?: string;
  /** 图片列表 */
  images?: string[];
  /** 视频信息 */
  video?: MomentVideo;
  /** 链接信息 */
  link?: MomentLink;
  /** 音乐信息 */
  music?: MomentMusic;
}

/**
 * 动态列表项
 */
export interface MomentListItem {
  /** 动态 ID */
  id: number;
  /** 是否发布 */
  is_publish: boolean;
  /** 发布时间 */
  publish_time: string;
  /** 动态内容 */
  content: MomentContent;
}

/**
 * 动态查询参数
 */
export interface MomentQuery extends PageQuery {}

// ==================== 用户相关类型 ====================

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
  /** 邮箱哈希 */
  email_hash: string;
  /** 是否为虚拟邮箱 */
  is_virtual_email: boolean;
  /** 头像 */
  avatar?: string;
  /** 铭牌标识 */
  badge?: string;
  /** 昵称 */
  nickname: string;
  /** 个人网站 */
  website?: string;
  /** 最后登录时间 */
  last_login?: string;
  /** 注册时间 */
  created_at: string;
  /** 角色 */
  role: UserRole;
  /** 是否有密码 */
  has_password: boolean;
  /** 已绑定的 OAuth */
  linked_oauths: string[];
}

// ==================== 请求相关类型 ====================

/**
 * 请求方法
 */
export type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

/**
 * 请求配置
 */
export interface RequestConfig {
  /** 请求地址 */
  url: string;
  /** 请求方法 */
  method?: RequestMethod;
  /** 请求参数 */
  data?: Record<string, unknown>;
  /** 查询参数 */
  params?: Record<string, unknown>;
  /** 请求头 */
  header?: Record<string, string>;
  /** 是否显示加载提示 */
  loading?: boolean;
  /** 加载提示文字 */
  loadingText?: string;
  /** 是否需要 token */
  needAuth?: boolean;
}

/**
 * API 错误类型
 */
export interface ApiError {
  /** 错误码 */
  code: number;
  /** 错误信息 */
  message: string;
  /** 请求路径 */
  url?: string;
  /** HTTP 状态码 */
  statusCode?: number;
}

/**
 * 请求拦截器
 * 在请求发送前调用，可用于修改配置或添加公共参数
 */
export type RequestInterceptor = (config: RequestConfig) => RequestConfig | Promise<RequestConfig>;

/**
 * 响应拦截器 - 成功回调
 */
export type ResponseInterceptor<T = unknown> = (
  response: ApiResponse<T>,
  config: RequestConfig
) => ApiResponse<T>;

/**
 * 响应拦截器 - 错误回调
 */
export type ErrorInterceptor = (error: ApiError, config: RequestConfig) => void | Promise<void>;
