package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const PanelURL = "https://panel.flec.top"

type PanelClient struct {
	baseURL string
	client  *http.Client
}

func NewPanelClient() *PanelClient {
	return &PanelClient{
		baseURL: PanelURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *PanelClient) GetLatestVersion() (string, error) {
	resp, err := c.client.Get(c.baseURL + "/api/versions")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("panel 返回错误状态码: %d", resp.StatusCode)
	}

	var response struct {
		Versions []struct {
			Version string `json:"version"`
		} `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Versions) == 0 || response.Versions[0].Version == "" {
		return "", fmt.Errorf("panel 未返回版本信息")
	}
	return response.Versions[0].Version, nil
}
