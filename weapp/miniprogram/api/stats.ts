/**
 * 统计相关 API 接口
 * 对接后端公开统计接口
 */

import type { SiteStats } from '../types';
import { get } from '../utils/request';

/**
 * 获取站点统计信息（公开接口）
 * @returns 站点统计
 */
export function getSiteStats(): Promise<SiteStats> {
  return get<SiteStats>('/stats/site');
}
