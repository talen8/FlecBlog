package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/random"
	"flec_blog/pkg/utils"
)

// FeedbackService 反馈服务
type FeedbackService struct {
	repo                *repository.FeedbackRepository
	notificationService *NotificationService
	fileService         *FileService
}

// NewFeedbackService 创建反馈服务实例
func NewFeedbackService(repo *repository.FeedbackRepository, notificationService *NotificationService, fileService *FileService) *FeedbackService {
	return &FeedbackService{
		repo:                repo,
		notificationService: notificationService,
		fileService:         fileService,
	}
}

// Submit 提交反馈
func (s *FeedbackService) Submit(ctx context.Context, req *dto.SubmitFeedbackRequest, ip, userAgent string) (*dto.FeedbackResponse, error) {
	formContent := &dto.FeedbackContent{
		Description:     req.Description,
		Reason:          req.Reason,
		AttachmentFiles: req.AttachmentFiles,
	}

	formContentJSON, err := json.Marshal(formContent)
	if err != nil {
		return nil, err
	}

	feedbackTime := utils.Now().Time

	feedback := &model.Feedback{
		TicketNo:     s.generateTicketNo(feedbackTime),
		ReportUrl:    req.ReportUrl,
		ReportType:   req.ReportType,
		FormContent:  string(formContentJSON),
		Email:        req.Email,
		Status:       "pending",
		IP:           ip,
		UserAgent:    userAgent,
		FeedbackTime: feedbackTime, // 自动设置反馈时间
	}

	if err := s.repo.Create(ctx, feedback); err != nil {
		return nil, err
	}

	go s.markFilesAsUsed(ctx, req)
	go s.notifyAdmins(ctx, feedback)

	return s.toDTO(feedback), nil
}

// List 获取反馈列表
func (s *FeedbackService) List(ctx context.Context, req *dto.FeedbackQueryRequest) ([]dto.FeedbackResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	feedbacks, total, err := s.repo.List(ctx, offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.FeedbackResponse, len(feedbacks))
	for i, feedback := range feedbacks {
		result[i] = *s.toDTO(&feedback)
	}

	return result, total, nil
}

// Get 获取反馈详情
func (s *FeedbackService) Get(ctx context.Context, id uint) (*dto.FeedbackResponse, error) {
	feedback, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toDTO(feedback), nil
}

// GetByTicketNo 根据工单号获取反馈详情（公开接口）
func (s *FeedbackService) GetByTicketNo(ctx context.Context, ticketNo string) (*dto.FeedbackResponse, error) {
	feedback, err := s.repo.GetByTicketNo(ctx, ticketNo)
	if err != nil {
		return nil, err
	}
	return s.toDTO(feedback), nil
}

// Update 更新反馈
func (s *FeedbackService) Update(ctx context.Context, id uint, req *dto.UpdateFeedbackRequest) error {
	feedback, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	feedback.Status = req.Status

	// 如果有管理员回复内容，且之前没有回复时间，则设置回复时间
	if req.AdminReply != "" && feedback.ReplyTime == nil {
		now := utils.Now().Time
		feedback.ReplyTime = &now
	}
	feedback.AdminReply = req.AdminReply

	return s.repo.Update(ctx, feedback)
}

// Delete 删除反馈
func (s *FeedbackService) Delete(ctx context.Context, id uint) error {
	// 获取反馈详情以清理附件文件
	feedback, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	// 删除反馈
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 标记附件文件为未使用
	if s.fileService != nil && feedback.FormContent != "" {
		var content dto.FeedbackContent
		if err := json.Unmarshal([]byte(feedback.FormContent), &content); err == nil {
			for _, fileUrl := range content.AttachmentFiles {
				_ = s.fileService.MarkAsUnused(fileUrl)
			}
		}
	}

	return nil
}

// notifyAdmins 通知管理员
func (s *FeedbackService) notifyAdmins(ctx context.Context, feedback *model.Feedback) {
	if s.notificationService == nil {
		return
	}

	_ = s.notificationService.NotifyFeedback(ctx, feedback)
}

// markFilesAsUsed 标记文件为使用中
func (s *FeedbackService) markFilesAsUsed(_ context.Context, req *dto.SubmitFeedbackRequest) {
	if s.fileService == nil || len(req.AttachmentFiles) == 0 {
		return
	}
	for _, fileUrl := range req.AttachmentFiles {
		_ = s.fileService.MarkAsUsed(fileUrl)
	}
}

// generateTicketNo 生成工单号
// 格式：FB + 日期(YYYYMMDD) + 随机3位数字
// 例如：FB20241108001
func (s *FeedbackService) generateTicketNo(t time.Time) string {
	dateStr := t.Format("20060102")
	randomSuffix := random.Digits(3)
	return fmt.Sprintf("FB%s%s", dateStr, randomSuffix)
}

// toDTO 转换为DTO
func (s *FeedbackService) toDTO(feedback *model.Feedback) *dto.FeedbackResponse {
	var formContent dto.FeedbackContent
	if feedback.FormContent != "" {
		_ = json.Unmarshal([]byte(feedback.FormContent), &formContent)
	}

	return &dto.FeedbackResponse{
		ID:           feedback.ID,
		TicketNo:     feedback.TicketNo,
		ReportUrl:    feedback.ReportUrl,
		ReportType:   feedback.ReportType,
		FormContent:  formContent,
		Email:        feedback.Email,
		Status:       feedback.Status,
		AdminReply:   feedback.AdminReply,
		ReplyTime:    utils.ToJSONTime(feedback.ReplyTime),
		UserAgent:    feedback.UserAgent,
		IP:           feedback.IP,
		FeedbackTime: utils.NewJSONTime(feedback.FeedbackTime),
	}
}
