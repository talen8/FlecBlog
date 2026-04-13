import { getArticlesForWeb, getArticleBySlug } from '@/composables/api/article'
import { getComments, createComment, deleteComment } from '@/composables/api/comment'
import { flattenComments } from '@/composables/useComment'
import { getMoments } from '@/composables/api/moment'
import { getNotifications, markAsRead, markAllAsRead } from '@/composables/api/notification'
import { getCategories } from '@/composables/api/category'
import { getTags } from '@/composables/api/tag'
import type { Article, ArticleQuery } from '@@/types/article'
import type { Category } from '@@/types/category'
import type { Comment, CreateCommentParams, CommentTargetType } from '@@/types/comment'
import type { Menu } from '@@/types/menu'
import type { Moment } from '@@/types/moment'
import type { Notification, GetNotificationsParams } from '@@/types/notification'
import type { SiteStats } from '@@/types/stats'
import type { Tag } from '@@/types/tag'

export function useArticles() {
  const articles = useState<Article[]>('articles', () => [])
  const total = useState<number>('articles-total', () => 0)
  const currentPage = useState<number>('articles-currentPage', () => 1)
  const pageSize = useState<number>('articles-pageSize', () => 10)

  const fetchArticles = async (query: ArticleQuery = {}, forceRefresh = false) => {
    if (query.page) currentPage.value = query.page
    if (!forceRefresh && articles.value.length && !Object.keys(query).length) return

    try {
      const { list, total: resTotal } = await getArticlesForWeb({
        page: currentPage.value,
        page_size: pageSize.value,
        ...query
      })
      articles.value = list
      total.value = resTotal
    } catch (error) {
      console.error('获取文章列表失败:', error)
    }
  }

  return {
    articles,
    total,
    currentPage,
    pageSize,
    fetchArticles,
    setPageSize: (size: number) => {
      pageSize.value = size
      currentPage.value = 1
    },
    resetPage: () => currentPage.value = 1
  }
}

export function useCurrentArticle() {
  const currentArticle = useState<Article | null>('currentArticle', () => null)

  return {
    currentArticle,
    setCurrentArticle: (article: Article | null) => currentArticle.value = article,
    clearCurrentArticle: () => currentArticle.value = null
  }
}

export function useCategories() {
  const categories = useState<Category[]>('categories', () => [])
  const total = useState<number>('categories-total', () => 0)

  const fetchCategories = async (forceRefresh = false) => {
    if (!forceRefresh && categories.value.length) return

    try {
      const { list, total: resTotal } = await getCategories()
      categories.value = list
      total.value = resTotal
    } catch (error) {
      console.error('获取分类列表失败:', error)
    }
  }

  return { categories, total, fetchCategories }
}

export function useComments() {
  const comments = useState<Comment[]>('comments', () => [])
  const currentTargetType = useState<CommentTargetType | null>('comments-targetType', () => null)
  const currentTargetKey = useState<string | number | null>('comments-targetKey', () => null)
  const { articles } = useArticles()

  const fetchComments = async (targetType: CommentTargetType, targetKey: string | number) => {
    if (!targetType || !targetKey) return

    currentTargetType.value = targetType
    currentTargetKey.value = targetKey

    try {
      const data = await getComments({ target_type: targetType, target_key: targetKey })
      comments.value = data.list
    } catch (error) {
      console.error('获取评论失败:', error)
      comments.value = []
    }
  }

  const addComment = async (params: CreateCommentParams) => {
    const newComment = await createComment(params)

    if (!params.parent_id) {
      comments.value.unshift(newComment)
    } else {
      const addReplyToComment = (commentList: Comment[]): boolean => {
        for (const comment of commentList) {
          if (comment.id === params.parent_id) {
            if (!comment.replies) comment.replies = []
            comment.replies.push(newComment)
            return true
          }
          if (comment.replies?.length && addReplyToComment(comment.replies)) return true
        }
        return false
      }
      addReplyToComment(comments.value)
    }

    // 更新首页文章列表中的评论数量
    if (params.target_type === 'article') {
      const article = articles.value.find(a => a.slug === params.target_key)
      if (article) {
        article.comment_count = (article.comment_count || 0) + 1
      }
      // 刷新首页的SSR缓存数据
      refreshNuxtData('articles-list')
    }

    return newComment
  }

  const removeComment = async (commentId: number) => {
    await deleteComment(commentId)

    const removeFromList = (commentList: Comment[]): boolean => {
      const index = commentList.findIndex(c => c.id === commentId)
      if (index !== -1) {
        commentList.splice(index, 1)
        return true
      }
      for (const comment of commentList) {
        if (comment.replies?.length && removeFromList(comment.replies)) return true
      }
      return false
    }
    removeFromList(comments.value)

    if (currentTargetType.value === 'article' && currentTargetKey.value) {
      const article = articles.value.find(a => a.slug === currentTargetKey.value)
      if (article) {
        article.comment_count = Math.max(0, (article.comment_count || 0) - 1)
      }
      refreshNuxtData('articles-list')
    }
  }

  return {
    comments,
    fetchComments,
    addComment,
    removeComment,
    resetComments: () => {
      comments.value = []
      currentTargetType.value = null
      currentTargetKey.value = null
    },
    flattenComments
  }
}

