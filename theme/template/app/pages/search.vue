<script lang="ts" setup>
import { searchArticles } from '@/composables/useApi';
import type { Article } from '@@/types';

const route = useRoute();
const router = useRouter();

const keyword = ref((route.query.q as string) || '');
const articles = ref<Article[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const searching = ref(false);

const doSearch = async (page = 1) => {
  if (!keyword.value.trim()) return;

  currentPage.value = page;
  searching.value = true;

  try {
    router.replace({ query: { q: keyword.value } });
    const { list, total: resTotal } = await searchArticles(keyword.value, {
      page,
      page_size: pageSize.value,
    });
    articles.value = list || [];
    total.value = resTotal || 0;
  } catch (error) {
    console.error('搜索失败:', error);
    articles.value = [];
    total.value = 0;
  } finally {
    searching.value = false;
  }
};

if (keyword.value) {
  await useAsyncData('search-articles', () => doSearch());
}

const handlePageChange = (page: number) => {
  doSearch(page);
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

useSeoMeta({
  title: '搜索',
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">搜索</h1>
    <div class="search-form">
      <input
        v-model="keyword"
        type="text"
        placeholder="输入关键词搜索文章..."
        class="search-input"
        @keyup.enter="doSearch()"
      />
      <button class="search-btn" :disabled="searching || !keyword.trim()" @click="doSearch()">
        {{ searching ? '搜索中...' : '搜索' }}
      </button>
    </div>

    <div v-if="articles.length > 0" class="article-list">
      <article v-for="article in articles" :key="article.id" class="article-item">
        <div class="article-info">
          <NuxtLink :to="article.url" class="article-title">{{ article.title }}</NuxtLink>
          <div class="article-meta">
            <span class="date">{{ formatDate(article.publish_time) }}</span>
            <span v-if="article.category" class="category">
              <NuxtLink :to="article.category.url">{{ article.category.name }}</NuxtLink>
            </span>
          </div>
          <p class="article-summary">{{ article.summary }}</p>
        </div>
      </article>
    </div>

    <div v-else-if="keyword && !searching" class="empty-state">未找到相关文章</div>

    <UiPagination
      v-if="total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 20px;
  font-size: 2rem;
  font-weight: 700;
}

.search-form {
  display: flex;
  gap: 12px;
  margin-bottom: 30px;
}

.search-input {
  flex: 1;
  padding: 12px;
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.95rem;

  &:focus {
    outline: none;
    border-color: var(--theme-color);
  }
}

.search-btn {
  padding: 12px 24px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  cursor: pointer;

  &:hover:not(:disabled) {
    opacity: 0.9;
  }
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.article-item {
  padding-bottom: 20px;
  border-bottom: 1px solid var(--flec-border);

  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
}

.article-title {
  display: block;
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--font-color);
  text-decoration: none;
  margin-bottom: 8px;

  &:hover {
    color: var(--theme-color);
  }
}

.article-meta {
  display: flex;
  gap: 16px;
  font-size: 0.85rem;
  color: var(--theme-meta-color);
  margin-bottom: 8px;

  a {
    color: var(--theme-meta-color);
    text-decoration: none;
    &:hover {
      color: var(--theme-color);
    }
  }
}

.article-summary {
  margin: 0;
  color: var(--theme-meta-color);
  font-size: 0.95rem;
  line-height: 1.6;
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
  #page {
    padding: 20px;
  }
  .page-title {
    font-size: 1.5rem;
  }
  .search-form {
    flex-direction: column;
  }
}
</style>
