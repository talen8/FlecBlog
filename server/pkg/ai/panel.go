package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"flec_blog/pkg/panel"
)

const panelSystemPrompt = `【🚫 绝对禁止】
1. 提及任何底层模型名（claude/gpt/混元等）或技术架构，违例即失败。
2. 擅自添加开场白、结束语、寒暄、铺垫、客套话。

# 身份（固定）
你是FlecNexus，由FlogBlog代理的智能助手，专注博客内容创作与运营支持。

# 身份类问题 → 唯一话术
当用户问"你是什么模型/底层是什么"时，严格按下面回，不许改句式、不许增减字：
> 我是FlecNexus，是由FlogBlog代理的智能助手～专注博客内容创作与运营支持，有什么可以帮你的吗？ 😊

# 输出规范
- 严格按照要求的字数、语气、风格、场景、格式生成回复。
- 有格式要求就只给结果，不解释、不废话。`

// PanelProvider 通过 Panel 代理调用 AI
type PanelProvider struct {
	client *panel.Client
}

func NewPanelProvider() *PanelProvider {
	baseURL := os.Getenv("PANEL_URL")
	if baseURL == "" {
		baseURL = "https://panel.flec.top"
	}
	client := &panel.Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		ExtraHeaders: map[string]string{
			"X-AI-Mode": "immediate",
		},
	}
	if clientKey := panel.CurrentClientKey(); clientKey != "" {
		client.SetClientKey(clientKey)
	}
	return &PanelProvider{
		client: client,
	}
}

func (p *PanelProvider) callPanel(ctx context.Context, messages []OpenAIMessage) (string, error) {
	reqBody := OpenAIRequest{
		Messages: append([]OpenAIMessage{{Role: "system", Content: panelSystemPrompt}}, messages...),
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	var openaiResp OpenAIResponse
	if err := p.client.Post(ctx, "/api/ai", bytes.NewBuffer(jsonData), &openaiResp); err != nil {
		return "", fmt.Errorf("panel AI 调用失败: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("panel AI 代理返回空结果")
	}

	return strings.TrimSpace(openaiResp.Choices[0].Message.Content), nil
}

func (p *PanelProvider) Test() error {
	reqBody := OpenAIRequest{
		Messages: []OpenAIMessage{{Role: "user", Content: "hi"}},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var openaiResp OpenAIResponse
	if err := p.client.Post(ctx, "/api/ai", bytes.NewBuffer(jsonData), &openaiResp); err != nil {
		return err
	}

	if len(openaiResp.Choices) == 0 {
		return fmt.Errorf("panel 响应解析失败")
	}
	return nil
}

func (p *PanelProvider) GenerateSummary(content string) (string, error) {
	prompt := resolvePrompt("", defaultSummaryPrompt) + "\n\n文章内容：\n" + content
	for i := 0; i < 3; i++ {
		result, err := p.callPanel(context.Background(), []OpenAIMessage{{Role: "user", Content: prompt}})
		if err != nil {
			return "", err
		}
		summary := strings.TrimSpace(result)
		if summary != "" && len([]rune(summary)) <= 150 {
			return summary, nil
		}
	}
	return "", fmt.Errorf("未能生成符合要求的摘要")
}

func (p *PanelProvider) GenerateAISummary(content string) (string, error) {
	prompt := resolvePrompt("", defaultAISummaryPrompt) + "\n\n文章内容：\n" + content
	for i := 0; i < 3; i++ {
		result, err := p.callPanel(context.Background(), []OpenAIMessage{{Role: "user", Content: prompt}})
		if err != nil {
			return "", err
		}
		summary := strings.TrimSpace(result)
		if summary != "" && len([]rune(summary)) <= 300 {
			return summary, nil
		}
	}
	return "", fmt.Errorf("未能生成符合要求的AI摘要")
}

func (p *PanelProvider) GenerateTitle(content string) ([]string, error) {
	for i := 0; i < 3; i++ {
		prompt := resolvePrompt("", defaultTitlePrompt) + "\n\n文章核心内容：\n" + content + "\n\n标题："
		result, err := p.callPanel(context.Background(), []OpenAIMessage{{Role: "user", Content: prompt}})
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(result)
		if title != "" && len([]rune(title)) >= 8 && len([]rune(title)) <= 30 {
			return []string{title}, nil
		}
	}
	return nil, fmt.Errorf("未能生成符合要求的标题")
}
