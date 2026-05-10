package core

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const DefaultInstallPath = "/opt/flecblog"

type Config struct {
	DeployMode  string `json:"deployMode"`
	InstallPath string `json:"installPath"`
	Version     string `json:"version"`
	ServerPort  int    `json:"serverPort"`
	BlogPort    int    `json:"blogPort"`
	AdminPort   int    `json:"adminPort"`
	Domain      string `json:"-"`
	DBPassword  string `json:"-"`
	JWTSecret   string `json:"-"`
	APIURL      string `json:"-"`
}

func NewConfig() *Config {
	return &Config{
		DeployMode:  "docker",
		InstallPath: DefaultInstallPath,
		ServerPort:  8080,
		BlogPort:    3000,
		AdminPort:   4000,
	}
}

func (c *Config) GetBlogImage() string {
	return fmt.Sprintf("talen8/flec-blog:%s", c.getVersionOrLatest())
}

func (c *Config) GetServerImage() string {
	return fmt.Sprintf("talen8/flec-server:%s", c.getVersionOrLatest())
}

func (c *Config) GetAdminImage() string {
	return fmt.Sprintf("talen8/flec-admin:%s", c.getVersionOrLatest())
}

func (c *Config) getVersionOrLatest() string {
	if c.Version != "" {
		return c.Version
	}
	return "latest"
}

func (c *Config) Validate() error {
	if c.DeployMode == "" {
		return fmt.Errorf("部署方式不能为空")
	}
	if c.Domain == "" {
		return fmt.Errorf("无法获取服务器地址")
	}
	if len(c.DBPassword) < 8 {
		return fmt.Errorf("数据库密码长度至少 8 位")
	}
	if c.InstallPath == "" {
		return fmt.Errorf("安装路径不能为空")
	}
	if c.DeployMode == "bare-metal" {
		return fmt.Errorf("裸机部署暂不可用，请使用 Docker 模式")
	}
	used := map[int]bool{}
	for _, p := range []int{c.ServerPort, c.BlogPort, c.AdminPort} {
		if used[p] {
			return fmt.Errorf("端口 %d 被重复使用", p)
		}
		used[p] = true
	}
	return nil
}

func (c *Config) BuildAPIURL() {
	if c.Domain == "" {
		client := &http.Client{Timeout: 5 * time.Second}
		services := []string{
			"https://icanhazip.com",
			"https://ipinfo.io/ip",
			"https://checkip.amazonaws.com",
		}
		for _, url := range services {
			resp, err := client.Get(url)
			if err != nil {
				continue
			}
			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				continue
			}
			if ip := strings.TrimSpace(string(body)); ip != "" {
				c.Domain = ip
				break
			}
		}
	}
	if c.Domain == "" {
		c.Domain = "localhost"
	}
	if net.ParseIP(c.Domain) != nil {
		c.APIURL = fmt.Sprintf("http://%s:%d/api/v1", c.Domain, c.ServerPort)
	} else {
		c.APIURL = fmt.Sprintf("https://%s/api/v1", c.Domain)
	}
}

func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %w", err)
	}
	return filepath.Join(homeDir, ".flecblog", "config.json"), nil
}

func (c *Config) Save() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return "", fmt.Errorf("创建配置目录失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return "", fmt.Errorf("保存配置失败: %w", err)
	}

	return configPath, nil
}

func LoadConfig() (*Config, string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, "", err
	}

	cfg, err := loadFromPath(configPath)
	if err != nil {
		return nil, "", fmt.Errorf("未找到安装配置")
	}

	return cfg, configPath, nil
}

func loadFromPath(configPath string) (*Config, error) {
	// #nosec G304 - configPath 来自受信任的源（用户主目录）
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ModeLabel(mode string) string {
	if mode == "bare-metal" {
		return "裸机部署 (systemd + Nginx)"
	}
	return "Docker 部署"
}

func DetectMode(installPath string) string {
	cfg, _, _ := LoadConfig()
	if cfg != nil && cfg.DeployMode != "" {
		return cfg.DeployMode
	}

	if _, err := os.Stat(getComposePath(installPath)); err == nil {
		return "docker"
	}
	return ""
}
