<script lang="ts" setup>
definePageMeta({
  typeHeader: 'post',
});

const route = useRoute();
const router = useRouter();
const { $tracker } = useNuxtApp();
const { setCurrentArticle, clearCurrentArticle } = useCurrentArticle();

// SSR 获取文章详情（与 PostHeader 共享同一请求）
const { data: article, error } = useArticle(() => route.params.slug as string);

// 同步到全局当前文章（供 TocCard / MobileToc / Menu 读取）
watch(article, a => setCurrentArticle(a ?? null), { immediate: true });

// 404 处理
watch(
  error,
  err => {
    const e = err as (Error & { response?: { status?: number }; statusCode?: number }) | null;
    if (e?.response?.status === 404 || e?.statusCode === 404) {
      router.replace('/404');
    }
  },
  { immediate: true }
);

// 动态页面标题和 SEO
useHead({
  title: () => article.value?.title,
});

useSeoMeta({
  title: () => article.value?.title,
  description: () => article.value?.summary || `${article.value?.title} - 阅读全文了解更多详情`,
  ogTitle: () => article.value?.title,
  ogDescription: () => article.value?.summary,
  ogImage: () => article.value?.cover,
  ogType: 'article',
  twitterTitle: () => article.value?.title,
  twitterDescription: () => article.value?.summary,
  twitterImage: () => article.value?.cover,
});

// 文章结构化数据
useSchemaOrg([
  defineArticle({
    headline: () => article.value?.title,
    description: () => article.value?.summary,
    image: () => article.value?.cover,
    datePublished: () => article.value?.publish_time,
    dateModified: () => article.value?.update_time,
  }),
]);

// 文章加载后（客户端）：埋点 + 锚点跳转
watch(
  article,
  a => {
    if (!a || !import.meta.client) return;
    $tracker?.setArticleId(a.id);
    $tracker?.trackPageView(undefined, a.id);
    nextTick(() => {
      if (route.hash) requestAnimationFrame(() => scrollToElement(route.hash, { block: 'start' }));
    });
  },
  { immediate: true }
);

// 监听 URL hash 变化，实现锚点跳转
watch(
  () => route.hash,
  hash => {
    if (hash) scrollToElement(hash, { block: 'start' });
  }
);

// 组件卸载时清除文章数据
onUnmounted(() => {
  clearCurrentArticle();
  $tracker?.setArticleId(undefined);
});
</script>

<template>
  <div v-if="article" id="post">
    <FeaturesArticleAISummary v-if="article.ai_summary" :summary="article.ai_summary" />

    <FeaturesArticleOutdatedNotice v-if="article.is_outdated" />

    <FeaturesArticleContent :content="article.content!" />

    <FeaturesArticleCopyright :article="article" />

    <FeaturesArticleTags :article="article" />

    <FeaturesArticleNavigation :prev="article.prev" :next="article.next" />

    <LazyFeaturesCommentComments target-type="article" :target-key="article.id!" />
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#post {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #post {
    padding: 30px;
  }
}

@media screen and (max-width: 768px) {
  #post {
    padding: 18px;
  }
}
</style>
