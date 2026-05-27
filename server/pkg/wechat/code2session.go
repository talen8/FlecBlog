package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

// Code2SessionResponse 微信 code2session 接口响应
type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Code2Session 使用临时 code 换取微信小程序 openid
func Code2Session(appID, secret, code string) (string, error) {
	params := url.Values{}
	params.Set("appid", appID)
	params.Set("secret", secret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(code2SessionURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("请求微信接口失败: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var result Code2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析微信响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误 [%d]: %s", result.ErrCode, result.ErrMsg)
	}

	if result.OpenID == "" {
		return "", fmt.Errorf("微信未返回 openid")
	}

	return result.OpenID, nil
}
