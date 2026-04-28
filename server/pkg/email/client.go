package email

import (
	"crypto/tls"
	"flec_blog/config"
	"flec_blog/pkg/logger"
	"fmt"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

// 邮件加密方式常量
const (
	emailSecureNone     = "none"     // 无加密 (端口 25)
	emailSecureSSL      = "ssl"      // SSL/TLS (端口 465)
	emailSecureSTARTTLS = "starttls" // STARTTLS (端口 587)
)

var (
	globalClient *Client
	globalMu     sync.RWMutex
)

// GetClient 获取全局邮件客户端实例
func GetClient() *Client {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalClient
}

// Client 邮件客户端
type Client struct {
	config      *config.Config // 全局配置对象引用（支持热重载）
	rateLimiter *RateLimiter
}

// RateLimiter 限流器
type RateLimiter struct {
	mu      sync.RWMutex
	records map[string][]time.Time
	limit   int
	window  time.Duration
}

// NewRateLimiter 创建限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{records: make(map[string][]time.Time), limit: limit, window: window}
}

// Allow 检查是否允许发送（滑动时间窗口限流）
func (r *RateLimiter) Allow(email string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	times := r.records[email]

	// 过滤有效时间窗口内的记录
	var validTimes []time.Time
	for _, t := range times {
		if now.Sub(t) < r.window {
			validTimes = append(validTimes, t)
		}
	}

	// 检查是否超过限制
	if len(validTimes) >= r.limit {
		return false
	}

	// 记录本次发送时间
	validTimes = append(validTimes, now)
	r.records[email] = validTimes
	return true
}

// Config 邮件配置
type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// Initialize 从全局配置创建邮件客户端
func Initialize(conf *config.Config) *Client {
	c := newClient(conf)

	globalMu.Lock()
	globalClient = c
	globalMu.Unlock()

	return c
}

// newClient 创建邮件客户端实例
func newClient(conf *config.Config) *Client {
	if conf == nil || conf.Notification.EmailHost == "" || conf.Notification.EmailUsername == "" {
		return nil
	}

	return &Client{
		config:      conf,
		rateLimiter: NewRateLimiter(5, time.Hour), // 5次/小时
	}
}

// Reload 重新加载邮件配置（热重载）
func Reload(conf *config.Config) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalClient = newClient(conf)
}

// SendEmail 发送邮件
func (c *Client) SendEmail(to, subject, htmlBody, fromName string) error {
	if c == nil || c.config == nil {
		return fmt.Errorf("邮件服务未初始化")
	}

	if !c.rateLimiter.Allow(to) {
		return fmt.Errorf("发送频率过高")
	}

	// 从全局配置读取最新的邮件配置（支持热重载）
	cfg := c.config.Notification

	// 如果没有指定发件人名称，使用配置中的博客标题
	if fromName == "" {
		fromName = c.config.Blog.Title
	}

	// 确定发件人邮箱地址
	fromEmail := cfg.EmailFrom
	if fromEmail == "" {
		fromEmail = cfg.EmailUsername
	}

	// 创建邮件
	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(fromEmail, fromName))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	// 发送邮件
	dialer := gomail.NewDialer(cfg.EmailHost, cfg.EmailPort, cfg.EmailUsername, cfg.EmailPassword)

	// 根据加密方式配置连接
	// #nosec G402 - 允许自签名证书，用于兼容某些邮件服务器配置
	switch cfg.EmailSecure {
	case emailSecureSSL:
		dialer.SSL = true
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	case emailSecureNone:
		dialer.SSL = false
		dialer.TLSConfig = nil
	case emailSecureSTARTTLS:
		fallthrough
	default:
		dialer.SSL = false
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := dialer.DialAndSend(msg); err != nil {
		logger.Error("邮件发送失败 to=%s err=%v", to, err)
		return err
	}

	logger.Info("邮件发送成功 to=%s subject=%s from=%s", to, subject, fromName)
	return nil
}

// HealthCheck 检查邮件服务可用性
func (c *Client) HealthCheck() error {
	if c == nil {
		return fmt.Errorf("未配置")
	}
	cfg := c.config.Notification
	dialer := gomail.NewDialer(cfg.EmailHost, cfg.EmailPort, cfg.EmailUsername, cfg.EmailPassword)

	// 根据加密方式配置连接
	// #nosec G402 - 允许自签名证书，用于兼容某些邮件服务器配置
	switch cfg.EmailSecure {
	case emailSecureSSL:
		dialer.SSL = true
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	case emailSecureNone:
		dialer.SSL = false
		dialer.TLSConfig = nil
	case emailSecureSTARTTLS:
		fallthrough
	default:
		dialer.SSL = false
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	conn, err := dialer.Dial()
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
