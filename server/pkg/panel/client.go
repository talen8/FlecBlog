package panel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://panel.flec.top"

// apiKey 由构建参数注入，为空时自部署用户无法使用 panel 功能
var apiKey string

// Client Panel API 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New 创建 Panel 客户端
func New() *Client {
	return &Client{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Version 版本信息
type Version struct {
	ID      int    `json:"id"`
	Version string `json:"version"`
	Date    string `json:"date"`
	Changes string `json:"changes"`
}

// Announcement 官方公告
type Announcement struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

// FetchVersions 获取已启用的版本列表
func (c *Client) FetchVersions(ctx context.Context) ([]Version, error) {
	var versions []Version
	if err := c.get(ctx, "/api/versions", &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// FetchLatestVersion 获取最新版本
func (c *Client) FetchLatestVersion(ctx context.Context) (*Version, error) {
	versions, err := c.FetchVersions(ctx)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("没有可用的版本信息")
	}
	return &versions[0], nil
}

// FetchAnnouncements 获取最近公告
func (c *Client) FetchAnnouncements(ctx context.Context) ([]Announcement, error) {
	var announcements []Announcement
	if err := c.get(ctx, "/api/announcements", &announcements); err != nil {
		return nil, err
	}
	return announcements, nil
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "FlecBlog-PanelClient")
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("该功能需要官方版本")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("panel API 错误: status=%d body=%s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
