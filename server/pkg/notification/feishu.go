package notification

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"flec_blog/config"
)

// FeishuChannel 飞书自定义机器人渠道
type FeishuChannel struct {
	cfg *config.Config
}

// NewFeishuChannel 创建飞书渠道
func NewFeishuChannel(cfg *config.Config) *FeishuChannel {
	return &FeishuChannel{cfg: cfg}
}

// ID 渠道标识，同时作为配置项前缀
func (f *FeishuChannel) ID() string { return "feishu" }

// IsEnabled 填了 Webhook 地址即视为启用
func (f *FeishuChannel) IsEnabled() bool {
	return f.cfg != nil && f.cfg.Notification.FeishuWebhookURL != ""
}

// Send 投递一条通知到飞书群
func (f *FeishuChannel) Send(ctx context.Context, msg Message) error {
	card, err := json.Marshal(f.buildCard(msg))
	if err != nil {
		return fmt.Errorf("序列化卡片失败: %w", err)
	}

	payload := map[string]any{
		"msg_type": "interactive",
		"card":     json.RawMessage(card),
	}
	if secret := f.cfg.Notification.FeishuSecret; secret != "" {
		ts := time.Now().Unix()
		payload["timestamp"] = strconv.FormatInt(ts, 10)
		payload["sign"] = hmacSign(ts, secret)
	}

	// payload 的成员只有字符串和刚序列化出的 card，不可能序列化失败
	body, _ := json.Marshal(payload)

	if err := postJSON(ctx, f.cfg.Notification.FeishuWebhookURL, body); err != nil {
		return fmt.Errorf("飞书推送失败: %w", err)
	}
	return nil
}

