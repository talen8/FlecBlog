import type { PaginationQuery } from './request';

// 订阅者实体
export interface Subscriber {
  id: number;
  email: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

// 订阅者查询参数
export type SubscriberQuery = PaginationQuery;

// 订阅者列表数据
export interface SubscriberListData {
  list: Subscriber[];
  total: number;
  page: number;
  page_size: number;
}
