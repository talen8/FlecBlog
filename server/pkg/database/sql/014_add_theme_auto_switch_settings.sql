-- 添加主题自动切换配置项
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('blog.theme_light_start', '06:00', 'blog', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;

INSERT INTO settings (key, value, "group", is_public, created_at, updated_at)
VALUES ('blog.theme_dark_start', '18:00', 'blog', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
