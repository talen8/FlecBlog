package service

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"flec_blog/config"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/auth"
	"flec_blog/pkg/email"
	"flec_blog/pkg/feishu"
	"flec_blog/pkg/random"

	"gorm.io/gorm"
)

// 配置键常量 - Basic 相关
const (
	KeyBasicAuthor       = "basic.author"        // 站长姓名
	KeyBasicAuthorEmail  = "basic.author_email"  // 站长邮箱
	KeyBasicAuthorDesc   = "basic.author_desc"   // 站长简介
	KeyBasicAuthorAvatar = "basic.author_avatar" // 站长头像
	KeyBasicAuthorPhoto  = "basic.author_photo"  // 站长形象
	KeyBasicICP          = "basic.icp"           // ICP备案号
	KeyBasicPoliceRecord = "basic.police_record" // 公安备案号
	KeyBasicAdminURL     = "basic.admin_url"     // 管理地址
	KeyBasicBlogURL      = "basic.blog_url"      // 博客地址
	KeyBasicHomeURL      = "basic.home_url"      // 主页地址
)

// 配置键常量 - Blog 相关
//
//goland:noinspection ALL
const (
	KeyBlogTitle             = "blog.title"               // 博客标题
	KeyBlogSubtitle          = "blog.subtitle"            // 博客副标题
	KeyBlogSlogan            = "blog.slogan"              // 博客标语
	KeyBlogDescription       = "blog.description"         // 博客描述
	KeyBlogKeywords          = "blog.keywords"            // 博客关键词
	KeyBlogEstablished       = "blog.established"         // 建站日期
	KeyBlogFavicon           = "blog.favicon"             // 网站Favicon
	KeyBlogBackgroundImage   = "blog.background_image"    // 背景图片
	KeyBlogScreenshot        = "blog.screenshot"          // 站点截图
	KeyBlogAnnouncement      = "blog.announcement"        // 公告内容
	KeyBlogTypingTexts       = "blog.typing_texts"        // 打字机效果文本（JSON数组）
	KeyBlogSidebarSocial     = "blog.sidebar_social"      // 侧边栏社交媒体（JSON数组）
	KeyBlogFooterSocial      = "blog.footer_social"       // 页脚社交媒体（JSON数组）
	KeyBlogFooterLinks       = "blog.footer_links"        // 页脚右侧链接（JSON数组）
	KeyBlogAboutDescribe     = "blog.about_describe"      // 个人描述
	KeyBlogAboutDescribeTips = "blog.about_describe_tips" // 描述提示
	KeyBlogAboutExhibition   = "blog.about_exhibition"    // 展览图片URL
	KeyBlogAboutProfile      = "blog.about_profile"       // 个人资料（JSON数组）
	KeyBlogAboutPersonality  = "blog.about_personality"   // 性格类型代码（如 INFJ-A）
	KeyBlogAboutMottoMain    = "blog.about_motto_main"    // 座右铭（JSON数组）
	KeyBlogAboutMottoSub     = "blog.about_motto_sub"     // 一言
	KeyBlogAboutSocialize    = "blog.about_socialize"     // 联系方式（JSON数组）
	KeyBlogAboutCreation     = "blog.about_creation"      // 创作平台（JSON数组）
	KeyBlogAboutVersions     = "blog.about_versions"      // 版本信息（JSON数组）
	KeyBlogAboutUnions       = "blog.about_unions"        // 站长联盟（JSON数组）
	KeyBlogAboutStory        = "blog.about_story"         // 心路历程
	KeyBlogCustomHead        = "blog.custom_head"         // 自定义 Head 代码
	KeyBlogCustomBody        = "blog.custom_body"         // 自定义 Body 代码
	KeyBlogEmojis            = "blog.emojis"              // 表情包配置
	KeyBlogFont              = "blog.font"                // 字体配置（URL|字体名称）
	KeyBlogMomentsSize       = "blog.moments_size"        // 动态列表每页数量
	KeyBlogMessageContent    = "blog.message_content"     // 留言信内容
	KeyBlogHomeLayout        = "blog.home_layout"         // 首页布局（waterfall/single_column）
)

