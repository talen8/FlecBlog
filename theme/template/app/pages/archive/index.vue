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

<script lang="ts" setup>
import { getArticlesForWeb } from '@/composables/useApi';

const PAGE_SIZE = 20;
const { articles, total, currentPage, fetchArticles } = useArticles();

useSeoMeta({
  title: '归档',
  description: '浏览所有文章归档',
});

const { data: initialData } = await useAsyncData('articles-list', async () => {
  const { list, total: resTotal } = await getArticlesForWeb({
    page: 1,
    page_size: 20,
  });
  return { list, total: resTotal };
});

if (initialData.value) {
  articles.value = initialData.value.list || [];
  total.value = initialData.value.total || 0;
  currentPage.value = 1;
}

const handlePageChange = async (page: number) => {
  await fetchArticles({ page, page_size: PAGE_SIZE });
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }
}
</style>
