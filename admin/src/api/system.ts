import request from '@/utils/request';
import type { SystemStatic, SystemDynamic } from '@/types/system';

/**
 * 获取系统静态信息
 * @returns Promise<SystemStatic>
 */
export function getSystemStatic(): Promise<SystemStatic> {
  return request.get('/admin/system/static');
}

/**
 * 获取系统动态信息
 * @returns Promise<SystemDynamic>
 */
export function getSystemDynamic(): Promise<SystemDynamic> {
  return request.get('/admin/system/dynamic');
}

/**
 * 版本信息
 */
export interface VersionInfo {
  id: number;
  version: string;
  date: string;
  changes: string;
}

/**
 * 检查更新响应
 */
export interface CheckUpdateResponse {
  has_update: boolean;
  current_version: string;
  latest_version: string;
  versions: VersionInfo[];
  last_check_error: string;
}

/**
 * 检查版本更新
 * @returns Promise<CheckUpdateResponse>
 */
export function checkUpdate(): Promise<CheckUpdateResponse> {
  return request.post('/admin/system/check-update');
}
