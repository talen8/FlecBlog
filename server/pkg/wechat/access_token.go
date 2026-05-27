package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token" //nolint:gosec // 不是凭据，是 API 地址

// accessTokenResponse 微信 access_token 响应
type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// tokenCache access_token 缓存
type tokenCache struct {
	mu        sync.RWMutex
	token     string
	expiresAt time.Time
	appID     string
	secret    string
}

var cache = &tokenCache{}

// GetAccessToken 获取微信接口调用凭据（带内存缓存）
func GetAccessToken(appID, secret string) (string, error) {
	cache.mu.RLock()
	if cache.token != "" && time.Now().Before(cache.expiresAt) && cache.appID == appID && cache.secret == secret {
		token := cache.token
		cache.mu.RUnlock()
		return token, nil
	}
	cache.mu.RUnlock()

	cache.mu.Lock()
	defer cache.mu.Unlock()

	// 双重检查
	if cache.token != "" && time.Now().Before(cache.expiresAt) && cache.appID == appID && cache.secret == secret {
		return cache.token, nil
	}

	params := fmt.Sprintf("?grant_type=client_credential&appid=%s&secret=%s", appID, secret)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(accessTokenURL + params)
	if err != nil {
		return "", fmt.Errorf("请求微信 access_token 接口失败: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var result accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析微信 access_token 响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误 [%d]: %s", result.ErrCode, result.ErrMsg)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("微信未返回 access_token")
	}

	// 提前 5 分钟过期
	cache.token = result.AccessToken
	cache.expiresAt = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)
	cache.appID = appID
	cache.secret = secret

	return result.AccessToken, nil
}
