import request from '@/utils/request';
import type {
  ThemeConfig,
  ThemeResponse,
  ThemeMenuItem,
  ThemeUpdateCheckResponse,
} from '@/types/theme';

export const getThemes = (): Promise<ThemeResponse[]> => {
  return request.get('/admin/themes');
};

export const getTheme = (slug: string): Promise<ThemeResponse> => {
  return request.get(`/admin/themes/${slug}`);
};

export const updateThemeConfig = (slug: string, config: ThemeConfig): Promise<ThemeConfig> => {
  return request.put(`/admin/themes/${slug}/config`, { config });
};

export const updateThemeMenus = (
  slug: string,
  menus: Record<string, ThemeMenuItem[]>
): Promise<Record<string, ThemeMenuItem[]>> => {
  return request.put(`/admin/themes/${slug}/menus`, { menus });
};

export const checkThemeUpdate = (slug: string): Promise<ThemeUpdateCheckResponse> => {
  return request.post(`/admin/themes/${slug}/check`);
};
