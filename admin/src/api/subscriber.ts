import request from '@/utils/request';
import type { SubscriberQuery, SubscriberListData } from '@/types/subscriber';

/**
 * 获取订阅者列表
 * @param params 查询参数
 * @returns Promise<SubscriberListData>
 */
export function getSubscribers(params: SubscriberQuery): Promise<SubscriberListData> {
  return request.get('/admin/subscribers', { params });
}

/**
 * 删除订阅者
 * @param id 订阅者ID
 * @returns Promise<void>
 */
export function deleteSubscriber(id: number): Promise<void> {
  return request.delete(`/admin/subscribers/${id}`);
}
