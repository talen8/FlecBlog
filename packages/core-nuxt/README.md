# @flecblog/core-nuxt

FlecBlog 主题开发核心包。创建主题时只需 `extends` 此包，即可获得自动导入的函数、开箱即用的文章排版、完整的 TypeScript 类型。

## 快速开始

```bash
npm install @flecblog/core-nuxt
```

在主题的 `nuxt.config.ts` 中：

```ts
export default defineNuxtConfig({
  extends: ['@flecblog/core-nuxt'],
})
```

## 函数

安装后在 `.vue` / `.ts` 文件中直接调用，无需 import。

### 文章

| 函数 | 说明 | 示例 |
|---|---|---|
| `useCurrentArticle` | 获取当前文章 | `const { currentArticle } = useCurrentArticle()` |
| `useArticle` | 按 slug 获取文章 | `const { data: article, refresh } = useArticle(slug)` |
| `useArticleList` | 分页获取文章列表 | `const { data: articles, page } = useArticleList()` |
| `useSearch` | 搜索文章 | `const { results, search } = useSearch()` |
| `useRandomArticle` | 获取随机文章 | `const { go } = useRandomArticle()` |
| `getArticleList` | API：文章列表 | `await getArticleList({ page: 1 })` |
| `getArticleBySlug` | API：按 slug 获取 | `await getArticleBySlug('hello-world')` |
| `searchArticles` | API：搜索 | `await searchArticles({ keyword: 'vue' })` |
| `getRandomArticleSlug` | API：随机 slug | `await getRandomArticleSlug()` |

### 认证

| 函数 | 说明 | 示例 |
|---|---|---|
| `isLoggedIn` | 响应式登录状态 | `v-if="isLoggedIn"` |
| `useAuth` | 获取登录状态 | `const isAuth = useAuth()` |
| `useAuthActions` | 登录/注册/找回密码 | `const { loginUser } = useAuthActions()` |
| `useLoginModal` | 控制登录弹窗 | `const { open } = useLoginModal()` |
| `useBindEmail` | 邮箱绑定提示 | `const { triggerGlobal } = useBindEmail()` |
| `accessToken` | 当前 token（ref） | `headers['Authorization'] = \`Bearer ${accessToken.value}\`` |
| `setAccessToken` | 设置 token | `setAccessToken(token)` |
| `getAccessToken` | 读取 token | `const t = getAccessToken()` |
| `clearAccessToken` | 清除 token | `clearAccessToken()` |
| `logoutUser` | 登出 | `logoutUser()` |
| `login` | API：登录 | `await login({ email, password })` |
| `register` | API：注册 | `await register({ email, nickname, password })` |
| `refreshToken` | API：刷新 token | `await refreshToken()` |
| `forgotPassword` | API：发送验证码 | `await forgotPassword({ email })` |
| `resetPassword` | API：重置密码 | `await resetPassword({ email, code, password })` |
| `logout` | API：服务端登出 | `await logout()` |
| `getOAuthUrl` | API：OAuth 登录地址 | `getOAuthUrl('github', redirect)` |
| `getOAuthBindUrl` | API：OAuth 绑定地址 | `getOAuthBindUrl('github', redirect)` |
| `getWechatQrcode` | API：微信扫码登录 | `const { blob, scene } = await getWechatQrcode()` |
| `getWechatScene` | API：轮询扫码状态 | `await getWechatScene(scene)` |

### 用户

| 函数 | 说明 | 示例 |
|---|---|---|
| `useUser` | 用户信息 + 拉取 | `const { userInfo, fetchUserInfo } = useUser()` |
| `useAvatar` | 头像 | `const { avatarUrl } = useAvatar()` |
| `getUserProfile` | API：获取用户信息 | `await getUserProfile()` |
| `updateUserProfile` | API：更新用户信息 | `await updateUserProfile({ nickname })` |
| `changePassword` | API：修改密码 | `await changePassword({ old, new })` |
| `setPassword` | API：首次设置密码 | `await setPassword({ password })` |
| `deactivateAccount` | API：注销账号 | `await deactivateAccount()` |
| `unbindOAuth` | API：解绑 OAuth | `await unbindOAuth('github')` |

### 评论

| 函数 | 说明 | 示例 |
|---|---|---|
| `useComments` | 获取评论列表 | `const { comments, fetchComments } = useComments()` |
| `fillComment` | 填充评论用户信息 | `fillComment(comment)` |
| `flattenComments` | 将嵌套评论扁平化 | `flattenComments(comments)` |
| `provideCommentContext` | 提供评论上下文 | `provideCommentContext(articleId)` |
| `useCommentContext` | 获取评论上下文 | `const ctx = useCommentContext()` |
| `getComments` | API：获取评论 | `await getComments({ target_type: 'article', target_key: 'slug' })` |
| `createComment` | API：发表评论 | `await createComment({ target_type: 'article', target_key: 'slug', content: '...' })` |
| `deleteComment` | API：删除评论 | `await deleteComment(id)` |

