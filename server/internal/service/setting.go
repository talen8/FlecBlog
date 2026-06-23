package service

import (
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
	"flec_blog/pkg/utils"

	"gorm.io/gorm"
)

// 配置键常量 - Basic 相关
const (
	KeyBasicAuthor        = "author"          // 站长姓名
	KeyBasicAuthorAvatar  = "author_avatar"   // 站长头像
	KeyBasicICP           = "icp"             // ICP备案号
	KeyBasicPoliceRecord  = "police_record"   // 公安备案号
	KeyBasicAdminURL      = "admin_url"       // 管理地址
	KeyBasicBlogURL       = "blog_url"        // 博客地址
	KeyBasicTitle         = "title"           // 博客标题
	KeyBasicSubtitle      = "subtitle"        // 博客副标题
	KeyBasicDescription   = "description"     // 博客描述
	KeyBasicKeywords      = "keywords"        // 博客关键词
	KeyBasicFavicon       = "favicon"         // 网站Favicon
	KeyBasicCustomHead    = "custom_head"     // 自定义 Head 代码
	KeyBasicCustomBody    = "custom_body"     // 自定义 Body 代码
	KeyBasicEmojis        = "emojis"          // 表情包配置
	KeyBasicMetingAPI     = "meting_api"      // Meting-API 地址
	KeyBasicCravatarURL   = "cravatar_url"    // 头像服务 URL（%s 为邮箱哈希）
	KeyBasicIPApiURL      = "ip_api_url"      // IP 归属地查询 URL（%s 为 IP）
	KeyBasicCoverMakerAPI = "cover_maker_api" // 封面制作图片源 API
)

// 配置键常量 - Notification 相关
const (
	KeyNotificationEmailHost     = "email_host"
	KeyNotificationEmailPort     = "email_port"
	KeyNotificationEmailSecure   = "email_secure"
	KeyNotificationEmailUsername = "email_username"
	KeyNotificationEmailFrom     = "email_from"
	KeyNotificationEmailPassword = "email_password"
	KeyNotificationFeishuAppID   = "feishu_app_id"
	KeyNotificationFeishuSecret  = "feishu_secret"
	KeyNotificationFeishuChatID  = "feishu_chat_id"
)

// 配置键常量 - Upload 相关
const (
	KeyUploadStorageType = "storage_type"
	KeyUploadMaxFileSize = "max_file_size"
	KeyUploadPathPattern = "path_pattern"
	KeyUploadAccessKey   = "access_key"
	KeyUploadSecretKey   = "secret_key"
	KeyUploadRegion      = "region"
	KeyUploadBucket      = "bucket"
	KeyUploadEndpoint    = "endpoint"
	KeyUploadDomain      = "domain"
	KeyUploadUseSSL      = "use_ssl"
)

// 配置键常量 - AI 相关
const (
	KeyAIBaseURL         = "base_url"
	KeyAIAPIKey          = "api_key"
	KeyAIModel           = "model"
	KeyAISummaryPrompt   = "summary_prompt"
	KeyAIAISummaryPrompt = "ai_summary_prompt"
	KeyAITitlePrompt     = "title_prompt"
	KeyAIMCPSecret       = "mcp_secret"
)

