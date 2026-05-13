import type { Article, ArticleQuery } from '@@/types/article';
import type { PaginationData } from '@@/types/request';
import { createApi } from './createApi';

const articleApi = createApi<Article>('/articles');

/** 获取文章列表 */
export const getArticlesForWeb = async (params: ArticleQuery = {}) => {
  return articleApi.getList(params);
};

/** 获取文章详情 */
export const getArticleBySlug = async (slug: string) => {
  return articleApi.getOne(slug);
};

/** 搜索文章 */
export const searchArticles = async (keyword: string, params: Partial<ArticleQuery> = {}) => {
  return articleApi.get<PaginationData<Article>>('/search', { keyword, ...params });
};

/** 随机文章 slug */
export const getRandomArticleSlug = async () => {
  return articleApi.get<string>('/random');
};
