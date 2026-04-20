package linkparser

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Metadata 链接元数据
type Metadata struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Favicon     string `json:"favicon"`
}

// Parse 解析URL获取元数据
func Parse(urlStr string) (*Metadata, error) {
	// 验证URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("无效的URL: %w", err)
	}

	// 发起HTTP请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	// 设置User-Agent，避免被某些网站拦截
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Bot/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	metadata := &Metadata{
		URL: urlStr,
	}

	// 遍历HTML节点提取元数据
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil && metadata.Title == "" {
					metadata.Title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				parseMeta(n, metadata)
			case "link":
				parseLink(n, metadata)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	// 如果没有找到favicon，使用默认路径
	if metadata.Favicon == "" {
		metadata.Favicon = fmt.Sprintf("%s://%s/favicon.ico", parsedURL.Scheme, parsedURL.Host)
	}

	// 确保Favicon是完整URL
	if metadata.Favicon != "" && !strings.HasPrefix(metadata.Favicon, "http") {
		metadata.Favicon = resolveURL(parsedURL, metadata.Favicon)
	}

	return metadata, nil
}

// parseMeta 解析meta标签
func parseMeta(n *html.Node, metadata *Metadata) {
	var name, property, content string
	for _, attr := range n.Attr {
		switch attr.Key {
		case "name":
			name = attr.Val
		case "property":
			property = attr.Val
		case "content":
			content = attr.Val
		}
	}

	// Open Graph 标签
	switch property {
	case "og:title":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "og:description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	}

	// Twitter Card 标签
	switch name {
	case "twitter:title":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "twitter:description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	}
}

// parseLink 解析link标签
func parseLink(n *html.Node, metadata *Metadata) {
	var rel, href string
	for _, attr := range n.Attr {
		switch attr.Key {
		case "rel":
			rel = attr.Val
		case "href":
			href = attr.Val
		}
	}

	// 查找favicon
	if strings.Contains(rel, "icon") && metadata.Favicon == "" {
		metadata.Favicon = href
	}
}

// resolveURL 解析相对URL为绝对URL
func resolveURL(base *url.URL, ref string) string {
	refURL, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return base.ResolveReference(refURL).String()
}
