package panel

import (
	"bytes"
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

// currentClientKey 用于跨实例共享 client_key
var currentClientKey string

// Client Panel API 客户端
type Client struct {
	BaseURL      string
	HTTPClient   *http.Client
	ExtraHeaders map[string]string
}

// New 创建 Panel 客户端
func New() *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetClientKey 设置已有的 client_key
func (c *Client) SetClientKey(key string) {
	if c.ExtraHeaders == nil {
		c.ExtraHeaders = make(map[string]string)
	}
	c.ExtraHeaders["X-Client-Key"] = key
	currentClientKey = key
}

// ClientKey 获取当前 client_key
func (c *Client) ClientKey() string {
	if c.ExtraHeaders == nil {
		return ""
	}
	return c.ExtraHeaders["X-Client-Key"]
}

// CurrentClientKey 获取全局最新的 client_key
func CurrentClientKey() string {
	return currentClientKey
}

// Register 注册或心跳
func (c *Client) Register(ctx context.Context, siteURL, version string) error {
	if c.ExtraHeaders == nil {
		c.ExtraHeaders = make(map[string]string)
	}
	c.ExtraHeaders["X-Site-Url"] = siteURL
	c.ExtraHeaders["X-Version"] = version

	key := c.ExtraHeaders["X-Client-Key"]
	if apiKey == "" && key == "" {
		return nil
	}

	body := map[string]any{"site_url": siteURL, "version": version, "client_key": key, "api_key": apiKey}

	jsonData, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.BaseURL+"/api/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "FlecBlog-PanelClient")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("request failed: status=%d body=%s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		ClientKey string `json:"client_key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil // 心跳返回 {status:"ok"}，无 client_key，正常
	}
	if result.ClientKey != "" {
		c.ExtraHeaders["X-Client-Key"] = result.ClientKey
		currentClientKey = result.ClientKey
	}
	return nil
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

func (c *Client) setCommonHeaders(req *http.Request) {
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}
	for k, v := range c.ExtraHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", "FlecBlog-PanelClient")
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	c.setCommonHeaders(req)
	return c.do(req, out)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setCommonHeaders(req)
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.HTTPClient.Do(req) //nolint:gosec // URL 由内部拼接 c.BaseURL+path，BaseURL 为预设常量
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return fmt.Errorf("%s", errResp.Error)
		}
		return fmt.Errorf("request failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
