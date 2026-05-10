# FlecBlog Theme Template

这是一个简洁的 Nuxt 博客主题模板，用于 FlecBlog 项目的二开。

## 快速开始

### 环境要求

- Node.js >= 18.x
- npm >= 9.x 或 pnpm >= 8.x

### 安装

1. 复制 `theme/template` 目录到你的项目
2. 安装依赖：

```bash
npm install
```

3. 配置环境变量，创建 `.env` 文件：

```env
NUXT_PUBLIC_API_URL=http://localhost:8080/api
```

4. 运行开发服务器：

```bash
npm run dev
```

5. 构建生产版本：

```bash
npm run build
```

## 目录结构

```
├── app/                        # 应用主目录
│   ├── assets/css/             # 样式文件
│   │   ├── color.css           # CSS 变量（主题色）
│   │   ├── global.scss         # 全局样式
│   │   ├── _mixins.scss        # SCSS 混入
│   │   └── _prose.scss         # 文章排版样式
│   ├── components/             # 组件
│   │   ├── features/           # 功能组件
│   │   │   ├── archive/        # 归档相关
│   │   │   ├── article/        # 文章相关
│   │   │   ├── comment/        # 评论相关
│   │   │   └── modals/         # 弹窗组件
│   │   ├── layout/             # 布局组件
│   │   │   ├── footer/         # 页脚
│   │   │   ├── navbar/         # 导航栏
│   │   │   └── sidebar/        # 侧边栏
│   │   └── ui/                 # UI 基础组件
│   ├── composables/            # 组合式函数
│   │   ├── useApi.ts           # API 请求函数
│   │   ├── useComment.ts       # 评论逻辑
│   │   ├── useStores.ts        # 全局状态
│   │   ├── useToast.ts         # 消息提示
│   │   └── useUser.ts          # 用户逻辑
│   ├── layouts/                # 布局模板
│   │   └── default.vue         # 默认布局
│   ├── pages/                  # 页面路由
│   ├── plugins/                # 插件
│   │   └── tracker.client.ts   # 访问统计插件
│   ├── utils/                  # 工具函数
│   │   ├── auth.ts             # 认证相关
│   │   ├── avatar.ts           # 头像处理
│   │   ├── date.ts             # 日期格式化
│   │   ├── json.ts             # JSON 解析
│   │   ├── markdown.ts         # Markdown 渲染
│   │   ├── request.ts          # HTTP 请求封装
│   │   └── upload.ts           # 文件上传
│   ├── app.vue                 # 根组件
│   └── error.vue               # 错误页面
├── server/                     # 服务端
│   └── routes/                 # 服务端路由
│       ├── atom.xml.ts         # Atom 订阅
│       └── rss.xml.ts          # RSS 订阅
├── types/                      # TypeScript 类型定义
│   └── index.ts                # 统一导出
├── .prettierrc                 # Prettier 配置
├── eslint.config.js            # ESLint 配置
├── nuxt.config.ts              # Nuxt 配置
├── package.json                # 依赖配置
└── tsconfig.json               # TypeScript 配置
```

## 页面列表

| 页面      | 路径                      | 说明         |
| --------- | ------------------------- | ------------ |
| 首页      | `/`                       | 文章列表     |
| 文章详情  | `/posts/[slug]`           | 文章内容     |
| 分类列表  | `/categories`             | 所有分类     |
| 分类详情  | `/category/[slug]`        | 分类文章     |
| 标签列表  | `/tags`                   | 所有标签     |
| 标签详情  | `/tag/[slug]`             | 标签文章     |
| 归档      | `/archive`                | 文章归档     |
| 月度归档  | `/archive/[year]/[month]` | 月度文章     |
| 动态      | `/moment`                 | 动态列表     |
| 友链      | `/friend`                 | 友情链接     |
| 关于      | `/about`                  | 关于页面     |
| 留言板    | `/message`                | 留言板       |
| 个人中心  | `/profile`                | 用户信息     |
| 通知中心  | `/notifications`          | 通知列表     |
| 反馈      | `/feedback`               | 反馈表单     |
| 查询反馈  | `/feedback/query`         | 查询反馈进度 |
| 统计      | `/statistics`             | 网站统计     |
| OAuth回调 | `/oauth/callback`         | 登录回调     |
| 订阅      | `/subscribe`              | RSS 订阅     |
| 搜索      | `/search`                 | 文章搜索     |

## 开发指南

### 代码规范

项目使用 ESLint + Prettier 进行代码规范检查：

```bash
# 检查代码规范
npm run lint

# 自动修复
npm run lint:fix

# 格式化代码
npm run format
```

### 添加新页面

1. 在 `app/pages/` 目录下创建 `.vue` 文件
2. 使用 `useSeoMeta` 配置页面 SEO 信息：

```vue
<script lang="ts" setup>
useSeoMeta({
  title: '页面标题',
  description: '页面描述',
});
</script>
```

### 添加新组件

1. 在 `app/components/` 对应目录下创建 `.vue` 文件
2. 组件会自动注册，无需手动导入

### 添加新 API

在 `app/composables/useApi.ts` 中添加新的 API 方法：

```ts
// 使用 createApi 工厂函数
const myApi = createApi<MyType>('/my-endpoint');

// 或直接导出函数
export const getMyData = async () => myApi.getList();
```

### 状态管理

使用 `useState` 进行全局状态管理，相关定义在 `app/composables/useStores.ts` 中：

```ts
// 用户相关
const isLoggedIn = useAuth(); // 来自 utils/auth.ts
const { userInfo, fetchUserInfo } = useUser();

// 配置相关
const { blogConfig, basicConfig } = useSysConfig();

// 其他状态
const { articles, fetchArticles } = useArticles();
const { categories, fetchCategories } = useCategories();
const { tags, fetchTags } = useTags();
```

## 技术栈

- [Nuxt 4](https://nuxt.com/) - Vue.js 框架
- [Vue 3](https://vuejs.org/) - 渐进式 JavaScript 框架
- [TypeScript](https://www.typescriptlang.org/) - 类型安全
- [SCSS](https://sass-lang.com/) - CSS 预处理器
- [VueUse](https://vueuse.org/) - Vue 组合式 API 工具集
- [markdown-it](https://github.com/markdown-it/markdown-it) - Markdown 解析器
- [DOMPurify](https://github.com/cure53/DOMPurify) - XSS 防护
- [dayjs](https://day.js.org/) - 日期处理库
- [Remix Icon](https://remixicon.com/) - 图标库
