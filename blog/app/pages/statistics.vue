<script setup lang="ts">
import type { SiteStats } from '@@/types/stats';
import { getSiteStats } from '@/composables/api/stats';

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '统计',
  description: '公开展示本站的文章、评论、友链、分类、标签、动态与访问情况等统计数据',
});

const { blogConfig } = useSysConfig();

const emptyStats: SiteStats = {
  total_words: '0',
  total_visitors: 0,
  total_page_views: 0,
  online_users: 0,
  total_articles: 0,
  total_comments: 0,
  total_friends: 0,
  total_moments: 0,
  total_categories: 0,
  total_tags: 0,
  today_visitors: 0,
  today_pageviews: 0,
  yesterday_visitors: 0,
  yesterday_pageviews: 0,
  month_pageviews: 0,
};

const { data } = await useAsyncData('site-stats-page', async () => {
  try {
    return await getSiteStats();
  } catch (error) {
    console.error('获取站点统计数据失败:', error);
    return emptyStats;
  }
});

const stats = computed<SiteStats>(() => ({
  ...emptyStats,
  ...(data.value ?? {}),
}));

const establishedDate = computed(() => blogConfig.value.established || '2024-01-01');

const runningDays = computed(() => {
  const startDate = new Date(establishedDate.value).getTime();
  return Math.max(0, Math.floor((Date.now() - startDate) / 86400000));
});

const formatNumber = (value: string | number) => {
  if (typeof value === 'string') return value;
  if (value >= 1000000) return `${(value / 1000000).toFixed(1)}M`;
  if (value >= 10000) return `${(value / 10000).toFixed(1)}w`;
  if (value >= 1000) return `${(value / 1000).toFixed(1)}k`;
  return `${value}`;
};

const overviewCards = computed(() => [
  {
    key: 'words',
    label: '字数',
    value: stats.value.total_words,
    icon: 'ri-file-text-line',
    hint: '已发布文章全文累计',
  },
  {
    key: 'pageviews',
    label: '浏览',
    value: stats.value.total_page_views,
    icon: 'ri-bar-chart-box-line',
    hint: '累计页面浏览次数',
  },
  {
    key: 'visitors',
    label: '访客',
    value: stats.value.total_visitors,
    icon: 'ri-user-line',
    hint: '累计独立访客数量',
  },
  {
    key: 'articles',
    label: '文章',
    value: stats.value.total_articles,
    icon: 'ri-article-line',
    hint: '已发布内容总量',
  },
  {
    key: 'categories',
    label: '分类',
    value: stats.value.total_categories,
    icon: 'ri-folder-2-line',
    hint: '已有内容的分类数量',
  },
  {
    key: 'tags',
    label: '标签',
    value: stats.value.total_tags,
    icon: 'ri-price-tag-3-line',
    hint: '可浏览标签总数',
  },
  {
    key: 'moments',
    label: '动态',
    value: stats.value.total_moments,
    icon: 'ri-quill-pen-line',
    hint: '公开动态记录数量',
  },
  {
    key: 'friends',
    label: '友链',
    value: stats.value.total_friends,
    icon: 'ri-links-line',
    hint: '有效站点链接数量',
  },
  {
    key: 'comments',
    label: '评论',
    value: stats.value.total_comments,
    icon: 'ri-message-3-line',
    hint: '公开可见评论数量',
  },
]);

const visitCards = computed(() => [
  {
    label: '今日访客',
    value: stats.value.today_visitors,
  },
  {
    label: '今日浏览',
    value: stats.value.today_pageviews,
  },
  {
    label: '当前在线',
    value: stats.value.online_users,
  },
  {
    label: '昨日访客',
    value: stats.value.yesterday_visitors,
  },
  {
    label: '昨日浏览',
    value: stats.value.yesterday_pageviews,
  },
  {
    label: '本月浏览',
    value: stats.value.month_pageviews,
  },
]);
</script>

