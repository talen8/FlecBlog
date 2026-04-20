/**
 * 菜单类型
 */
export type MenuType = 'navigation' | 'footer' | 'aggregate';

/**
 * 菜单项接口
 */
export interface Menu {
  id: number;
  type: MenuType;
  parent_id: number | null;
  title: string;
  url: string;
  icon: string;
  sort: number;
  children?: Menu[];
}

/**
 * 菜单响应接口
 */
export interface MenusResponse {
  code: number;
  message: string;
  data: Menu[];
}
