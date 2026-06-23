import type { Category } from '@@/types/category';
import { createApi } from './createApi';

const categoryApi = createApi<Category>('/categories');

/** 获取分类列表 */
export const getCategories = async () => {
  return categoryApi.getList();
};

/** 获取分类详情（By ID） */
export const getCategoryById = async (id: number) => {
  return categoryApi.getOne(id);
};

/** 获取分类详情（By Slug） */
export const getCategoryBySlug = async (slug: string) => {
  return categoryApi.getOne(slug);
};
