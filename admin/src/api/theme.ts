import request from '@/utils/request';
import type { ThemeConfig, ThemeConfigResponse, ThemeMenuItem, ThemePageItem } from '@/types/theme';

/**
 * 获取主题配置
 * @returns Promise<ThemeConfigResponse>
 */
export const getThemeConfig = (): Promise<ThemeConfigResponse> => {
  return request.get('/admin/themes');
};

/**
 * 更新主题配置
 * @param config 主题配置数据
 * @returns Promise<ThemeConfig>
 */
export const updateThemeConfig = (config: ThemeConfig): Promise<ThemeConfig> => {
  return request.put('/admin/themes/config', { config });
};

/**
 * 更新主题菜单
 * @param menus 菜单数据
 * @returns Promise<Record<string, ThemeMenuItem[]>>
 */
export const updateThemeMenus = (
  menus: Record<string, ThemeMenuItem[]>
): Promise<Record<string, ThemeMenuItem[]>> => {
  return request.put('/admin/themes/menus', { menus });
};

/**
 * 更新主题页面
 * @param pages 页面数据
 * @returns Promise<Record<string, ThemePageItem>>
 */
export const updateThemePages = (
  pages: Record<string, ThemePageItem>
): Promise<Record<string, ThemePageItem>> => {
  return request.put('/admin/themes/pages', { pages });
};
