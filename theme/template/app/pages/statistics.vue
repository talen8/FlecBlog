<template>
  <div id="page">
    <h1 class="page-title">统计</h1>
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_articles }}</div>
        <div class="stat-label">文章</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_comments }}</div>
        <div class="stat-label">评论</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_visitors }}</div>
        <div class="stat-label">访客</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_page_views }}</div>
        <div class="stat-label">浏览量</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_categories }}</div>
        <div class="stat-label">分类</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_tags }}</div>
        <div class="stat-label">标签</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_moments }}</div>
        <div class="stat-label">动态</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ siteStats.total_friends }}</div>
        <div class="stat-label">友链</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { getSiteStats } from '@/composables/useApi';

const { siteStats } = useStats();

useSeoMeta({
  title: '统计',
  description: '网站统计数据',
});

const { data: initialData } = await useAsyncData('site-stats', async () => {
  const stats = await getSiteStats();
  return stats;
});

if (initialData.value) {
  siteStats.value = initialData.value;
}
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 20px;
}

.stat-card {
  padding: 24px;
  background: var(--flec-heavy-bg);
  border-radius: 12px;
  text-align: center;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--theme-color);
  margin-bottom: 8px;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--theme-meta-color);
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-value {
    font-size: 1.5rem;
  }
}
</style>
