<script setup lang="ts">
const { siteStats } = useStats();
const { blogConfig } = useSysConfig();

const runningDays = computed(() => {
  const established = blogConfig.value.established || '2024-01-01';
  const startDate = new Date(established).getTime();
  const now = Date.now();
  return Math.floor((now - startDate) / 86400000);
});
</script>

<template>
  <div class="card-widget card-webinfo">
    <div class="item-headline">
      <i class="ri-line-chart-fill" />
      <span>网站信息</span>
    </div>
    <div class="webinfo">
      <div class="webinfo-item">
        <div class="item-name">本站总字数 :</div>
        <div class="item-count">{{ siteStats.total_words }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">本站访客量:</div>
        <div class="item-count">{{ siteStats.total_visitors }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">本站总浏览量 :</div>
        <div class="item-count">{{ siteStats.total_page_views }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">当前在线人数 :</div>
        <div class="item-count">{{ siteStats.online_users }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">网站运行天数 :</div>
        <div class="item-count">{{ runningDays }}</div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.webinfo-item {
  display: flex;
  align-items: center;
  padding: 2px 10px 0;

  .item-name {
    flex: 1;
    padding-right: 20px;
  }
}

@media screen and (max-width: 900px) {
  .webinfo-item {
    padding: 2px 6px 0;
    font-size: 0.95rem;

    .item-name {
      padding-right: 12px;
    }
  }
}
</style>
