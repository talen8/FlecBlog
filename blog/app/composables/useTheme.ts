import type { ThemeMenuItem, ThemePageItem } from '../../types/theme';
import { getThemeConfig } from './api/theme';

/**
 * 获取主题配置、菜单与页面文案（支持 SSR）
 * @returns themeConfig - 主题配置对象
 * @returns getMenus(type) - 按类型获取菜单项列表
 * @returns getPage(path) - 按路径获取页面项
 */
export function useTheme() {
  const { data } = useAsyncData('theme-fetch', () => getThemeConfig());

  const allMenus = computed(() => data.value?.menus ?? {});

  const allPages = computed(() => data.value?.pages ?? {});

  const themeConfig = computed(() => data.value?.config ?? {});

  const getMenus = (type: string): ThemeMenuItem[] =>
    (allMenus.value[type] as ThemeMenuItem[]) ?? [];

  const getPage = (path: string): ThemePageItem | undefined => allPages.value[path];

  return {
    themeConfig,
    getMenus,
    getPage,
  };
}
