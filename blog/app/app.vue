<script setup lang="ts">
const { toasts } = useToast();
const { showLoginModal } = useLoginModal();
const { showBindEmailModal, triggerGlobal, onBindSuccess } = useBindEmail();

// 全局数据
const { basicConfig } = useSysConfig();
const { themeConfig } = useTheme();

// 全局路由切换时触发邮箱绑定提示
const router = useRouter();
router.afterEach(() => {
  triggerGlobal();
});

// 背景图片
const bgImage = computed(() => themeConfig.value.background_image || '/bg.webp');

// 刷新时恢复滚动位置
onMounted(() => {
  const key = 'scroll-y';
  const nav = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
  if (nav?.type === 'reload') {
    const y = +(sessionStorage.getItem(key) || 0);
    if (y > 0) setTimeout(() => window.scrollTo(0, y), 100);
  }
  let t: ReturnType<typeof setTimeout>;
  const save = () => sessionStorage.setItem(key, '' + window.scrollY);
  window.addEventListener(
    'scroll',
    () => {
      clearTimeout(t);
      t = setTimeout(save, 200);
    },
    { passive: true }
  );
  window.addEventListener('pagehide', save);
});

// SEO Meta
useSeoMeta({
  description: () => basicConfig.value.description,
  keywords: () => basicConfig.value.keywords,
  author: () => basicConfig.value.author,
  // Open Graph
  ogTitle: () => basicConfig.value.title,
  ogDescription: () => basicConfig.value.description,
  ogImage: () => basicConfig.value.favicon,
  ogType: 'website',
  ogSiteName: () => basicConfig.value.title,
  // Twitter Card
  twitterCard: 'summary_large_image',
  twitterTitle: () => basicConfig.value.title,
  twitterDescription: () => basicConfig.value.description,
  twitterImage: () => basicConfig.value.favicon,
});

// 页面标题模板和 favicon
const route = useRoute();
const siteTitle = computed(() => basicConfig.value.title);

useHead({
  titleTemplate: (title): string | null => {
    // 首页特殊处理：显示"网站标题 - 网站副标题"
    if (route.path === '/') {
      const subtitle = basicConfig.value.subtitle;
      return subtitle ? `${siteTitle.value} - ${subtitle}` : siteTitle.value || null;
    }

    // 其他页面：显示"页面标题 | 网站标题"
    const pageTitle = title || (route.meta.title as string);
    if (pageTitle) return `${pageTitle} | ${siteTitle.value}`;
    return siteTitle.value || null;
  },
  link: [
    { rel: 'icon', href: basicConfig.value.favicon || '/favicon.ico' },
    // PWA Manifest
    { rel: 'manifest', href: '/manifest.json' },
    // RSS/Atom 订阅
    {
      rel: 'alternate',
      type: 'application/rss+xml',
      title: `${basicConfig.value.title} - RSS 2.0 Feed`,
      href: '/rss.xml',
    },
    {
      rel: 'alternate',
      type: 'application/atom+xml',
      title: `${basicConfig.value.title} - Atom Feed`,
      href: '/atom.xml',
    },
  ],
  meta: computed(() => [
    { name: 'description', content: basicConfig.value.description },
    { name: 'keywords', content: basicConfig.value.keywords },
    { name: 'author', content: basicConfig.value.author },
    // PWA 主题色
    { name: 'theme-color', content: '#f7f7f7' },
    { name: 'mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' },
  ]),
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: basicConfig.value.title,
        description: basicConfig.value.description,
      }),
    },
  ],
});
</script>

<template>
  <!-- 背景图片 -->
  <div class="web_bg" :style="{ backgroundImage: `url(${bgImage})` }" />

  <!-- Nuxt 布局和页面系统 -->
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>

  <!-- Toast 消息提示 -->
  <UiToast
    v-for="toast in toasts"
    :key="toast.id"
    :message="toast.message"
    :type="toast.type"
    :show="toast.show"
  />

  <!-- 登录弹窗 -->
  <FeaturesModalsLoginModal v-model="showLoginModal" />

  <!-- 邮箱绑定弹窗 -->
  <FeaturesModalsBindEmailModal v-model="showBindEmailModal" @success="onBindSuccess" />

  <!-- 右键菜单 -->
  <UiContextMenu />
</template>

<style scoped>
.web_bg {
  position: fixed;
  width: 100%;
  height: 100%;
  z-index: -50;
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
}

[data-theme='dark'] .web_bg::before {
  position: absolute;
  width: 100%;
  height: 100%;
  background-color: #121212b0;
  content: '';
}
</style>
