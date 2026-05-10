<script lang="ts" setup>
const { siteStats, fetchStats } = useStats();
const { blogConfig } = useSysConfig();

const runningDays = computed(() => {
  const established = blogConfig.value['blog.established'] || '2024-01-01';
  return Math.floor((Date.now() - new Date(established).getTime()) / 86400000);
});

onMounted(() => fetchStats());
</script>

<template>
  <div class="sidebar-card">
    <div class="headline"><i class="ri-line-chart-fill" /> 网站信息</div>
    <div class="info-row">
      <span>运行天数</span><span>{{ runningDays }} 天</span>
    </div>
    <div class="info-row">
      <span>访客量</span><span>{{ siteStats.total_visitors }}</span>
    </div>
    <div class="info-row">
      <span>浏览量</span><span>{{ siteStats.total_page_views }}</span>
    </div>
    <NuxtLink to="/statistics" class="stats-link">
      查看详细统计 <i class="ri-arrow-right-s-line" />
    </NuxtLink>
  </div>
</template>

<style lang="scss" scoped>
.headline {
  font-weight: 600;
  margin-bottom: 6px;
}
.headline i {
  margin-right: 4px;
}
.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  padding: 2px 0;
}
.info-row span:first-child {
  color: var(--theme-meta-color);
}
.stats-link {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 8px;
  padding: 4px 0;
  color: var(--theme-color);
  text-decoration: none;
  font-size: 0.85rem;
}
</style>
