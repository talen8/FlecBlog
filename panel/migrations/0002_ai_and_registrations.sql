-- AI 提供商表
CREATE TABLE IF NOT EXISTS ai_providers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  base_url TEXT NOT NULL,
  api_key TEXT NOT NULL,
  model TEXT NOT NULL,
  priority_immediate INTEGER NOT NULL DEFAULT 999,
  priority_deferred INTEGER NOT NULL DEFAULT 999,
  call_count INTEGER NOT NULL DEFAULT 0,
  fail_count INTEGER NOT NULL DEFAULT 0,
  total_latency_ms INTEGER NOT NULL DEFAULT 0,
  total_prompt_tokens INTEGER NOT NULL DEFAULT 0,
  total_completion_tokens INTEGER NOT NULL DEFAULT 0,
  last_call_at TEXT,
  last_error_message TEXT,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- 部署注册表
CREATE TABLE IF NOT EXISTS registrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  client_key TEXT NOT NULL UNIQUE,
  site_url TEXT NOT NULL UNIQUE,
  site_url_history TEXT NOT NULL DEFAULT '[]',
  version TEXT NOT NULL DEFAULT '',
  first_seen_at TEXT NOT NULL DEFAULT (datetime('now')),
  last_seen_at TEXT NOT NULL DEFAULT (datetime('now')),
  banned INTEGER NOT NULL DEFAULT 0,
  call_counts TEXT NOT NULL DEFAULT '{}',
  remark TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_registrations_client_key ON registrations(client_key);
CREATE INDEX IF NOT EXISTS idx_registrations_site_url ON registrations(site_url);
