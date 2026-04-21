package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultSummaryPrompt   = "你是一位博客作者，请根据文章内容生成中文摘要。要求：1. 以作者视角介绍文章；2. 控制在50到100字之间；3. 只输出摘要正文，不要附加解释；4. 内容简洁准确，覆盖核心信息；5. 不要空泛拔高文章意义。"
	defaultAISummaryPrompt = "你是一名AI助手，请根据文章内容生成中文总结。要求：1. 以旁观者视角总结并推荐文章；2. 控制在150到200字之间；3. 只输出总结正文，不要附加解释；4. 保持语言自然、信息完整；5. 采用“这篇文章...”的表述方式。"
	defaultTitlePrompt     = "你是一位资深技术作者，请根据文章内容生成1个中文标题。要求：1. 突出主题亮点和核心价值；2. 控制在15到25字之间；3. 尽量不用标点符号；4. 只返回标题本身，不要解释。"
)

func resolvePrompt(customPrompt, defaultPrompt string) string {
	if strings.TrimSpace(customPrompt) != "" {
		return strings.TrimSpace(customPrompt)
	}
	return defaultPrompt
}

// callOpenAI 通用OpenAI API调用函数
func (c *OpenAIClient) callOpenAI(prompt string) (string, error) {
	reqBody := OpenAIRequest{
		Model: c.Model,
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	fullURL := c.BaseURL + "/chat/completions"

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("API 返回空结果")
	}

	result := strings.TrimSpace(openaiResp.Choices[0].Message.Content)
	if result == "" {
		return "", fmt.Errorf("生成的内容为空")
	}

	return result, nil
}

// Test 测试AI配置是否可用
func (c *OpenAIClient) Test() error {
	client := &http.Client{Timeout: 15 * time.Second}

	reqBody := OpenAIRequest{
		Model: c.Model,
		Messages: []OpenAIMessage{
			{Role: "user", Content: "hi"},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API 返回错误 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil || len(openaiResp.Choices) == 0 {
		return fmt.Errorf("响应解析失败，请检查模型名称是否正确")
	}
	return nil
}

// GenerateSummary 生成文章摘要（50-100字，创作者角度）
func (c *OpenAIClient) GenerateSummary(content string) (string, error) {
	prompt := resolvePrompt(c.SummaryPrompt, defaultSummaryPrompt) + "\n\n文章内容：\n" + content
	return c.callOpenAI(prompt)
}

// GenerateAISummary 生成AI摘要（150-200字，旁观者角度）
func (c *OpenAIClient) GenerateAISummary(content string) (string, error) {
	prompt := resolvePrompt(c.AISummaryPrompt, defaultAISummaryPrompt) + "\n\n文章内容：\n" + content
	return c.callOpenAI(prompt)
}

// GenerateTitle 生成标题
func (c *OpenAIClient) GenerateTitle(content string) ([]string, error) {
	for i := 0; i < 3; i++ { // 最多重试3次
		prompt := resolvePrompt(c.TitlePrompt, defaultTitlePrompt) + "\n\n文章核心内容：\n" + content + "\n\n标题："
		result, err := c.callOpenAI(prompt)
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
