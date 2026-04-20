<script lang="ts" setup>
import { getArticleBySlug } from '@/composables/api/article';
import type { Article } from '@@/types/article';

definePageMeta({
  typeHeader: 'post',
});

const route = useRoute();
const router = useRouter();
const article = ref<Article | null>(null);
const { setCurrentArticle, clearCurrentArticle } = useCurrentArticle();
const { $tracker } = useNuxtApp();

// 使用SSR获取文章详情
const { data: initialData } = await useAsyncData(`post-${route.params.slug}`, async () => {
  const slug = route.params.slug as string;

  try {
    const articleData = await getArticleBySlug(slug);
    setCurrentArticle(articleData);
    return { article: articleData };
  } catch (error: unknown) {
    const err = error as Error & { response?: { status?: number } };
    if (err.response?.status === 404) {
      router.replace('/404');
    }
    return null;
  }
});

// 初始化本地 article ref
article.value = initialData.value?.article ?? null;

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

const fetchArticle = async () => {
  const slug = route.params.slug as string;

  try {
    article.value = await getArticleBySlug(slug);
    setCurrentArticle(article.value);

    // 设置当前文章ID用于埋点追踪
    if (article.value) {
      $tracker?.setArticleId(article.value.id);
      // 发送包含 article_id 的页面访问埋点
      $tracker?.trackPageView(undefined, article.value.id);

      // 处理 URL hash 锚点跳转
      nextTick(() => {
        if (route.hash) {
          requestAnimationFrame(() => scrollToElement(route.hash, { block: 'start' }));
        }
      });
    }
  } catch (error: unknown) {
    const err = error as Error & { response?: { status?: number } };
    clearCurrentArticle();
    $tracker?.setArticleId(undefined);

    // 如果是404错误，替换到404页面（不保留历史记录，避免循环）
    if (err.response?.status === 404) {
      router.replace('/404');
    }
  }
};

// 监听路由参数变化
watch(() => route.params.slug, fetchArticle);

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

    <LazyFeaturesCommentComments target-type="article" :target-key="article.slug!" />
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
