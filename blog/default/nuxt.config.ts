export default defineNuxtConfig({
  // Extend from core-nuxt layer
  extends: ['@flecblog/core-nuxt'],

  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  // 启用 SSR
  ssr: true,

  // 内联 SSR 样式到 HTML
  features: {
    inlineStyles: true,
  },

  // 应用配置
  app: {
    head: {
      htmlAttrs: { lang: 'zh-CN' },
      script: [
        {
          innerHTML: `
            (function() {
              var theme = localStorage.getItem('theme');
              var isDark = theme === 'dark' || (!theme && window.matchMedia('(prefers-color-scheme: dark)').matches);
              document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
            })();
          `,
          type: 'text/javascript',
          tagPosition: 'head',
        },
      ],
    },
  },

  // 主题专属模块
  modules: [
    '@vite-pwa/nuxt',
    [
      '@nuxtjs/critters',
      {
        config: {
          preload: 'swap',
          inlineFonts: false,
          pruneSource: false,
        },
      },
    ],
  ],

  // CSS 配置（主题专属样式，在包默认样式之后加载以覆盖）
  css: ['@/assets/css/color.css', '@/assets/css/global.scss'],

  // SEO 配置
  site: {
    url: '',
  },

  // 运行时配置
  runtimeConfig: {
    public: {
      apiUrl: '',
      appVersion: process.env.FLECBLOG_VERSION || '1.0.0',
    },
  },

  // PWA 配置
  pwa: {
    registerType: 'autoUpdate',
    manifest: false, // 使用自定义的动态 manifest
    workbox: {
      navigateFallback: '/favicon.ico',
      navigateFallbackAllowlist: [/$^/],
      globPatterns: ['**/*.{js,css,png,ico,webp,woff,woff2}'],
      globIgnores: ['**/remixicon*.svg'],
      maximumFileSizeToCacheInBytes: 3 * 1024 * 1024,
      runtimeCaching: [
        {
          urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp|ico)$/i,
          handler: 'CacheFirst',
          options: {
            cacheName: 'images',
            expiration: {
              maxEntries: 100,
              maxAgeSeconds: 60 * 60 * 24 * 30, // 30 天
            },
          },
        },
      ],
    },
    client: {
      installPrompt: true,
      periodicSyncForUpdates: 3600, // 每小时检查更新
    },
    devOptions: {
      enabled: true,
      type: 'module',
      navigateFallbackAllowlist: [/$^/],
    },
  },

  // Vite 配置
  vite: {
    build: {
      rollupOptions: {
        output: {
          // 细粒度的代码分割策略
          manualChunks(id) {
            // 核心框架（首屏必需）
            if (id.includes('node_modules/vue/') || id.includes('node_modules/@vue/')) {
              return 'vue-core';
            }
            if (id.includes('node_modules/vue-router')) {
              return 'vue-router';
            }

            // 日期处理库
            if (id.includes('node_modules/dayjs')) {
              return 'dayjs';
            }

            // Markdown 渲染生态
            if (
              id.includes('node_modules/markdown-it') ||
              id.includes('node_modules/dompurify') ||
              id.includes('node_modules/isomorphic-dompurify')
            ) {
              return 'markdown-renderer';
            }

            // KaTeX 数学公式渲染
            if (id.includes('node_modules/katex')) {
              return 'katex';
            }

            // 代码高亮（较大，独立分割）
            if (id.includes('node_modules/highlight.js')) {
              return 'highlight';
            }

            // VueUse 工具库
            if (id.includes('node_modules/@vueuse')) {
              return 'vueuse';
            }
          },
        },
      },
      // 调整 chunk 大小警告阈值
      chunkSizeWarningLimit: 600,
      // CSS 代码分割
      cssCodeSplit: true,
      // 关闭生产环境 sourcemap
      sourcemap: false,
      cssMinify: true,
    },
  },

  // 路由规则（主题专属）
  routeRules: {
    '/oauth/**': { sitemap: false },
    '/profile': { sitemap: false },
    '/notifications': { sitemap: false },
    '/feedback': { sitemap: false },
    '/subscribe': { sitemap: false },
  },
});
