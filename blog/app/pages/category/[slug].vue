<script lang="ts" setup>
import { getCategoryBySlug } from '@/composables/api/category';
import { getArticlesForWeb } from '@/composables/api/article';
import type { Category } from '@@/types/category';
import type { Article } from '@@/types/article';
definePageMeta({});

const route = useRoute();
const router = useRouter();
const category = ref<Category | null>(null);
const articles = ref<Article[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

// 使用SSR获取分类详情和文章列表
const { data: initialData } = await useAsyncData(`category-${route.params.slug}`, async () => {
  const slug = route.params.slug as string;

  try {
    const [categoryData, articlesData] = await Promise.all([
      getCategoryBySlug(slug),
      getArticlesForWeb({
        category: slug,
        page: 1,
        page_size: pageSize.value,
      }),
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

// 初始化数据
if (initialData.value) {
  category.value = initialData.value.category;
  articles.value = initialData.value.articles;
  total.value = initialData.value.total;
  currentPage.value = 1;
}

// 动态页面标题
useHead({
  title: () => (category.value ? `分类:${category.value.name}` : undefined),
});

useSeoMeta({
  title: () => (category.value ? `分类 - ${category.value.name}` : '分类'),
  description: () =>
    category.value
      ? `浏览 ${category.value.name} 分类下的 ${total.value} 篇文章，探索更多相关内容`
      : '浏览分类下的文章',
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

// 处理分页变化
const handlePageChange = (page: number) => {
  fetchData(page);
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

// 监听路由参数变化
watch(
  () => route.params.slug,
  () => {
    currentPage.value = 1;
    fetchData(1);
  }
);
</script>

<template>
  <div id="page">
    <FeaturesArchiveArticleList
      v-if="category"
      :articles="articles"
      :title="`分类 - ${category.name}`"
      :total="total"
    />

    <!-- 分页 -->
    <UiPagination
      v-if="category && total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
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
