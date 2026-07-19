package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config 应用配置
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	JWT          JWTConfig
	Basic        BasicConfig        // 从数据库加载
	Notification NotificationConfig // 从数据库加载
	Upload       UploadConfig       // 从数据库加载
	AI           AIConfig           // 从数据库加载
	OAuth        OAuthConfig        // 从数据库加载
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int
	AllowOrigins []string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string
}

// BasicConfig 基本配置（从数据库动态加载）
type BasicConfig struct {
	Author        string // 站长姓名
	AuthorAvatar  string // 站长头像
	ICP           string // ICP备案号
	PoliceRecord  string // 公安备案号
	AdminURL      string // 管理地址
	BlogURL       string // 博客地址
	Title         string // 博客标题
	Subtitle      string // 博客副标题
	Description   string // 博客描述
	Keywords      string // 博客关键词
	Favicon       string // 网站Favicon
	CustomHead    string // 自定义 Head 代码
	CustomBody    string // 自定义 Body 代码
	Emojis        string // 表情包配置
	MetingAPI     string // Meting-API 地址
	CravatarURL   string // 头像服务 URL（%s 为邮箱哈希）
	IPApiURL      string // IP 归属地查询 URL（%s 为 IP）
	CoverMakerAPI string // 封面制作图片源 API
}

// NotificationConfig 通知配置（从数据库动态加载）
type NotificationConfig struct {
	EmailHost     string // SMTP服务器地址
	EmailPort     int    // SMTP服务器端口
	EmailSecure   string // 加密方式: none/ssl/starttls
	EmailUsername string // SMTP登录账号
	EmailFrom     string // 发件人邮箱地址
	EmailPassword string // 邮箱密码
	FeishuAppID   string // 飞书应用ID
	FeishuSecret  string // 飞书应用Secret
	FeishuChatID  string // 飞书群聊ID
}

// UploadConfig 上传配置（从数据库动态加载）
type UploadConfig struct {
	StorageType string                 // 存储类型: local/s3/cos/oss/kodo/minio
	MaxFileSize int64                  // 最大文件大小(MB)
	PathPattern string                 // 路径生成模式
	AccessKey   string                 // 访问密钥
	SecretKey   string                 // 秘密密钥
	Region      string                 // 地域
	Bucket      string                 // 存储桶名称
	Endpoint    string                 // 服务端点
	Domain      string                 // 访问域名
	UseSSL      bool                   // 是否使用HTTPS
	Extra       map[string]interface{} // 额外配置
}

// AIConfig AI服务配置（从数据库动态加载）
type AIConfig struct {
	BaseURL         string // API 端点 (OpenAI 兼容格式)
	APIKey          string // API 密钥
	Model           string // 模型名称
	SummaryPrompt   string // 文章摘要提示词
	AISummaryPrompt string // AI 总结提示词
	TitlePrompt     string // 标题生成提示词
	MCPSecret       string // MCP 专用鉴权密钥（系统自动生成）
}

// OAuthConfig OAuth配置（从数据库动态加载）
type OAuthConfig struct {
	SessionSecret string // Session加密密钥（自动生成）
	WorkerProxy   string // OAuth Worker 代理地址
	Github        OAuthProviderConfig
	Google        OAuthProviderConfig
	QQ            OAuthProviderConfig
	Microsoft     OAuthProviderConfig
	OIDC          OIDCConfig
}

// OIDCConfig 通用 OIDC 提供商配置
type OIDCConfig struct {
	Enabled      bool   // 开关
	ClientID     string // Client ID
	ClientSecret string // Client Secret
	RedirectURL  string // 回调地址
	IssuerURL    string // OIDC Issuer URL
}

// OAuthProviderConfig 单个OAuth提供商配置
type OAuthProviderConfig struct {
	Enabled      bool   // 开关
	ClientID     string // Client ID
	ClientSecret string // Client Secret
	RedirectURL  string // 回调地址
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			AllowOrigins: getEnvAsSlice("SERVER_ALLOW_ORIGINS", []string{"*"}),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "postgres"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
	}

	// 验证必需的配置
	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD 环境变量未设置")
	}
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET 环境变量未设置")
	}

	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// getEnvAsInt 获取整数类型的环境变量
func getEnvAsInt(key string, defaultVal int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// getEnvAsSlice 获取切片类型的环境变量（逗号分隔）
func getEnvAsSlice(key string, defaultVal []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultVal
	}
	// 分割字符串并去除空格
	values := strings.Split(valueStr, ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	return values
}
