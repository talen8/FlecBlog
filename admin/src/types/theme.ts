export type ThemeConfig = Record<string, unknown>;

export interface ThemeResponse {
  slug: string;
  name: string;
  version: string;
  author: string;
  description: string;
  license: string;
  repo: string;
  schema?: Record<string, unknown>;
  is_active: boolean;
  config?: ThemeConfig;
  menus?: Record<string, ThemeMenuItem[]>;
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
  title?: string;
  maxDepth?: number;
  defaults?: Partial<ThemeMenuItem>[];
}

export interface ThemeSchema {
  $menus?: Record<string, MenuSlot>;
}

export interface ThemeUpdateCheckResponse {
  has_update: boolean;
  current_version: string;
  latest_version: string;
  release_url: string;
}
