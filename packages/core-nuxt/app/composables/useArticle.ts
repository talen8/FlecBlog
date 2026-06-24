import type { Article, ArticleQuery } from '../../types/article';
import {
  getArticleList,
  getArticleBySlug,
  searchArticles,
  getRandomArticleSlug,
} from './api/article';

/**
 * 获取当前正在浏览的文章（跨页面共享状态）
 * @returns currentArticle - 当前文章对象
 * @returns setCurrentArticle - 设置当前文章
 * @returns clearCurrentArticle - 清除当前文章
 */
export function useCurrentArticle() {
  const currentArticle = useState<Article | null>('currentArticle', () => null);

  return {
    currentArticle,
    setCurrentArticle: (article: Article | null) => (currentArticle.value = article),
    clearCurrentArticle: () => (currentArticle.value = null),
  };
}

/**
 * 按 slug 获取文章详情（支持 SSR）
 * @param slug - 文章 slug，支持 ref / computed / getter
 * @returns data - 文章对象，refresh - 刷新方法
 */
export function useArticle(slug: MaybeRefOrGetter<string>) {
  return useAsyncData(`article:${toValue(slug)}`, () => getArticleBySlug(toValue(slug)), {
    watch: [() => toValue(slug)],
  });
}

/**
 * 分页获取文章列表（支持 SSR）
 * @param query - 筛选条件 { page_size, category, tag 等 }，支持 ref / computed / getter
 * @returns data - 文章分页数据
 * @returns page - 当前页码（ref）
 * @returns refresh - 刷新方法
 */
export function useArticleList(query: MaybeRefOrGetter<ArticleQuery> = () => ({})) {
  const page = ref(1);

  watch(
    () => JSON.stringify(toValue(query)),
    () => {
      page.value = 1;
    }
  );

  const result = useAsyncData(
    `articles:${JSON.stringify(toValue(query))}`,
    () => getArticleList({ page_size: 10, ...toValue(query), page: page.value }),
    { watch: [page, () => JSON.stringify(toValue(query))] }
  );

  return { ...result, page };
}

/**
 * 搜索文章（客户端交互式）
 * @returns results - 搜索结果列表
 * @returns total - 结果总数
 * @returns loading - 加载状态
 * @returns search(keyword, page?, pageSize?) - 执行搜索
 * @returns reset - 重置结果
 */
export function useSearch() {
  const results = useState<Article[]>('search-results', () => []);
  const total = useState<number>('search-total', () => 0);
  const loading = useState<boolean>('search-loading', () => false);

  const search = async (keyword: string, page = 1, pageSize = 5) => {
    const kw = keyword.trim();
    if (!kw) {
      results.value = [];
      total.value = 0;
      return;
    }
    loading.value = true;
    try {
      const data = await searchArticles(kw, { page, page_size: pageSize });
      results.value = data.list;
      total.value = data.total;
    } catch {
      results.value = [];
      total.value = 0;
    } finally {
      loading.value = false;
    }
  };

  const reset = () => {
    results.value = [];
    total.value = 0;
  };

  return { results, total, loading, search, reset };
}

/**
 * 随机跳转到一篇文章
 * @returns loading - 跳转中状态，go - 执行随机跳转
 */
export function useRandomArticle() {
  const loading = useState<boolean>('random-article-loading', () => false);

  const go = async () => {
    if (loading.value) return;
    loading.value = true;
    try {
      const slug = await getRandomArticleSlug();
      await navigateTo(`/posts/${slug}`);
    } finally {
      loading.value = false;
    }
  };

  return { loading, go };
}
