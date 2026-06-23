-- 主题系统迁移：settings/menu 收敛到 default 主题实例

CREATE TABLE IF NOT EXISTS theme_instances (
    slug TEXT PRIMARY KEY,
    name TEXT DEFAULT '',
    version TEXT DEFAULT '',
    author TEXT DEFAULT '',
    description TEXT DEFAULT '',
    license TEXT DEFAULT '',
    repo TEXT DEFAULT '',
    "schema" JSON,
    is_active BOOLEAN DEFAULT FALSE,
    config JSON,
    menus JSON DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_theme_instances_active
    ON theme_instances (is_active)
    WHERE is_active = TRUE;

INSERT INTO theme_instances (slug, name, "schema", is_active, config, menus)
SELECT 'default', '默认主题', '{}'::json, TRUE, '{}'::json, '{}'::json
WHERE NOT EXISTS (SELECT 1 FROM theme_instances);

WITH s AS (
    SELECT key, value FROM settings
),
config AS (
    SELECT jsonb_build_object(
        'author_email', COALESCE((SELECT value FROM s WHERE key = 'basic.author_email'), ''),
        'author_desc', COALESCE((SELECT value FROM s WHERE key = 'basic.author_desc'), ''),
        'author_photo', COALESCE((SELECT value FROM s WHERE key = 'basic.author_photo'), ''),
        'about_creation', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.about_creation'), '[]'::jsonb),
        'about_describe', COALESCE((SELECT value FROM s WHERE key = 'blog.about_describe'), ''),
        'about_describe_tips', COALESCE((SELECT value FROM s WHERE key = 'blog.about_describe_tips'), ''),
        'about_exhibition', COALESCE((SELECT value FROM s WHERE key = 'blog.about_exhibition'), ''),
        'about_motto_main', COALESCE((
            SELECT jsonb_agg(jsonb_build_object('value', item))
            FROM jsonb_array_elements(COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.about_motto_main'), '[]'::jsonb)) item
        ), '[]'::jsonb),
        'about_motto_sub', COALESCE((SELECT value FROM s WHERE key = 'blog.about_motto_sub'), ''),
        'about_personality', COALESCE((SELECT value FROM s WHERE key = 'blog.about_personality'), ''),
        'about_profile', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.about_profile'), '[]'::jsonb),
        'about_socialize', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.about_socialize'), '[]'::jsonb),
        'about_story', COALESCE((SELECT value FROM s WHERE key = 'blog.about_story'), ''),
        'about_unions', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.about_unions'), '[]'::jsonb),
        'about_versions', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.about_versions'), '[]'::jsonb),
        'slogan', COALESCE((SELECT value FROM s WHERE key = 'blog.slogan'), ''),
        'background_image', COALESCE((SELECT value FROM s WHERE key = 'blog.background_image'), ''),
        'screenshot', COALESCE((SELECT value FROM s WHERE key = 'blog.screenshot'), ''),
        'announcement', COALESCE((SELECT value FROM s WHERE key = 'blog.announcement'), ''),
        'typing_texts', COALESCE((
            SELECT jsonb_agg(jsonb_build_object('value', item))
            FROM jsonb_array_elements(COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.typing_texts'), '[]'::jsonb)) item
        ), '[]'::jsonb),
        'font', COALESCE((SELECT value FROM s WHERE key = 'blog.font'), ''),
        'sidebar_social', COALESCE((
            SELECT jsonb_agg(CASE
                WHEN item->>'icon' LIKE 'ri-%' THEN item
                ELSE jsonb_set(item, '{icon}', to_jsonb('ri-' || (item->>'icon')))
            END)
            FROM jsonb_array_elements(COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.sidebar_social'), '[]'::jsonb)) item
        ), '[]'::jsonb),
        'footer_social', COALESCE((
            SELECT jsonb_agg(CASE
                WHEN item->>'icon' LIKE 'ri-%' THEN item
                ELSE jsonb_set(item, '{icon}', to_jsonb('ri-' || (item->>'icon')))
            END)
            FROM jsonb_array_elements(COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.footer_social'), '[]'::jsonb)) item
        ), '[]'::jsonb),
        'footer_links', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.footer_links'), '[]'::jsonb),
        'home_layout', COALESCE((SELECT value FROM s WHERE key = 'blog.home_layout'), ''),
        'moments_size', to_jsonb(COALESCE((SELECT value FROM s WHERE key = 'blog.moments_size'), '30')::int),
        'theme_light_start', regexp_replace(COALESCE((SELECT value FROM s WHERE key = 'blog.theme_light_start'), ''), '^([0-9]{2}:[0-9]{2})$', '\1:00'),
        'theme_dark_start', regexp_replace(COALESCE((SELECT value FROM s WHERE key = 'blog.theme_dark_start'), ''), '^([0-9]{2}:[0-9]{2})$', '\1:00'),
        'donation_methods', COALESCE((SELECT value::jsonb FROM s WHERE key = 'blog.donation_methods'), '[]'::jsonb),
        'message_content', COALESCE((SELECT value FROM s WHERE key = 'blog.message_content'), ''),
        'wechat_qrcode', COALESCE((SELECT value FROM s WHERE key = 'blog.wechat_qrcode'), ''),
        'wechat_name', COALESCE((SELECT value FROM s WHERE key = 'blog.wechat_name'), ''),
        'home_url', COALESCE((SELECT value FROM s WHERE key = 'basic.home_url'), ''),
        'established', COALESCE((SELECT value FROM s WHERE key = 'blog.established'), '')
    ) AS data
)
UPDATE theme_instances
SET config = (config.data || COALESCE(theme_instances.config::jsonb, '{}'::jsonb))::json,
    updated_at = CURRENT_TIMESTAMP
FROM config
WHERE slug = 'default';

WITH menu_tree AS (
    SELECT jsonb_object_agg(type, items) AS data
    FROM (
        SELECT p.type,
               jsonb_agg(
                   jsonb_build_object(
                       'id', p.id,
                       'title', p.title,
                       'url', COALESCE(p.url, ''),
                       'icon', COALESCE(p.icon, ''),
                       'sort', p.sort,
                       'is_enabled', p.is_enabled,
                       'children', COALESCE(children.items, '[]'::jsonb)
                   )
                   ORDER BY p.sort, p.id
               ) AS items
        FROM menus p
        LEFT JOIN LATERAL (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'id', c.id,
                    'title', c.title,
                    'url', COALESCE(c.url, ''),
                    'icon', COALESCE(c.icon, ''),
                    'sort', c.sort,
                    'is_enabled', c.is_enabled,
                    'children', '[]'::jsonb
                )
                ORDER BY c.sort, c.id
            ) AS items
            FROM menus c
            WHERE c.parent_id = p.id
        ) children ON TRUE
        WHERE p.parent_id IS NULL
        GROUP BY p.type
    ) grouped
)
UPDATE theme_instances
SET menus = (COALESCE(menu_tree.data, '{}'::jsonb) || COALESCE(theme_instances.menus::jsonb, '{}'::jsonb))::json,
    updated_at = CURRENT_TIMESTAMP
FROM menu_tree
WHERE slug = 'default';

WITH file_urls AS (
    SELECT value AS url FROM settings
    WHERE key IN (
        'basic.author_photo',
        'blog.background_image',
        'blog.screenshot',
        'blog.about_exhibition',
        'blog.wechat_qrcode'
    )
    UNION
    SELECT item->>'qrcode'
    FROM settings, jsonb_array_elements(value::jsonb) item
    WHERE key = 'blog.donation_methods'
    UNION
    SELECT icon FROM menus
)
UPDATE files
SET upload_type = 'default',
    status = 1,
    updated_at = CURRENT_TIMESTAMP
WHERE file_url IN (SELECT url FROM file_urls WHERE COALESCE(url, '') <> '');

DELETE FROM settings
WHERE key IN (
    'basic.author_email',
    'basic.author_desc',
    'basic.author_photo',
    'basic.home_url',
    'blog.established',
    'blog.about_creation',
    'blog.about_describe',
    'blog.about_describe_tips',
    'blog.about_exhibition',
    'blog.about_motto_main',
    'blog.about_motto_sub',
    'blog.about_personality',
    'blog.about_profile',
    'blog.about_socialize',
    'blog.about_story',
    'blog.about_unions',
    'blog.about_versions',
    'blog.slogan',
    'blog.background_image',
    'blog.screenshot',
    'blog.announcement',
    'blog.typing_texts',
    'blog.font',
    'blog.sidebar_social',
    'blog.footer_social',
    'blog.footer_links',
    'blog.home_layout',
    'blog.moments_size',
    'blog.theme_light_start',
    'blog.theme_dark_start',
    'blog.donation_methods',
    'blog.message_content',
    'blog.wechat_qrcode',
    'blog.wechat_name'
);

ALTER TABLE settings DROP CONSTRAINT IF EXISTS settings_key_key;

DELETE FROM settings blog_setting
USING settings basic_setting
WHERE blog_setting."group" = 'blog'
    AND basic_setting."group" = 'basic'
    AND replace(blog_setting.key, 'blog.', '') = replace(basic_setting.key, 'basic.', '');

UPDATE settings
SET key = regexp_replace(key, '^(basic|blog)\.', ''),
    "group" = 'basic',
    updated_at = CURRENT_TIMESTAMP
WHERE "group" IN ('basic', 'blog');

-- notification 分组去前缀
UPDATE settings
SET key = regexp_replace(key, '^notification\.', ''),
    "group" = 'notification',
    updated_at = CURRENT_TIMESTAMP
WHERE "group" = 'notification';

-- upload 分组去前缀
UPDATE settings
SET key = regexp_replace(key, '^upload\.', ''),
    "group" = 'upload',
    updated_at = CURRENT_TIMESTAMP
WHERE "group" = 'upload';

-- ai 分组去前缀
UPDATE settings
SET key = regexp_replace(key, '^ai\.', ''),
    "group" = 'ai',
    updated_at = CURRENT_TIMESTAMP
WHERE "group" = 'ai';

-- oauth 分组去前缀
UPDATE settings
SET key = regexp_replace(key, '^oauth\.', ''),
    "group" = 'oauth',
    updated_at = CURRENT_TIMESTAMP
WHERE "group" = 'oauth';

CREATE UNIQUE INDEX IF NOT EXISTS idx_settings_group_key ON settings ("group", key);

DROP TABLE menus;

-- 移除 moment 作为独立的评论目标类型（实际使用 page + target_key='moment'）
ALTER TABLE comments DROP CONSTRAINT IF EXISTS chk_comments_target_type;
ALTER TABLE comments ADD CONSTRAINT chk_comments_target_type CHECK (target_type IN ('article', 'page'));

-- 统一文件用途名称：网站Favicon → 博客图标
UPDATE files SET upload_type = '博客图标', updated_at = CURRENT_TIMESTAMP WHERE upload_type = '网站Favicon';
