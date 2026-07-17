-- 添加 MCP 静态认证用户管理操作身份配置项；空值表示自动选择唯一启用的超级管理员。
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('mcp_static_operator_user_id', '', 'ai', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT ("group", key) DO NOTHING;

-- 添加 MCP 管理员工具开关；默认关闭。
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('mcp_admin_tools_enabled', 'false', 'ai', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT ("group", key) DO NOTHING;
