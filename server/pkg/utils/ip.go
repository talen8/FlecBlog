package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const ipAPIURL = "http://ip-api.com/json/%s?lang=zh-CN"

var httpClient = &http.Client{Timeout: 5 * time.Second}

type ipAPIResponse struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
	Message    string `json:"message"`
}

// InitIPSearcher 初始化 IP 地址搜索器（在线模式无需初始化）
func InitIPSearcher(_ string) error {
	return nil
}

// CloseIPSearcher 关闭 IP 搜索器（在线模式无需关闭）
func CloseIPSearcher() {}

// GetIPLocation 获取 IP 地址的地理位置
func GetIPLocation(ip string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return "未知"
	}

	if ipAddr.IsLoopback() || ipAddr.IsPrivate() {
		return "本地"
	}

	resp, err := httpClient.Get(fmt.Sprintf(ipAPIURL, ip))
	if err != nil {
		return "未知"
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "未知"
	}

	if result.Status != "success" {
		return "未知"
	}

	var parts []string
	if result.Country != "" {
		parts = append(parts, result.Country)
	}
	if result.RegionName != "" && result.RegionName != result.Country {
		parts = append(parts, result.RegionName)
	}
	if result.City != "" && result.City != result.RegionName {
		parts = append(parts, result.City)
	}

	if len(parts) == 0 {
		return "未知"
	}
	return strings.Join(parts, " ")
}
