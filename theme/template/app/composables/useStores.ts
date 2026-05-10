/**
 * 全局状态管理 Composables
 *
 * 本文件包含所有全局状态的组合式函数，用于管理应用的核心数据状态。
 * 使用 Nuxt 的 useState 实现跨组件的状态共享。
 *
 * @module composables/useStores
 */

import {
  getArticlesForWeb,
  getComments,
  createComment,
  deleteComment,
  getMoments,
  getNotifications,
  markAsRead,
  markAllAsRead,
  getCategories,
  getTags,
  getSiteStats,
} from '@/composables/useApi';
import type {
  Article,
  ArticleQuery,
  Category,
  Comment,
  CreateCommentParams,
  CommentTargetType,
  Menu,
  Moment,
  Notification,
  GetNotificationsParams,
  SiteStats,
  Tag,
} from '@@/types';

/**
 * 文章状态管理
 *
 * 管理文章列表、分页信息等状态，提供文章数据获取和操作方法。
 *
 * @returns {Object} 文章状态和方法
 * @returns {Ref<Article[]>} articles - 文章列表
 * @returns {Ref<number>} total - 文章总数
 * @returns {Ref<number>} currentPage - 当前页码
 * @returns {Ref<number>} pageSize - 每页文章数
 * @returns {Function} fetchArticles - 获取文章列表
 * @returns {Function} setPageSize - 设置每页文章数
 * @returns {Function} resetPage - 重置页码
 *
 * @example
 * ```ts
 * const { articles, total, fetchArticles } = useArticles();
 * await fetchArticles({ page: 1, category: 'tech' });
 * ```
 */
export function useArticles() {
  const articles = useState<Article[]>('articles', () => []);
  const total = useState<number>('articles-total', () => 0);
  const currentPage = useState<number>('articles-currentPage', () => 1);
  const pageSize = useState<number>('articles-pageSize', () => 10);

  /**
   * 获取文章列表
   * @param query - 查询参数
   * @param forceRefresh - 是否强制刷新缓存
   */
  const fetchArticles = async (query: ArticleQuery = {}, forceRefresh = false) => {
    if (query.page) currentPage.value = query.page;
    if (!forceRefresh && articles.value.length && !Object.keys(query).length) return;

    try {
      const { list, total: resTotal } = await getArticlesForWeb({
        page: currentPage.value,
        page_size: pageSize.value,
        ...query,
      });
      articles.value = list || [];
      total.value = resTotal || 0;
    } catch (error) {
      console.error('获取文章列表失败:', error);
      articles.value = [];
      total.value = 0;
    }
  };

  return {
    articles,
    total,
    currentPage,
    pageSize,
    fetchArticles,
    setPageSize: (size: number) => {
      pageSize.value = size;
      currentPage.value = 1;
    },
    resetPage: () => (currentPage.value = 1),
  };
}

/**
 * 分类状态管理
 *
 * 管理分类列表状态，提供分类数据获取方法。
 *
 * @returns {Object} 分类状态和方法
 * @returns {Ref<Category[]>} categories - 分类列表
 * @returns {Ref<number>} total - 分类总数
 * @returns {Function} fetchCategories - 获取分类列表
 *
 * @example
 * ```ts
 * const { categories, fetchCategories } = useCategories();
 * await fetchCategories();
 * ```
 */
export function useCategories() {
  const categories = useState<Category[]>('categories', () => []);
  const total = useState<number>('categories-total', () => 0);

  /**
   * 获取分类列表
   * @param forceRefresh - 是否强制刷新缓存
   */
  const fetchCategories = async (forceRefresh = false) => {
    if (!forceRefresh && categories.value.length) return;

    try {
      const { list, total: resTotal } = await getCategories();
      categories.value = list || [];
      total.value = resTotal || 0;
    } catch (error) {
      console.error('获取分类列表失败:', error);
      categories.value = [];
      total.value = 0;
    }
  };

  return { categories, total, fetchCategories };
}

/**
 * 评论状态管理
 *
 * 管理评论列表状态，提供评论的增删改查方法。
 * 支持嵌套评论（回复）的管理。
 *
 * @returns {Object} 评论状态和方法
 * @returns {Ref<Comment[]>} comments - 评论列表
 * @returns {Function} fetchComments - 获取评论列表
 * @returns {Function} addComment - 添加评论
 * @returns {Function} removeComment - 删除评论
 * @returns {Function} resetComments - 重置评论状态
 *
 * @example
 * ```ts
 * const { comments, addComment, removeComment } = useComments();
 * await addComment({
 *   target_type: 'article',
 *   target_key: 'my-article-slug',
 *   content: '这是一条评论',
 * });
 * ```
 */
