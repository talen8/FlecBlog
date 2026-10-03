import type { ThemeMenuItem } from '../../types/theme';
import { getThemeConfig } from './api/theme';

/**
 * 获取主题配置及其菜单（支持 SSR）
 * @returns themeConfig - 主题配置对象
 * @returns getMenus(type) - 按类型获取菜单项列表
 */
export function useTheme() {
  const { data } = useAsyncData('theme-fetch', () => getThemeConfig());

  const allMenus = computed(() => data.value?.menus ?? {});

  const themeConfig = computed(() => data.value?.config ?? {});

  const getMenus = (type: string): ThemeMenuItem[] =>
    (allMenus.value[type] as ThemeMenuItem[]) ?? [];

  return {
    themeConfig,
    getMenus,
  };
}