// hmacSign 计算签名
func hmacSign(timestamp int64, secret string) string {
	h := hmac.New(sha256.New, []byte(strconv.FormatInt(timestamp, 10)+"\n"+secret))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// buildCard 按通知类型构建卡片
func (f *FeishuChannel) buildCard(msg Message) *feishuCard {
	switch msg.Type {
	case "comment_new":
		return f.commentCard(msg)
	case "friend_apply":
		return f.friendApplyCard(msg)
	case "feedback_new":
		return f.feedbackCard(msg)
	case "rss_feed_daily":
		return f.rssFeedCard(msg)
	default:
		return f.fallbackCard(msg)
	}
}

// commentCard 评论通知卡片
func (f *FeishuChannel) commentCard(msg Message) *feishuCard {
	content := fmt.Sprintf("来源：%s\n评论者：%s\n评论内容：%s",
		getString(msg.Data, "page_title"),
		msg.SenderName,
		getString(msg.Data, "comment_content"),
	)
	return newFeishuCard("📬 收到了新的评论通知", "blue", []any{
		newFeishuMarkdownElement(content),
		newFeishuButtonElement("查看", "primary", f.cfg.Basic.BlogURL+getString(msg.Data, "page_link")),
		newFeishuButtonElement("详情", "default", f.cfg.Basic.AdminURL+"/comments"),
	})
}

// friendApplyCard 友链申请卡片
func (f *FeishuChannel) friendApplyCard(msg Message) *feishuCard {
	name := getString(msg.Data, "site_name")
	link := getString(msg.Data, "site_url")
	logo := getString(msg.Data, "site_logo")
	screenshot := getString(msg.Data, "site_screenshot")
	content := fmt.Sprintf("名称：%s\n地址：[%s](%s)\n描述：%s\n头像：[%s](%s)\n截图：[%s](%s)",
		name, link, link, getString(msg.Data, "site_description"), logo, logo, screenshot, screenshot)
	return newFeishuCard("🔗 收到了新的友链申请", "green", []any{
		newFeishuMarkdownElement(content),
		newFeishuButtonElement("详情", "primary", f.cfg.Basic.AdminURL+"/friends"),
	})
}

// feedbackCard 反馈投诉卡片
func (f *FeishuChannel) feedbackCard(msg Message) *feishuCard {
	content := fmt.Sprintf("工单号：%s\n举报类型：%s",
		getString(msg.Data, "ticket_no"),
		getString(msg.Data, "report_type_text"),
	)
	return newFeishuCard("📝 收到了新的反馈投诉", "orange", []any{
		newFeishuMarkdownElement(content),
		newFeishuButtonElement("详情", "primary",
			fmt.Sprintf("%s/feedback/%d", f.cfg.Basic.AdminURL, getUint(msg.Data, "feedback_id"))),
	})
}

// rssFeedCard RSS 订阅日报卡片
func (f *FeishuChannel) rssFeedCard(msg Message) *feishuCard {
	articles, unread := rssDigest(msg)

	var content strings.Builder
	fmt.Fprintf(&content, "未读文章：%d 篇\n\n", unread)
	for i, article := range articles {
		fmt.Fprintf(&content, "%d. 《[%s](%s)》 %s\n\n", i+1, article.Title, article.Link, article.FriendName)
	}

	return newFeishuCard("📰 RSS订阅日报", "purple", []any{
		newFeishuMarkdownElement(content.String()),
		newFeishuButtonElement("查看详情", "primary", f.cfg.Basic.AdminURL+"/rssfeed"),
	})
}

// fallbackCard 未知类型的通用文本卡片
func (f *FeishuChannel) fallbackCard(msg Message) *feishuCard {
	title := msg.Title
	if title == "" {
		title = "📢 站点通知"
	}
	content := msg.Content
	if msg.Link != "" {
		content += fmt.Sprintf("\n\n[查看详情](%s)", msg.Link)
	}
	return newFeishuCard(title, "blue", []any{newFeishuMarkdownElement(content)})
}

// newFeishuCard 组装一张卡片 2.0
func newFeishuCard(title, template string, elements []any) *feishuCard {
	return &feishuCard{
		Schema: "2.0",
		Config: defaultCardConfig(),
		Header: defaultCardHeader(title, template),
		Body:   defaultCardBody(elements),
	}
}

// defaultCardConfig 默认卡片配置
func defaultCardConfig() *feishuCardConfig {
	return &feishuCardConfig{
		UpdateMulti: true,
		Style: &feishuCardStyle{
			TextSize: &feishuTextSizeStyle{
				NormalV2: &feishuSizeConfig{
					Default: "normal",
					PC:      "normal",
					Mobile:  "heading",
				},
			},
		},
	}
}

// defaultCardHeader 默认卡片头部
func defaultCardHeader(title, template string) *feishuCardHeader {
	return &feishuCardHeader{
		Template: template,
		Title: &feishuTextObject{
			Content: title,
			Tag:     "plain_text",
		},
		Subtitle: &feishuTextObject{
			Content: "",
			Tag:     "plain_text",
		},
		Padding: "12px 12px 12px 12px",
	}
}

// defaultCardBody 默认卡片主体
func defaultCardBody(elements []any) *feishuCardBody {
	return &feishuCardBody{
		Direction: "vertical",
		Padding:   "12px 12px 12px 12px",
		Elements:  elements,
	}
}

// newFeishuMarkdownElement 创建 markdown 元素
func newFeishuMarkdownElement(content string) *feishuMarkdownElement {
	return &feishuMarkdownElement{Tag: "markdown", Content: content}
}

// newFeishuButtonElement 创建带跳转链接的 button 元素
func newFeishuButtonElement(content, btnType, url string) *feishuButtonElement {
	return &feishuButtonElement{
		Tag:      "button",
		Text:     &feishuTextObject{Content: content, Tag: "plain_text"},
		Type:     btnType,
		Width:    "fill",
		Size:     "medium",
		MultiURL: &feishuURLElement{URL: url, PCURL: url},
	}
}

// feishuCard 飞书卡片 2.0 结构
type feishuCard struct {
	Schema string            `json:"schema"`
	Config *feishuCardConfig `json:"config,omitempty"`
	Header *feishuCardHeader `json:"header,omitempty"`
	Body   *feishuCardBody   `json:"body,omitempty"`
}

// feishuCardConfig 卡片配置
type feishuCardConfig struct {
	UpdateMulti bool             `json:"update_multi,omitempty"`
	Style       *feishuCardStyle `json:"style,omitempty"`
}

// feishuCardStyle 卡片样式
type feishuCardStyle struct {
	TextSize *feishuTextSizeStyle `json:"text_size,omitempty"`
}

// feishuTextSizeStyle 文本大小样式
type feishuTextSizeStyle struct {
	NormalV2 *feishuSizeConfig `json:"normal_v2,omitempty"`
}

// feishuSizeConfig 大小配置
type feishuSizeConfig struct {
	Default string `json:"default,omitempty"`
	PC      string `json:"pc,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
}

// feishuCardHeader 卡片头部
type feishuCardHeader struct {
	Template string            `json:"template"`
	Title    *feishuTextObject `json:"title"`
	Subtitle *feishuTextObject `json:"subtitle,omitempty"`
	Padding  string            `json:"padding,omitempty"`
}

// feishuCardBody 卡片主体
type feishuCardBody struct {
	Direction string `json:"direction,omitempty"`
	Padding   string `json:"padding,omitempty"`
	Elements  []any  `json:"elements"`
}

// feishuTextObject 文本对象
type feishuTextObject struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

// feishuMarkdownElement markdown 元素
type feishuMarkdownElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// feishuButtonElement 按钮元素
type feishuButtonElement struct {
	Tag      string            `json:"tag"`
	Text     *feishuTextObject `json:"text"`
	Type     string            `json:"type,omitempty"`
	Width    string            `json:"width,omitempty"`
	Size     string            `json:"size,omitempty"`
	MultiURL *feishuURLElement `json:"multi_url,omitempty"`
}

// feishuURLElement URL跳转元素
type feishuURLElement struct {
	URL   string `json:"url,omitempty"`
	PCURL string `json:"pc_url,omitempty"`
}
