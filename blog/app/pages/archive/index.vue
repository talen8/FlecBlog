<script lang="ts" setup>
import { getArticlesForWeb } from '@/composables/api/article';

definePageMeta({});

const PAGE_SIZE = 20;
const { articles, total, currentPage, fetchArticles } = useArticles();

useSeoMeta({
  title: '归档',
  description: () => `浏览所有文章归档，共 ${total.value} 篇文章，按时间顺序查看历史文章`,
});

// 使用SSR获取归档数据
const { data: initialData } = await useAsyncData('articles-list', async () => {
  const { list, total: resTotal } = await getArticlesForWeb({
    page: 1,
    page_size: 20,
  });
  return { list, total: resTotal };
});

// 初始化数据
if (initialData.value) {
  articles.value = initialData.value.list || [];
  total.value = initialData.value.total || 0;
  currentPage.value = 1;
}

// 加载数据（总览页：按年分组显示）
const loadData = async (page: number = 1) => {
  await fetchArticles({
    page,
    page_size: PAGE_SIZE,
  });
};

// 处理分页变化
const handlePageChange = async (page: number) => {
  await loadData(page);
  if (process.client) {
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
      v-if="total > PAGE_SIZE"
      :total="total"
      :current-page="currentPage"
      :page-size="PAGE_SIZE"
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
