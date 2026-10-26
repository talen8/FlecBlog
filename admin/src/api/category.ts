import request from '@/utils/request';
import type { Category, CategoryListData } from '@/types/category';
import type { PaginationQuery } from '@/types/request';

/**
 * 获取分类列表
 * @param params 查询参数
 * @returns Promise<CategoryListData>
 */
export function getCategories(params?: PaginationQuery): Promise<CategoryListData> {
  return request.get('/admin/categories', { params });
}

/**
 * 创建分类
 * @param data 分类数据
 * @returns Promise<Category>
 */
export function createCategory(data: Partial<Category>): Promise<Category> {
  return request.post('/admin/categories', data);
}

/**
 * 更新分类
 * @param id 分类ID
 * @param data 分类数据
 * @returns Promise<Category>
 */
export function updateCategory(id: number, data: Partial<Category>): Promise<Category> {
  return request.put(`/admin/categories/${id}`, data);
}

/**
 * 删除分类
 * @param id 分类ID
 * @returns Promise<void>
 */
export function deleteCategory(id: number): Promise<void> {
  return request.delete(`/admin/categories/${id}`);
}
