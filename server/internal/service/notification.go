package service

import (
	"context"
	"encoding/json"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/logger"
	notifier "flec_blog/pkg/notification"
	"flec_blog/pkg/utils"
)

// NotificationService 通知服务
type NotificationService struct {
	repo            *repository.NotificationRepository
	notificationSvc *notifier.Service
}

// NewNotificationService 创建通知服务实例
func NewNotificationService(repo *repository.NotificationRepository, notificationSvc *notifier.Service) *NotificationService {
	return &NotificationService{
		repo:            repo,
		notificationSvc: notificationSvc,
	}
}

// ============ 对外通知方法 ============

// NotifyCommentReply 评论回复通知（发送给被回复者）
func (s *NotificationService) NotifyCommentReply(ctx context.Context, senderID, receiverID uint, comment *model.Comment, pageTitle string) error {
	if senderID == receiverID {
		return nil
	}

	senderName := s.getUserNickname(ctx, senderID)
	receiverName := s.getUserNickname(ctx, receiverID)
	pageLink := buildCommentLink(comment.TargetType, comment.TargetKey, comment.ID)

	title := fmt.Sprintf("在《%s》收到了来自 @%s 的评论回复", pageTitle, senderName)
	data := map[string]interface{}{
		"page_title":        pageTitle,
		"page_link":         pageLink,
		"comment_id":        comment.ID,
		"comment_content":   comment.Content,
		"parent_comment_id": comment.ParentID,
		"sender_name":       senderName,
		"receiver_name":     receiverName,
	}

	return s.send(ctx, model.TypeCommentReply, title, comment.Content, pageLink, data, &senderID, &comment.ID, []uint{receiverID})
}

// NotifyCommentToAdmins 新评论通知管理员
func (s *NotificationService) NotifyCommentToAdmins(ctx context.Context, senderID uint, comment *model.Comment, pageTitle string, excludeUserID *uint) error {
	senderName := s.getUserNickname(ctx, senderID)
	pageLink := buildCommentLink(comment.TargetType, comment.TargetKey, comment.ID)

	title := "收到了新的评论通知"
	content := fmt.Sprintf("%s：%s", senderName, comment.Content)
	link := "/comments"
	data := map[string]interface{}{
		"comment_id":      comment.ID,
		"page_title":      pageTitle,
		"page_link":       pageLink,
		"comment_content": comment.Content,
		"sender_name":     senderName,
	}

	return s.sendToAdmins(ctx, model.TypeCommentNew, title, content, link, data, nil, &comment.ID, excludeUserID)
}

// NotifyFeedback 反馈投诉通知管理员
func (s *NotificationService) NotifyFeedback(ctx context.Context, feedback *model.Feedback) error {
	reportTypeText := getReportTypeText(feedback.ReportType)

	title := "收到了新的反馈投诉"
	content := fmt.Sprintf("工单号：%s，类型：%s", feedback.TicketNo, reportTypeText)
	link := fmt.Sprintf("/feedback/%d", feedback.ID)
	data := map[string]interface{}{
		"feedback_id":      feedback.ID,
		"ticket_no":        feedback.TicketNo,
		"report_url":       feedback.ReportUrl,
		"report_type":      feedback.ReportType,
		"report_type_text": reportTypeText,
		"form_content":     feedback.FormContent,
		"feedback_content": feedback.FormContent,
		"status":           feedback.Status,
	}

	return s.sendToAdmins(ctx, model.TypeFeedbackNew, title, content, link, data, nil, &feedback.ID, nil)
}

// NotifyFriendApply 友链申请通知管理员
func (s *NotificationService) NotifyFriendApply(ctx context.Context, friendID uint, siteName, siteURL, description, avatar, screenshot string, applicantID *uint) error {
	applicantName := "匿名用户"
	if applicantID != nil {
		applicantName = s.getUserNickname(ctx, *applicantID)
	}

	title := "收到了新的友链申请"
	content := fmt.Sprintf("%s：%s - %s", applicantName, siteName, siteURL)
	link := "/friends"
	data := map[string]interface{}{
		"friend_id":        friendID,
		"site_name":        siteName,
		"site_url":         siteURL,
		"site_description": description,
		"site_logo":        avatar,
		"site_screenshot":  screenshot,
		"applicant_name":   applicantName,
	}

	return s.sendToAdmins(ctx, model.TypeFriendApply, title, content, link, data, applicantID, nil, nil)
}

