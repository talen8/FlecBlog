import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')

  // 从 VITE_API_URL 提取后端基础地址，提供默认值
  const apiBaseUrl = env.VITE_API_URL || 'http://localhost:8080/api/v1'
  const backendBaseUrl = apiBaseUrl.replace(/\/api\/v\d+$/, '') || 'http://localhost:8080'

  return {
    plugins: [
      vue(),
      VitePWA({
        registerType: 'autoUpdate',
        manifest: {
          name: 'Flec 管理系统',
          short_name: 'Flec Admin',
          description: '后台管理系统',
          theme_color: '#f7f7f7',
          background_color: '#ffffff',
          display: 'standalone',
          icons: [
            {
              src: '/pwa-192x192.png',
              sizes: '192x192',
              type: 'image/png'
            },
            {
              src: '/pwa-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'any maskable'
            }
          ]
        },
        includeAssets: ['favicon.ico', 'pwa-192x192.png', 'pwa-512x512.png'],
        devOptions: {
          enabled: true
        },
        workbox: {
          // 只缓存静态资源
          globPatterns: ['**/*.{js,css,html,ico,png,woff,woff2}'],
          // 排除大文件
          globIgnores: ['**/remixicon*.svg'],
          // 不缓存 API 请求
          navigateFallbackDenylist: [/^\/api/]
        }
      })
    ],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src')
      }
    },
    server: {
      port: 5174,
      proxy: {
        '/api': {
          target: backendBaseUrl,
          changeOrigin: true
        },
        '/uploads': {
          target: backendBaseUrl,
          changeOrigin: true
        }
      }
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            // Vue 核心
            'vue-vendor': ['vue', 'vue-router'],
            // Element Plus UI 框架
            'element-plus': ['element-plus'],
            // CodeMirror 编辑器核心
            'codemirror': [
              '@codemirror/autocomplete',
              '@codemirror/commands',
              '@codemirror/lang-markdown',
              '@codemirror/language',
              '@codemirror/lint',
              '@codemirror/search',
              '@codemirror/state',
              '@codemirror/view',
              'codemirror'
            ],
            // Mermaid 图表库（体积较大，单独分割）
            'mermaid': ['mermaid'],
            // Markdown 解析器及插件
            'markdown': [
              'markdown-it',
              'markdown-it-anchor',
              'markdown-it-kbd',
              'markdown-it-link-attributes',
              'markdown-it-mark',
              'markdown-it-plugin-underline',
              'markdown-it-sub',
              'markdown-it-sup',
              'markdown-it-task-lists'
            ],
            // KaTeX 数学公式渲染
            'katex': ['katex'],
            // 其他工具库
            'utils': ['axios', 'dayjs', 'dompurify', '@vueuse/core'],
            // 图表库
            'echarts': ['echarts', 'echarts-wordcloud']
          }
        }
      }
    }
  }
})
