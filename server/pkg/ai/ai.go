package ai

import (
	"flec_blog/config"
	"fmt"
	"os"
)

// Provider AI服务提供商接口
type Provider interface {
	GenerateSummary(content string) (string, error)
	GenerateAISummary(content string) (string, error)
	GenerateTitle(content string) ([]string, error)
	Test() error
}

// GetProvider 根据配置获取AI服务提供商
func GetProvider(cfg *config.AIConfig) (Provider, error) {
	hasLocalCfg := cfg != nil && cfg.BaseURL != "" && cfg.APIKey != "" && cfg.Model != ""

	if hasLocalCfg {
		return NewOpenAIClientWithConfig(cfg), nil
	}

	if os.Getenv("PANEL_API_KEY") != "" {
		return NewPanelProvider(), nil
	}

	return nil, fmt.Errorf("AI配置不完整或非官方版本")
}
