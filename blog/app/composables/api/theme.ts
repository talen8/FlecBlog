import type { ThemeConfigResponse } from '../../../types/theme';
import { createApi } from './createApi';

const themeApi = createApi<ThemeConfigResponse>('/themes');

/** 获取主题配置与菜单 */
export const getThemeConfig = async () => {
  return themeApi.get('');
};