// 配置键常量 - OAuth 相关
// #nosec G101 - 这些是配置项键名，不是实际凭证
const (
	KeyOAuthGithubEnabled         = "github.enabled"
	KeyOAuthGithubClientID        = "github.client_id"
	KeyOAuthGithubClientSecret    = "github.client_secret" // #nosec G101 - 配置键名
	KeyOAuthGithubRedirectURL     = "github.redirect_url"
	KeyOAuthGoogleEnabled         = "google.enabled"
	KeyOAuthGoogleClientID        = "google.client_id"
	KeyOAuthGoogleClientSecret    = "google.client_secret" // #nosec G101 - 配置键名
	KeyOAuthGoogleRedirectURL     = "google.redirect_url"
	KeyOAuthQQEnabled             = "qq.enabled"
	KeyOAuthQQClientID            = "qq.client_id"     // QQ AppID
	KeyOAuthQQClientSecret        = "qq.client_secret" // QQ AppKey #nosec G101 - 配置键名
	KeyOAuthQQRedirectURL         = "qq.redirect_url"
	KeyOAuthMicrosoftEnabled      = "microsoft.enabled"
	KeyOAuthMicrosoftClientID     = "microsoft.client_id"
	KeyOAuthMicrosoftClientSecret = "microsoft.client_secret"
	KeyOAuthMicrosoftRedirectURL  = "microsoft.redirect_url"
	KeyOAuthOIDCEnabled           = "oidc.enabled"
	KeyOAuthOIDCIssuerURL         = "oidc.issuer_url"
	KeyOAuthOIDCClientID          = "oidc.client_id"
	KeyOAuthOIDCClientSecret      = "oidc.client_secret"
	KeyOAuthOIDCRedirectURL       = "oidc.redirect_url"
	KeyOAuthSessionSecret         = "session_secret" // Session 加密密钥
	KeyOAuthWechatEnabled         = "wechat.enabled" // 微信小程序是否启用
	KeyOAuthWechatAppID           = "wechat.appid"   // 微信小程序 AppID
	KeyOAuthWechatSecret          = "wechat.secret"  // 微信小程序 AppSecret
	KeyOAuthWorkerProxy           = "worker_proxy"   // OAuth Worker 代理地址
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
	var oldSettings map[string]string
	if s.fileService != nil && group == model.SettingGroupBasic {
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
			handleImageChange(KeyBasicFavicon)
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
		if v, ok := basicSettings[KeyBasicAuthorAvatar]; ok && v != "" {
			cfg.Basic.AuthorAvatar = v
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
		if v, ok := basicSettings[KeyBasicTitle]; ok && v != "" {
			cfg.Basic.Title = v
		}
		if v, ok := basicSettings[KeyBasicSubtitle]; ok && v != "" {
			cfg.Basic.Subtitle = v
		}
		if v, ok := basicSettings[KeyBasicDescription]; ok && v != "" {
			cfg.Basic.Description = v
		}
		if v, ok := basicSettings[KeyBasicKeywords]; ok && v != "" {
			cfg.Basic.Keywords = v
		}
		if v, ok := basicSettings[KeyBasicFavicon]; ok && v != "" {
			cfg.Basic.Favicon = v
		}
		if v, ok := basicSettings[KeyBasicCustomHead]; ok && v != "" {
			cfg.Basic.CustomHead = v
		}
		if v, ok := basicSettings[KeyBasicCustomBody]; ok && v != "" {
			cfg.Basic.CustomBody = v
		}
		if v, ok := basicSettings[KeyBasicEmojis]; ok && v != "" {
			cfg.Basic.Emojis = v
		}
		if v, ok := basicSettings[KeyBasicMetingAPI]; ok && v != "" {
			cfg.Basic.MetingAPI = v
		}
		if v, ok := basicSettings[KeyBasicCravatarURL]; ok && v != "" {
			cfg.Basic.CravatarURL = v
			utils.SetCravatarURL(v)
		}
		if v, ok := basicSettings[KeyBasicIPApiURL]; ok && v != "" {
			cfg.Basic.IPApiURL = v
			utils.SetIPApiURL(v)
		}
		if v, ok := basicSettings[KeyBasicCoverMakerAPI]; ok && v != "" {
			cfg.Basic.CoverMakerAPI = v
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
		// 默认使用 SSL
		cfg.Notification.EmailSecure = "ssl"
		if v, ok := notificationSettings[KeyNotificationEmailSecure]; ok && v != "" {
			cfg.Notification.EmailSecure = v
		}
		if v, ok := notificationSettings[KeyNotificationEmailUsername]; ok && v != "" {
			cfg.Notification.EmailUsername = v
		}
		if v, ok := notificationSettings[KeyNotificationEmailFrom]; ok && v != "" {
			cfg.Notification.EmailFrom = v
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

		// OIDC
		if v, ok := oauthSettings[KeyOAuthOIDCEnabled]; ok && v != "" {
			if enabled, err := strconv.ParseBool(v); err == nil {
				cfg.OAuth.OIDC.Enabled = enabled
			}
		}
		if v, ok := oauthSettings[KeyOAuthOIDCIssuerURL]; ok && v != "" {
			cfg.OAuth.OIDC.IssuerURL = v
		}
		if v, ok := oauthSettings[KeyOAuthOIDCClientID]; ok && v != "" {
			cfg.OAuth.OIDC.ClientID = v
		}
		if v, ok := oauthSettings[KeyOAuthOIDCClientSecret]; ok && v != "" {
			cfg.OAuth.OIDC.ClientSecret = v
		}
		if v, ok := oauthSettings[KeyOAuthOIDCRedirectURL]; ok && v != "" {
			cfg.OAuth.OIDC.RedirectURL = v
		}
		if v, ok := oauthSettings[KeyOAuthWorkerProxy]; ok && v != "" {
			cfg.OAuth.WorkerProxy = v
			auth.SetWorkerProxy(v)
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
