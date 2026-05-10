export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  ssr: true,

  features: {
    inlineStyles: true,
  },

  app: {
    head: {
      htmlAttrs: { lang: 'zh-CN' },
    },
  },

  modules: ['@vueuse/nuxt'],

  css: ['@/assets/css/color.css', '@/assets/css/global.scss'],

  runtimeConfig: {
    public: {
      apiUrl: '',
    },
  },

  router: {
    options: {
      scrollBehaviorType: 'smooth',
    },
  },
});
