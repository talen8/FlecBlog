package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// GenerateExcerpt 生成包含关键词的文章摘录
func GenerateExcerpt(content, keyword string, maxLength int) string {
	if maxLength == 0 {
		maxLength = 200
	}

	// 移除 Markdown 标记
	plainText := stripMarkdown(content)

	// 如果正文为空，返回空字符串
	if plainText == "" {
		return ""
	}

	// 查找关键词在正文中的位置（不区分大小写）
	lowerContent := strings.ToLower(plainText)
	lowerKeyword := strings.ToLower(keyword)
	keywordPos := strings.Index(lowerContent, lowerKeyword)

	// 如果没有找到关键词，返回空字符串
	if keywordPos == -1 {
		return ""
	}

	// 计算提取片段的起始和结束位置
	keywordLength := utf8.RuneCountInString(keyword)
	contextSize := (maxLength - keywordLength) / 2

	// 转换为 rune 数组以正确处理中文字符
	runes := []rune(plainText)
	totalLength := len(runes)

	// 计算关键词在 rune 数组中的位置
	runeKeywordPos := utf8.RuneCountInString(plainText[:keywordPos])

	// 计算起始和结束位置
	start := runeKeywordPos - contextSize
	if start < 0 {
		start = 0
	}

	end := runeKeywordPos + keywordLength + contextSize
	if end > totalLength {
		end = totalLength
	}

	// 如果片段太长，调整范围
	if end-start > maxLength {
		// 优先保证关键词前后各有一定的上下文
		halfMax := maxLength / 2
		start = runeKeywordPos - halfMax
		end = runeKeywordPos + keywordLength + halfMax

		if start < 0 {
			start = 0
			end = min(maxLength, totalLength)
		}
		if end > totalLength {
			end = totalLength
			start = max(0, totalLength-maxLength)
		}
	}

	// 提取片段
	excerpt := string(runes[start:end])
	excerpt = strings.TrimSpace(excerpt)

	// 添加省略号
	if start > 0 {
		excerpt = "..." + excerpt
	}
	if end < totalLength {
		excerpt += "..."
	}

	return excerpt
}

// stripMarkdown 移除 Markdown 标记，保留纯文本
func stripMarkdown(content string) string {
	// 移除代码块
	re := regexp.MustCompile("(?s)```.*?```")
	content = re.ReplaceAllString(content, "")

	// 移除行内代码
	re = regexp.MustCompile("`[^`]+`")
	content = re.ReplaceAllString(content, "")

	// 移除图片
	re = regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`)
	content = re.ReplaceAllString(content, "$1")

	// 移除链接，保留文本
	re = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	content = re.ReplaceAllString(content, "$1")

	// 移除标题符号
	re = regexp.MustCompile(`(?m)^#+\s+`)
	content = re.ReplaceAllString(content, "")

	// 移除粗体和斜体
	re = regexp.MustCompile(`\*\*([^\*]+)\*\*`)
	content = re.ReplaceAllString(content, "$1")
	re = regexp.MustCompile(`\*([^\*]+)\*`)
	content = re.ReplaceAllString(content, "$1")

	// 移除删除线
	re = regexp.MustCompile(`~~([^~]+)~~`)
	content = re.ReplaceAllString(content, "$1")

	// 移除列表符号
	re = regexp.MustCompile(`(?m)^[\*\-\+]\s+`)
	content = re.ReplaceAllString(content, "")

	// 移除有序列表
	re = regexp.MustCompile(`(?m)^\d+\.\s+`)
	content = re.ReplaceAllString(content, "")

	// 移除引用符号
	re = regexp.MustCompile(`(?m)^>\s+`)
	content = re.ReplaceAllString(content, "")

	// 移除水平线
	re = regexp.MustCompile(`(?m)^[-\*_]{3,}$`)
	content = re.ReplaceAllString(content, "")

	// 移除 HTML 标签
	re = regexp.MustCompile(`<[^>]+>`)
	content = re.ReplaceAllString(content, "")

	// 合并多个空白为一个空格
	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	// 移除首尾空白
	content = strings.TrimSpace(content)

	// 返回处理后的纯文本
	return content
}
