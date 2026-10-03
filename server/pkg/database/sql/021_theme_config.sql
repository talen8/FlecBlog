-- 主题系统简化迁移：theme_instances 收敛为 theme_configs 键值配置表

CREATE TABLE IF NOT EXISTS theme_configs (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 迁移激活主题行的 config/menus（无激活行则取最早一行）
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM information_schema.tables
        WHERE table_schema = 'public'
        AND table_name = 'theme_instances'
    ) THEN
        -- 展开 config 为 key-value 行，menus 单独一行；字符串值去引号，其余保留 JSON 文本
        WITH chosen AS (
            SELECT *
            FROM theme_instances
            ORDER BY is_active DESC, created_at ASC
            LIMIT 1
        ),
        kv AS (
            SELECT e.key,
                   CASE WHEN jsonb_typeof(e.value) = 'string'
                        THEN e.value ->> 0
                        ELSE e.value::text
                   END AS value
            FROM chosen,
                 LATERAL jsonb_each(COALESCE(chosen.config, '{}'::json)::jsonb) e
            UNION ALL
            SELECT 'menus', COALESCE(chosen.menus::text, '{}')
            FROM chosen
        )
        INSERT INTO theme_configs (key, value)
        SELECT key, value
        FROM kv
        ON CONFLICT (key) DO NOTHING;

        -- 还原文件用途：菜单图标 → 「菜单图标」
        WITH chosen AS (
            SELECT slug, menus
            FROM theme_instances
            ORDER BY is_active DESC, created_at ASC
            LIMIT 1
        ),
        menu_icons AS (
            SELECT DISTINCT i.value ->> 'icon' AS url
            FROM chosen,
                 LATERAL jsonb_each(COALESCE(chosen.menus, '{}'::json)::jsonb) g,
                 LATERAL jsonb_array_elements(g.value) i
            UNION
            SELECT DISTINCT c.value ->> 'icon' AS url
            FROM chosen,
                 LATERAL jsonb_each(COALESCE(chosen.menus, '{}'::json)::jsonb) g,
                 LATERAL jsonb_array_elements(g.value) i,
                 LATERAL jsonb_array_elements(COALESCE(i.value -> 'children', '[]'::jsonb)) c
        )
        UPDATE files f
        SET upload_type = '菜单图标',
            updated_at = CURRENT_TIMESTAMP
        FROM chosen, menu_icons m
        WHERE f.upload_type = chosen.slug
          AND f.file_url = m.url
          AND COALESCE(m.url, '') <> '';

        -- 还原文件用途：配置图片 → 「主题图片」
        WITH chosen AS (
            SELECT slug, config
            FROM theme_instances
            ORDER BY is_active DESC, created_at ASC
            LIMIT 1
        ),
        config_images AS (
            SELECT DISTINCT e.value #>> '{}' AS url
            FROM chosen,
                 LATERAL jsonb_each(COALESCE(chosen.config, '{}'::json)::jsonb) e
            WHERE e.key IN ('author_photo', 'background_image', 'screenshot',
                            'about_exhibition', 'wechat_qrcode')
              AND jsonb_typeof(e.value) = 'string'
            UNION
            SELECT DISTINCT d.value ->> 'qrcode' AS url
            FROM chosen,
                 LATERAL jsonb_array_elements(
                     COALESCE(chosen.config -> 'donation_methods', '[]'::json)::jsonb
                 ) d
        )
        UPDATE files f
        SET upload_type = '主题图片',
            updated_at = CURRENT_TIMESTAMP
        FROM chosen, config_images ci
        WHERE f.upload_type = chosen.slug
          AND f.file_url = ci.url
          AND COALESCE(ci.url, '') <> '';

        -- 统一文件用途名称：default → 主题图片
        UPDATE files
        SET upload_type = '主题图片',
            updated_at = CURRENT_TIMESTAMP
        WHERE upload_type = 'default';

        DROP TABLE theme_instances;
    END IF;
END $$;

-- 站点基础信息回迁 settings（basic 组）
INSERT INTO settings ("group", key, value, is_public)
SELECT 'basic', target.key, config.value, TRUE
FROM theme_configs config
JOIN (VALUES ('author_email'), ('author_desc'), ('home_url'), ('established')) AS target(key)
  ON config.key = target.key
WHERE COALESCE(config.value, '') <> ''
ON CONFLICT ("group", key) DO NOTHING;

DELETE FROM theme_configs
WHERE key IN ('author_email', 'author_desc', 'home_url', 'established');

-- 预置站点信息配置项：settings 写入非 upsert，缺行会导致保存报错
INSERT INTO settings (key, value, "group", is_public)
VALUES
('author_email', '', 'basic', TRUE),
('author_desc', '', 'basic', TRUE),
('home_url', '', 'basic', TRUE),
('established', '', 'basic', TRUE)
ON CONFLICT ("group", key) DO NOTHING;
