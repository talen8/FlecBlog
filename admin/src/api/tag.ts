import request from '@/utils/request';
import type { Tag, TagListData } from '@/types/tag';
import type { PaginationQuery } from '@/types/request';

/**
 * 获取标签列表
 * @param params 查询参数
 * @returns Promise<TagListData>
 */
export function getTags(params?: PaginationQuery): Promise<TagListData> {
  return request.get('/admin/tags', { params });
}

/**
 * 创建标签
 * @param data 标签数据
 * @returns Promise<Tag>
 */
export function createTag(data: Partial<Tag>): Promise<Tag> {
  return request.post('/admin/tags', data);
}

/**
 * 更新标签
 * @param id 标签ID
 * @param data 标签数据
 * @returns Promise<Tag>
 */
export function updateTag(id: number, data: Partial<Tag>): Promise<Tag> {
  return request.put(`/admin/tags/${id}`, data);
}

/**
 * 删除标签
 * @param id 标签ID
 * @returns Promise<void>
 */
export function deleteTag(id: number): Promise<void> {
  return request.delete(`/admin/tags/${id}`);
}
