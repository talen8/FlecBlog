/**
 * 分类类型定义
 */
export interface Category {
  id: number;
  name: string;
  slug: string;
  url: string;
  description: string;
  count: number;
  sort: number;
}
