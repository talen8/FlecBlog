/**
 * 类型定义
 * 包含所有业务实体类型和环境类型
 */

/** 版本信息 */
export interface Version {
  id?: number;
  version: string;
  date: string;
  changes: string;
  enabled: boolean;
  created_at?: string;
}

/** 公告信息 */
export interface Announcement {
  id: number;
  title: string;
  content: string;
  link?: string;
  created_at: string;
}

/** 环境变量 */
export interface Env {
  DB: D1Database;
  ASSETS?: Fetcher;
}

/** 设置项 */
export interface Setting {
  key: string;
  value: string;
  updated_at?: string;
}
