/**
 * API 请求函数
 *
 * 本文件包含所有与后端 API 交互的函数，提供统一的 API 调用接口。
 * 使用工厂模式创建 API 实例，支持 CRUD 操作和自定义请求。
 *
 * @module composables/useApi
 */

import type {
  ApiResponse,
  PaginationData,
  PaginationQuery,
  Article,
  ArticleQuery,
  Category,
  Tag,
  Comment,
  CommentTargetType,
  CreateCommentParams,
  Moment,
  FriendGroupedResponse,
  FriendQueryParams,
  FriendApplyRequest,
  Feedback,
  SubmitFeedbackParams,
  Menu,
  NotificationListResponse,
  GetNotificationsParams,
  SiteStats,
  SettingGroupType,
  LoginParams,
  LoginResponse,
  RegisterParams,
  RegisterResponse,
  UserInfo,
  UpdateProfileParams,
  ForgotPasswordParams,
  ResetPasswordParams,
  ChangePasswordParams,
  DeactivateAccountParams,
  RefreshTokenResponse,
} from '@@/types';
import { get, post, put, patch, del } from '@/utils/request';

/**
 * API 工厂配置选项
 */
interface ApiFactoryOptions {
  /** 是否将 target_key 转换为字符串 */
  stringifyTargetKey?: boolean;
  /**
   * 参数转换函数
   * @param params - 原始参数
   * @returns 转换后的参数
   */
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 转换函数，数据类型不确定
  transformParams?: (params: any) => any;
  /**
   * 请求体转换函数
   * @param body - 原始请求体
   * @returns 转换后的请求体
   */
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 转换函数，数据类型不确定
  transformBody?: (body: any) => any;
}

/**
 * 处理请求数据
 *
 * 根据配置选项对数据进行预处理，包括类型转换和自定义转换。
 *
 * @param data - 原始数据
 * @param options - 工厂配置选项
 * @param isBody - 是否为请求体（true 为请求体，false 为查询参数）
 * @returns 处理后的数据
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用数据处理函数
function processData(data: any, options: ApiFactoryOptions, isBody: boolean) {
  let processed = { ...data };

  if (options.stringifyTargetKey && processed.target_key !== undefined) {
    processed.target_key = String(processed.target_key);
  }

  const transformFn = isBody ? options.transformBody : options.transformParams;
  if (transformFn) {
    processed = transformFn(processed);
  }

  return processed;
}

/**
 * 创建 API 实例
 *
 * 工厂函数，为指定端点创建完整的 CRUD API 实例。
 *
 * @template T - 数据实体类型
 * @param endpoint - API 端点路径
 * @param options - 配置选项
 * @returns API 实例对象
 *
 * @example
 * ```ts
 * const articleApi = createApi<Article>('/articles');
 *
 * // 获取列表
 * const { list, total } = await articleApi.getList({ page: 1, page_size: 10 });
 *
 * // 获取单个
 * const article = await articleApi.getOne('my-article-slug');
 *
 * // 创建
 * const newArticle = await articleApi.create({ title: '新文章', content: '...' });
 *
 * // 更新
 * await articleApi.update(1, { title: '更新标题' });
 *
 * // 删除
 * await articleApi.delete(1);
 * ```
 */
