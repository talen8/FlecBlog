-- users 表添加 OIDC ID 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS oidc_id VARCHAR(100) DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_users_oidc_id ON users (oidc_id) WHERE oidc_id != '' AND deleted_at IS NULL;

-- OIDC 配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('oauth.oidc.enabled', 'false', 'oauth', TRUE),
('oauth.oidc.issuer_url', '', 'oauth', FALSE),
('oauth.oidc.client_id', '', 'oauth', FALSE),
('oauth.oidc.client_secret', '', 'oauth', FALSE),
('oauth.oidc.redirect_url', '', 'oauth', FALSE)
ON CONFLICT (key) DO NOTHING;
