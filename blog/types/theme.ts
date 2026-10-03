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

/** 主题配置响应 */
export interface ThemeConfigResponse {
  config: Record<string, string>;
  menus: Record<string, ThemeMenuItem[]>;
}