export function useComments() {
  const comments = useState<Comment[]>('comments', () => []);
  const currentTargetType = useState<CommentTargetType | null>('comments-targetType', () => null);
  const currentTargetKey = useState<string | number | null>('comments-targetKey', () => null);
  const { articles } = useArticles();

  /**
   * 获取评论列表
   * @param targetType - 评论目标类型（article/page/moment）
   * @param targetKey - 评论目标标识
   */
  const fetchComments = async (targetType: CommentTargetType, targetKey: string | number) => {
    if (!targetType || !targetKey) return;

    currentTargetType.value = targetType;
    currentTargetKey.value = targetKey;

    try {
      const data = await getComments({
        target_type: targetType,
        target_key: targetKey,
      });
      comments.value = data.list || [];
    } catch (error) {
      console.error('获取评论失败:', error);
      comments.value = [];
    }
  };

  /**
   * 添加评论
   * @param params - 评论参数
   * @returns 新创建的评论
   */
  const addComment = async (params: CreateCommentParams) => {
    const newComment = await createComment(params);

    if (!params.parent_id) {
      comments.value.unshift(newComment);
    } else {
      const addReplyToComment = (commentList: Comment[]): boolean => {
        for (const comment of commentList) {
          if (comment.id === params.parent_id) {
            if (!comment.replies) comment.replies = [];
            comment.replies.push(newComment);
            return true;
          }
          if (comment.replies?.length && addReplyToComment(comment.replies)) return true;
        }
        return false;
      };
      addReplyToComment(comments.value);
    }

    if (params.target_type === 'article') {
      const article = articles.value.find(a => a.slug === params.target_key);
      if (article) {
        article.comment_count = (article.comment_count || 0) + 1;
      }
      refreshNuxtData('articles-list');
    }

    return newComment;
  };

  /**
   * 删除评论
   * @param commentId - 评论ID
   */
  const removeComment = async (commentId: number) => {
    await deleteComment(commentId);

    const removeFromList = (commentList: Comment[]): boolean => {
      const index = commentList.findIndex(c => c.id === commentId);
      if (index !== -1) {
        commentList.splice(index, 1);
        return true;
      }
      for (const comment of commentList) {
        if (comment.replies?.length && removeFromList(comment.replies)) return true;
      }
      return false;
    };
    removeFromList(comments.value);

    if (currentTargetType.value === 'article' && currentTargetKey.value) {
      const article = articles.value.find(a => a.slug === currentTargetKey.value);
      if (article) {
        article.comment_count = Math.max(0, (article.comment_count || 0) - 1);
      }
      refreshNuxtData('articles-list');
    }
  };

  return {
    comments,
    fetchComments,
    addComment,
    removeComment,
    resetComments: () => {
      comments.value = [];
      currentTargetType.value = null;
      currentTargetKey.value = null;
    },
  };
}

/**
 * 菜单状态管理
 *
 * 管理导航菜单、页脚菜单、聚合菜单等状态。
 * 支持菜单的层级结构和扁平化访问。
 *
 * @returns {Object} 菜单状态和方法
 * @returns {Ref<Menu[]>} menus - 所有菜单
 * @returns {ComputedRef<Menu[]>} navigationMenus - 导航菜单
 * @returns {ComputedRef<Menu[]>} footerMenus - 页脚菜单
 * @returns {ComputedRef<Menu[]>} aggregateMenus - 聚合菜单
 * @returns {ComputedRef<Menu[]>} flatNavigationMenus - 扁平化的导航菜单
 *
 * @example
 * ```ts
 * const { navigationMenus, footerMenus } = useMenus();
 * ```
 */
export function useMenus() {
  const menus = useState<Menu[]>('menus', () => []);

  const filterByType = (type: string) =>
    menus.value.filter(menu => menu.type === type).sort((a, b) => a.sort - b.sort);

  const flatNavigationMenus = computed(() => {
    const result: Menu[] = [];
    const flatten = (items: Menu[]) => {
      items.forEach(item => {
        if (item.type === 'navigation') {
          result.push(item);
          if (item.children?.length) flatten(item.children);
        }
      });
    };
    flatten(filterByType('navigation'));
    return result;
  });

  return {
    menus,
    navigationMenus: computed(() => filterByType('navigation')),
    footerMenus: computed(() => filterByType('footer')),
    aggregateMenus: computed(() => filterByType('aggregate')),
    flatNavigationMenus,
  };
}

