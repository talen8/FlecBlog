package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const getUnlimitedURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"

// qrcodeRequest 小程序码请求参数
type qrcodeRequest struct {
	Scene      string `json:"scene"`
	Page       string `json:"page,omitempty"`
	Width      int    `json:"width,omitempty"`
	EnvVersion string `json:"env_version,omitempty"`
}

// GetUnlimitedQRCode 生成小程序码（getUnlimited）
func GetUnlimitedQRCode(appID, secret, scene string, envVersion ...string) ([]byte, error) {
	accessToken, err := GetAccessToken(appID, secret)
	if err != nil {
		return nil, err
	}

	reqBody := qrcodeRequest{
		Scene:      scene,
		Page:       "pages/wechat-auth/wechat-auth",
		Width:      280,
		EnvVersion: "release", // develop: 开发版, trial: 体验版, release: 正式版
	}
	if len(envVersion) > 0 && envVersion[0] != "" {
		reqBody.EnvVersion = envVersion[0]
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败: %w", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(
		fmt.Sprintf("%s?access_token=%s", getUnlimitedURL, accessToken),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("请求微信小程序码接口失败: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信响应失败: %w", err)
	}

	// 微信返回图片时 Content-Type 为 image/jpeg，返回 JSON 时表示出错
	if resp.Header.Get("Content-Type") == "application/json" || resp.Header.Get("Content-Type") == "text/plain" {
		var errResp struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(data, &errResp); err == nil && errResp.ErrCode != 0 {
			return nil, fmt.Errorf("微信接口错误 [%d]: %s", errResp.ErrCode, errResp.ErrMsg)
		}
	}

	return data, nil
}
