<template>
  <div id="page">
    <h1 class="page-title">{{ title }}</h1>
    <p v-if="total" class="page-subtitle">共 {{ total }} 篇文章</p>

    <template v-if="groupByYear && groupedArticles.length > 0">
      <div v-for="group in groupedArticles" :key="group.year" class="year-group">
        <h2 class="year-title">{{ group.year }}</h2>
        <div class="article-list">
          <article v-for="article in group.articles" :key="article.id" class="article-item">
            <div class="article-info">
              <NuxtLink :to="article.url" class="article-title">{{ article.title }}</NuxtLink>
              <div class="article-meta">
                <span class="date">{{ formatDate(article.publish_time) }}</span>
                <span v-if="article.category" class="category">
                  <NuxtLink :to="article.category.url">{{ article.category.name }}</NuxtLink>
                </span>
              </div>
            </div>
            <img
              v-if="article.cover"
              :src="article.cover"
              :alt="article.title"
              class="article-cover"
              loading="lazy"
            />
          </article>
        </div>
      </div>
    </template>

    <template v-else-if="articles.length > 0">
      <div class="article-list">
        <article v-for="article in articles" :key="article.id" class="article-item">
          <div class="article-info">
            <NuxtLink :to="article.url" class="article-title">{{ article.title }}</NuxtLink>
            <div class="article-meta">
              <span class="date">{{ formatDate(article.publish_time) }}</span>
              <span v-if="article.category" class="category">
                <NuxtLink :to="article.category.url">{{ article.category.name }}</NuxtLink>
              </span>
            </div>
          </div>
          <img
            v-if="article.cover"
            :src="article.cover"
            :alt="article.title"
            class="article-cover"
            loading="lazy"
          />
        </article>
      </div>
    </template>

    <div v-else class="empty-state">暂无文章</div>
  </div>
</template>

<script lang="ts" setup>
import type { Article } from '@@/types';

const props = defineProps<{
  articles: Article[];
  title: string;
  total?: number;
  groupByYear?: boolean;
}>();

const groupedArticles = computed(() => {
  if (!props.groupByYear || !props.articles.length) return [];

  const groups: { year: string; articles: Article[] }[] = [];
  const map = new Map<string, Article[]>();

  for (const article of props.articles) {
    const year = new Date(article.publish_time).getFullYear().toString();
    if (!map.has(year)) {
      map.set(year, []);
    }
    map.get(year)!.push(article);
  }

  for (const [year, articles] of map) {
    groups.push({ year, articles });
  }

  groups.sort((a, b) => b.year.localeCompare(a.year));
  return groups;
});
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 10px;
  font-size: 2rem;
  font-weight: 700;
}

.page-subtitle {
  margin: 0 0 30px;
  color: var(--theme-meta-color);
  font-size: 0.9rem;
}

.year-group {
  margin-bottom: 30px;
  &:last-child {
    margin-bottom: 0;
  }
}

.year-title {
  font-size: 1.4rem;
  font-weight: 700;
  margin: 0 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--flec-border);
  color: var(--theme-color);
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.article-item {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--flec-border);
  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
}

.article-info {
  flex: 1;
  min-width: 0;
}

.article-title {
  display: block;
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--font-color);
  text-decoration: none;
  margin-bottom: 8px;
}

.article-meta {
  display: flex;
  gap: 16px;
  font-size: 0.85rem;
  color: var(--theme-meta-color);
  a {
    color: var(--theme-meta-color);
    text-decoration: none;
  }
}

.article-cover {
  width: 160px;
  height: 100px;
  object-fit: cover;
  flex-shrink: 0;
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
  .article-item {
    flex-direction: column;
    gap: 12px;
  }
  .article-cover {
    width: 100%;
    height: auto;
    aspect-ratio: 16 / 9;
  }
}
</style>
