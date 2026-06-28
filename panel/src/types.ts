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
  PANEL_API_KEY?: string;
}

/** 部署注册 */
export interface Registration {
  id: number;
  client_key: string;
  site_url: string;
  site_url_history: string;
  version: string;
  first_seen_at: string;
  last_seen_at: string;
  banned: number;
  call_counts: string;
  remark: string;
}

/** 编辑部署请求体 */
export interface RegistrationUpdate {
  banned?: number;
  remark?: string;
}

/** 设置项 */
export interface Setting {
  key: string;
  value: string;
  updated_at?: string;
}

/** AI 提供商 */
export interface AIProvider {
  id: number;
  name: string;
  base_url: string;
  api_key: string;
  model: string;
  priority_immediate: number;
  priority_deferred: number;
  call_count: number;
  fail_count: number;
  total_latency_ms: number;
  total_prompt_tokens: number;
  total_completion_tokens: number;
  last_call_at: string | null;
  last_error_message: string | null;
  created_at: string;
  updated_at: string;
}
