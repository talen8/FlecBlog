package notification

import (
	"context"
	"log/slog"

	"flec_blog/config"
	"flec_blog/pkg/email"
	"flec_blog/pkg/logsetup"
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
func (s *Service) SendEmail(ctx context.Context, email string, msg Message) error {
	if email != "" && s.emailNotifier != nil {
		return s.emailNotifier.Send(ctx, email, msg)
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
			defer logsetup.RecoverContext(ctx, "通知投递", slog.String("channel", ch.ID()), slog.String("type", msg.Type))
			if err := ch.Send(ctx, msg); err != nil {
				slog.WarnContext(ctx, "通知投递失败",
					slog.String("channel", ch.ID()),
					slog.String("type", msg.Type),
					slog.Any("err", err),
				)
			}
		}()
	}
}
