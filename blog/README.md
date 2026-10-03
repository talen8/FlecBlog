# FlecBlog 博客前端

> 基于 Nuxt 4 + Vue 3 的现代化博客前端应用

## 技术栈

- **框架**: [Nuxt 4](https://nuxt.com) - Vue.js 全栈框架
- **文章渲染**: markdown-it, Highlight.js, KaTeX, Mermaid
- **样式**: SCSS
- **SEO**: @nuxtjs/seo, Sitemap, RSS / Atom Feed
- **PWA**: @vite-pwa/nuxt
- **图片**: @nuxt/image, medium-zoom
- **其他**: TypeScript, VueUse, dayjs, DOMPurify, Remix Icon

## 文件结构

```
blog/
├── app/                  # 应用主目录
│   ├── assets/           # 静态资源（样式表、字体）
│   ├── components/       # Vue 组件
│   ├── composables/      # 组合式函数（含 api/ 请求封装）
│   ├── layouts/          # 页面布局
│   ├── pages/            # 页面路由
│   ├── plugins/          # Nuxt 插件（暗色模式切换、自定义代码、埋点等）
│   ├── utils/            # 工具函数（Markdown 渲染、日期、滚动等）
│   ├── app.vue           # 根组件
│   └── error.vue         # 错误页面
├── public/               # 公共静态文件
├── server/               # 服务端路由（RSS、Atom、动态 manifest）
├── types/                # TypeScript 类型定义
├── nuxt.config.ts        # Nuxt 配置
└── Dockerfile            # Docker 配置
```

详细文档请查看 [项目主 README](../README.md)
