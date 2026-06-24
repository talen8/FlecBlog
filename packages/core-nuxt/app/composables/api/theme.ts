import type { ActiveTheme } from '../../../types/theme';
import { createApi } from './createApi';

const themeApi = createApi<ActiveTheme>('/themes');

/** 获取当前激活主题 */
export const getActiveTheme = async () => {
  return themeApi.get('');
};