### 主题与配置

| 函数 | 说明 | 示例 |
|---|---|---|
| `useTheme` | 获取激活主题及其菜单 | `const { themeConfig, getMenus } = useTheme()` |
| `useSysConfig` | 系统配置（标题/描述/OAuth） | `const { basicConfig, oauthConfig } = useSysConfig()` |
| `getActiveTheme` | API：获取激活主题 | `await getActiveTheme()` |
| `getSettingGroup` | API：获取配置组 | `await getSettingGroup('basic')` |

### 暗色模式

| 函数 | 说明 | 示例 |
|---|---|---|
| `useDarkMode` | 暗色模式控制 | `const { isDark, toggleTheme, initAutoSwitch } = useDarkMode()` |

### 分类

| 函数 | 说明 | 示例 |
|---|---|---|
| `useCategory` | 按 slug 获取分类 | `const { category } = useCategory(slug)` |
| `useCategories` | 获取全部分类 | `const { categories } = useCategories()` |
| `getCategories` | API：分类列表 | `await getCategories()` |
| `getCategoryById` | API：按 ID 获取 | `await getCategoryById(1)` |
| `getCategoryBySlug` | API：按 slug 获取 | `await getCategoryBySlug('frontend')` |

### 标签

| 函数 | 说明 | 示例 |
|---|---|---|
| `useTag` | 按 slug 获取标签 | `const { tag } = useTag(slug)` |
| `useTags` | 获取全部标签 | `const { tags } = useTags()` |
| `getTags` | API：标签列表 | `await getTags()` |
| `getTagById` | API：按 ID 获取 | `await getTagById(1)` |
| `getTagBySlug` | API：按 slug 获取 | `await getTagBySlug('vue')` |

### 友链

| 函数 | 说明 | 示例 |
|---|---|---|
| `useFriends` | 友链和友链分组 | `const { data: friendGroups, refresh } = useFriends()` |
| `useFriendApply` | 提交友链申请 | `const { apply } = useFriendApply()` |
| `getFriends` | API：友链列表 | `await getFriends()` |
| `applyFriend` | API：申请友链 | `await applyFriend({ name, url, description, avatar })` |

### 统计

| 函数 | 说明 | 示例 |
|---|---|---|
| `useStats` | 站点统计 | `const { siteStats } = useStats()` |
| `useArchives` | 归档统计 | `const { archives } = useArchives()` |
| `getSiteStats` | API：站点统计 | `await getSiteStats()` |
| `getArchiveStats` | API：归档统计 | `await getArchiveStats()` |

### 动态

| 函数 | 说明 | 示例 |
|---|---|---|
| `useMomentList` | 动态列表 | `const { data: moments } = useMomentList()` |
| `getMoments` | API：动态列表 | `await getMoments()` |

### 订阅

| 函数 | 说明 | 示例 |
|---|---|---|
| `useSubscribe` | 邮箱订阅/退订 | `const { subscribe, unsubscribe } = useSubscribe()` |
| `subscribe` | API：订阅 | `await subscribe({ email })` |
| `unsubscribe` | API：退订 | `await unsubscribe({ email })` |

### 通知

| 函数 | 说明 | 示例 |
|---|---|---|
| `useNotifications` | 通知列表 + 未读数 | `const { unreadCount, notifications } = useNotifications()` |
| `getNotifications` | API：通知列表 | `await getNotifications()` |
| `markAsRead` | API：标记已读 | `await markAsRead(id)` |
| `markAllAsRead` | API：全部已读 | `await markAllAsRead()` |

### 反馈

| 函数 | 说明 | 示例 |
|---|---|---|
| `useFeedback` | 提交反馈 | `const { submit } = useFeedback()` |
| `submitFeedback` | API：提交反馈 | `await submitFeedback({ reportUrl, reportType: 'bug', description })` |
| `getFeedbackByTicketNo` | API：按工单号查询 | `await getFeedbackByTicketNo('T2025...')` |

### Markdown

| 函数 | 说明 | 示例 |
|---|---|---|
| `renderMarkdown` | 完整渲染（代码高亮/Katex） | `const html = renderMarkdown(content)` |
| `renderSimpleMarkdown` | 简单渲染（不含高亮） | `const html = renderSimpleMarkdown(content)` |
| `extractToc` | 提取文章目录 | `const toc = extractToc(content)` |
| `countWords` | 统计字数 | `countWords(content)` |
| `estimateReadingTime` | 预估阅读时间（分钟） | `estimateReadingTime(content, 300)` |