export function useMenus() {
  const menus = useState<Menu[]>('menus', () => [])

  const filterByType = (type: string) =>
    menus.value.filter(menu => menu.type === type).sort((a, b) => a.sort - b.sort)

  const flatNavigationMenus = computed(() => {
    const result: Menu[] = []
    const flatten = (items: Menu[]) => {
      items.forEach(item => {
        if (item.type === 'navigation') {
          result.push(item)
          if (item.children?.length) flatten(item.children)
        }
      })
    }
    flatten(filterByType('navigation'))
    return result
  })

  return {
    menus,
    navigationMenus: computed(() => filterByType('navigation')),
    footerMenus: computed(() => filterByType('footer')),
    aggregateMenus: computed(() => filterByType('aggregate')),
    flatNavigationMenus
  }
}

export function useMoments() {
  const moments = useState<Moment[]>('moments', () => [])
  const total = useState<number>('moments-total', () => 0)
  const currentPage = useState<number>('moments-currentPage', () => 1)
  const pageSize = useState<number>('moments-pageSize', () => 30)

  const fetchMoments = async (page: number = 1, forceRefresh = false) => {
    if (page) currentPage.value = page
    if (!forceRefresh && moments.value.length) return

    try {
      const { list, total: resTotal, page: resPage, page_size: resPageSize } = await getMoments({
        page: currentPage.value,
        page_size: pageSize.value
      })
      moments.value = list
      total.value = resTotal
      currentPage.value = resPage
      pageSize.value = resPageSize
    } catch (error) {
      console.error('获取动态列表失败:', error)
      moments.value = []
      total.value = 0
    }
  }

  return { moments, total, currentPage, pageSize, fetchMoments }
}

export function useNotifications() {
  const notifications = useState<Notification[]>('notifications', () => [])
  const total = useState<number>('notifications-total', () => 0)
  const currentPage = useState<number>('notifications-currentPage', () => 1)
  const pageSize = useState<number>('notifications-pageSize', () => 10)
  const unreadCount = useState<number>('notifications-unreadCount', () => 0)
  const loading = useState<boolean>('notifications-loading', () => false)

  const fetchNotifications = async (params?: Partial<GetNotificationsParams>) => {
    loading.value = true
    try {
      const response = await getNotifications({
        page: params?.page ?? currentPage.value,
        page_size: params?.page_size ?? pageSize.value
      })
      notifications.value = response.list
      total.value = response.total
      unreadCount.value = response.unread_count
      params?.page && (currentPage.value = params.page)
    } catch (error) {
      console.error('获取通知列表失败:', error)
      notifications.value = []
      total.value = 0
      unreadCount.value = 0
    } finally {
      loading.value = false
    }
  }

  const markNotificationAsRead = async (id: number) => {
    try {
      await markAsRead(id)
      const notification = notifications.value.find(n => n.id === id)
      if (notification?.is_read === false) {
        notification.is_read = true
        notification.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (error) {
      console.error('标记通知已读失败:', error)
      throw error
    }
  }

  const markAllNotificationsAsRead = async () => {
    try {
      await markAllAsRead()
      notifications.value.forEach(n => {
        n.is_read = true
        n.read_at = new Date().toISOString()
      })
      unreadCount.value = 0
    } catch (error) {
      console.error('标记所有通知已读失败:', error)
      throw error
    }
  }

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
    resetPage: () => currentPage.value = 1,
    clearNotifications: () => {
      notifications.value = []
      total.value = 0
      currentPage.value = 1
      unreadCount.value = 0
    }
  }
}

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
    month_pageviews: 0
  }))

  return { siteStats }
}

export function useSysConfig() {
  const basicConfig = useState<Record<string, string>>('sysconfig-basic', () => ({
    'author': '',
    'author_email': '',
    'author_desc': '',
    'author_avatar': '',
    'author_photo': '',
    'icp': '',
    'police_record': '',
    'admin_url': '',
    'blog_url': '',
    'home_url': ''
  }))

  const blogConfig = useState<Record<string, string>>('sysconfig-blog', () => ({
    'title': 'FlecBLOG',
    'subtitle': 'FlecBLOG',
    'slogan': '',
    'description': '',
    'keywords': '',
    'established': '',
    'favicon': '',
    'background_image': '',
    'screenshot': '',
    'announcement': '',
    'typing_texts': '',
    'sidebar_social': '',
    'footer_social': '',
    'about_describe': '',
    'about_describe_tips': '',
    'about_exhibition': '',
    'about_profile': '',
    'about_personality': '',
    'about_motto_main': '',
    'about_motto_sub': '',
    'about_socialize': '',
    'about_creation': '',
    'about_versions': '',
    'about_unions': '',
    'about_story': '',
    'custom_head': '',
    'custom_body': '',
    'emojis': '',
    'font': ''
  }))

  const oauthConfig = useState<Record<string, string>>('sysconfig-oauth', () => ({
    'github.enabled': 'false',
    'google.enabled': 'false',
    'qq.enabled': 'false',
    'microsoft.enabled': 'false'
  }))

  const uploadConfig = useState<Record<string, string>>('sysconfig-upload', () => ({
    'max_file_size': '5'
  }))

  return {
    basicConfig,
    blogConfig,
    oauthConfig,
    uploadConfig
  }
}

export function useTags() {
  const tags = useState<Tag[]>('tags', () => [])
  const total = useState<number>('tags-total', () => 0)

  const fetchTags = async (forceRefresh = false) => {
    if (!forceRefresh && tags.value.length) return

    try {
      const { list, total: resTotal } = await getTags()
      tags.value = list
      total.value = resTotal
    } catch (error) {
      console.error('获取标签列表失败:', error)
    }
  }

  return { tags, total, fetchTags }
}
