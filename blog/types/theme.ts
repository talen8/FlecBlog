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

/** 主题页面项 */
export interface ThemePageItem {
  title?: string;
  description?: string;
}

/** 主题配置响应 */
export interface ThemeConfigResponse {
  config: Record<string, string>;
  menus: Record<string, ThemeMenuItem[]>;
  pages?: Record<string, ThemePageItem>;
}

/** 版权协议标识 */
export type LicenseKey =
  'cc-by' | 'cc-by-sa' | 'cc-by-nd' | 'cc-by-nc' | 'cc-by-nc-sa' | 'cc-by-nc-nd' | 'cc0';

/** 协议条款标识 */
export type LicenseClauseCode = 'by' | 'nc' | 'nd' | 'sa' | 'cc0';

/** 版权协议完整信息 */
export interface CopyrightLicense {
  short: string;
  name: string;
  url: string;
  clauses: LicenseClauseCode[];
}
