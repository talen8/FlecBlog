<template>
  <div id="page">
    <FeaturesArchiveArticleList
      v-if="category"
      :articles="articles"
      :title="`分类 - ${category.name}`"
      :total="total"
    />
    <UiPagination
      v-if="category && total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" setup>
import { getCategoryBySlug, getArticlesForWeb } from '@/composables/useApi';
import type { Category, Article } from '@@/types';

const route = useRoute();
const router = useRouter();
const category = ref<Category | null>(null);
const articles = ref<Article[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

const { data: initialData } = await useAsyncData(`category-${route.params.slug}`, async () => {
  const slug = route.params.slug as string;
  try {
    const [categoryData, articlesData] = await Promise.all([
      getCategoryBySlug(slug),
      getArticlesForWeb({ category: slug, page: 1, page_size: pageSize.value }),
    ]);
    return { category: categoryData, articles: articlesData.list, total: articlesData.total };
  } catch (error: unknown) {
    const err = error as Error & { response?: { status?: number } };
    if (err.response?.status === 404) {
      router.replace('/404');
    }
    return null;
  }
});

if (initialData.value) {
  category.value = initialData.value.category;
  articles.value = initialData.value.articles;
  total.value = initialData.value.total;
  currentPage.value = 1;
}

useHead({
  title: () => (category.value ? `分类:${category.value.name}` : undefined),
});

const fetchData = async (page = 1) => {
  const slug = route.params.slug as string;
  currentPage.value = page;
  try {
    const articlesData = await getArticlesForWeb({
      category: slug,
      page,
      page_size: pageSize.value,
    });
    articles.value = articlesData.list;
    total.value = articlesData.total;
  } catch (error: unknown) {
    const err = error as Error & { response?: { status?: number } };
    if (err.response?.status === 404) {
      router.replace('/404');
    }
  }
};

const handlePageChange = (page: number) => {
  fetchData(page);
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

watch(
  () => route.params.slug,
  () => {
    currentPage.value = 1;
    fetchData(1);
  }
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
