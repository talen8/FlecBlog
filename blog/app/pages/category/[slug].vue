<script lang="ts" setup>
definePageMeta({});

const route = useRoute();
const router = useRouter();
const slug = () => route.params.slug as string;

// SSR 获取分类详情 + 该分类下文章列表
const { data: category, error } = useCategory(slug);
const { data: list, page } = useArticleList(() => ({ category: slug() }));

const articles = computed(() => list.value?.list ?? []);
const total = computed(() => list.value?.total ?? 0);

// 分类不存在 → 404
watch(
  error,
  e => {
    const err = e as { statusCode?: number; response?: { status?: number } } | null;
    if (err?.statusCode === 404 || err?.response?.status === 404) router.replace('/404');
  },
  { immediate: true }
);

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
      v-if="category"
      :articles="articles"
      :title="`分类 - ${category.name}`"
      :total="total"
    />

    <!-- 分页 -->
    <UiPagination
      v-if="category && total > 10"
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