### 日期

| 函数 | 说明 | 示例 |
|---|---|---|
| `formatDateTime` | `2025-10-03 13:46:59` | `formatDateTime(date)` |
| `formatDate` | `2025-10-03` | `formatDate(date)` |
| `formatRelativeTime` | `2小时前` | `formatRelativeTime(date)` |
| `formatFriendly` | `2025年10月3日` | `formatFriendly(date)` |
| `formatMomentTime` | `刚刚` / `3天前` / `10月3日` | `formatMomentTime(date)` |
| `formatForBackend` | 转后端格式 | `formatForBackend(date)` |
| `parseBackendDate` | 解析后端日期 | `parseBackendDate('2025-10-03 13:46:59')` |
| `isValidDate` | 日期是否有效 | `isValidDate(date)` |

### 表情

| 函数 | 说明 | 示例 |
|---|---|---|
| `loadEmojiGroups` | 加载表情分组 | `const groups = await loadEmojiGroups(url)` |
| `loadEmojiMap` | 加载 `:key:` → URL 映射 | `const map = await loadEmojiMap(url)` |
| `replaceEmojisInText` | 替换文本中的表情占位符 | `replaceEmojisInText(text, map)` |

### 滚动

| 函数 | 说明 | 示例 |
|---|---|---|
| `scrollToTop` | 平滑回到顶部 | `scrollToTop()` |
| `scrollToElement` | 滚动到指定元素 | `scrollToElement('#comments')` |

### UI

| 函数 | 说明 | 示例 |
|---|---|---|
| `useToast` | Toast 消息 | `const { success, error } = useToast()` |
| `useExpandable` | 内容展开/收起 | `const { isExpanded, toggle } = useExpandable()` |

### 上传

| 函数 | 说明 | 示例 |
|---|---|---|
| `useUpload` | 上传文件 | `const { upload } = useUpload()` |
| `uploadFile` | API：上传单文件 | `await uploadFile({ file })` |

### 工具

| 函数 | 说明 | 示例 |
|---|---|---|
| `parseJSON` | 安全解析 JSON | `parseJSON(str)` |
| `createApi` | 创建 API 工厂实例 | `const api = createApi('/resource')` |
| `apiRequest` | 通用请求（自动 token/Auth） | `await apiRequest('/path')` |

## 类型

TypeScript 类型需要显式 import。统一从 `@flecblog/core-nuxt` 导入，自动索引所有类型。

```ts
import type { Article, ArticleQuery, ApiResponse, PaginationData, UserInfo, Comment } from '@flecblog/core-nuxt';
```

可用类型：

| 模块 | 类型 |
|---|---|
| 文章 | `Article` `ArticleNav` `ArticleQuery` |
| 认证 | `LoginParams` `LoginResponse` `RegisterParams` `RegisterResponse` `ForgotPasswordParams` `ResetPasswordParams` `RefreshTokenResponse` |
| 分类 | `Category` |
| 评论 | `Comment` `CreateCommentParams` `GetCommentsParams` `FlatComment` `EmojiItem` `EmojiGroup` `CommentTargetType` |
| 反馈 | `Feedback` `SubmitFeedbackParams` `ReportType` `FeedbackStatus` |
| 友链 | `Friend` `FriendGroup` `FriendGroupedResponse` `FriendApplyRequest` |
| 动态 | `Moment` `MomentContent` `MomentMusic` `AudioTrack` `LyricLine` `MomentListResponse` |
| 通知 | `Notification` `NotificationType` `GetNotificationsParams` |
| 请求 | `ApiResponse<T>` `PaginationQuery` `PaginationData<T>` |
| 统计 | `SiteStats` `ArchiveItem` `ArchiveStats` |
| 配置 | `SysConfigData` `SettingGroupType` |
| 标签 | `Tag` |
| 主题 | `ThemeMenuItem` `ActiveTheme` |
| 上传 | `UploadResponse` `UploadType` |
| 用户 | `UserInfo` `UpdateProfileParams` `ChangePasswordParams` `UserRole` |
| Markdown | `TocItem` |

## 样式

包提供 `base.css` 和 `prose.css`，会在 `extends` 时自动加载，可在主题中创建同名选择器覆盖。

- **`base.css`**：CSS 重置、RemixIcon 字体图标
- **`prose.css`**：文章内容完整排版（亮/暗双主题），包含标题、代码块、表格、引用、自定义卡片等 20+ 种元素样式

## 许可证

MIT
