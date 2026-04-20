// API响应数据结构
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

// 分页查询参数
export interface PaginationQuery {
  page: number;
  page_size: number;
}
