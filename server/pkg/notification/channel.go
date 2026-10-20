package notification

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Channel 出站通知渠道
type Channel interface {
	ID() string                                  // 渠道唯一标识
	IsEnabled() bool                             // 是否启用
	Send(ctx context.Context, msg Message) error // 投递一条通知
}

// Message 投递内容的中间表示
type Message struct {
	Type       string         // 通知类型
	Title      string         // 标题
	Content    string         // 正文
	Link       string         // 跳转链接
	SenderName string         // 触发者名称
	Data       map[string]any // 原始结构化数据
}

// RssArticleItem RSS文章项
type RssArticleItem struct {
	Title      string
	Link       string
	FriendName string
}

// rssDigest 取出 RSS 日报数据：文章列表与未读总数
func rssDigest(msg Message) ([]RssArticleItem, uint) {
	articles, _ := msg.Data["articles"].([]RssArticleItem)
	return articles, getUint(msg.Data, "unread_count")
}

// getString 从 Data 中取出字符串值
func getString(data map[string]any, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

// getUint 从 Data 中取出无符号数值（ID、计数等）
func getUint(data map[string]any, key string) uint {
	switch v := data[key].(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case float64:
		return uint(v)
	}
	return 0
}

// Registry 渠道注册表
type Registry struct {
	channels map[string]Channel
}

// NewRegistry 创建渠道注册表
func NewRegistry() *Registry {
	return &Registry{channels: make(map[string]Channel)}
}

// Register 登记一个渠道
func (r *Registry) Register(c Channel) {
	id := c.ID()
	if _, exists := r.channels[id]; exists {
		panic(fmt.Sprintf("notification: duplicate channel id %q", id))
	}
	r.channels[id] = c
}

// Enabled 返回所有启用的渠道
func (r *Registry) Enabled() []Channel {
	out := make([]Channel, 0, len(r.channels))
	for _, c := range r.channels {
		if c.IsEnabled() {
			out = append(out, c)
		}
	}
	return out
}

// webhookTimeout 单次 Webhook 投递的超时上限
const webhookTimeout = 5 * time.Second

// postJSON 向目标地址 POST 一段 JSON
func postJSON(ctx context.Context, target string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, webhookTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(snippet))
	}
	return nil
}
