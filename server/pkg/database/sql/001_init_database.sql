-- =============================================
-- FlecBLOG - 完整数据库初始化脚本
-- =============================================
-- 数据库：PostgreSQL 12+
-- 编码：UTF-8
-- 作用：创建表结构、索引、触发器、函数和示例数据
-- 最后更新：2025-11-25
-- =============================================


-- =============================================
-- 第一部分：表结构定义
-- =============================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) DEFAULT '',
    has_password BOOLEAN DEFAULT FALSE,
    nickname VARCHAR(50) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',
    website VARCHAR(255) DEFAULT '',
    badge VARCHAR(50) DEFAULT '',
    is_enabled BOOLEAN DEFAULT TRUE,
    role VARCHAR(20) DEFAULT 'user',
    last_login TIMESTAMP,
    github_id VARCHAR(50) DEFAULT '',
    google_id VARCHAR(50) DEFAULT '',
    qq_id VARCHAR(50) DEFAULT '',
    feishu_open_id VARCHAR(50) DEFAULT '',
    microsoft_id VARCHAR(50) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT chk_users_role CHECK (role IN ('super_admin', 'admin', 'user', 'guest'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_date ON users (DATE(created_at)) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_github_id ON users (github_id) WHERE github_id != '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users (google_id) WHERE google_id != '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_qq_id ON users (qq_id) WHERE qq_id != '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_feishu_open_id ON users (feishu_open_id) WHERE feishu_open_id != '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_microsoft_id ON users (microsoft_id) WHERE microsoft_id != '' AND deleted_at IS NULL;


CREATE TABLE IF NOT EXISTS token_blacklists (
    id BIGSERIAL PRIMARY KEY,
    token_hash VARCHAR(64) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_token_blacklists_expires_at ON token_blacklists(expires_at);
CREATE INDEX IF NOT EXISTS idx_token_blacklists_user_id ON token_blacklists(user_id);


-- 邮箱验证码表
CREATE TABLE IF NOT EXISTS verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    failed_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_verifications_failed_count CHECK (failed_count >= 0 AND failed_count <= 10)
);

CREATE INDEX IF NOT EXISTS idx_verifications_email_code ON verifications(email, code) WHERE used = FALSE;
CREATE INDEX IF NOT EXISTS idx_verifications_expires_at ON verifications(expires_at) WHERE used = FALSE;


-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    count INTEGER DEFAULT 0,
    sort INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_categories_count CHECK (count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_categories_sort ON categories (sort DESC);


-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_tags_count CHECK (count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_tags_count ON tags (count DESC);


-- 文章表
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    summary VARCHAR(500) DEFAULT '',
    ai_summary TEXT DEFAULT '',
    cover VARCHAR(255) DEFAULT '',
    location VARCHAR(100) DEFAULT '',
    is_publish BOOLEAN DEFAULT FALSE,
    is_top BOOLEAN DEFAULT FALSE,
    is_essence BOOLEAN DEFAULT FALSE,
    is_outdated BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    category_id BIGINT,
    publish_time TIMESTAMP NULL,
    update_time TIMESTAMP NULL,
    search_vector TSVECTOR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_article_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    CONSTRAINT chk_articles_view_count CHECK (view_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_articles_category_id ON articles (category_id) WHERE category_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_articles_is_publish ON articles (is_publish);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_publish_time ON articles (publish_time DESC) WHERE publish_time IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_articles_publish_top_time ON articles (is_publish, is_top DESC, publish_time DESC) WHERE is_publish = TRUE;
CREATE INDEX IF NOT EXISTS idx_articles_publish_year_month ON articles (EXTRACT(YEAR FROM publish_time), EXTRACT(MONTH FROM publish_time)) WHERE is_publish = TRUE;
CREATE INDEX IF NOT EXISTS idx_articles_hot ON articles (view_count DESC, publish_time DESC) WHERE is_publish = TRUE;


-- 文章标签关联表
CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (article_id, tag_id),
    CONSTRAINT fk_at_article FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    CONSTRAINT fk_at_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_article_tags_tag_id ON article_tags(tag_id);


-- 评论表（支持多态评论和嵌套回复）
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_key VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL,
    parent_id BIGINT,
    root_id BIGINT,
    reply_to BIGINT,
    status INTEGER DEFAULT 1,
    ip VARCHAR(45) DEFAULT '',
    location VARCHAR(100) DEFAULT '',
    browser VARCHAR(50) DEFAULT '',
    os VARCHAR(50) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_comment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_reply_to FOREIGN KEY (reply_to) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_comments_status CHECK (status IN (0, 1)),
    CONSTRAINT chk_comments_target_type CHECK (target_type IN ('article', 'page', 'moment'))
);

CREATE INDEX IF NOT EXISTS idx_comments_target_root ON comments (target_type, target_key, root_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_root_created ON comments (root_id, created_at ASC) WHERE deleted_at IS NULL AND root_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_reply_to ON comments (reply_to, created_at DESC) WHERE deleted_at IS NULL AND reply_to IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_comments_status_created ON comments (status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_latest ON comments (created_at DESC) WHERE deleted_at IS NULL AND status = 1;
CREATE INDEX IF NOT EXISTS idx_comments_created_date ON comments (DATE(created_at)) WHERE deleted_at IS NULL;


-- 文件上传表
CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT DEFAULT 0,
    file_type VARCHAR(100) DEFAULT '',
    upload_type VARCHAR(20) DEFAULT '',
    storage_type VARCHAR(20) DEFAULT 'local',
    user_id BIGINT,
    file_url VARCHAR(500) DEFAULT '',
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_file_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_files_status CHECK (status IN (0, 1)),
    CONSTRAINT chk_files_file_size CHECK (file_size >= 0)
);

CREATE INDEX IF NOT EXISTS idx_files_user ON files (user_id);
CREATE INDEX IF NOT EXISTS idx_files_type ON files (upload_type) WHERE upload_type != '';
CREATE INDEX IF NOT EXISTS idx_files_status ON files (status);
CREATE INDEX IF NOT EXISTS idx_files_user_type ON files (user_id, upload_type);
CREATE INDEX IF NOT EXISTS idx_files_storage_type ON files (storage_type);


-- 访问记录表（PV/UV统计）
CREATE TABLE IF NOT EXISTS visits (
    id BIGSERIAL PRIMARY KEY,
    visitor_id VARCHAR(64) NOT NULL,
    ip VARCHAR(45) NOT NULL,
    page_url VARCHAR(500) DEFAULT '',
    article_id BIGINT,
    user_agent TEXT DEFAULT '',
    location VARCHAR(100) DEFAULT '',
    browser VARCHAR(50) DEFAULT '',
    os VARCHAR(50) DEFAULT '',
    referer VARCHAR(500) DEFAULT '',
    visit_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_visit_article FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_visits_date ON visits (visit_date DESC);
CREATE INDEX IF NOT EXISTS idx_visits_visitor_date ON visits (visitor_id, visit_date);
CREATE INDEX IF NOT EXISTS idx_visits_article_date ON visits (article_id, visit_date) WHERE article_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_visits_recent_online ON visits (created_at DESC, visitor_id);


-- 友链类型表
CREATE TABLE IF NOT EXISTS friend_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    sort INTEGER DEFAULT 0,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_friend_types_sort ON friend_types (sort DESC, id ASC);
CREATE INDEX IF NOT EXISTS idx_friend_types_visible ON friend_types (is_visible) WHERE is_visible = TRUE;

-- 友链表
CREATE TABLE IF NOT EXISTS friends (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',
    screenshot VARCHAR(255) DEFAULT '',
    sort INTEGER DEFAULT 5,
    is_invalid BOOLEAN DEFAULT FALSE,
    is_pending BOOLEAN DEFAULT FALSE,
    type INTEGER REFERENCES friend_types(id) ON DELETE SET NULL,
    rss_url VARCHAR(500) DEFAULT '',
    rss_latime TIMESTAMP,
    accessible INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_friends_type ON friends (type);
CREATE INDEX IF NOT EXISTS idx_friends_sort ON friends (sort DESC, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_friends_pending ON friends (is_pending) WHERE is_pending = TRUE;
CREATE INDEX IF NOT EXISTS idx_friends_status_sort ON friends (is_pending, is_invalid, sort DESC, created_at DESC);


-- RSS文章表
CREATE TABLE IF NOT EXISTS rss_articles (
    id BIGSERIAL PRIMARY KEY,
    friend_id BIGINT NOT NULL,
    title VARCHAR(500) NOT NULL,
    link VARCHAR(1000) NOT NULL,
    published_at TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_rss_article_friend FOREIGN KEY (friend_id) REFERENCES friends(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_rss_articles_friend_id ON rss_articles(friend_id);
CREATE INDEX IF NOT EXISTS idx_rss_articles_published_at ON rss_articles(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_rss_articles_is_read ON rss_articles(is_read);
CREATE UNIQUE INDEX IF NOT EXISTS idx_rss_articles_link ON rss_articles(link);


-- 动态表
CREATE TABLE IF NOT EXISTS moments (
    id BIGSERIAL PRIMARY KEY,
    content JSONB NOT NULL,
    is_publish BOOLEAN DEFAULT TRUE,
    publish_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_moments_is_publish ON moments (is_publish);
CREATE INDEX IF NOT EXISTS idx_moments_publish_time ON moments (publish_time DESC);
CREATE INDEX IF NOT EXISTS idx_moments_created_at ON moments (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_moments_publish_filter ON moments (is_publish, COALESCE(publish_time, created_at) DESC) WHERE is_publish = TRUE;
CREATE INDEX IF NOT EXISTS idx_moments_content ON moments USING GIN (content);


-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    parent_id BIGINT DEFAULT NULL,
    title VARCHAR(100) NOT NULL,
    url VARCHAR(500) DEFAULT NULL,
    icon VARCHAR(500) DEFAULT NULL,
    sort INTEGER DEFAULT 5,
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_menus_parent FOREIGN KEY (parent_id) REFERENCES menus (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_menus_type ON menus (type);
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus (parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_sort ON menus (sort);


-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    link VARCHAR(500) NOT NULL DEFAULT '',
    data JSON NOT NULL DEFAULT '{}',
    sender_id INTEGER,
    target_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_sender_id ON notifications(sender_id);
CREATE INDEX IF NOT EXISTS idx_notifications_target_id ON notifications(target_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_title ON notifications USING gin(to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS idx_notifications_content ON notifications USING gin(to_tsvector('simple', content));


-- 用户通知关联表
CREATE TABLE IF NOT EXISTS user_notifications (
    id SERIAL PRIMARY KEY,
    notification_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_notifications_user_id ON user_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_user_notifications_notification_id ON user_notifications(notification_id);
CREATE INDEX IF NOT EXISTS idx_user_notifications_user_read ON user_notifications(user_id, is_read);


-- 反馈投诉表
CREATE TABLE IF NOT EXISTS feedbacks (
    id SERIAL PRIMARY KEY,
    ticket_no VARCHAR(50) UNIQUE NOT NULL,
    report_url VARCHAR(500) NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    form_content TEXT NOT NULL DEFAULT '{}',
    email VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    admin_reply TEXT,
    reply_time TIMESTAMP,
    user_agent VARCHAR(500),
    ip VARCHAR(45),
    feedback_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feedbacks_ticket_no ON feedbacks(ticket_no);
CREATE INDEX IF NOT EXISTS idx_feedbacks_report_type ON feedbacks(report_type);
CREATE INDEX IF NOT EXISTS idx_feedbacks_status ON feedbacks(status);
CREATE INDEX IF NOT EXISTS idx_feedbacks_email ON feedbacks(email);
CREATE INDEX IF NOT EXISTS idx_feedbacks_feedback_time ON feedbacks(feedback_time);


-- 设置表
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    "group" VARCHAR(50) NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_settings_group ON settings ("group");


-- 订阅者表
CREATE TABLE IF NOT EXISTS subscribers (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    active BOOLEAN DEFAULT TRUE,
    token VARCHAR(64) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscribers_email ON subscribers(email);
CREATE INDEX IF NOT EXISTS idx_subscribers_token ON subscribers(token);
CREATE INDEX IF NOT EXISTS idx_subscribers_active ON subscribers(active);


-- =============================================
-- 第二部分：全文搜索优化
-- =============================================

-- 创建触发器函数：自动更新搜索向量
CREATE OR REPLACE FUNCTION articles_search_trigger() RETURNS trigger AS $$
BEGIN
  new.search_vector :=
    setweight(to_tsvector('simple', coalesce(new.title, '')), 'A') ||
    setweight(to_tsvector('simple', coalesce(new.content, '')), 'B');
  RETURN new;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
DROP TRIGGER IF EXISTS tsvector_update ON articles;
CREATE TRIGGER tsvector_update
  BEFORE INSERT OR UPDATE OF title, content
  ON articles
  FOR EACH ROW
  EXECUTE FUNCTION articles_search_trigger();

-- 创建 GIN 索引
CREATE INDEX IF NOT EXISTS idx_articles_search_vector 
  ON articles USING GIN(search_vector);


-- =============================================
-- 第三部分：维护函数
-- =============================================

-- 自动更新 updated_at 的触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为 friends 表添加触发器
DROP TRIGGER IF EXISTS trigger_update_friends_updated_at ON friends;
CREATE TRIGGER trigger_update_friends_updated_at
    BEFORE UPDATE ON friends
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 为 settings 表添加触发器
DROP TRIGGER IF EXISTS trigger_update_settings_updated_at ON settings;
CREATE TRIGGER trigger_update_settings_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 为 rss_articles 表添加触发器
DROP TRIGGER IF EXISTS trigger_update_rss_articles_updated_at ON rss_articles;
CREATE TRIGGER trigger_update_rss_articles_updated_at
    BEFORE UPDATE ON rss_articles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- 第四部分：默认数据
-- =============================================

-- 插入默认超级管理员用户
-- 邮箱: admin@example.com
-- 密码: 123456
INSERT INTO users (email, password, has_password, nickname, avatar, is_enabled, role, github_id, google_id, qq_id)
SELECT 'admin@example.com', '$2a$10$he9ko8GhgnC.b2/GZ8oUHem6pE0HbwfQsnbcodpO5WaUlgsb0wyBS', TRUE, '超级管理员', '', TRUE, 'super_admin', '', '', ''
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@example.com' AND deleted_at IS NULL);

-- 插入示例分类
INSERT INTO categories (name, slug, description, count, sort)
SELECT '技术分享', '技术分享', 'Go、数据库、架构等技术文章', 1, 1
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE slug = '技术分享');

-- 插入示例标签
INSERT INTO tags (name, slug, description, count)
SELECT 'Go语言', 'Go语言', 'Go语言相关技术', 1
WHERE NOT EXISTS (SELECT 1 FROM tags WHERE slug = 'Go语言');

INSERT INTO tags (name, slug, description, count)
SELECT 'PostgreSQL', 'PostgreSQL', 'PostgreSQL数据库', 1
WHERE NOT EXISTS (SELECT 1 FROM tags WHERE slug = 'PostgreSQL');

-- 插入示例文章
INSERT INTO articles (title, slug, content, summary, cover, is_publish, is_top, view_count, category_id, publish_time, update_time, created_at, updated_at)
SELECT 
  'FLEC 博客系统使用指南',
  'a7k9m2x5',
  E'# FLEC 博客系统使用指南\n\n## 简介\n\n欢迎使用 FLEC 博客系统！这是一个基于 Go + Vue 3 开发的现代化博客平台。\n\n## 主要特性\n\n### 1. 技术栈\n- **后端**: Go + Gin + GORM + PostgreSQL\n- **前端**: Vue 3 + TypeScript + Vite\n- **缓存**: Redis（推荐）\n\n### 2. 核心功能\n- ✅ 文章管理（支持 Markdown）\n- ✅ 分类和标签\n- ✅ 多态评论系统\n- ✅ 文件上传\n- ✅ 友链管理\n- ✅ 动态发布\n- ✅ 访问统计\n- ✅ 全文搜索\n\n### 3. 性能优化\n- 全文搜索使用 PostgreSQL GIN 索引\n- 合理的索引设计\n- 支持 Redis 缓存\n\n## 快速开始\n\n### 默认管理员账号\n- 邮箱: `admin@example.com`\n- 密码: `123456`\n- **重要**: 请在首次登录后立即修改密码！\n\n### 使用建议\n1. 先完善个人资料\n2. 创建分类和标签\n3. 开始写作\n4. 配置 Redis 提升性能\n\n## 技术亮点\n\n本系统在数据库层面做了大量优化：\n- CHECK 约束保证数据完整性\n- 外键级联操作\n- 部分索引减少存储空间\n- 复合索引优化查询性能\n- 全文搜索 GIN 索引\n\n祝您使用愉快！',
  '欢迎使用 FLEC 博客系统，这是一个现代化的博客平台，支持 Markdown、全文搜索、多态评论等功能。',
  '',
  TRUE,
  TRUE,
  0,
  c.id,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM categories c
WHERE c.slug = '技术分享'
AND NOT EXISTS (SELECT 1 FROM articles WHERE slug = 'a7k9m2x5');

-- 插入文章标签关联
INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id
FROM articles a, tags t
WHERE a.slug = 'a7k9m2x5' AND t.slug = 'Go语言'
AND NOT EXISTS (SELECT 1 FROM article_tags WHERE article_id = a.id AND tag_id = t.id);

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id
FROM articles a, tags t
WHERE a.slug = 'a7k9m2x5' AND t.slug = 'PostgreSQL'
AND NOT EXISTS (SELECT 1 FROM article_tags WHERE article_id = a.id AND tag_id = t.id);

-- 插入示例评论
INSERT INTO comments (content, target_type, target_key, user_id, parent_id, root_id, reply_to, status, created_at, updated_at)
SELECT 
  '欢迎使用 FLEC 博客系统！如有问题请随时反馈。',
  'article',
  a.slug,
  u.id,
  NULL,
  NULL,
  NULL,
  1,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM articles a, users u
WHERE a.slug = 'a7k9m2x5' AND u.email = 'admin@example.com'
AND NOT EXISTS (
  SELECT 1 FROM comments 
  WHERE target_type = 'article' 
    AND target_key = a.slug
    AND user_id = u.id
);

-- 插入友链类型初始数据
INSERT INTO friend_types (name, sort, is_visible)
SELECT '友情链接', 5, TRUE
WHERE NOT EXISTS (SELECT 1 FROM friend_types WHERE name = '友情链接');

-- 插入示例友链
INSERT INTO friends (name, url, description, avatar, screenshot, sort, is_invalid, is_pending, type)
SELECT '爱吃猫的鱼', 'https://talen.top/', '前景可待 未来可期', 'https://image.talen.top/20251229184705_wgsva7g9.png', 'https://image.talen.top/20251229175030_heptu6ds.png', 5, FALSE, FALSE, 1
WHERE NOT EXISTS (SELECT 1 FROM friends WHERE url = 'https://blog.talen.top/');

-- 插入 settings 初始数据

-- 基本配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('basic.author', 'talen8', 'basic', TRUE),
('basic.author_email', 'your-email@example.com', 'basic', TRUE),
('basic.author_desc', '前景可待 未来可期', 'basic', TRUE),
('basic.author_avatar', '', 'basic', TRUE),
('basic.author_photo', '', 'basic', TRUE),
('basic.icp', '', 'basic', TRUE),
('basic.police_record', '', 'basic', TRUE),
('basic.admin_url', 'https://admin.your-domain.com', 'basic', TRUE),
('basic.blog_url', 'https://blog.your-domain.com', 'basic', TRUE),
('basic.home_url', 'https://your-domain.com', 'basic', TRUE)
ON CONFLICT (key) DO NOTHING;

-- 博客配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('blog.title', 'FlecBLOG', 'blog', TRUE),
('blog.subtitle', 'FlecBLOG', 'blog', TRUE),
('blog.slogan', '分享技术，记录生活', 'blog', TRUE),
('blog.description', 'FlecBLOG - 基于 Go + Vue 的现代化博客系统', 'blog', TRUE),
('blog.keywords', '博客,技术,编程,Go,Vue', 'blog', TRUE),
('blog.established', '2025-01-01', 'blog', TRUE),
('blog.favicon', '', 'blog', TRUE),
('blog.background_image', '', 'blog', TRUE),
('blog.screenshot', '', 'blog', TRUE),
('blog.typing_texts', '["前景可待 未来可期","笔耕不辍 文思不竭","笔学不怠 精进不休"]', 'blog', TRUE),
('blog.sidebar_social', '[{"name":"GitHub","url":"","icon":"github-line"},{"name":"Email","url":"","icon":"mail-line"},{"name":"X","url":"","icon":"twitter-x-line"}]', 'blog', TRUE),
('blog.footer_social', '[{"name":"Email","url":"","icon":"mail-line","position":"left"},{"name":"X","url":"","icon":"twitter-x-line","position":"left"},{"name":"GitHub","url":"","icon":"github-line","position":"right"},{"name":"Bilibili","url":"","icon":"bilibili-line","position":"right"}]', 'blog', TRUE),
-- 关于页面配置
('blog.about_describe', '持续学习的终身践行者，用代码和文字构建小小世界。在这里记录技术洞见、思考碎片和生活灵感。保持好奇，步履不停。', 'blog', TRUE),
('blog.about_describe_tips', '开发者 · 创作者 · 终身学习者', 'blog', TRUE),
('blog.about_exhibition', '', 'blog', TRUE),
('blog.about_profile', '[{"label":"姓名","value":"君の名は。","color":"#43a6c6"},{"label":"生于","value":"2000","color":"#c69043"},{"label":"故乡","value":"中国","color":"#d44040"},{"label":"职业","value":"开发者","color":"#b04fe6"},{"label":"兴趣","value":"编程与创作","color":"#43c66f"},{"label":"梦想","value":"创造价值","color":"#c643b3"}]', 'blog', TRUE),
('blog.about_personality', 'INFJ-A', 'blog', TRUE),
('blog.about_motto_main', '["前景可待，","未来可期。"]', 'blog', TRUE),
('blog.about_motto_sub', '用代码和文字记录生活，分享知识。', 'blog', TRUE),
('blog.about_socialize', '[{"name":"GitHub","url":"https://github.com"},{"name":"Email","url":"mailto:contact@example.com"}]', 'blog', TRUE),
('blog.about_creation', '[{"name":"掘金","url":"https://juejin.cn"},{"name":"知乎","url":"https://zhihu.com"}]', 'blog', TRUE),
('blog.about_versions', '[{"name":"FlecBLOG","version":"x.y.z"},{"name":"Vue","version":"3.5.0"},{"name":"Go","version":"1.25.0"}]', 'blog', TRUE),
('blog.about_unions', '[{"name":"BlogFinder","url":"https://bf.zzxworld.com"},{"name":"个站商店","url":"https://storeweb.cn"},{"name":"博友圈","url":"https://www.boyouquan.com"}]', 'blog', TRUE),
('blog.about_story', '欢迎来到我的个人博客！这里是我分享技术、记录生活、积累知识的小天地。希望这些内容能对你有所帮助。如果你有任何想法或建议，欢迎随时与我交流！', 'blog', TRUE),
('blog.custom_head', '', 'blog', TRUE),
('blog.custom_body', '', 'blog', TRUE),
('blog.emojis', '', 'blog', TRUE)
ON CONFLICT (key) DO NOTHING;

-- 通知配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('notification.email_host', 'smtp.qq.com', 'notification', FALSE),
('notification.email_port', '465', 'notification', FALSE),
('notification.email_username', 'your_email@qq.com', 'notification', FALSE),
('notification.email_password', 'your_smtp_authorization_code', 'notification', FALSE),
('notification.feishu_app_id', '', 'notification', FALSE),
('notification.feishu_secret', '', 'notification', FALSE),
('notification.feishu_chat_id', '', 'notification', FALSE)
ON CONFLICT (key) DO NOTHING;

-- 上传配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('upload.storage_type', 'local', 'upload', FALSE),
('upload.max_file_size', '5', 'upload', TRUE),
('upload.path_pattern', '{timestamp}_{random}{ext}', 'upload', FALSE),
('upload.access_key', '', 'upload', FALSE),
('upload.secret_key', '', 'upload', FALSE),
('upload.region', '', 'upload', FALSE),
('upload.bucket', '', 'upload', FALSE),
('upload.endpoint', '', 'upload', FALSE),
('upload.domain', '', 'upload', FALSE),
('upload.use_ssl', 'true', 'upload', FALSE)
ON CONFLICT (key) DO NOTHING;

-- OAuth 配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('oauth.github.enabled', 'false', 'oauth', TRUE),
('oauth.github.client_id', '', 'oauth', FALSE),
('oauth.github.client_secret', '', 'oauth', FALSE),
('oauth.github.redirect_url', '', 'oauth', FALSE),
('oauth.google.enabled', 'false', 'oauth', TRUE),
('oauth.google.client_id', '', 'oauth', FALSE),
('oauth.google.client_secret', '', 'oauth', FALSE),
('oauth.google.redirect_url', '', 'oauth', FALSE),
('oauth.qq.enabled', 'false', 'oauth', TRUE),
('oauth.qq.client_id', '', 'oauth', FALSE),
('oauth.qq.client_secret', '', 'oauth', FALSE),
('oauth.qq.redirect_url', '', 'oauth', FALSE),
('oauth.microsoft.enabled', 'false', 'oauth', TRUE),
('oauth.microsoft.client_id', '', 'oauth', FALSE),
('oauth.microsoft.client_secret', '', 'oauth', FALSE),
('oauth.microsoft.redirect_url', '', 'oauth', FALSE),
('oauth.session_secret', '', 'oauth', FALSE)
ON CONFLICT (key) DO NOTHING;

-- AI 配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('ai.base_url', 'https://api.deepseek.com', 'ai', FALSE),
('ai.api_key', '', 'ai', FALSE),
('ai.model', 'deepseek-chat', 'ai', FALSE)
ON CONFLICT (key) DO NOTHING;

-- 微信公众号配置
INSERT INTO settings (key, value, "group", is_public) VALUES
('wechat.app_id', '', 'wechat', FALSE),
('wechat.app_secret', '', 'wechat', FALSE),
('wechat.token_url', '', 'wechat', FALSE)
ON CONFLICT (key) DO NOTHING;

-- 插入示例动态
INSERT INTO moments (content, is_publish, publish_time, created_at)
SELECT '{"text":"欢迎使用 FLEC 博客系统！这是一条示例动态。"}'::jsonb, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM moments WHERE content->>'text' LIKE '%欢迎使用 FLEC 博客系统%');

-- 插入示例菜单
DO $$
DECLARE
    aggregate_id BIGINT;
    nav_id BIGINT;
    zhida_id BIGINT;
    service_id BIGINT;
    protocol_id BIGINT;
    support_id BIGINT;
BEGIN
    -- 网站聚合菜单 (aggregate)
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled)
    VALUES ('aggregate', NULL, '个站', '', '', 1, TRUE)
    RETURNING id INTO aggregate_id;
    
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('aggregate', aggregate_id, '主页', 'https://talen.top', 'ri-home-line', 1, TRUE),
    ('aggregate', aggregate_id, '博客', 'https://blog.talen.top', 'ri-article-line', 2, TRUE),
    ('aggregate', aggregate_id, 'GitHub', 'https://github.com/talen8', 'ri-github-fill', 3, TRUE);

    -- 顶部导航菜单 (navigation)
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('navigation', NULL, '归档', '/archive', 'ri-archive-line', 1, TRUE),
    ('navigation', NULL, '分类', '/categories', 'ri-folder-line', 2, TRUE),
    ('navigation', NULL, '标签', '/tags', 'ri-price-tag-3-line', 3, TRUE),
    ('navigation', NULL, '友链', '/friend', 'ri-links-line', 4, TRUE),
    ('navigation', NULL, '动态', '/moment', 'ri-chat-3-line', 5, TRUE),
    ('navigation', NULL, '留言', '/message', 'ri-quill-pen-line', 6, TRUE),
    ('navigation', NULL, '关于', '/about', 'ri-user-line', 7, TRUE);

    -- 页脚菜单 (footer) - 导航
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled)
    VALUES ('footer', NULL, '导航', '', '', 1, TRUE)
    RETURNING id INTO nav_id;
    
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('footer', nav_id, '文章归档', '/archive', '', 1, TRUE),
    ('footer', nav_id, '文章分类', '/categories', '', 2, TRUE),
    ('footer', nav_id, '文章标签', '/tags', '', 3, TRUE),
    ('footer', nav_id, '即时动态', '/moment', '', 4, TRUE);

    -- 页脚菜单 (footer) - 直达
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled)
    VALUES ('footer', NULL, '直达', '', '', 2, TRUE)
    RETURNING id INTO zhida_id;
    
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('footer', zhida_id, '申请友链', '/friend#apply', '', 1, TRUE);

    -- 页脚菜单 (footer) - 服务
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled)
    VALUES ('footer', NULL, '服务', '', '', 3, TRUE)
    RETURNING id INTO service_id;
    
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('footer', service_id, '服务状态', '/status', '', 2, TRUE),
    ('footer', service_id, '反馈投诉', '/feedback', '', 3, TRUE);

    -- 页脚菜单 (footer) - 协议
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled)
    VALUES ('footer', NULL, '协议', '', '', 4, TRUE)
    RETURNING id INTO protocol_id;
    
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('footer', protocol_id, 'Cookies', '/cookies', '', 1, TRUE),
    ('footer', protocol_id, '隐私协议', '/privacy', '', 2, TRUE),
    ('footer', protocol_id, '版权协议', '/copyright', '', 3, TRUE),
    ('footer', protocol_id, '提问须知', '/ask', '', 4, TRUE);

    -- 页脚菜单 (footer) - 支持
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled)
    VALUES ('footer', NULL, '支持', '', '', 5, TRUE)
    RETURNING id INTO support_id;
    
    INSERT INTO menus (type, parent_id, title, url, icon, sort, is_enabled) VALUES
    ('footer', support_id, '腾讯云', 'https://cloud.tencent.com', '', 1, TRUE);

END $$;


-- =============================================
-- 初始化完成提示
-- =============================================

DO $$
BEGIN
  RAISE NOTICE '==============================================';
  RAISE NOTICE 'FlecBLOG 系统数据库初始化完成！';
  RAISE NOTICE '==============================================';
  RAISE NOTICE '默认管理员账号:';
  RAISE NOTICE '  邮箱: admin@example.com';
  RAISE NOTICE '  密码: 123456';
  RAISE NOTICE '  重要: 请在首次登录后立即修改密码！';
  RAISE NOTICE '==============================================';
  RAISE NOTICE '已自动配置数据库权限';
  RAISE NOTICE '==============================================';
END $$;