// NotifyFriendAbnormal 异常友链通知管理员（仅站内信）
func (s *NotificationService) NotifyFriendAbnormal(ctx context.Context, friendID uint, siteName string, abnormalCount int) error {
	title := "收到了异常友链提醒"
	content := fmt.Sprintf("%s已连续 %d 次检测异常", siteName, abnormalCount)
	link := "/friends"
	data := map[string]interface{}{
		"friend_id":      friendID,
		"site_name":      siteName,
		"abnormal_count": abnormalCount,
	}

	adminIDs, err := s.repo.GetAllAdmins(ctx)
	if err != nil {
		return err
	}
	if len(adminIDs) == 0 {
		return nil
	}

	return s.sendInApp(ctx, model.TypeFriendAbnormal, title, content, link, data, nil, &friendID, adminIDs)
}

// NotifyRssFeedDaily RSS订阅日报通知管理员
func (s *NotificationService) NotifyRssFeedDaily(unreadCount int, articles interface{}) error {
	if s.notificationSvc == nil {
		return nil
	}

	data := map[string]interface{}{
		"unread_count": unreadCount,
		"articles":     articles,
	}

	return s.notificationSvc.SendFeishu(notifier.Data{
		Type: "rss_feed_daily",
		Data: data,
	})
}

// NotifyVersionUpdateToSuperAdmins 版本更新通知超级管理员（仅站内信）
func (s *NotificationService) NotifyVersionUpdateToSuperAdmins(ctx context.Context, currentVersion, latestVersion, releaseURL string) error {
	superAdminIDs, err := s.repo.GetAllSuperAdmins(ctx)
	if err != nil {
		return err
	}
	if len(superAdminIDs) == 0 {
		return nil
	}

	content := fmt.Sprintf("当前版本 %s，发现新版本 %s", currentVersion, latestVersion)

	data := map[string]interface{}{
		"alert_type":      model.AlertTypeVersionUpdate,
		"message":         content,
		"severity":        "info",
		"current_version": currentVersion,
		"latest_version":  latestVersion,
		"release_url":     releaseURL,
	}

	return s.sendInApp(ctx, model.TypeSystemAlert, "发现新版本", content, releaseURL, data, nil, nil, superAdminIDs)
}

// HasVersionUpdateNotification 检查指定版本是否已经通知过
func (s *NotificationService) HasVersionUpdateNotification(ctx context.Context, latestVersion string) (bool, error) {
	return s.repo.ExistsVersionUpdateNotification(ctx, latestVersion)
}

// ============ 查询方法 ============

// ListForWeb 获取前台用户通知列表（仅评论回复）
func (s *NotificationService) ListForWeb(ctx context.Context, userID uint, req *dto.NotificationListRequest) (*dto.NotificationListResponse, error) {
	return s.list(ctx, userID, req, []model.NotificationType{model.TypeCommentReply})
}

// List 获取后台管理员通知列表
func (s *NotificationService) List(ctx context.Context, userID uint, req *dto.NotificationListRequest) (*dto.NotificationListResponse, error) {
	return s.list(ctx, userID, req, []model.NotificationType{
		model.TypeCommentNew,
		model.TypeFeedbackNew,
		model.TypeFriendApply,
		model.TypeFriendAbnormal,
		model.TypeSystemAlert,
	})
}

// MarkAsRead 标记为已读
func (s *NotificationService) MarkAsRead(ctx context.Context, id, userID uint) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

// MarkAllAsRead 标记所有为已读
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// ============ 内部方法 ============

