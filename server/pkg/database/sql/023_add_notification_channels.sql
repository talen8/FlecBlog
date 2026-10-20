-- 通知渠道配置：清旧键 + 新增各渠道凭证

-- 1. 清理应用机器人时代的旧配置
DELETE FROM settings
WHERE "group" = 'notification'
  AND key IN ('feishu_app_id', 'feishu_secret', 'feishu_chat_id');

-- 2. 新增各渠道凭证
INSERT INTO settings (key, value, "group", is_public) VALUES
('feishu_webhook_url', '', 'notification', FALSE),
('feishu_secret', '', 'notification', FALSE),
('telegram_bot_token', '', 'notification', FALSE),
('telegram_chat_id', '', 'notification', FALSE)
ON CONFLICT ("group", key) DO NOTHING;

-- 3. worker_proxy 由 oauth 分组迁到 basic 分组
UPDATE settings
SET "group" = 'basic',
    updated_at = CURRENT_TIMESTAMP
WHERE "group" = 'oauth'
  AND key = 'worker_proxy';

-- 4. 删除 users 表的飞书 OpenID 列
DROP INDEX IF EXISTS idx_users_feishu_open_id;
ALTER TABLE users DROP COLUMN IF EXISTS feishu_open_id;
