import { createResolver } from '@nuxt/kit';

// @flecblog/core-nuxt Nuxt Layer Configuration
const { resolve } = createResolver(import.meta.url);

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',

  future: {
    compatibilityVersion: 4,
  },

  // 模块
  modules: [
    '@vueuse/nuxt',
    '@nuxt/image',
    '@nuxtjs/seo',
  ],

  // Auto-import 配置
  imports: {
    dirs: [
      resolve('./app/composables'),
      resolve('./app/lib'),
      resolve('./app/utils'),
    ],
  },

  // CSS 配置
  css: [
    'remixicon/fonts/remixicon.css',
    resolve('./app/assets/css/base.css'),
    resolve('./app/assets/css/prose.css'),
  ],

  // SEO 配置
  site: {
    defaultLocale: 'zh-CN',
  },

  // Sitemap 配置
  sitemap: {
    strictNuxtContentPaths: true,
  },

  // Robots 配置
  robots: {
    allow: '/',
  },

  // 禁用 OG Image 自动生成
  ogImage: {
    enabled: false,
  },

  // Vite 配置
  vite: {
    // Vite 7 CJS 需逐个预构建转 ESM
    optimizeDeps: {
      include: [
        'highlight.js/lib/core',
        'highlight.js/lib/languages/javascript',
        'highlight.js/lib/languages/typescript',
        'highlight.js/lib/languages/python',
        'highlight.js/lib/languages/go',
        'highlight.js/lib/languages/java',
        'highlight.js/lib/languages/sql',
        'highlight.js/lib/languages/xml',
        'highlight.js/lib/languages/css',
        'highlight.js/lib/languages/json',
        'highlight.js/lib/languages/yaml',
        'highlight.js/lib/languages/markdown',
        'highlight.js/lib/languages/bash',
        'highlight.js/lib/languages/shell',
        'highlight.js/lib/languages/dockerfile',
        'highlight.js/lib/languages/diff',
        'markdown-it',
        'markdown-it-anchor',
        'markdown-it-task-lists',
        'markdown-it-mark',
        'markdown-it-link-attributes',
        'markdown-it-kbd',
        'markdown-it-sub',
        'markdown-it-sup',
        'markdown-it-plugin-underline',
        '@traptitech/markdown-it-katex',
      ],
    },
  },

  // 路由配置
  router: {
    options: {
      scrollBehaviorType: 'smooth',
    },
  },

  // 注册插件
  plugins: [
    resolve('./plugins/tracker.client.ts'),
    resolve('./plugins/console-banner.client.ts'),
    resolve('./plugins/custom-code.ts'),
  ],
});
