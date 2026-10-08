export type ThemeConfig = Record<string, unknown>;

export interface ThemeConfigResponse {
  config?: ThemeConfig;
  menus?: Record<string, ThemeMenuItem[]>;
  pages?: Record<string, ThemePageItem>;
}

export interface ThemeMenuItem {
  id: number;
  title: string;
  url: string;
  icon: string;
  sort: number;
  is_enabled: boolean;
  children?: ThemeMenuItem[];
}

export interface SchemaField {
  type?: string;
  title?: string;
  description?: string;
  default?: unknown;
  enum?: Array<{ label: string; value: string } | string>;
  format?: string;
  placeholder?: string;
  width?: number;
  height?: number;
  min?: number;
  max?: number;
  'x-item-fields'?: Array<string | Record<string, unknown>>;
}

export interface SchemaGroup {
  name: string;
  label: string;
  fields: Record<string, SchemaField>;
}

export interface MenuSlot {
  label?: string;
  maxDepth?: number;
  defaults?: Partial<ThemeMenuItem>[];
}

export interface ThemePageItem {
  title?: string;
  description?: string;
}

export interface ThemeSchema {
  $menus?: Record<string, MenuSlot>;
  $pages?: string[];
}
