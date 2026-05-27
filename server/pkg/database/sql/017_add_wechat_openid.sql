-- 微信小程序登录支持
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_open_id VARCHAR(100) DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_users_wechat_open_id ON users(wechat_open_id) WHERE wechat_open_id != '';

-- 微信小程序配置项
INSERT INTO settings (key, value, "group", is_public, created_at, updated_at) VALUES
    ('oauth.wechat.enabled', 'false', 'oauth', true, NOW(), NOW()),
    ('oauth.wechat.appid', '', 'oauth', false, NOW(), NOW()),
    ('oauth.wechat.secret', '', 'oauth', false, NOW(), NOW())
ON CONFLICT (key) DO NOTHING;