/**
 * 动态状态管理
 *
 * 管理动态（说说）列表状态，提供动态数据获取方法。
 *
 * @returns {Object} 动态状态和方法
 * @returns {Ref<Moment[]>} moments - 动态列表
 * @returns {Ref<number>} total - 动态总数
 * @returns {Ref<number>} currentPage - 当前页码
 * @returns {Ref<number>} pageSize - 每页动态数
 * @returns {Function} fetchMoments - 获取动态列表
 *
 * @example
 * ```ts
 * const { moments, fetchMoments } = useMoments();
 * await fetchMoments(1);
 * ```
 */
export function useMoments() {
  const moments = useState<Moment[]>('moments', () => []);
  const total = useState<number>('moments-total', () => 0);
  const currentPage = useState<number>('moments-currentPage', () => 1);
  const { blogConfig } = useSysConfig();
  const pageSize = useState<number>('moments-pageSize', () => {
    const configSize = parseInt(blogConfig.value['blog.moments_size'] || '30');
    return configSize > 0 ? configSize : 30;
  });

  /**
   * 获取动态列表
   * @param page - 页码
   * @param forceRefresh - 是否强制刷新缓存
   */
  const fetchMoments = async (page: number = 1, forceRefresh = false) => {
    if (page) currentPage.value = page;
    if (!forceRefresh && moments.value.length) return;

    try {
      const {
        list,
        total: resTotal,
        page: resPage,
        page_size: resPageSize,
      } = await getMoments({
        page: currentPage.value,
        page_size: pageSize.value,
      });
      moments.value = list || [];
      total.value = resTotal || 0;
      currentPage.value = resPage || 1;
      pageSize.value = resPageSize || 30;
    } catch (error) {
      console.error('获取动态列表失败:', error);
      moments.value = [];
      total.value = 0;
    }
  };

  return { moments, total, currentPage, pageSize, fetchMoments };
}

/**
 * 通知状态管理
 *
 * 管理用户通知列表状态，提供通知获取和已读标记方法。
 *
 * @returns {Object} 通知状态和方法
 * @returns {Ref<Notification[]>} notifications - 通知列表
 * @returns {Ref<number>} total - 通知总数
 * @returns {Ref<number>} currentPage - 当前页码
 * @returns {Ref<number>} pageSize - 每页通知数
 * @returns {Ref<number>} unreadCount - 未读通知数
 * @returns {Ref<boolean>} loading - 加载状态
 * @returns {Function} fetchNotifications - 获取通知列表
 * @returns {Function} markNotificationAsRead - 标记单条通知已读
 * @returns {Function} markAllNotificationsAsRead - 标记所有通知已读
 *
 * @example
 * ```ts
 * const { notifications, unreadCount, markNotificationAsRead } = useNotifications();
 * await markNotificationAsRead(123);
 * ```
 */
