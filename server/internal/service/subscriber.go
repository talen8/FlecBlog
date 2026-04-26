package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/email"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/utils"

	"gorm.io/gorm"
)

// SubscriberService 订阅者服务
type SubscriberService struct {
	repo        *repository.SubscriberRepository
	emailClient *email.Client
	config      *config.Config
}

// NewSubscriberService 创建订阅者服务
func NewSubscriberService(repo *repository.SubscriberRepository, emailClient *email.Client, config *config.Config) *SubscriberService {
	return &SubscriberService{
		repo:        repo,
		emailClient: emailClient,
		config:      config,
	}
}

// Subscribe 订阅
func (s *SubscriberService) Subscribe(ctx context.Context, email string) error {
	sub, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		if sub.Active {
			return errors.New("该邮箱已订阅")
		}
		// 重新激活订阅
		sub.Active = true
		if err := s.repo.Update(ctx, sub); err != nil {
			return err
		}
		// 异步发送欢迎邮件，邮件失败不影响订阅成功
		go func() {
			if err := s.sendWelcomeEmail(sub); err != nil {
				logger.Warn("发送欢迎邮件失败 (邮箱: %s): %v", sub.Email, err)
			}
		}()
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建新订阅者
	newSub := &model.Subscriber{
		Email:  email,
		Active: true,
		Token:  generateToken(),
	}

	if err := s.repo.Create(ctx, newSub); err != nil {
		return err
	}

	// 异步发送欢迎邮件，邮件失败不影响订阅成功
	go func() {
		if err := s.sendWelcomeEmail(newSub); err != nil {
			logger.Warn("发送欢迎邮件失败 (邮箱: %s): %v", newSub.Email, err)
		}
	}()
	return nil
}

// Unsubscribe 退订
func (s *SubscriberService) Unsubscribe(ctx context.Context, token string) error {
	sub, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无效的退订链接")
		}
		return err
	}

	if !sub.Active {
		return errors.New("该邮箱已退订")
	}

	sub.Active = false
	return s.repo.Update(ctx, sub)
}

// sendWelcomeEmail 发送欢迎邮件
func (s *SubscriberService) sendWelcomeEmail(sub *model.Subscriber) error {
	if s.emailClient == nil {
		return errors.New("邮件服务未配置")
	}

	siteURL := s.config.Basic.BlogURL
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	unsubscribeURL := fmt.Sprintf("%s/subscribe?action=unsubscribe&token=%s", siteURL, sub.Token)

	subject := fmt.Sprintf("订阅成功 - %s", s.config.Blog.Title)
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>订阅成功</h2>
			<p>感谢您订阅 <strong>%s</strong>！</p>
			<p>您将会收到本站的最新文章推送。</p>
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
			<p style="color: #999; font-size: 12px;">
				如果这不是您本人的操作，或您想退订，请点击：<a href="%s">退订链接</a>
			</p>
		</div>
	`, s.config.Blog.Title, unsubscribeURL)

	return s.emailClient.SendEmail(sub.Email, subject, htmlBody, "")
}

// SendArticleNotification 发送文章推送通知（并发发送）
func (s *SubscriberService) SendArticleNotification(ctx context.Context, article *model.Article) error {
	// 获取所有活跃订阅者
	subscribers, err := s.repo.GetActiveSubscribers(ctx)
	if err != nil {
		return fmt.Errorf("获取订阅者列表失败: %w", err)
	}

	if len(subscribers) == 0 {
		return nil
	}

	// 并发控制：最多10个并发
	semaphore := make(chan struct{}, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0

	// 并发发送邮件
	for _, sub := range subscribers {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(subscriber *model.Subscriber) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			if err := s.sendArticleEmail(subscriber, article); err != nil {
				logger.Warn("发送文章推送失败 (邮箱: %s): %v", subscriber.Email, err)
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(sub)
	}

	// 等待所有邮件发送完成
	wg.Wait()

	logger.Info("文章推送完成: 成功 %d/%d", successCount, len(subscribers))
	return nil
}

// sendArticleEmail 发送文章推送邮件
func (s *SubscriberService) sendArticleEmail(sub *model.Subscriber, article *model.Article) error {
	if s.emailClient == nil {
		return errors.New("邮件服务未配置")
	}

	siteURL := s.config.Basic.BlogURL
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	articleURL := fmt.Sprintf("%s/posts/%s", siteURL, article.Slug)
	unsubscribeURL := fmt.Sprintf("%s/subscribe?action=unsubscribe&token=%s", siteURL, sub.Token)

	subject := fmt.Sprintf("新文章推送：%s - %s", article.Title, s.config.Blog.Title)
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">%s</h2>
			<p style="color: #666; line-height: 1.6;">%s</p>
			<div style="margin: 30px 0; text-align: center;">
				<a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 30px; text-decoration: none; border-radius: 4px; display: inline-block;">
					阅读全文
				</a>
			</div>
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
			<p style="color: #999; font-size: 12px;">
				这是来自 <strong>%s</strong> 的文章推送。<br>
				如需退订，请点击：<a href="%s">退订链接</a>
			</p>
		</div>
	`, article.Title, article.Summary, articleURL, s.config.Blog.Title, unsubscribeURL)

	return s.emailClient.SendEmail(sub.Email, subject, htmlBody, "")
}

// generateToken 生成随机令牌
func generateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		logger.Error("生成随机令牌失败: %v", err)
	}
	return hex.EncodeToString(b)
}

// List 获取订阅者列表（后台管理）
func (s *SubscriberService) List(ctx context.Context, req *dto.SubscriberQueryRequest) ([]dto.SubscriberResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	subscribers, total, err := s.repo.List(ctx, offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.SubscriberResponse, len(subscribers))
	for i, sub := range subscribers {
		result[i] = dto.SubscriberResponse{
			ID:        sub.ID,
			Email:     sub.Email,
			Active:    sub.Active,
			CreatedAt: utils.NewJSONTime(sub.CreatedAt),
			UpdatedAt: utils.NewJSONTime(sub.UpdatedAt),
		}
	}

	return result, total, nil
}

// Delete 删除订阅者（后台管理）
func (s *SubscriberService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
