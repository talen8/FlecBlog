package core

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// PromptDeployMode 提示用户选择部署方式
func PromptDeployMode(reader *bufio.Reader) string {
	fmt.Println("  部署方式:")
	fmt.Println("    [1] Docker 部署 (推荐)")
	fmt.Println("    [2] 裸机部署 (systemd + Nginx) [暂不可用]")

	for {
		fmt.Print("  请选择 [1]: ")
		s, _ := reader.ReadString('\n')
		s = strings.TrimSpace(s)

		if s == "" || s == "1" {
			return "docker"
		}

		if s == "2" {
			fmt.Println("  裸机部署暂不可用，请选择 Docker 部署")
			continue
		}

		fmt.Println("  无效的选择，请输入 1 或 2")
	}
}

// PromptConfig 收集部署配置
func PromptConfig(reader *bufio.Reader, deployMode string) (*Config, error) {
	var err error

	cfg := NewConfig()
	cfg.DeployMode = deployMode

	fmt.Println("\n  端口配置（直接回车使用默认值）：")
	PromptPorts(reader, cfg)

	cfg.DBPassword, err = GenerateSecret()
	if err != nil {
		return nil, fmt.Errorf("生成数据库密码失败: %w", err)
	}

	cfg.JWTSecret, err = GenerateSecret()
	if err != nil {
		return nil, fmt.Errorf("生成 JWT 密钥失败: %w", err)
	}

	fmt.Printf("  安装路径 [%s]: ", DefaultInstallPath)
	input, _ := reader.ReadString('\n')
	if ip := strings.TrimSpace(input); ip != "" {
		cfg.InstallPath = ip
	}

	cfg.BuildAPIURL()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// PromptPorts 配置端口
func PromptPorts(reader *bufio.Reader, cfg *Config) {
	cfg.ServerPort = promptInt(reader, "Server 端口", cfg.ServerPort)
	cfg.BlogPort = promptInt(reader, "Blog 端口", cfg.BlogPort)
	cfg.AdminPort = promptInt(reader, "Admin 端口", cfg.AdminPort)
}

// promptInt 提示输入整数
func promptInt(r *bufio.Reader, label string, defaultVal int) int {
	fmt.Printf("  %s [%d]: ", label, defaultVal)
	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(input)
	if err != nil || v <= 0 || v > 65535 {
		fmt.Printf("  无效端口号，使用默认值 %d\n", defaultVal)
		return defaultVal
	}
	return v
}