export function useNotifications() {
  const notifications = useState<Notification[]>('notifications', () => []);
  const total = useState<number>('notifications-total', () => 0);
  const currentPage = useState<number>('notifications-currentPage', () => 1);
  const pageSize = useState<number>('notifications-pageSize', () => 10);
  const unreadCount = useState<number>('notifications-unreadCount', () => 0);
  const loading = useState<boolean>('notifications-loading', () => false);

  /**
   * 获取通知列表
   * @param params - 查询参数
   */
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

  /**
   * 标记单条通知为已读
   * @param id - 通知ID
   */
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

  /**
   * 标记所有通知为已读
   */
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

/**
 * 网站统计状态管理
 *
 * 管理网站统计数据状态，如文章数、访客数、浏览量等。
 *
 * @returns {Object} 统计状态和方法
 * @returns {Ref<SiteStats>} siteStats - 网站统计数据
 * @returns {Function} fetchStats - 获取统计数据
 *
 * @example
 * ```ts
 * const { siteStats, fetchStats } = useStats();
 * await fetchStats();
 * console.log(siteStats.value.total_articles);
 * ```
 */
export function useStats() {
  const siteStats = useState<SiteStats>('siteStats', () => ({
    total_words: '0',
    total_visitors: 0,
    total_page_views: 0,
    online_users: 0,
    total_articles: 0,
    total_comments: 0,
    total_friends: 0,
    total_moments: 0,
    total_categories: 0,
    total_tags: 0,
    today_visitors: 0,
    today_pageviews: 0,
    yesterday_visitors: 0,
    yesterday_pageviews: 0,
    month_pageviews: 0,
  }));

  /**
   * 获取站点统计数据
   * @param forceRefresh - 是否强制刷新
   */
  const fetchStats = async (forceRefresh = false) => {
    if (!forceRefresh && siteStats.value.total_visitors > 0) return;

    try {
      const stats = await getSiteStats();
      siteStats.value = stats;
    } catch (error) {
      console.error('获取站点统计失败:', error);
    }
  };

  return { siteStats, fetchStats };
}

/**
 * 系统配置状态管理
 *
 * 管理系统各类配置项状态，包括基础配置、博客配置、OAuth配置、上传配置等。
 * 配置项以键值对形式存储，通过配置键获取对应值。
 *
 * @returns {Object} 配置状态
 * @returns {Ref<Record<string, string>>} basicConfig - 基础配置
 * @returns {Ref<Record<string, string>>} blogConfig - 博客配置
 * @returns {Ref<Record<string, string>>} oauthConfig - OAuth 配置
 * @returns {Ref<Record<string, string>>} uploadConfig - 上传配置
 *
 * @example
 * ```ts
 * const { blogConfig, basicConfig } = useSysConfig();
 * const siteTitle = blogConfig.value['blog.title'];
 * const authorName = basicConfig.value['basic.author'];
 * ```
 */
export function useSysConfig() {
  const basicConfig = useState<Record<string, string>>('sysconfig-basic', () => ({
    'basic.author': '',
    'basic.author_email': '',
    'basic.author_desc': '',
    'basic.author_avatar': '',
    'basic.author_photo': '',
    'basic.icp': '',
    'basic.police_record': '',
    'basic.admin_url': '',
    'basic.blog_url': '',
    'basic.home_url': '',
  }));

  const blogConfig = useState<Record<string, string>>('sysconfig-blog', () => ({
    'blog.title': 'FlecBLOG',
    'blog.subtitle': 'FlecBLOG',
    'blog.slogan': '',
    'blog.description': '',
    'blog.keywords': '',
    'blog.established': '',
    'blog.favicon': '',
    'blog.background_image': '',
    'blog.screenshot': '',
    'blog.announcement': '',
    'blog.typing_texts': '',
    'blog.sidebar_social': '',
    'blog.footer_social': '',
    'blog.footer_links': '',
    'blog.about_describe': '',
    'blog.about_describe_tips': '',
    'blog.about_exhibition': '',
    'blog.about_profile': '',
    'blog.about_personality': '',
    'blog.about_motto_main': '',
    'blog.about_motto_sub': '',
    'blog.about_socialize': '',
    'blog.about_creation': '',
    'blog.about_versions': '',
    'blog.about_unions': '',
    'blog.about_story': '',
    'blog.custom_head': '',
    'blog.custom_body': '',
    'blog.emojis': '',
    'blog.font': '',
    'blog.moments_size': '30',
    'blog.message_content': '',
    'blog.home_layout': 'waterfall',
  }));

  const oauthConfig = useState<Record<string, string>>('sysconfig-oauth', () => ({
    'oauth.github.enabled': 'false',
    'oauth.google.enabled': 'false',
    'oauth.qq.enabled': 'false',
    'oauth.microsoft.enabled': 'false',
  }));

  const uploadConfig = useState<Record<string, string>>('sysconfig-upload', () => ({
    'upload.max_file_size': '5',
  }));

  return {
    basicConfig,
    blogConfig,
    oauthConfig,
    uploadConfig,
  };
}

/**
 * 标签状态管理
 *
 * 管理标签列表状态，提供标签数据获取方法。
 *
 * @returns {Object} 标签状态和方法
 * @returns {Ref<Tag[]>} tags - 标签列表
 * @returns {Ref<number>} total - 标签总数
 * @returns {Function} fetchTags - 获取标签列表
 *
 * @example
 * ```ts
 * const { tags, fetchTags } = useTags();
 * await fetchTags();
 * ```
 */
export function useTags() {
  const tags = useState<Tag[]>('tags', () => []);
  const total = useState<number>('tags-total', () => 0);

  /**
   * 获取标签列表
   * @param forceRefresh - 是否强制刷新缓存
   */
  const fetchTags = async (forceRefresh = false) => {
    if (!forceRefresh && tags.value.length) return;

    try {
      const { list, total: resTotal } = await getTags();
      tags.value = list || [];
      total.value = resTotal || 0;
    } catch (error) {
      console.error('获取标签列表失败:', error);
      tags.value = [];
      total.value = 0;
    }
  };

  return { tags, total, fetchTags };
}
