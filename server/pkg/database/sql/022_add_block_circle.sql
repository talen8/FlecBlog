-- 友圈功能：屏蔽开关与发布时间索引
ALTER TABLE friends ADD COLUMN IF NOT EXISTS block_circle BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE rss_articles ADD COLUMN IF NOT EXISTS block_circle BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX IF NOT EXISTS idx_rss_articles_published_at ON rss_articles (published_at DESC);
