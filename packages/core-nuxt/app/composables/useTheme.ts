import type { ThemeMenuItem } from '../../types/theme';
import { getActiveTheme } from './api/theme';

/**
 * 获取当前激活主题及其菜单（支持 SSR）
 * @returns themeConfig - 主题配置对象，getMenus(type) - 按类型获取菜单项列表
 */
export function useTheme() {
  const { data } = useAsyncData('theme-fetch', () => getActiveTheme());

  const allMenus = computed(() => data.value?.menus ?? {});

  const getMenus = (type: string): ThemeMenuItem[] =>
    (allMenus.value[type] as ThemeMenuItem[]) ?? [];

  return {
    themeConfig: computed(() => data.value?.config ?? {}),
    getMenus,
  };
}
