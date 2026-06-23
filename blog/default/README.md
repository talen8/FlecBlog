# FlecBlog 默认主题

> 基于 Nuxt 4 + Vue 3 的现代化博客默认主题，通过 extends `@flecblog/core-nuxt` 获得核心能力。

## 技术栈

- **核心**: [@flecblog/core-nuxt](https://www.npmjs.com/package/@flecblog/core-nuxt)
- **框架**: [Nuxt 4](https://nuxt.com)
- **文章渲染**: markdown-it、Highlight.js、KaTeX
- **样式**: SCSS
- **SEO**: @nuxtjs/seo、Sitemap、Atom Feed
- **PWA**: @vite-pwa/nuxt
- **平滑滚动**: Lenis
- **图片预览**: medium-zoom
- **音乐播放**: APlayer
- **其他**: TypeScript、VueUse、dayjs

## 文件结构

```
blog/default/
├── app/                  # 应用主目录
│   ├── assets/           # 主题专属样式（global.scss、color.css）
│   ├── components/       # Vue 组件（48 个）
│   ├── composables/      # 主题专属组合式函数
│   ├── layouts/          # 页面布局
│   ├── pages/            # 页面路由（22 个）
│   ├── plugins/          # 主题专属插件（暗色模式自动切换等）
│   ├── app.vue           # 根组件
│   └── error.vue         # 错误页面
├── public/               # 公共静态文件
├── nuxt.config.ts        # Nuxt 配置（核心层 + 主题专属模块）
├── package.json
└── Dockerfile
```

详细文档请查看 [项目主 README](../../README.md)
