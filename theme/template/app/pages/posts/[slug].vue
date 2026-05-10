<script lang="ts" setup>
import { getArticleBySlug } from '@/composables/useApi';
import type { Article } from '@@/types';

const route = useRoute();
const router = useRouter();
const article = ref<Article | null>(null);

const { data: initialData } = await useAsyncData(`post-${route.params.slug}`, async () => {
  const slug = route.params.slug as string;
  try {
    const articleData = await getArticleBySlug(slug);
    return { article: articleData };
  } catch (error: unknown) {
    const err = error as Error & { response?: { status?: number } };
    if (err.response?.status === 404) {
      router.replace('/404');
    }
    return null;
  }
});

article.value = initialData.value?.article ?? null;

useHead({
  title: () => article.value?.title,
});

useSeoMeta({
  title: () => article.value?.title,
  description: () => article.value?.summary,
});

watch(
  () => route.params.slug,
  async () => {
    const slug = route.params.slug as string;
    try {
      article.value = await getArticleBySlug(slug);
    } catch (error: unknown) {
      const err = error as Error & { response?: { status?: number } };
      if (err.response?.status === 404) {
        router.replace('/404');
      }
    }
  }
);
</script>

<template>
  <div v-if="article" id="post">
    <h1 class="post-title">{{ article.title }}</h1>
    <div class="post-meta">
      <span><i class="ri-calendar-line" /> {{ formatDate(article.publish_time) }}</span>
      <span v-if="article.category">
        <i class="ri-folder-line" />
        <NuxtLink :to="article.category.url">{{ article.category.name }}</NuxtLink>
      </span>
      <span><i class="ri-eye-line" /> {{ article.view_count || 0 }} 次阅读</span>
      <span v-if="article.location"><i class="ri-map-pin-line" /> {{ article.location }}</span>
    </div>

    <div v-if="article.is_outdated" class="outdated-notice">
      <i class="ri-error-warning-line" />
      <span>本文内容可能已过时，请谨慎参考。</span>
    </div>
    <FeaturesArticleAISummary v-if="article.ai_summary" :ai-summary="article.ai_summary" />

    <FeaturesArticleContent :content="article.content!" />

    <FeaturesArticleCopyright :article="article" />

    <FeaturesArticleTags :article="article" />

    <FeaturesArticleNavigation :prev="article.prev" :next="article.next" />

    <LazyFeaturesCommentComments target-type="article" :target-key="article.slug!" />
  </div>
</template>

<style lang="scss" scoped>
#post {
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.post-title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 16px;
  line-height: 1.3;
}

.post-meta {
  display: flex;
  gap: 20px;
  margin-bottom: 30px;
  font-size: 0.9rem;
  color: var(--theme-meta-color);
  flex-wrap: wrap;
  i {
    margin-right: 4px;
  }
  a {
    color: var(--theme-meta-color);
    text-decoration: none;
  }
}

.outdated-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 20px 0;
  padding: 12px 20px;
  background: rgba(250, 173, 20, 0.1);
  border: 1px solid rgba(250, 173, 20, 0.3);
  color: #d48806;
  font-size: 0.95rem;
}

@media screen and (max-width: 768px) {
  #post {
    padding: 20px;
  }
  .post-title {
    font-size: 1.5rem;
  }
}
</style>
