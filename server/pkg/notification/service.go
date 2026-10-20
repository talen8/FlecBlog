package notification

import (
	"context"

	"flec_blog/config"
	"flec_blog/pkg/email"
	"flec_blog/pkg/logger"
)

// Service 通知服务（多渠道）
type Service struct {
	emailNotifier *EmailNotifier
	registry      *Registry
}

// NewService 创建通知服务
func NewService(emailClient *email.Client, cfg *config.Config) *Service {
	var emailNotifier *EmailNotifier
	if emailClient != nil && cfg != nil {
		emailNotifier = NewEmailNotifier(emailClient, cfg)
	}

	registry := NewRegistry()
	if cfg != nil {
		registry.Register(NewFeishuChannel(cfg))
		registry.Register(NewTelegramChannel(cfg))
	}

	return &Service{
		emailNotifier: emailNotifier,
		registry:      registry,
	}
}

// SendEmail 发送邮件通知
func (s *Service) SendEmail(email string, msg Message) error {
	if email != "" && s.emailNotifier != nil {
		return s.emailNotifier.Send(email, msg)
	}
	return nil
}

// Broadcast 向所有已启用渠道投递通知
func (s *Service) Broadcast(ctx context.Context, msg Message) {
	if s.registry == nil {
		return
	}

	for _, ch := range s.registry.Enabled() {
		go func() {
			if err := ch.Send(ctx, msg); err != nil {
				logger.Warn("通知渠道 %s 投递失败 (类型: %s): %v", ch.ID(), msg.Type, err)
			}
		}()
	}
}
