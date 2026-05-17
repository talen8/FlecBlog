-- 添加赞赏方式配置项
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('blog.donation_methods', '[]', 'blog', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
