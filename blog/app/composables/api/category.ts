import type { Category } from '../../../types/category';
import { createApi } from './createApi';

const categoryApi = createApi<Category>('/categories');

/** 获取分类列表 */
export const getCategories = async () => {
  return categoryApi.getList();
};

/** 获取分类详情 */
export const getCategoryBySlug = async (slug: string) => {
  return categoryApi.getOne(slug);
};
