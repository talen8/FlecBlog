export interface SystemStatic {
  cpu_core: number;
  cpu_model: string;
  cpu_arch: string;
  hostname: string;
  os: string;
  server_ip: string;
  timezone: string;
  db_type: string;
  memory_total: number;
  swap_total: number;
  disk_total: number;
  db_tables: number;
  storage_status: string;
  email_status: string;
  feishu_status: string;
  panel_status: string;
  app_version: string;
}

export interface SystemDynamic {
  cpu_usage: number;
  load_1: number;
  load_5: number;
  load_15: number;
  memory_used: number;
  memory_available: number;
  swap_used: number;
  host_uptime: number;
  disk_used: number;
  disk_free: number;
  db_status: string;
  db_size: number;
  db_conn_count: number;
  version_latest_version: string;
  version_last_check_error: string;
}