// 配置键常量 - Notification 相关
const (
	KeyNotificationEmailHost     = "notification.email_host"
	KeyNotificationEmailPort     = "notification.email_port"
	KeyNotificationEmailUsername = "notification.email_username"
	KeyNotificationEmailPassword = "notification.email_password"
	KeyNotificationFeishuAppID   = "notification.feishu_app_id"
	KeyNotificationFeishuSecret  = "notification.feishu_secret"
	KeyNotificationFeishuChatID  = "notification.feishu_chat_id"
)

// 配置键常量 - Upload 相关
const (
	KeyUploadStorageType = "upload.storage_type"
	KeyUploadMaxFileSize = "upload.max_file_size"
	KeyUploadPathPattern = "upload.path_pattern"
	KeyUploadAccessKey   = "upload.access_key"
	KeyUploadSecretKey   = "upload.secret_key"
	KeyUploadRegion      = "upload.region"
	KeyUploadBucket      = "upload.bucket"
	KeyUploadEndpoint    = "upload.endpoint"
	KeyUploadDomain      = "upload.domain"
	KeyUploadUseSSL      = "upload.use_ssl"
)

// 配置键常量 - AI 相关
const (
	KeyAIBaseURL         = "ai.base_url"
	KeyAIAPIKey          = "ai.api_key"
	KeyAIModel           = "ai.model"
	KeyAISummaryPrompt   = "ai.summary_prompt"
	KeyAIAISummaryPrompt = "ai.ai_summary_prompt"
	KeyAITitlePrompt     = "ai.title_prompt"
	KeyAIMCPSecret       = "ai.mcp_secret"
)

// 配置键常量 - OAuth 相关
// #nosec G101 - 这些是配置项键名，不是实际凭证
const (
	KeyOAuthGithubEnabled         = "oauth.github.enabled"
	KeyOAuthGithubClientID        = "oauth.github.client_id"
	KeyOAuthGithubClientSecret    = "oauth.github.client_secret" // #nosec G101 - 配置键名
	KeyOAuthGithubRedirectURL     = "oauth.github.redirect_url"
	KeyOAuthGoogleEnabled         = "oauth.google.enabled"
	KeyOAuthGoogleClientID        = "oauth.google.client_id"
	KeyOAuthGoogleClientSecret    = "oauth.google.client_secret" // #nosec G101 - 配置键名
	KeyOAuthGoogleRedirectURL     = "oauth.google.redirect_url"
	KeyOAuthQQEnabled             = "oauth.qq.enabled"
	KeyOAuthQQClientID            = "oauth.qq.client_id"     // QQ AppID
	KeyOAuthQQClientSecret        = "oauth.qq.client_secret" // QQ AppKey #nosec G101 - 配置键名
	KeyOAuthQQRedirectURL         = "oauth.qq.redirect_url"
	KeyOAuthMicrosoftEnabled      = "oauth.microsoft.enabled"
	KeyOAuthMicrosoftClientID     = "oauth.microsoft.client_id"
	KeyOAuthMicrosoftClientSecret = "oauth.microsoft.client_secret"
	KeyOAuthMicrosoftRedirectURL  = "oauth.microsoft.redirect_url"
	KeyOAuthSessionSecret         = "oauth.session_secret" // Session 加密密钥
)

// SettingService 配置服务
type SettingService struct {
	repo        *repository.SettingRepository
	config      *config.Config // 全局配置对象引用（支持热重载）
	db          *gorm.DB
	mu          sync.RWMutex // 保护配置重载的并发安全
	fileService *FileService // 文件服务（用于标记文件状态）
}

// NewSettingService 创建配置服务
func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{repo: repository.NewSettingRepository(db), db: db}
}

// SetFileService 设置文件服务（用于依赖注入）
func (s *SettingService) SetFileService(fileService *FileService) {
	s.fileService = fileService
}

// SetConfig 设置全局配置对象（用于热重载）
func (s *SettingService) SetConfig(cfg *config.Config) {
	s.config = cfg
}

// GetByGroup 获取某个分组的所有配置
func (s *SettingService) GetByGroup(group string, isPublicOnly ...bool) (map[string]string, error) {
	return s.repo.GetByGroup(group, isPublicOnly...)
}

