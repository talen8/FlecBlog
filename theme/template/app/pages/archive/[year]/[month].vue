<template>
  <div id="page">
    <FeaturesArchiveArticleList
      :articles="articles"
      :group-by-year="false"
      :title="listTitle"
      :total="total"
    />
    <UiPagination
      v-if="total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" setup>
import { getArticlesForWeb } from '@/composables/useApi';

const route = useRoute();
const { articles, total, currentPage, pageSize, fetchArticles } = useArticles();

useHead({
  title: () => `${route.params.year}${route.params.month}归档`,
});

useSeoMeta({
  title: () => `${route.params.year}年${route.params.month}月归档`,
  description: () => `浏览 ${route.params.year}年${route.params.month}月发布的所有文章`,
});

const listTitle = computed(() => `归档 - ${route.params.year}年${route.params.month}月`);

const { data: initialData } = await useAsyncData(
  `archive-${route.params.year}-${route.params.month}`,
  async () => {
    const articlesData = await getArticlesForWeb({
      year: route.params.year as string,
      month: route.params.month as string,
      page: 1,
      page_size: pageSize.value,
    });
    return {
      articles: articlesData.list,
      total: articlesData.total,
    };
  }
);

if (initialData.value) {
  articles.value = initialData.value.articles || [];
  total.value = initialData.value.total || 0;
  currentPage.value = 1;
}

const handlePageChange = async (page: number) => {
  await fetchArticles({
    year: route.params.year as string,
    month: route.params.month as string,
    page,
  });
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

watch(
  () => [route.params.year, route.params.month],
  () => {
    fetchArticles({
      year: route.params.year as string,
      month: route.params.month as string,
      page: 1,
    });
  },
  { deep: true }
);
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
