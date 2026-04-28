-- 添加邮件加密方式和发件人邮箱配置项

-- 添加邮件加密方式配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('notification.email_secure', 'ssl', 'notification', FALSE)
ON CONFLICT (key) DO NOTHING;

-- 添加发件人邮箱配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('notification.email_from', '', 'notification', FALSE)
ON CONFLICT (key) DO NOTHING;
