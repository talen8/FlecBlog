import type { Tag } from '@@/types/tag';
import { createApi } from './createApi';

const tagApi = createApi<Tag>('/tags');

/** 获取标签列表 */
export const getTags = async () => {
  return tagApi.getList();
};

/** 根据ID获取标签详情 */
export const getTagById = async (id: number) => {
  return tagApi.getOne(id);
};

/** 根据Slug获取标签详情 */
export const getTagBySlug = async (slug: string) => {
  return tagApi.getOne(slug);
};
