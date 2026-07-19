import request from '@/utils/request';
import type { SettingGroupType } from '@/types/sysconfig';

/**
 * 重置MCP Secret响应
 */
export interface ResetMCPSecretResponse {
  secret: string;
}

/**
 * 获取指定分组的配置
 * @param group 分组类型
 * @returns Promise<Record<string, string>>
 */
export const getSettingGroup = (group: SettingGroupType): Promise<Record<string, string>> => {
  return request.get(`/admin/settings/${group}`);
};

/**
 * 更新指定分组的配置
 * @param group 分组类型
 * @param data 配置数据
 * @returns Promise<void>
 */
export const updateSettingGroup = (group: SettingGroupType, data: Record<string, string>) => {
  return request.patch(`/admin/settings/${group}`, data);
};

/**
 * 重置 MCP Secret
 * @returns Promise<ResetMCPSecretResponse>
 */
export const resetMCPSecret = (): Promise<ResetMCPSecretResponse> => {
  return request.put('/admin/settings/ai/mcp-secret/reset');
};
