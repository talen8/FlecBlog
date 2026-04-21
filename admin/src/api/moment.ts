import request from '@/utils/request';
import type {
  Moment,
  MomentListData,
  CreateMomentRequest,
  UpdateMomentRequest,
  MomentListQuery,
} from '@/types/moment';

/**
 * 获取动态列表
 * @param params 查询参数
 * @returns Promise<MomentListData>
 */
export function getMoments(params: MomentListQuery): Promise<MomentListData> {
  return request.get('/admin/moments', { params });
}

/**
 * 获取动态详情
 * @param id 动态ID
 * @returns Promise<Moment>
 */
export function getMoment(id: number): Promise<Moment> {
  return request.get(`/admin/moments/${id}`);
}

/**
 * 创建动态
 * @param data 动态数据
 * @returns Promise<Moment>
 */
export function createMoment(data: CreateMomentRequest): Promise<Moment> {
  return request.post('/admin/moments', data);
}

/**
 * 更新动态
 * @param id 动态ID
 * @param data 动态数据
 * @returns Promise<Moment>
 */
export function updateMoment(id: number, data: UpdateMomentRequest): Promise<Moment> {
  return request.put(`/admin/moments/${id}`, data);
}

/**
 * 删除动态
 * @param id 动态ID
 * @returns Promise<void>
 */
export function deleteMoment(id: number): Promise<void> {
  return request.delete(`/admin/moments/${id}`);
}
