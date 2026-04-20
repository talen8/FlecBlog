/**
 * API响应数据结构
 * @template T 响应数据类型
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，允许调用方不指定类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

/**
 * 分页查询参数
 */
export interface PaginationQuery {
  page?: number;
  page_size?: number;
}

/**
 * 分页响应数据结构
 * @template T 列表项数据类型
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，允许调用方不指定类型
export interface PaginationData<T = any> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}
