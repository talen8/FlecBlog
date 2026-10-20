package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"flec_blog/config"
)

// telegramAPIBase 直连时的 Telegram Bot API 根地址
const telegramAPIBase = "https://api.telegram.org"

// telegramProxyPath Worker 上转发到官方域名的路径前缀
const telegramProxyPath = "/telegram"

// TelegramChannel Telegram Bot 渠道
type TelegramChannel struct {
	cfg *config.Config
}

// NewTelegramChannel 创建 Telegram 渠道
func NewTelegramChannel(cfg *config.Config) *TelegramChannel {
	return &TelegramChannel{cfg: cfg}
}

// ID 渠道标识，同时作为配置项前缀
func (t *TelegramChannel) ID() string { return "telegram" }

// IsEnabled Bot Token 与会话 ID 都填了才启用
func (t *TelegramChannel) IsEnabled() bool {
	return t.cfg != nil &&
		t.cfg.Notification.TelegramBotToken != "" && t.cfg.Notification.TelegramChatID != ""
}

// Send 投递一条通知到 Telegram
func (t *TelegramChannel) Send(ctx context.Context, msg Message) error {
	body, _ := json.Marshal(map[string]any{
		"chat_id": t.cfg.Notification.TelegramChatID,
		"text":    t.renderTelegramText(msg),
	})

	target := t.apiBase() + "/bot" + t.cfg.Notification.TelegramBotToken + "/sendMessage"
	if err := postJSON(ctx, target, body); err != nil {
		return fmt.Errorf("推送到 Telegram 失败: %w", err)
	}
	return nil
}

// apiBase 返回 API 根地址
func (t *TelegramChannel) apiBase() string {
	if t.cfg != nil && t.cfg.Basic.WorkerProxy != "" {
		return strings.TrimSuffix(t.cfg.Basic.WorkerProxy, "/") + telegramProxyPath
	}
	return telegramAPIBase
}

// renderTelegramText 按通知类型渲染纯文本
func (t *TelegramChannel) renderTelegramText(msg Message) string {
	switch msg.Type {
	case "comment_new":
		return t.renderCommentText(msg)
	case "friend_apply":
		return t.renderFriendApplyText(msg)
	case "feedback_new":
		return t.renderFeedbackText(msg)
	case "rss_feed_daily":
		return t.renderRssFeedText(msg)
	default:
		return t.renderFallbackText(msg)
	}
}

// renderCommentText 评论通知
func (t *TelegramChannel) renderCommentText(msg Message) string {
	return joinTelegramBlocks(
		"📬 收到了新的评论通知",
		joinTelegramLines(
			telegramField("来源", getString(msg.Data, "page_title")),
			telegramField("评论者", msg.SenderName),
			telegramField("评论内容", getString(msg.Data, "comment_content")),
		),
		joinTelegramLines(
			telegramField("查看", t.cfg.Basic.BlogURL+getString(msg.Data, "page_link")),
			telegramField("详情", t.cfg.Basic.AdminURL+"/comments"),
		),
	)
}

// renderFriendApplyText 友链申请
func (t *TelegramChannel) renderFriendApplyText(msg Message) string {
	return joinTelegramBlocks(
		"🔗 收到了新的友链申请",
		joinTelegramLines(
			telegramField("名称", getString(msg.Data, "site_name")),
			telegramField("地址", getString(msg.Data, "site_url")),
			telegramField("描述", getString(msg.Data, "site_description")),
		),
		telegramField("详情", t.cfg.Basic.AdminURL+"/friends"),
	)
}

// renderFeedbackText 反馈投诉
func (t *TelegramChannel) renderFeedbackText(msg Message) string {
	return joinTelegramBlocks(
		"📝 收到了新的反馈投诉",
		joinTelegramLines(
			telegramField("工单号", getString(msg.Data, "ticket_no")),
			telegramField("举报类型", getString(msg.Data, "report_type_text")),
		),
		telegramField("详情", fmt.Sprintf("%s/feedback/%d", t.cfg.Basic.AdminURL, getUint(msg.Data, "feedback_id"))),
	)
}

// renderRssFeedText RSS 日报
func (t *TelegramChannel) renderRssFeedText(msg Message) string {
	articles, unread := rssDigest(msg)

	var list strings.Builder
	fmt.Fprintf(&list, "未读文章：%d 篇", unread)
	for i, article := range articles {
		fmt.Fprintf(&list, "\n%d. 《%s》\n   %s", i+1, article.Title, article.Link)
		if article.FriendName != "" {
			list.WriteString(" —— " + article.FriendName)
		}
	}

	return joinTelegramBlocks(
		"📰 RSS订阅日报",
		list.String(),
		telegramField("查看详情", t.cfg.Basic.AdminURL+"/rssfeed"),
	)
}

// renderFallbackText 未知类型的通用文本
func (t *TelegramChannel) renderFallbackText(msg Message) string {
	title := msg.Title
	if title == "" {
		title = "📢 站点通知"
	}
	return joinTelegramBlocks(title, msg.Content, t.absoluteURL(msg.Link))
}

// absoluteURL 把后台相对路径补全成可点击的绝对 URL
func (t *TelegramChannel) absoluteURL(link string) string {
	if link == "" || strings.HasPrefix(link, "http") {
		return link
	}
	return t.cfg.Basic.AdminURL + link
}

// joinTelegramBlocks 用空行连接非空段落
func joinTelegramBlocks(blocks ...string) string {
	return joinTelegram("\n\n", blocks...)
}

// joinTelegramLines 用换行连接非空行
func joinTelegramLines(lines ...string) string {
	return joinTelegram("\n", lines...)
}

// joinTelegram 用 sep 连接非空片段
func joinTelegram(sep string, parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, sep)
}

// telegramField 拼「标签：值」，值为空时整段省略
func telegramField(label, value string) string {
	if value == "" {
		return ""
	}
	return label + "：" + value
}
