<script lang="ts" setup>
import { getTagBySlug } from '@/composables/api/tag';
import { getArticlesForWeb } from '@/composables/api/article';
import type { Tag } from '@@/types/tag';
import type { Article } from '@@/types/article';
definePageMeta({});

const route = useRoute();
const router = useRouter();
const tag = ref<Tag | null>(null);
const articles = ref<Article[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

// 使用SSR获取标签详情和文章列表
const { data: initialData } = await useAsyncData(`tag-${route.params.slug}`, async () => {
  const slug = route.params.slug as string;

  try {
    const [tagData, articlesData] = await Promise.all([
      getTagBySlug(slug),
      getArticlesForWeb({
        tag: slug,
        page: 1,
        page_size: pageSize.value,
      }),
    ]);
    return { tag: tagData, articles: articlesData.list, total: articlesData.total };
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
  tag.value = initialData.value.tag;
  articles.value = initialData.value.articles;
  total.value = initialData.value.total;
  currentPage.value = 1;
}

// 动态页面标题
useHead({
  title: () => (tag.value ? `标签:${tag.value.name}` : undefined),
});

useSeoMeta({
  title: () => (tag.value ? `标签 - ${tag.value.name}` : '标签'),
  description: () =>
    tag.value
      ? `浏览 ${tag.value.name} 标签下的 ${total.value} 篇文章，发现更多相关内容`
      : '浏览标签下的文章',
});

const fetchData = async (page = 1) => {
  const slug = route.params.slug as string;
  currentPage.value = page;

  try {
    const articlesData = await getArticlesForWeb({
      tag: slug,
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
      v-if="tag"
      :articles="articles"
      :title="`标签 - ${tag.name}`"
      :total="total"
    />

    <!-- 分页 -->
    <UiPagination
      v-if="tag && total > pageSize"
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