// GetAIConfig 获取AI配置
func (s *SettingService) GetAIConfig() (*config.AIConfig, error) {
	aiSettings, err := s.repo.GetByGroup(model.SettingGroupAI)
	if err != nil {
		return nil, err
	}

	cfg := &config.AIConfig{}
	if v, ok := aiSettings[KeyAIBaseURL]; ok && v != "" {
		cfg.BaseURL = v
	}
	if v, ok := aiSettings[KeyAIAPIKey]; ok && v != "" {
		cfg.APIKey = v
	}
	if v, ok := aiSettings[KeyAIModel]; ok && v != "" {
		cfg.Model = v
	}
	if v, ok := aiSettings[KeyAISummaryPrompt]; ok {
		cfg.SummaryPrompt = v
	}
	if v, ok := aiSettings[KeyAIAISummaryPrompt]; ok {
		cfg.AISummaryPrompt = v
	}
	if v, ok := aiSettings[KeyAITitlePrompt]; ok {
		cfg.TitlePrompt = v
	}
	if v, ok := aiSettings[KeyAIMCPSecret]; ok && v != "" {
		cfg.MCPSecret = v
	}

	return cfg, nil
}

// ResetMCPSecret 重新生成 MCP Secret 并持久化
func (s *SettingService) ResetMCPSecret() (string, error) {
	secret := random.String(32)
	if err := s.repo.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPSecret: secret,
	}); err != nil {
		return "", err
	}

	if s.config != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		if err := s.ApplyDatabaseConfig(s.config); err != nil {
			return "", err
		}
	}

	return secret, nil
}

// UpdateGroup 更新某个分组的配置（patch 方式），更新后自动重载
func (s *SettingService) UpdateGroup(group string, updates map[string]string) error {
	if err := validateSettingGroupUpdates(group, updates); err != nil {
		return err
	}

	var oldSettings map[string]string
	if s.fileService != nil && (group == model.SettingGroupBasic || group == model.SettingGroupBlog) {
		settings, err := s.repo.GetByGroup(group)
		if err == nil {
			oldSettings = settings
		}
	}

	// 更新数据库
	if err := s.repo.UpdateGroup(group, updates); err != nil {
		return err
	}

	if s.fileService != nil && oldSettings != nil {
		handleImageChange := func(key string) {
			newValue, ok := updates[key]
			if !ok {
				return
			}
			oldValue := oldSettings[key]
			if oldValue == newValue {
				return
			}
			if oldValue != "" {
				_ = s.fileService.MarkAsUnused(oldValue)
			}
			if newValue != "" {
				_ = s.fileService.MarkAsUsed(newValue)
			}
		}

		if group == model.SettingGroupBasic {
			handleImageChange(KeyBasicAuthorAvatar)
			handleImageChange(KeyBasicAuthorPhoto)
		}
		if group == model.SettingGroupBlog {
			handleImageChange(KeyBlogFavicon)
			handleImageChange(KeyBlogBackgroundImage)
			handleImageChange(KeyBlogAboutExhibition)
			handleImageChange(KeyBlogScreenshot)
		}
	}

	// 自动重载配置到内存（热重载）
	if s.config != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.ApplyDatabaseConfig(s.config)
	}

	return nil
}

// validateSettingGroupUpdates 校验配置项是否都属于指定分组
func validateSettingGroupUpdates(group string, updates map[string]string) error {
	for key := range updates {
		if !strings.HasPrefix(key, group+".") {
			return fmt.Errorf("配置项 %s 不属于分组 %s", key, group)
		}
	}
	return nil
}

