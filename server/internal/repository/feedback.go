package repository

import (
	"context"
	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// FeedbackRepository 反馈仓储
type FeedbackRepository struct {
	db *gorm.DB
}

// NewFeedbackRepository 创建反馈仓储实例
func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

// Create 创建反馈
func (r *FeedbackRepository) Create(ctx context.Context, feedback *model.Feedback) error {
	return r.db.WithContext(ctx).Create(feedback).Error
}

// Get 获取反馈详情
func (r *FeedbackRepository) Get(ctx context.Context, id uint) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.db.WithContext(ctx).First(&feedback, id).Error
	return &feedback, err
}

// List 获取反馈列表
func (r *FeedbackRepository) List(
	ctx context.Context,
	offset, limit int,
	keyword, reportType, status string,
	startTime, endTime string,
) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Feedback{})

	// 关键词搜索（工单号、投诉地址）
	if keyword != "" {
		searchKeyword := "%" + keyword + "%"
		query = query.Where("ticket_no ILIKE ? OR report_url ILIKE ?", searchKeyword, searchKeyword)
	}

	// 反馈类型筛选
	if reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 反馈时间范围筛选
	if startTime != "" {
		query = query.Where("feedback_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("feedback_time <= ?", endTime+" 23:59:59")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedbacks).Error
	return feedbacks, total, err
}

// ExistsByAttachmentURL 检查是否有反馈附件引用该文件
func (r *FeedbackRepository) ExistsByAttachmentURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Feedback{}).Where("form_content LIKE ?", "%"+url+"%").Count(&count).Error
	return count > 0, err
}

// Update 更新反馈
func (r *FeedbackRepository) Update(ctx context.Context, feedback *model.Feedback) error {
	return r.db.WithContext(ctx).Save(feedback).Error
}

// Delete 删除反馈
func (r *FeedbackRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Feedback{}, id).Error
}

// GetByTicketNo 根据工单号获取反馈
func (r *FeedbackRepository) GetByTicketNo(ctx context.Context, ticketNo string) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.db.WithContext(ctx).Where("ticket_no = ?", ticketNo).First(&feedback).Error
	return &feedback, err
}