function createApi<T>(endpoint: string, options: ApiFactoryOptions = {}) {
  return {
    /**
     * 获取分页列表
     * @param params - 分页查询参数
     * @returns 分页数据
     */
    getList: async (params?: Partial<PaginationQuery>): Promise<PaginationData<T>> => {
      const response = await get<ApiResponse<PaginationData<T>>>(endpoint, {
        params: processData(params, options, false),
      });
      return response.data;
    },

    /**
     * 获取单个实体
     * @param id - 实体 ID 或别名
     * @returns 实体数据
     */
    getOne: async (id: number | string): Promise<T> => {
      const response = await get<ApiResponse<T>>(`${endpoint}/${id}`);
      return response.data;
    },

    /**
     * 创建实体
     * @param data - 实体数据
     * @returns 创建的实体
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用 CRUD 方法，数据类型由调用方决定
    create: async (data: any): Promise<T> => {
      const response = await post<ApiResponse<T>>(endpoint, processData(data, options, true));
      return response.data;
    },

    /**
     * 更新实体（全量更新）
     * @param id - 实体 ID
     * @param data - 更新数据
     * @returns 更新后的实体
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用 CRUD 方法，数据类型由调用方决定
    update: async (id: number | string, data: any): Promise<T> => {
      const response = await put<ApiResponse<T>>(
        `${endpoint}/${id}`,
        processData(data, options, true)
      );
      return response.data;
    },

    /**
     * 部分更新实体
     * @param id - 实体 ID
     * @param data - 更新数据
     * @returns 更新后的实体
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用 CRUD 方法，数据类型由调用方决定
    patch: async (id: number | string, data: any): Promise<T> => {
      const response = await patch<ApiResponse<T>>(
        `${endpoint}/${id}`,
        processData(data, options, true)
      );
      return response.data;
    },

    /**
     * 删除实体
     * @param id - 实体 ID
     */
    delete: async (id: number | string): Promise<void> => {
      await del(`${endpoint}/${id}`);
    },

    /**
     * 自定义 GET 请求
     * @template R - 响应数据类型
     * @param url - 相对 URL 路径
     * @param params - 查询参数
     * @returns 响应数据
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用请求方法，参数类型不确定
    get: async <R = T>(url: string, params?: any): Promise<R> => {
      const response = await get<ApiResponse<R>>(`${endpoint}${url}`, {
        params: processData(params, options, false),
      });
      return response.data;
    },

    /**
     * 自定义 POST 请求
     * @template R - 响应数据类型
     * @param url - 相对 URL 路径
     * @param data - 请求体数据
     * @returns 响应数据
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用请求方法，数据类型由调用方决定
    post: async <R = T>(url: string, data?: any): Promise<R> => {
      const response = await post<ApiResponse<R>>(
        `${endpoint}${url}`,
        processData(data, options, true)
      );
      return response.data;
    },

    /**
     * 自定义 PUT 请求
     * @template R - 响应数据类型
     * @param url - 相对 URL 路径
     * @param data - 请求体数据
     * @returns 响应数据
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用请求方法，数据类型由调用方决定
    put: async <R = T>(url: string, data?: any): Promise<R> => {
      const response = await put<ApiResponse<R>>(
        `${endpoint}${url}`,
        processData(data, options, true)
      );
      return response.data;
    },

    /**
     * 自定义 PATCH 请求
     * @template R - 响应数据类型
     * @param url - 相对 URL 路径
     * @param data - 请求体数据
     * @returns 响应数据
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用请求方法，数据类型由调用方决定
    patchRequest: async <R = T>(url: string, data?: any): Promise<R> => {
      const response = await patch<ApiResponse<R>>(
        `${endpoint}${url}`,
        processData(data, options, true)
      );
      return response.data;
    },

    /**
     * 自定义 DELETE 请求
     * @param url - 相对 URL 路径
     * @param data - 请求体数据
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 通用请求方法，数据类型由调用方决定
    deleteRequest: async (url: string, data?: any): Promise<void> => {
      await del(`${endpoint}${url}`, {
        body: processData(data, options, true),
      });
    },
  };
}

/* ============================================
 * API 实例创建
 * ============================================ */

const articleApi = createApi<Article>('/articles');
const categoryApi = createApi<Category>('/categories');
const tagApi = createApi<Tag>('/tags');
const commentApi = createApi<Comment>('/comments', { stringifyTargetKey: true });
const momentApi = createApi<Moment>('/moments');
const friendApi = createApi<FriendGroupedResponse>('');
const feedbackApi = createApi<Feedback>('');
const statsApi = createApi<SiteStats>('');
const menuApi = createApi<Menu[]>('');
const notificationApi = createApi<NotificationListResponse>('');
const settingApi = createApi<Record<string, string>>('');
const authApi = createApi<LoginResponse>('/auth');
const userApi = createApi<UserInfo>('');

/* ============================================
 * Composable API（推荐使用）
 * ============================================ */

/**
 * API Composable
 *
 * 提供所有 API 方法的统一入口，按模块组织。
 * 推荐在组件中使用此 Composable。
 *
 * @returns 按模块组织的 API 方法对象
 *
 * @example
 * ```ts
 * const api = useApi();
 *
 * // 文章相关
 * const articles = await api.article.getList({ page: 1 });
 * const article = await api.article.getBySlug('my-article');
 *
 * // 评论相关
 * await api.comment.create({
 *   target_type: 'article',
 *   target_key: 'my-article',
 *   content: '这是一条评论',
 * });
 * ```
 */
export function useApi() {
  return {
    /** 文章 API */
    article: {
      /** 获取文章列表 */
      getList: async (params: ArticleQuery = {}) => articleApi.getList(params),
      /** 根据别名获取文章 */
      getBySlug: async (slug: string) => articleApi.getOne(slug),
      /** 搜索文章 */
      search: async (keyword: string, params: Partial<ArticleQuery> = {}) =>
        articleApi.get<PaginationData<Article>>('/search', { keyword, ...params }),
    },

    /** 分类 API */
    category: {
      /** 获取分类列表 */
      getList: async () => categoryApi.getList(),
      /** 根据别名获取分类 */
      getBySlug: async (slug: string) => categoryApi.getOne(slug),
    },

    /** 标签 API */
    tag: {
      /** 获取标签列表 */
      getList: async () => tagApi.getList(),
      /** 根据别名获取标签 */
      getBySlug: async (slug: string) => tagApi.getOne(slug),
    },

    /** 评论 API */
    comment: {
      /** 获取评论列表 */
      getList: async (
        params: { target_type: CommentTargetType; target_key: string | number } & PaginationQuery
      ) => commentApi.getList(params),
      /** 创建评论 */
      create: async (params: CreateCommentParams) => commentApi.create(params),
      /** 删除评论 */
      delete: async (id: number) => commentApi.delete(id),
    },

    /** 动态 API */
    moment: {
      /** 获取动态列表 */
      getList: async (params: PaginationQuery = {}) => momentApi.getList(params),
    },

    /** 友链 API */
    friend: {
      /** 获取友链列表 */
      getList: async (params?: FriendQueryParams) => friendApi.get('/friends', params),
      /** 申请友链 */
      apply: async (data: FriendApplyRequest) => friendApi.post('/friends/apply', data),
    },

    /** 反馈 API */
    feedback: {
      /** 提交反馈 */
      submit: async (data: SubmitFeedbackParams): Promise<Feedback> =>
        feedbackApi.post<Feedback>('/feedback', data),
      /** 根据工单号查询反馈 */
      getByTicketNo: async (ticketNo: string): Promise<Feedback> =>
        feedbackApi.get<Feedback>(`/feedback/ticket/${ticketNo}`),
    },

    /** 统计 API */
    stats: {
      /** 获取网站统计数据 */
      getSite: async () => statsApi.get<SiteStats>('/stats/site'),
    },

    /** 菜单 API */
    menu: {
      /** 获取菜单列表 */
      getList: async () => menuApi.get<Menu[]>('/menus'),
    },

    /** 通知 API */
    notification: {
      /** 获取通知列表 */
      getList: async (params: GetNotificationsParams) =>
        notificationApi.get<NotificationListResponse>('/notifications', params),
      /** 标记单条通知已读 */
      markAsRead: async (id: number) => notificationApi.put(`/notifications/${id}/read`),
      /** 标记所有通知已读 */
      markAllAsRead: async () => notificationApi.put('/notifications/read-all'),
    },

    /** 配置 API */
    setting: {
      /** 获取指定分组的配置 */
      getGroup: async (group: SettingGroupType) =>
        settingApi.get<Record<string, string>>(`/settings/${group}`),
    },

    /** 订阅 API */
    subscribe: {
      /** 邮箱订阅 */
      email: async (email: string) => post<ApiResponse<void>>('/subscribe', { email }),
    },

    /** 认证 API */
    auth: {
      /** 登录 */
      login: async (data: LoginParams) => authApi.post<LoginResponse>('/login', data),
      /** 注册 */
      register: async (data: RegisterParams) => authApi.post<RegisterResponse>('/register', data),
      /** 刷新 Token */
      refreshToken: async () => authApi.post<RefreshTokenResponse>('/refresh'),
      /** 忘记密码 */
      forgotPassword: async (data: ForgotPasswordParams) => authApi.post('/forgot-password', data),
      /** 重置密码 */
      resetPassword: async (data: ResetPasswordParams) => authApi.post('/reset-password', data),
    },

    /** 用户 API */
    user: {
      /** 获取用户信息 */
      getProfile: async () => userApi.get<UserInfo>('/user/profile'),
      /** 更新用户信息 */
      updateProfile: async (data: UpdateProfileParams) =>
        userApi.patchRequest<UserInfo>('/user/profile', data),
      /** 修改密码 */
      changePassword: async (data: ChangePasswordParams) => userApi.put('/user/password', data),
      /** 注销账号 */
      deactivate: async (data: DeactivateAccountParams) =>
        userApi.deleteRequest('/user/deactivate', data),
      /** 解绑 OAuth */
      unbindOAuth: async (provider: string) => userApi.deleteRequest(`/user/oauth/${provider}`),
    },
  };
}

/* ============================================
 * 独立导出函数（兼容旧代码）
 * ============================================ */

/** 获取文章列表（前台展示用） */
export const getArticlesForWeb = async (params: ArticleQuery = {}) => articleApi.getList(params);
/** 根据别名获取文章 */
export const getArticleBySlug = async (slug: string) => articleApi.getOne(slug);
/** 搜索文章 */
export const searchArticles = async (keyword: string, params: Partial<ArticleQuery> = {}) =>
  articleApi.get<PaginationData<Article>>('/search', { keyword, ...params });

/** 获取分类列表 */
export const getCategories = async () => categoryApi.getList();
/** 根据别名获取分类 */
export const getCategoryBySlug = async (slug: string) => categoryApi.getOne(slug);

/** 获取标签列表 */
export const getTags = async () => tagApi.getList();
/** 根据别名获取标签 */
export const getTagBySlug = async (slug: string) => tagApi.getOne(slug);

/** 获取评论参数 */
interface GetCommentsParams extends PaginationQuery {
  target_type: CommentTargetType;
  target_key: string | number;
}
/** 获取评论列表 */
export const getComments = async (params: GetCommentsParams) => commentApi.getList(params);
/** 创建评论 */
export const createComment = async (params: CreateCommentParams) => commentApi.create(params);
/** 删除评论 */
export const deleteComment = async (id: number) => commentApi.delete(id);

/** 获取动态列表 */
export const getMoments = async (params: PaginationQuery = {}) => momentApi.getList(params);

/** 获取友链列表 */
export const getFriends = async (params?: FriendQueryParams) => friendApi.get('/friends', params);
/** 申请友链 */
export const applyFriend = async (data: FriendApplyRequest) =>
  friendApi.post('/friends/apply', data);

/** 提交反馈 */
export const submitFeedback = async (data: SubmitFeedbackParams): Promise<Feedback> =>
  feedbackApi.post<Feedback>('/feedback', data);
/** 根据工单号查询反馈 */
export const getFeedbackByTicketNo = async (ticketNo: string): Promise<Feedback> =>
  feedbackApi.get<Feedback>(`/feedback/ticket/${ticketNo}`);

/** 获取网站统计数据 */
export const getSiteStats = async () => statsApi.get<SiteStats>('/stats/site');

/** 获取菜单列表 */
export const getMenus = async () => menuApi.get<Menu[]>('/menus');

/** 获取通知列表 */
export const getNotifications = async (params: GetNotificationsParams) =>
  notificationApi.get<NotificationListResponse>('/notifications', params);
/** 标记单条通知已读 */
export const markAsRead = async (id: number) => notificationApi.put(`/notifications/${id}/read`);
/** 标记所有通知已读 */
export const markAllAsRead = async () => notificationApi.put('/notifications/read-all');

/** 获取指定分组的配置 */
export const getSettingGroup = async (group: SettingGroupType) =>
  settingApi.get<Record<string, string>>(`/settings/${group}`);

/** 邮箱订阅 */
export const subscribe = async (email: string) => post<ApiResponse<void>>('/subscribe', { email });

/** 登录 */
export const login = async (data: LoginParams) => authApi.post<LoginResponse>('/login', data);
/** 注册 */
export const register = async (data: RegisterParams) =>
  authApi.post<RegisterResponse>('/register', data);
/** 刷新 Token */
export const refreshToken = async () => authApi.post<RefreshTokenResponse>('/refresh');
/** 忘记密码 */
export const forgotPassword = async (data: ForgotPasswordParams) =>
  authApi.post('/forgot-password', data);
/** 重置密码 */
export const resetPassword = async (data: ResetPasswordParams) =>
  authApi.post('/reset-password', data);

/** 获取用户信息 */
export const getUserProfile = async () => userApi.get<UserInfo>('/user/profile');
/** 更新用户信息 */
export const updateUserProfile = async (data: UpdateProfileParams) =>
  userApi.patchRequest<UserInfo>('/user/profile', data);
/** 修改密码 */
export const changePassword = async (data: ChangePasswordParams) =>
  userApi.put('/user/password', data);
/** 注销账号 */
export const deactivateAccount = async (data: DeactivateAccountParams) =>
  userApi.deleteRequest('/user/deactivate', data);
/** 解绑 OAuth */
export const unbindOAuth = async (provider: string) =>
  userApi.deleteRequest(`/user/oauth/${provider}`);
