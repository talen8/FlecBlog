-- 添加 MCP 管理员工具开关；默认关闭。
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('mcp_admin_tools_enabled', 'false', 'ai', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT ("group", key) DO NOTHING;
