package ai

import (
	"flec_blog/config"
	"strings"
)

// OpenAIClient OpenAI 兼容 API 客户端
type OpenAIClient struct {
	BaseURL         string
	APIKey          string
	Model           string
	SummaryPrompt   string
	AISummaryPrompt string
	TitlePrompt     string
}

// OpenAIRequest OpenAI API 请求结构
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// OpenAIMessage 消息结构
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse OpenAI API 响应结构
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAIClient 创建 OpenAI 兼容客户端
func NewOpenAIClient(cfg *config.AIConfig) *OpenAIClient {
	if cfg == nil {
		return &OpenAIClient{}
	}

	return &OpenAIClient{
		BaseURL:         strings.TrimRight(cfg.BaseURL, "/"),
		APIKey:          cfg.APIKey,
		Model:           cfg.Model,
		SummaryPrompt:   cfg.SummaryPrompt,
		AISummaryPrompt: cfg.AISummaryPrompt,
		TitlePrompt:     cfg.TitlePrompt,
	}
}