// send 发送通知（站内+邮件）
func (s *NotificationService) send(ctx context.Context, notifType model.NotificationType, title, content, link string, data any, senderID, targetID *uint, receiverIDs []uint) error {
	if err := s.sendInApp(ctx, notifType, title, content, link, data, senderID, targetID, receiverIDs); err != nil {
		return err
	}

	// 异步发送邮件
	if s.notificationSvc != nil {
		notification := &model.Notification{
			Type:     notifType,
			Title:    title,
			Content:  content,
			Link:     link,
			SenderID: senderID,
			TargetID: targetID,
		}
		dataJSON, err := json.Marshal(data)
		if err != nil {
			return err
		}
		notification.Data = string(dataJSON)
		go s.sendEmails(notification, receiverIDs)
	}

	return nil
}

// sendInApp 仅发送站内通知
func (s *NotificationService) sendInApp(ctx context.Context, notifType model.NotificationType, title, content, link string, data any, senderID, targetID *uint, receiverIDs []uint) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	notification := &model.Notification{
		Type:     notifType,
		Title:    title,
		Content:  content,
		Link:     link,
		Data:     string(dataJSON),
		SenderID: senderID,
		TargetID: targetID,
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return err
	}

	return s.createUserNotifications(ctx, notification.ID, senderID, receiverIDs)
}

// sendToAdmins 发送给所有管理员（站内+邮件+飞书）
func (s *NotificationService) sendToAdmins(ctx context.Context, notifType model.NotificationType, title, content, link string, data any, senderID, targetID, excludeUserID *uint) error {
	adminIDs, err := s.repo.GetAllAdmins(ctx)
	if err != nil {
		return err
	}

	// 过滤排除的用户
	if excludeUserID != nil {
		filtered := make([]uint, 0, len(adminIDs))
		for _, id := range adminIDs {
			if id != *excludeUserID {
				filtered = append(filtered, id)
			}
		}
		adminIDs = filtered
	}

	if len(adminIDs) == 0 {
		return nil
	}

	if err := s.send(ctx, notifType, title, content, link, data, senderID, targetID, adminIDs); err != nil {
		return err
	}

	// 异步发送飞书
	if s.notificationSvc != nil {
		go s.sendFeishu(notifType, title, content, link, data)
	}

	return nil
}

// createUserNotifications 创建用户通知关联
func (s *NotificationService) createUserNotifications(ctx context.Context, notificationID uint, senderID *uint, receiverIDs []uint) error {
	if len(receiverIDs) == 0 {
		return nil
	}

	userNotifications := make([]model.UserNotification, 0, len(receiverIDs))
	for _, userID := range receiverIDs {
		if senderID != nil && *senderID == userID {
			continue
		}
		userNotifications = append(userNotifications, model.UserNotification{
			NotificationID: notificationID,
			UserID:         userID,
			IsRead:         false,
		})
	}

	if len(userNotifications) == 0 {
		return nil
	}

	return s.repo.BatchCreateUserNotifications(ctx, userNotifications)
}

// sendEmails 异步发送邮件通知
func (s *NotificationService) sendEmails(notification *model.Notification, receiverIDs []uint) {
	var data map[string]any
	if err := json.Unmarshal([]byte(notification.Data), &data); err != nil {
		return
	}

	ctx := context.Background()
	for _, userID := range receiverIDs {
		if notification.SenderID != nil && *notification.SenderID == userID {
			continue
		}

		user, err := s.repo.GetUserByID(ctx, userID)
		if err != nil || user.Email == "" {
			continue
		}

		senderName, _ := data["sender_name"].(string)
		notifData := notifier.Data{
			Title:      notification.Title,
			Content:    notification.Content,
			Link:       notification.Link,
			Type:       string(notification.Type),
			SenderName: senderName,
			Data:       data,
		}

		if err := s.notificationSvc.SendEmail(user.Email, notifData); err != nil {
			logger.Warn("发送邮件通知失败 (用户ID: %d): %v", userID, err)
		}
	}
}

