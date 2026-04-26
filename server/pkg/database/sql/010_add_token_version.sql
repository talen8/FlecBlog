-- 添加 token_version 字段，用于撤销所有 token
ALTER TABLE users ADD COLUMN IF NOT EXISTS token_version INT DEFAULT 0 NOT NULL;
