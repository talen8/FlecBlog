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

// truncateContent 截断内容到指定长度
func truncateContent(content string, maxLength int) string {
	if len(content) > maxLength {
		return content[:maxLength] + "\n\n... (内容已截断)"
	}
	return content
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	const fixedPrompt = "你是一个博客作者，需要为自己的文章生成摘要。严格要求：1.摘要字数必须控制在50-100字之间，绝不能超出100字上限；2.以创作者角度介绍这一篇文章；3.只生成摘要内容，不要任何额外说明；4.用中文作答，去除特殊字符与冗余空格；5.请生成简洁、精准且涵盖核心信息的摘要；6.不要总结有什么意义或者价值！！。"

	content = truncateContent(content, 10000)

	prompt := fixedPrompt + "\n\n文章内容：\n" + content
	return c.callOpenAI(prompt)
}

// GenerateAISummary 生成AI摘要（150-200字，旁观者角度）
func (c *OpenAIClient) GenerateAISummary(content string) (string, error) {
	const fixedPrompt = "你是一个AI助手，需要为博客文章生成摘要。严格要求：1.摘要字数必须控制在150-200汉字（含标点符号）之间，绝不能超出200字上限；2.以旁观者AI角度总结介绍并推荐这一篇文章；3.只生成摘要内容，不要任何额外说明；4.用中文作答，去除特殊字符与冗余空格；5.输出格式固定为：这篇文章[摘要内容]；6.需严格校验最终摘要字数，确保完全符合区间要求，请生成简洁、精准且涵盖核心信息的摘要。"

	content = truncateContent(content, 10000)

	prompt := fixedPrompt + "\n\n文章内容：\n" + content
	return c.callOpenAI(prompt)
}

// GenerateTitle 生成标题
func (c *OpenAIClient) GenerateTitle(content string) ([]string, error) {
	const fixedPrompt = "你是一位资深技术博客作者，需要为文章创作一个吸引人的标题。要求 1.生成1个最具吸引力的标题 2.突出技术亮点和核心价值 3.使用简洁有力的语言，字数控制在15-25字 4.采用纯文字表达，禁止使用任何标点符号包括感叹号、问号、冒号、逗号等 5.直接返回标题内容，不要任何前缀、后缀或解释说明"

	content = truncateContent(content, 3000)

	for i := 0; i < 3; i++ { // 最多重试3次
		prompt := fixedPrompt + "\n\n文章核心内容：\n" + content + "\n\n标题："
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
