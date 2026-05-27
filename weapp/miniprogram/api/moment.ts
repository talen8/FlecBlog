/**
 * 动态相关 API 接口
 * 对接后端 /moments 接口
 */

import type {
  MomentListItem,
  PageQuery,
  PageResult,
} from '../types';
import { get } from '../utils/request';

/**
 * 获取动态列表
 * @param params - 查询参数
 * @returns 分页动态列表
 */
export function getMoments(params: PageQuery = {}): Promise<PageResult<MomentListItem>> {
  return get<PageResult<MomentListItem>>('/moments', params);
}
