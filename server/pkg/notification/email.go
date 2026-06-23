package notification

import (
	"fmt"

	"flec_blog/config"
	"flec_blog/pkg/email"
)

// EmailNotifier 邮件通知器
type EmailNotifier struct {
	client *email.Client
	config *config.Config // 全局配置对象引用（支持热重载）
}

// NewEmailNotifier 创建邮件通知器
func NewEmailNotifier(client *email.Client, cfg *config.Config) *EmailNotifier {
	return &EmailNotifier{client: client, config: cfg}
}

// Send 发送邮件通知
func (e *EmailNotifier) Send(to string, data Data) error {
	if e.client == nil || e.config == nil {
		return nil // 邮件服务未配置，静默返回
	}

	// 根据通知类型选择地址（管理员通知用后台地址）
	siteURL := e.config.Basic.BlogURL
	if e.isAdminNotification(data.Type) {
		siteURL = e.config.Basic.AdminURL
	}

	// 准备模板数据
	templateData := map[string]interface{}{
		"Title":    data.Title,
		"Content":  data.Content,
		"Link":     data.Link,
		"SiteURL":  siteURL,
		"SiteName": e.config.Basic.Title,
		"Data":     data.Data,
	}

	// 使用全局模板管理器渲染
	templateName := e.getTemplateName(data.Type)
	tmplMgr := email.GetGlobalTemplateManager()
	htmlBody, err := tmplMgr.Render(templateName, templateData)
	if err != nil {
		return fmt.Errorf("渲染邮件模板失败: %w", err)
	}

	// 生成邮件主题
	subject := e.generateSubject(data)

	// 确定发件人名称（评论通知用回复者昵称，其他用站点名）
	fromName := templateData["SiteName"].(string)
	if data.SenderName != "" && (data.Type == "comment_reply" || data.Type == "comment_new") {
		fromName = data.SenderName
	}

	// 发送邮件
	return e.client.SendEmail(to, subject, htmlBody, fromName)
}

// generateSubject 生成邮件主题
func (e *EmailNotifier) generateSubject(data Data) string {
	siteName := e.config.Basic.Title

	switch data.Type {
	case "comment_reply":
		senderName := data.SenderName
		if senderName == "" {
			senderName = "匿名用户"
		}
		return fmt.Sprintf("[%s] 您收到了来自 %s 的评论回复", siteName, senderName)
	case "comment_new":
		return fmt.Sprintf("[%s] 收到了新的评论通知", siteName)
	case "feedback_new":
		return fmt.Sprintf("[%s] 收到了新的反馈投诉", siteName)
	case "friend_apply":
		return fmt.Sprintf("[%s] 收到了新的友链申请", siteName)
	default:
		return fmt.Sprintf("[%s] 新通知", siteName)
	}
}

// isAdminNotification 判断是否为管理员通知（需要使用后台地址）
func (e *EmailNotifier) isAdminNotification(notifType string) bool {
	adminTypes := []string{
		"comment_new",  // 新评论通知管理员
		"feedback_new", // 反馈投诉通知管理员
		"friend_apply", // 友链申请通知管理员
	}

	for _, t := range adminTypes {
		if t == notifType {
			return true
		}
	}
	return false
}

// getTemplateName 根据通知类型获取模板名称
func (e *EmailNotifier) getTemplateName(notifType string) string {
	switch notifType {
	case "comment_reply", "comment_new", "feedback_new", "friend_apply":
		return notifType
	default:
		return "default"
	}
}
