<script setup lang="ts">
import { getArticlesForWeb } from '@/composables/useApi';

const { articles, total, currentPage, pageSize, fetchArticles } = useArticles();
const { blogConfig } = useSysConfig();

const { data: initialData } = await useAsyncData('articles-list', async () => {
  const { list, total: resTotal } = await getArticlesForWeb({
    page: 1,
    page_size: pageSize.value,
  });
  return { list, total: resTotal };
});

if (initialData.value) {
  articles.value = initialData.value.list || [];
  total.value = initialData.value.total || 0;
  currentPage.value = 1;
}

const handlePageChange = async (page: number) => {
  await fetchArticles({ page });
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

useSeoMeta({
  title: '首页',
  description: () => blogConfig.value?.['blog.description'] || '欢迎来到我的博客',
});
</script>

<template>
  <div id="home">
    <div v-if="articles.length > 0" class="article-list">
      <article v-for="article in articles" :key="article.id" class="article-card">
        <div v-if="article.cover" class="article-cover">
          <NuxtLink :to="article.url">
            <img :src="article.cover" :alt="article.title" loading="lazy" />
          </NuxtLink>
        </div>
        <div class="article-info">
          <div class="title-row">
            <NuxtLink class="article-title" :to="article.url">{{ article.title }}</NuxtLink>
            <span v-if="article.is_top" class="badge">置顶</span>
            <span v-else-if="article.is_essence" class="badge">精选</span>
          </div>
          <div class="article-meta">
            <span class="date"
              ><i class="ri-calendar-line" /> {{ formatDate(article.publish_time) }}</span
            >
            <span v-if="article.category" class="category">
              <i class="ri-folder-line" />
              <NuxtLink :to="article.category.url">{{ article.category.name }}</NuxtLink>
            </span>
            <span class="comments"
              ><i class="ri-message-3-line" /> {{ article.comment_count }}条评论</span
            >
          </div>
          <p class="article-summary">{{ article.summary }}</p>
        </div>
      </article>
    </div>
    <div v-else class="empty-state">暂无文章</div>

    <UiPagination
      v-if="articles.length > 0"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<style lang="scss" scoped>
#home {
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  padding: 30px;
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.article-card {
  display: flex;
  gap: 20px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--flec-border);
  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
}

.article-cover {
  flex-shrink: 0;
  width: 200px;
  height: 130px;
  overflow: hidden;
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.article-info {
  flex: 1;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.article-title {
  color: var(--font-color);
  text-decoration: none;
  font-size: 1.3rem;
  font-weight: 600;
}

.badge {
  padding: 2px 8px;
  font-size: 0.75rem;
  border: 1px solid var(--theme-color);
  color: var(--theme-color);
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 12px;
  font-size: 0.85rem;
  color: var(--theme-meta-color);
  i {
    margin-right: 4px;
  }
  a {
    color: var(--theme-meta-color);
    text-decoration: none;
  }
}

.article-summary {
  color: var(--theme-meta-color);
  font-size: 0.95rem;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--theme-meta-color);
}

@media screen and (max-width: 768px) {
  #home {
    padding: 20px;
  }
  .article-card {
    flex-direction: column;
    gap: 12px;
  }
  .article-cover {
    width: 100%;
    height: auto;
    aspect-ratio: 16 / 9;
  }
  .article-title {
    font-size: 1.1rem;
  }
}
</style>
