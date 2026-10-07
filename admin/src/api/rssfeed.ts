import request from '@/utils/request';
import type { RssArticleQuery, RssArticleListData } from '@/types/rssfeed';

/**
 * 全部标记已读响应
 */
export interface MarkAllReadResponse {
  affected: number;
}

/**
 * 获取RSS文章列表
 * @param params 查询参数
 * @returns Promise<RssArticleListData>
 */
export const getRssArticles = async (params?: RssArticleQuery): Promise<RssArticleListData> => {
  return request.get('/admin/rssfeed', { params });
};

/**
 * 标记文章已读
 * @param id 文章ID
 * @returns Promise<void>
 */
export const markRssArticleRead = async (id: number): Promise<void> => {
  await request.put(`/admin/rssfeed/${id}/read`);
};

/**
 * 设置文章的友圈屏蔽状态
 * @param id 文章ID
 * @param block true 屏蔽，false 取消屏蔽
 * @returns Promise<void>
 */
export const setRssArticleBlockCircle = async (id: number, block: boolean): Promise<void> => {
  await request.put(`/admin/rssfeed/${id}/block-circle`, { block_circle: block });
};

/**
 * 全部标记已读
 * @returns Promise<MarkAllReadResponse>
 */
export const markAllRssArticlesRead = async (): Promise<MarkAllReadResponse> => {
  return request.put('/admin/rssfeed/read-all');
};
