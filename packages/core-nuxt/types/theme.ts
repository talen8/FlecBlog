/** 主题菜单项 */
export interface ThemeMenuItem {
  id: number;
  title: string;
  url: string;
  icon: string;
  sort: number;
  is_enabled?: boolean;
  children?: ThemeMenuItem[];
}

/** 激活主题响应 */
export interface ActiveTheme {
  slug: string;
  name: string;
  version: string;
  author: string;
  description: string;
  license: string;
  repo: string;
  is_active: boolean;
  config: Record<string, string>;
  menus: Record<string, ThemeMenuItem[]>;
}
