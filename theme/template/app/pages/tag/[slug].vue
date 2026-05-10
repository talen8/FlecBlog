<template>
  <div id="page">
    <FeaturesArchiveArticleList
      v-if="tag"
      :articles="articles"
      :title="`标签 - ${tag.name}`"
      :total="total"
    />
    <UiPagination
      v-if="tag && total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" setup>
import { getTagBySlug, getArticlesForWeb } from '@/composables/useApi';
import type { Tag, Article } from '@@/types';

const route = useRoute();
const router = useRouter();
const tag = ref<Tag | null>(null);
const articles = ref<Article[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

const { data: initialData } = await useAsyncData(`tag-${route.params.slug}`, async () => {
  const slug = route.params.slug as string;
  try {
    const [tagData, articlesData] = await Promise.all([
      getTagBySlug(slug),
      getArticlesForWeb({ tag: slug, page: 1, page_size: pageSize.value }),
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

if (initialData.value) {
  tag.value = initialData.value.tag;
  articles.value = initialData.value.articles;
  total.value = initialData.value.total;
  currentPage.value = 1;
}

useHead({
  title: () => (tag.value ? `标签:${tag.value.name}` : undefined),
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