<template>
  <div id="statistics-page">
    <h1 class="page-title">统计</h1>

    <section id="overview" class="content-section">
      <div class="overview-grid">
        <article v-for="item in overviewCards" :key="item.key" class="stat-card">
          <div class="stat-card__icon">
            <i :class="item.icon" />
          </div>
          <div class="stat-card__content">
            <span class="stat-card__label">{{ item.label }}</span>
            <strong class="stat-card__value">{{ formatNumber(item.value) }}</strong>
            <span class="stat-card__hint">{{ item.hint }}</span>
          </div>
        </article>
      </div>
    </section>

    <section class="content-section split-section">
      <article id="traffic" class="panel-block">
        <div class="section-head">
          <div>
            <span class="section-kicker">Traffic</span>
            <h2>访问概览</h2>
          </div>
        </div>

        <div class="traffic-list">
          <div v-for="item in visitCards" :key="item.label" class="traffic-item">
            <span class="traffic-item__label">{{ item.label }}</span>
            <strong class="traffic-item__value">
              {{ formatNumber(item.value) }}
            </strong>
          </div>
        </div>
      </article>

      <article class="panel-block note-panel">
        <div class="section-head">
          <div>
            <span class="section-kicker">Runtime</span>
            <h2>站点信息</h2>
          </div>
        </div>

        <div class="runtime-cards">
          <article class="runtime-card">
            <span class="runtime-card__label">运行时长</span>
            <strong class="runtime-card__value">{{ runningDays }} 天</strong>
          </article>
          <article class="runtime-card">
            <span class="runtime-card__label">建站日期</span>
            <strong class="runtime-card__value">{{ establishedDate }}</strong>
          </article>
        </div>
      </article>
    </section>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#statistics-page {
  @extend .cardHover;
  width: 100%;
  padding: 40px;
  display: flex;
  flex-direction: column;
  gap: 32px;

  .page-title {
    margin: 0 0 6px;
    font-size: 2rem;
    font-weight: 700;
    line-height: 1.2;
  }

  .stat-card__hint,
  .traffic-item__label,
  .runtime-card__label {
    color: var(--theme-meta-color);
  }

  .section-kicker {
    display: inline-block;
    margin-bottom: 10px;
    color: var(--theme-color);
    font-size: 0.82rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
  }

  .section-head,
  .stat-card {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
  }

  .content-section {
    padding-top: 4px;
  }

  .split-section {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 18px;
  }

  .panel-block {
    padding: 24px;
    border-radius: 22px;
    background:
      linear-gradient(180deg, rgba(73, 177, 245, 0.08), transparent 120px), var(--flec-card-bg);
    border: 1px solid var(--flec-border-color);
  }

  .section-head {
    margin-bottom: 22px;

    h2 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 700;
    }
  }

  .overview-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 14px;
  }

  .stat-card {
    min-width: 0;
    padding: 20px;
    border-radius: 20px;
    background: var(--flec-card-bg);
    border: 1px solid var(--flec-border-color);
  }

  .stat-card__icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 52px;
    height: 52px;
    border-radius: 16px;
    background: linear-gradient(135deg, rgba(73, 177, 245, 0.16), rgba(73, 177, 245, 0.04));
    color: var(--theme-color);
    flex-shrink: 0;

    i {
      font-size: 1.45rem;
    }
  }

  .stat-card__content {
    flex: 1;
    min-width: 0;
  }

  .stat-card__label {
    display: block;
    color: var(--theme-meta-color);
  }

  .stat-card__value {
    display: block;
    margin: 8px 0 6px;
    font-size: 1.8rem;
    line-height: 1.1;
  }

  .traffic-list {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 14px;
  }

  .runtime-cards {
    display: grid;
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .traffic-item,
  .runtime-card {
    padding: 18px;
    border-radius: 18px;
    background: var(--flec-heavy-bg);
  }

  .runtime-card {
    min-height: 96px;
  }

  .runtime-card__label {
    display: block;
  }

  .traffic-item__value,
  .runtime-card__value {
    display: block;
    margin-top: 8px;
    font-size: 1.7rem;
    line-height: 1.1;
  }
}

@media screen and (max-width: 1024px) {
  #statistics-page {
    padding: 30px;
    gap: 26px;

    .split-section {
      grid-template-columns: 1fr;
    }

    .overview-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
}

@media screen and (max-width: 768px) {
  #statistics-page {
    padding: 18px;
    gap: 22px;

    .page-title {
      font-size: 1.4rem;
    }

    .panel-block,
    .stat-card {
      padding: 18px;
    }

    .overview-grid,
    .traffic-list {
      grid-template-columns: 1fr;
    }

    .section-head,
    .stat-card {
      flex-direction: column;
    }

    .stat-card__value,
    .traffic-item__value,
    .runtime-card__value {
      font-size: 1.45rem;
    }
  }
}
</style>
