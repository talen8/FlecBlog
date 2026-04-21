import request from '@/utils/request';
import type {
  Feedback,
  FeedbackListData,
  FeedbackListQuery,
  FeedbackUpdateRequest,
} from '@/types/feedback';

/**
 * 获取反馈列表
 * @param params 查询参数
 * @returns Promise<FeedbackListData>
 */
export function getFeedbackList(params: FeedbackListQuery): Promise<FeedbackListData> {
  return request.get('/admin/feedback', { params });
}

/**
 * 获取反馈详情
 * @param id 反馈ID
 * @returns Promise<Feedback>
 */
export function getFeedbackDetail(id: number): Promise<Feedback> {
  return request.get(`/admin/feedback/${id}`);
}

/**
 * 更新反馈
 * @param id 反馈ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export function updateFeedback(id: number, data: FeedbackUpdateRequest): Promise<void> {
  return request.put(`/admin/feedback/${id}`, data);
}

/**
 * 删除反馈
 * @param id 反馈ID
 * @returns Promise<void>
 */
export function deleteFeedback(id: number): Promise<void> {
  return request.delete(`/admin/feedback/${id}`);
}