// ApplyDatabaseConfig 从数据库加载配置并应用到 Config 对象
func (s *SettingService) ApplyDatabaseConfig(cfg *config.Config) error {
	if cfg == nil {
		return nil
	}

	// 加载 Basic 配置
	basicSettings, err := s.repo.GetByGroup(model.SettingGroupBasic)
	if err != nil {
		return err
	}
	if len(basicSettings) > 0 {
		if v, ok := basicSettings[KeyBasicAuthor]; ok && v != "" {
			cfg.Basic.Author = v
		}
		if v, ok := basicSettings[KeyBasicAuthorEmail]; ok && v != "" {
			cfg.Basic.AuthorEmail = v
		}
		if v, ok := basicSettings[KeyBasicAuthorDesc]; ok && v != "" {
			cfg.Basic.AuthorDesc = v
		}
		if v, ok := basicSettings[KeyBasicAuthorAvatar]; ok && v != "" {
			cfg.Basic.AuthorAvatar = v
		}
		if v, ok := basicSettings[KeyBasicAuthorPhoto]; ok && v != "" {
			cfg.Basic.AuthorPhoto = v
		}
		if v, ok := basicSettings[KeyBasicICP]; ok && v != "" {
			cfg.Basic.ICP = v
		}
		if v, ok := basicSettings[KeyBasicPoliceRecord]; ok && v != "" {
			cfg.Basic.PoliceRecord = v
		}
		if v, ok := basicSettings[KeyBasicAdminURL]; ok && v != "" {
			cfg.Basic.AdminURL = strings.TrimRight(v, "/")
		}
		if v, ok := basicSettings[KeyBasicBlogURL]; ok && v != "" {
			cfg.Basic.BlogURL = strings.TrimRight(v, "/")
		}
		if v, ok := basicSettings[KeyBasicHomeURL]; ok && v != "" {
			cfg.Basic.HomeURL = strings.TrimRight(v, "/")
		}
	}

	// 加载 Blog 配置
	blogSettings, err := s.repo.GetByGroup(model.SettingGroupBlog)
	if err != nil {
		return err
	}
	if len(blogSettings) > 0 {
		if v, ok := blogSettings[KeyBlogTitle]; ok && v != "" {
			cfg.Blog.Title = v
		}
		if v, ok := blogSettings[KeyBlogSubtitle]; ok && v != "" {
			cfg.Blog.Subtitle = v
		}
		if v, ok := blogSettings[KeyBlogSlogan]; ok && v != "" {
			cfg.Blog.Slogan = v
		}
		if v, ok := blogSettings[KeyBlogDescription]; ok && v != "" {
			cfg.Blog.Description = v
		}
		if v, ok := blogSettings[KeyBlogKeywords]; ok && v != "" {
			cfg.Blog.Keywords = v
		}
		if v, ok := blogSettings[KeyBlogEstablished]; ok && v != "" {
			cfg.Blog.Established = v
		}
		if v, ok := blogSettings[KeyBlogFavicon]; ok && v != "" {
			cfg.Blog.Favicon = v
		}
		if v, ok := blogSettings[KeyBlogBackgroundImage]; ok && v != "" {
			cfg.Blog.BackgroundImage = v
		}
		if v, ok := blogSettings[KeyBlogScreenshot]; ok && v != "" {
			cfg.Blog.Screenshot = v
		}
		if v, ok := blogSettings[KeyBlogAnnouncement]; ok {
			cfg.Blog.Announcement = v
		}
		if v, ok := blogSettings[KeyBlogCustomHead]; ok && v != "" {
			cfg.Blog.CustomHead = v
		}
		if v, ok := blogSettings[KeyBlogCustomBody]; ok && v != "" {
			cfg.Blog.CustomBody = v
		}
		if v, ok := blogSettings[KeyBlogEmojis]; ok && v != "" {
			cfg.Blog.Emojis = v
		}
		if v, ok := blogSettings[KeyBlogFont]; ok && v != "" {
			cfg.Blog.Font = v
		}
		cfg.Blog.MomentsSize = 30
		if v, ok := blogSettings[KeyBlogMomentsSize]; ok && v != "" {
			if size, err := strconv.Atoi(v); err == nil && size > 0 {
				cfg.Blog.MomentsSize = size
			}
		}
		if v, ok := blogSettings[KeyBlogMessageContent]; ok && v != "" {
			cfg.Blog.MessageContent = v
		}
		cfg.Blog.HomeLayout = "waterfall"
		if v, ok := blogSettings[KeyBlogHomeLayout]; ok && v != "" {
			cfg.Blog.HomeLayout = v
		}
	}

	// 加载 Notification 配置
	notificationSettings, err := s.repo.GetByGroup(model.SettingGroupNotification)
	if err != nil {
		return err
	}
	if len(notificationSettings) > 0 {
		if v, ok := notificationSettings[KeyNotificationEmailHost]; ok && v != "" {
			cfg.Notification.EmailHost = v
		}
		if v, ok := notificationSettings[KeyNotificationEmailPort]; ok && v != "" {
			if port, err := strconv.Atoi(v); err == nil {
				cfg.Notification.EmailPort = port
			}
		}
		if v, ok := notificationSettings[KeyNotificationEmailUsername]; ok && v != "" {
			cfg.Notification.EmailUsername = v
		}
		if v, ok := notificationSettings[KeyNotificationEmailPassword]; ok && v != "" {
			cfg.Notification.EmailPassword = v
		}
		if v, ok := notificationSettings[KeyNotificationFeishuAppID]; ok && v != "" {
			cfg.Notification.FeishuAppID = v
		}
		if v, ok := notificationSettings[KeyNotificationFeishuSecret]; ok && v != "" {
			cfg.Notification.FeishuSecret = v
		}
		if v, ok := notificationSettings[KeyNotificationFeishuChatID]; ok && v != "" {
			cfg.Notification.FeishuChatID = v
		}
	}

	// 加载 Upload 配置
	uploadSettings, err := s.repo.GetByGroup(model.SettingGroupUpload)
	if err != nil {
		return err
	}
	if len(uploadSettings) > 0 {
		if v, ok := uploadSettings[KeyUploadStorageType]; ok && v != "" {
			cfg.Upload.StorageType = v
		}
		if v, ok := uploadSettings[KeyUploadMaxFileSize]; ok && v != "" {
			if size, err := strconv.ParseInt(v, 10, 64); err == nil {
				cfg.Upload.MaxFileSize = size
			}
		}
		if v, ok := uploadSettings[KeyUploadPathPattern]; ok && v != "" {
			cfg.Upload.PathPattern = v
		}
		if v, ok := uploadSettings[KeyUploadAccessKey]; ok && v != "" {
			cfg.Upload.AccessKey = v
		}
		if v, ok := uploadSettings[KeyUploadSecretKey]; ok && v != "" {
			cfg.Upload.SecretKey = v
		}
		if v, ok := uploadSettings[KeyUploadRegion]; ok && v != "" {
			cfg.Upload.Region = v
		}
		if v, ok := uploadSettings[KeyUploadBucket]; ok && v != "" {
			cfg.Upload.Bucket = v
		}
		if v, ok := uploadSettings[KeyUploadEndpoint]; ok && v != "" {
			cfg.Upload.Endpoint = v
		}
		if v, ok := uploadSettings[KeyUploadDomain]; ok && v != "" {
			cfg.Upload.Domain = v
		}
		if v, ok := uploadSettings[KeyUploadUseSSL]; ok && v != "" {
			if useSSL, err := strconv.ParseBool(v); err == nil {
				cfg.Upload.UseSSL = useSSL
			}
		}
	}

	// 加载 AI 配置
	aiSettings, err := s.repo.GetByGroup(model.SettingGroupAI)
	if err != nil {
		return err
	}
	if len(aiSettings) > 0 {
		if v, ok := aiSettings[KeyAIBaseURL]; ok && v != "" {
			cfg.AI.BaseURL = v
		}
		if v, ok := aiSettings[KeyAIAPIKey]; ok && v != "" {
			cfg.AI.APIKey = v
		}
		if v, ok := aiSettings[KeyAIModel]; ok && v != "" {
			cfg.AI.Model = v
		}
		if v, ok := aiSettings[KeyAISummaryPrompt]; ok {
			cfg.AI.SummaryPrompt = v
		}
		if v, ok := aiSettings[KeyAIAISummaryPrompt]; ok {
			cfg.AI.AISummaryPrompt = v
		}
		if v, ok := aiSettings[KeyAITitlePrompt]; ok {
			cfg.AI.TitlePrompt = v
		}
		if v, ok := aiSettings[KeyAIMCPSecret]; ok && v != "" {
			cfg.AI.MCPSecret = v
		} else {
			cfg.AI.MCPSecret = random.String(32)
			_ = s.repo.UpdateGroup(model.SettingGroupAI, map[string]string{
				KeyAIMCPSecret: cfg.AI.MCPSecret,
			})
		}
	}

	// 加载 OAuth 配置
	oauthSettings, err := s.repo.GetByGroup(model.SettingGroupOAuth)
	if err != nil {
		return err
	}

	// 确保 Session Secret 存在
	var sessionSecret string
	if v, ok := oauthSettings[KeyOAuthSessionSecret]; ok && v != "" {
		sessionSecret = v
	} else {
		// 自动生成并保存
		sessionSecret = random.String(32)
		_ = s.repo.UpdateGroup(model.SettingGroupOAuth, map[string]string{
			KeyOAuthSessionSecret: sessionSecret,
		})
	}
	cfg.OAuth.SessionSecret = sessionSecret

	if len(oauthSettings) > 0 {
		// GitHub
		if v, ok := oauthSettings[KeyOAuthGithubEnabled]; ok && v != "" {
			if enabled, err := strconv.ParseBool(v); err == nil {
				cfg.OAuth.Github.Enabled = enabled
			}
		}
		if v, ok := oauthSettings[KeyOAuthGithubClientID]; ok && v != "" {
			cfg.OAuth.Github.ClientID = v
		}
		if v, ok := oauthSettings[KeyOAuthGithubClientSecret]; ok && v != "" {
			cfg.OAuth.Github.ClientSecret = v
		}
		if v, ok := oauthSettings[KeyOAuthGithubRedirectURL]; ok && v != "" {
			cfg.OAuth.Github.RedirectURL = v
		}

		// Google
		if v, ok := oauthSettings[KeyOAuthGoogleEnabled]; ok && v != "" {
			if enabled, err := strconv.ParseBool(v); err == nil {
				cfg.OAuth.Google.Enabled = enabled
			}
		}
		if v, ok := oauthSettings[KeyOAuthGoogleClientID]; ok && v != "" {
			cfg.OAuth.Google.ClientID = v
		}
		if v, ok := oauthSettings[KeyOAuthGoogleClientSecret]; ok && v != "" {
			cfg.OAuth.Google.ClientSecret = v
		}
		if v, ok := oauthSettings[KeyOAuthGoogleRedirectURL]; ok && v != "" {
			cfg.OAuth.Google.RedirectURL = v
		}

		// QQ
		if v, ok := oauthSettings[KeyOAuthQQEnabled]; ok && v != "" {
			if enabled, err := strconv.ParseBool(v); err == nil {
				cfg.OAuth.QQ.Enabled = enabled
			}
		}
		if v, ok := oauthSettings[KeyOAuthQQClientID]; ok && v != "" {
			cfg.OAuth.QQ.ClientID = v
		}
		if v, ok := oauthSettings[KeyOAuthQQClientSecret]; ok && v != "" {
			cfg.OAuth.QQ.ClientSecret = v
		}
		if v, ok := oauthSettings[KeyOAuthQQRedirectURL]; ok && v != "" {
			cfg.OAuth.QQ.RedirectURL = v
		}

		// Microsoft
		if v, ok := oauthSettings[KeyOAuthMicrosoftEnabled]; ok && v != "" {
			if enabled, err := strconv.ParseBool(v); err == nil {
				cfg.OAuth.Microsoft.Enabled = enabled
			}
		}
		if v, ok := oauthSettings[KeyOAuthMicrosoftClientID]; ok && v != "" {
			cfg.OAuth.Microsoft.ClientID = v
		}
		if v, ok := oauthSettings[KeyOAuthMicrosoftClientSecret]; ok && v != "" {
			cfg.OAuth.Microsoft.ClientSecret = v
		}
		if v, ok := oauthSettings[KeyOAuthMicrosoftRedirectURL]; ok && v != "" {
			cfg.OAuth.Microsoft.RedirectURL = v
		}
	}

	// 应用 OAuth 配置到 Goth (热重载)
	auth.UpdateConfig(&cfg.OAuth)

	// 应用邮件配置 (热重载)
	email.Reload(cfg)

	// 应用 Feishu 配置 (热重载)
	feishu.Reload(cfg.Notification.FeishuAppID, cfg.Notification.FeishuSecret, cfg.Notification.FeishuChatID)

	return nil
}
