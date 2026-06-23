<script lang="ts" setup>
definePageMeta({});

const route = useRoute();

// SSR 获取年月归档列表 + 分页
const { data: list, page } = useArticleList(() => ({
  year: route.params.year as string,
  month: route.params.month as string,
}));

const articles = computed(() => list.value?.list ?? []);
const total = computed(() => list.value?.total ?? 0);

// 列表标题：归档 - YYYY年MM月
const listTitle = computed(() => `归档 - ${route.params.year}年${route.params.month}月`);

// 标签页标题：YYYYMM归档
useHead({
  title: () => `${route.params.year}${route.params.month}归档`,
});

useSeoMeta({
  title: () => `${route.params.year}年${route.params.month}月归档`,
  description: () =>
    `浏览 ${route.params.year}年${route.params.month}月发布的所有文章，共 ${total.value} 篇`,
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
      :group-by-year="false"
      :title="listTitle"
      :total="total"
    />

    <UiPagination
      v-if="total > 10"
      :total="total"
      :current-page="page"
      :page-size="10"
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
