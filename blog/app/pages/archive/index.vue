<script lang="ts" setup>
definePageMeta({});

// SSR 获取归档列表（每页 20，按年分组展示）
const { data: list, page } = useArticleList(() => ({ page_size: 20 }));
const articles = computed(() => list.value?.list ?? []);
const total = computed(() => list.value?.total ?? 0);

useSeoMeta({
  title: '归档',
  description: () => `浏览所有文章归档，共 ${total.value} 篇文章，按时间顺序查看历史文章`,
});

// 处理分页变化
const handlePageChange = (p: number) => {
  page.value = p;
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};
</script>

<template>
  <div id="page">
    <FeaturesArchiveArticleList
      :articles="articles"
      :group-by-year="true"
      title="归档"
      :total="total"
    />

    <UiPagination
      v-if="total > 20"
      :total="total"
      :current-page="page"
      :page-size="20"
      @change="handlePageChange"
    />
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #page {
    padding: 30px;
  }
}

@media screen and (max-width: 768px) {
  #page {
    padding: 18px;
  }
}
</style>