// sendFeishu 异步发送飞书通知
func (s *NotificationService) sendFeishu(notifType model.NotificationType, title, content, link string, data any) {
	dataMap, _ := data.(map[string]any)
	senderName, _ := dataMap["sender_name"].(string)

	notifData := notifier.Data{
		Title:      title,
		Content:    content,
		Link:       link,
		Type:       string(notifType),
		SenderName: senderName,
		Data:       dataMap,
	}

	if err := s.notificationSvc.SendFeishu(notifData); err != nil {
		logger.Warn("发送飞书通知失败: %v", err)
	}
}

// list 获取通知列表
func (s *NotificationService) list(ctx context.Context, userID uint, req *dto.NotificationListRequest, allowedTypes []model.NotificationType) (*dto.NotificationListResponse, error) {
	offset := (req.Page - 1) * req.PageSize

	userNotifications, total, err := s.repo.List(ctx, userID, allowedTypes, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.repo.GetUnreadCount(ctx, userID, allowedTypes)
	if err != nil {
		return nil, err
	}

	list := make([]dto.NotificationResponse, len(userNotifications))
	for i, un := range userNotifications {
		list[i] = s.toDTO(&un)
	}

	return &dto.NotificationListResponse{
		List:        list,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		UnreadCount: unreadCount,
	}, nil
}

// toDTO 转换为DTO
func (s *NotificationService) toDTO(un *model.UserNotification) dto.NotificationResponse {
	notifType := string(un.Notification.Type)
	resp := dto.NotificationResponse{
		ID:        un.ID,
		Type:      notifType,
		TypeText:  getNotificationTypeText(notifType),
		Title:     un.Notification.Title,
		Content:   un.Notification.Content,
		Link:      un.Notification.Link,
		Data:      json.RawMessage(un.Notification.Data),
		TargetID:  un.Notification.TargetID,
		IsRead:    un.IsRead,
		CreatedAt: utils.NewJSONTime(un.CreatedAt),
	}

	if un.ReadAt != nil {
		readAt := utils.NewJSONTime(*un.ReadAt)
		resp.ReadAt = &readAt
	}

	if un.Notification.Sender != nil {
		resp.Sender = &un.Notification.Sender.Nickname
	}

	return resp
}

// getUserNickname 获取用户昵称
func (s *NotificationService) getUserNickname(ctx context.Context, userID uint) string {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil || user.Nickname == "" {
		return "匿名用户"
	}
	return user.Nickname
}

// getReportTypeText 获取反馈类型文本
func getReportTypeText(reportType string) string {
	typeMap := map[string]string{
		"copyright":     "版权侵权内容投诉",
		"inappropriate": "不当内容举报投诉",
		"summary":       "文章摘要问题反馈",
		"suggestion":    "功能建议优化反馈",
	}
	if text, ok := typeMap[reportType]; ok {
		return text
	}
	return reportType
}

// getNotificationTypeText 获取通知类型文本
func getNotificationTypeText(notifType string) string {
	typeMap := map[string]string{
		"comment_reply":   "评论回复",
		"comment_new":     "新评论",
		"feedback_new":    "反馈投诉",
		"friend_apply":    "友链申请",
		"friend_abnormal": "异常友链",
		"system_alert":    "系统通知",
	}
	if text, ok := typeMap[notifType]; ok {
		return text
	}
	return "通知"
}

func buildCommentLink(targetType, targetKey string, commentID uint) string {
	return fmt.Sprintf("%s#comment-%d", buildCommentTargetPath(targetType, targetKey), commentID)
}

func buildCommentTargetPath(targetType, targetKey string) string {
	switch targetType {
	case "article":
		return fmt.Sprintf("/posts/%s", targetKey)
	case "page":
		if targetKey == "" {
			return "/"
		}
		return fmt.Sprintf("/%s", targetKey)
	default:
		if targetKey == "" {
			return "/"
		}
		return fmt.Sprintf("/%s", targetKey)
	}
}
