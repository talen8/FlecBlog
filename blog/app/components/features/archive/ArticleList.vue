<script lang="ts" setup>
import type { Article } from '@@/types/article';

interface Props {
  articles: Article[];
  groupByYear?: boolean; // 是否按年份分组
  title?: string; // 页面标题
  total?: number; // 总文章数（分页时使用）
}

const props = withDefaults(defineProps<Props>(), {
  groupByYear: false,
});

// 显示的文章数量：优先使用 total，否则使用 articles.length
const displayTotal = computed(() => {
  return props.total ?? props.articles.length;
});

// 按年份对文章进行分组
const groupedArticles = computed(() => {
  if (!props.groupByYear || !props.articles) {
    return null;
  }

  const groups = new Map<string, Article[]>();

  props.articles.forEach(article => {
    if (article.publish_time) {
      const year = new Date(article.publish_time).getFullYear().toString();
      if (!groups.has(year)) {
        groups.set(year, []);
      }
      groups.get(year)!.push(article);
    }
  });

  // 转换为按年份排序的数组
  return Array.from(groups.entries())
    .sort((a, b) => Number(b[0]) - Number(a[0]))
    .map(([year, articles]) => ({
      year,
      articles: articles.sort((a, b) => {
        const bTime = b.publish_time ? new Date(b.publish_time).getTime() : 0;
        const aTime = a.publish_time ? new Date(a.publish_time).getTime() : 0;
        return bTime - aTime;
      }),
    }));
});
</script>

<template>
  <div class="article-sort">
    <!-- 页面头部 -->
    <div v-if="title" class="article-sort-header">
      <h1 class="article-sort-title">{{ title }}</h1>
      <div class="article-sort-meta">
        <i class="ri-file-list-line"></i>
        <span>共 {{ displayTotal }} 篇文章</span>
      </div>
    </div>

    <!-- 按年份分组显示 -->
    <template v-if="groupByYear && groupedArticles">
      <template v-for="group in groupedArticles" :key="group.year">
        <div class="article-sort-item year">{{ group.year }}</div>
        <div v-for="article in group.articles" :key="article.id" class="article-sort-item">
          <NuxtLink v-if="article.cover" :to="article.url" class="article-sort-item-img">
            <NuxtImg :src="article.cover" :alt="article.title" loading="lazy" />
          </NuxtLink>
          <div class="article-sort-item-info">
            <div class="article-sort-item-time">
              <i class="ri-calendar-2-fill"></i>
              <span>{{ formatDate(article.publish_time) }}</span>
            </div>
            <NuxtLink :to="article.url" class="article-sort-item-title">{{
              article.title
            }}</NuxtLink>
          </div>
        </div>
      </template>
    </template>

    <!-- 直接列表显示 -->
    <template v-else>
      <div v-for="article in articles" :key="article.id" class="article-sort-item">
        <NuxtLink v-if="article.cover" :to="article.url" class="article-sort-item-img">
          <NuxtImg :src="article.cover" :alt="article.title" loading="lazy" />
        </NuxtLink>
        <div class="article-sort-item-info">
          <div class="article-sort-item-time">
            <i class="ri-calendar-2-fill"></i>
            <span>{{ formatDate(article.publish_time) }}</span>
          </div>
          <NuxtLink :to="article.url" class="article-sort-item-title">{{ article.title }}</NuxtLink>
        </div>
      </div>
    </template>
  </div>
</template>

<style lang="scss">
@use '@/assets/css/mixins' as *;

.article-sort {
  .article-sort-header {
    margin-bottom: 30px;

    .article-sort-title {
      margin: 0 0 10px;
      font-weight: bold;
      font-size: 2rem;
    }

    .article-sort-meta {
      color: var(--theme-meta-color);
      font-size: 0.9rem;

      i {
        margin-right: 6px;
      }
    }
  }

  .article-sort-item {
    position: relative;
    display: flex;
    align-items: center;
    margin-bottom: 20px;

    &.year {
      margin-bottom: 10px;
      font-size: 1.2rem;
      font-weight: 600;
    }

    .article-sort-item-img {
      overflow: hidden;
      width: 120px;
      height: 70px;
      border-radius: 6px;

      img {
        @extend .imgHover;
      }
    }

    .article-sort-item-info {
      flex: 1;
      padding: 0 16px;

      .article-sort-item-time {
        color: var(--theme-meta-color);
        font-size: 0.85rem;

        span {
          padding-left: 6px;
        }
      }

      .article-sort-item-title {
        color: var(--font-color);
        font-size: 1.05rem;
        transition: all 0.3s;

        &:hover {
          color: var(--theme-color);
        }
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .article-sort {
    .article-sort-header {
      margin-bottom: 20px;

      .article-sort-title {
        font-size: 1.4rem;
      }

      .article-sort-meta {
        font-size: 0.85rem;
      }
    }

    .article-sort-item {
      margin-bottom: 15px;

      &.year {
        font-size: 1.05rem;
      }

      .article-sort-item-img {
        width: 90px;
        height: 55px;
      }

      .article-sort-item-info {
        padding: 0 12px;

        .article-sort-item-time {
          font-size: 0.78rem;
        }

        .article-sort-item-title {
          font-size: 0.92rem;
        }
      }
    }
  }
}
</style>
