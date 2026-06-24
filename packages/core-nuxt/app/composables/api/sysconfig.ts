import type { SettingGroupType } from '../../../types/sysconfig';
import { createApi } from './createApi';

const settingApi = createApi<Record<string, string>>('/settings');

/** 获取指定分组的配置 */
export const getSettingGroup = async (group: SettingGroupType) => {
  return settingApi.get<Record<string, string>>(`/${group}`);
};
