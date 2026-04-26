package service

import (
	"errors"
	"time"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/email"
	"flec_blog/pkg/random"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// VerificationService 验证服务
type VerificationService struct {
	repo        *repository.VerificationRepository
	userRepo    *repository.UserRepository
	emailClient *email.Client
	config      *config.Config
}

// NewVerificationService 创建验证服务实例
func NewVerificationService(repo *repository.VerificationRepository, userRepo *repository.UserRepository, emailClient *email.Client, cfg *config.Config) *VerificationService {
	return &VerificationService{repo: repo, userRepo: userRepo, emailClient: emailClient, config: cfg}
}

// ============ 通用服务 ============

// VerifyCode 验证验证码
func (s *VerificationService) VerifyCode(email, code string) (*model.Verification, error) {
	verification, err := s.repo.GetByEmailAndCode(email, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 增加失败次数
			if latest, _ := s.repo.GetLatestByEmail(email); latest != nil {
				_ = s.repo.IncrementFailedCount(latest.ID)
			}
			return nil, errors.New("验证码错误")
		}
		return nil, err
	}

	// 验证码状态检查
	if verification.Used {
		return nil, errors.New("验证码已被使用")
	}
	if time.Now().After(verification.ExpiresAt) {
		return nil, errors.New("验证码已过期")
	}
	if verification.FailedCount >= 5 {
		return nil, errors.New("验证失败次数过多，请重新获取验证码")
	}

	return verification, nil
}

// SendPasswordReset 发送密码重置邮件
func (s *VerificationService) SendPasswordReset(req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("该邮箱未注册")
		}
		return err
	}

	// 高权限账号不支持此方式重置密码
	if user.Role == model.RoleAdmin || user.Role == model.RoleSuperAdmin {
		return errors.New("管理员及超级管理员账号不支持此方式重置密码")
	}

	// 生成6位验证码
	code := random.Digits(6)

	// 使用模板管理器渲染邮件内容
	tmplMgr := email.GetGlobalTemplateManager()
	htmlBody, err := tmplMgr.Render("password_reset", map[string]interface{}{
		"Code":     code,
		"SiteName": s.config.Blog.Title,
		"SiteURL":  s.config.Basic.BlogURL,
	})
	if err != nil {
		return err
	}

	// 发送邮件
	subject := "[" + s.config.Blog.Title + "] 重置密码"
	if err := s.emailClient.SendEmail(req.Email, subject, htmlBody, ""); err != nil {
		return err
	}

	// 邮件发送成功后，再保存验证码
	verification := &model.Verification{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	return s.repo.Create(verification)
}

// ResetPassword 重置密码
func (s *VerificationService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// 验证验证码
	verification, err := s.VerifyCode(req.Email, req.Code)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码并递增 TokenVersion 使旧 token 失效
	if err := s.userRepo.UpdatePasswordAndIncrementVersion(user.ID, string(hashedPassword)); err != nil {
		return err
	}

	// 标记验证码为已使用
	return s.repo.MarkAsUsed(verification.ID)
}

// CleanExpiredVerifications 清理过期验证码
func (s *VerificationService) CleanExpiredVerifications() error {
	return s.repo.CleanExpired()
}
