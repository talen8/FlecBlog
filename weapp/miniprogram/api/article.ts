/**
 * 文章相关 API 接口
 * 对接后端 /articles 接口
 */

import type {
  ArticleListItem,
  ArticleDetail,
  ArticleQuery,
  PageResult,
} from '../types';
import { get } from '../utils/request';

/**
 * 获取文章列表
 * @param params - 查询参数
 * @returns 分页文章列表
 */
export function getArticles(params: ArticleQuery = {}): Promise<PageResult<ArticleListItem>> {
  return get<PageResult<ArticleListItem>>('/articles', params);
}

/**
 * 获取文章详情
 * @param slug - 文章别名
 * @returns 文章详情
 */
export function getArticleBySlug(slug: string): Promise<ArticleDetail> {
  return get<ArticleDetail>(`/articles/${slug}`);
}

