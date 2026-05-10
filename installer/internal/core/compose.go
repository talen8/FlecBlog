package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const composeTmpl = `services:
  postgres:
    image: postgres:15-alpine
    container_name: flec_postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: postgres
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - {{ .InstallPath }}/postgres_data:/var/lib/postgresql/data
    networks:
      - flec-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  server:
    image: {{ .GetServerImage }}
    container_name: flec_server
    restart: unless-stopped
    environment:
      DB_HOST: postgres
      DB_PORT: "5432"
      DB_NAME: postgres
      DB_USER: postgres
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
    ports:
      - "{{ .ServerPort }}:8080"
    volumes:
      - {{ .InstallPath }}/data:/app/data
    networks:
      - flec-network
    depends_on:
      postgres:
        condition: service_healthy

  blog:
    image: {{ .GetBlogImage }}
    container_name: flec_blog
    restart: unless-stopped
    environment:
      NUXT_PUBLIC_API_URL: ${API_URL}
    ports:
      - "{{ .BlogPort }}:3000"
    networks:
      - flec-network
    depends_on:
      - server

  admin:
    image: {{ .GetAdminImage }}
    container_name: flec_admin
    restart: unless-stopped
    environment:
      API_URL: ${API_URL}
    ports:
      - "{{ .AdminPort }}:4000"
    networks:
      - flec-network
    depends_on:
      - server

networks:
  flec-network:
    driver: bridge
`

func getComposePath(installPath string) string {
	return filepath.Join(installPath, "docker-compose.yml")
}

func getEnvPath(installPath string) string {
	return filepath.Join(installPath, ".env")
}

func getDataPath(installPath string) string {
	return filepath.Join(installPath, "data")
}

func getPostgresDataPath(installPath string) string {
	return filepath.Join(installPath, "postgres_data")
}

// WriteEnvFile 写入环境变量文件
func WriteEnvFile(installPath, dbPassword, jwtSecret, apiURL string) error {
	path := getEnvPath(installPath)
	var b strings.Builder

	fmt.Fprintf(&b, "DB_PASSWORD=%s\n", dbPassword)
	fmt.Fprintf(&b, "JWT_SECRET=%s\n", jwtSecret)
	fmt.Fprintf(&b, "API_URL=%s\n", apiURL)

	return os.WriteFile(path, []byte(b.String()), 0600)
}

// WriteComposeFile 写入 docker-compose.yml 文件
func WriteComposeFile(installPath string, cfg *Config) error {
	tmpl, err := template.New("compose").Parse(composeTmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return err
	}
	path := getComposePath(installPath)
	// #nosec G306 - docker-compose.yml 需要 0644 权限以便 Docker 读取
	return os.WriteFile(path, buf.Bytes(), 0644)
}

// composeFile 用于解析 docker-compose.yml
type composeFile struct {
	Services map[string]struct {
		Ports       []string          `yaml:"ports"`
		Environment map[string]string `yaml:"environment"`
	} `yaml:"services"`
}

// ParseComposePorts 解析 docker-compose.yml 中的端口配置
func ParseComposePorts(installPath string) (serverPort, blogPort, adminPort int, apiURL string, err error) {
	composePath := getComposePath(installPath)
	// #nosec G304 - 文件路径来自受信任的源（安装目录）
	data, err := os.ReadFile(composePath)
	if err != nil {
		return 0, 0, 0, "", err
	}

	var compose composeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return 0, 0, 0, "", err
	}

	serverPort = parsePort(compose.Services["server"].Ports)
	blogPort = parsePort(compose.Services["blog"].Ports)
	adminPort = parsePort(compose.Services["admin"].Ports)
	apiURL = compose.Services["blog"].Environment["NUXT_PUBLIC_API_URL"]

	// 从 .env 文件读取 API_URL（优先级更高）
	envPath := getEnvPath(installPath)
	// #nosec G304 - 文件路径来自受信任的源（安装目录）
	if envData, err := os.ReadFile(envPath); err == nil {
		envLines := strings.Split(string(envData), "\n")
		for _, envLine := range envLines {
			envLine = strings.TrimSpace(envLine)
			if strings.HasPrefix(envLine, "API_URL=") {
				apiURL = strings.TrimPrefix(envLine, "API_URL=")
				apiURL = strings.TrimSpace(apiURL)
				apiURL = strings.Trim(apiURL, "\"'")
				break
			}
		}
	}

	return serverPort, blogPort, adminPort, apiURL, nil
}

func parsePort(ports []string) int {
	if len(ports) == 0 {
		return 0
	}
	// 格式: "8080:8080" 或 "8080:8080/tcp"
	parts := strings.Split(ports[0], ":")
	if len(parts) >= 1 {
		if p, err := strconv.Atoi(parts[0]); err == nil {
			return p
		}
	}
	return 0
}

// UpdateComposePorts 更新 docker-compose.yml 中的端口配置
func UpdateComposePorts(installPath string, serverPort, blogPort, adminPort int) error {
	composePath := getComposePath(installPath)
	// #nosec G304 - 文件路径来自受信任的源（安装目录）
	data, err := os.ReadFile(composePath)
	if err != nil {
		return fmt.Errorf("读取 docker-compose.yml 失败: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("解析 docker-compose.yml 失败: %w", err)
	}

	if len(doc.Content) == 0 {
		return fmt.Errorf("docker-compose.yml 为空")
	}
	root := doc.Content[0]
	services := findMappingValue(root, "services")
	if services == nil {
		return fmt.Errorf("docker-compose.yml 中未找到 services 配置")
	}

	portMap := map[string]int{
		"server": serverPort,
		"blog":   blogPort,
		"admin":  adminPort,
	}
	for name, port := range portMap {
		svc := findMappingValue(services, name)
		if svc == nil {
			continue
		}
		ports := findMappingValue(svc, "ports")
		if ports == nil || len(ports.Content) == 0 {
			continue
		}
		firstPort := ports.Content[0]
		if firstPort.Kind == yaml.ScalarNode {
			parts := strings.SplitN(firstPort.Value, ":", 2)
			if len(parts) == 2 {
				firstPort.Value = fmt.Sprintf("%d:%s", port, parts[1])
			}
		}
	}

	// 备份原文件
	backupPath := composePath + ".backup"
	// #nosec G703 - backupPath 来自受信任的源（安装目录）
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("备份 docker-compose.yml 失败: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&doc); err != nil {
		return fmt.Errorf("序列化 docker-compose.yml 失败: %w", err)
	}
	_ = encoder.Close()

	// #nosec G306 - docker-compose.yml 需要 0644 权限以便 Docker 读取
	if err := os.WriteFile(composePath, buf.Bytes(), 0644); err != nil {
		_ = os.Rename(backupPath, composePath)
		return fmt.Errorf("写入 docker-compose.yml 失败: %w", err)
	}

	return nil
}

// findMappingValue 在 YAML Mapping 节点中查找指定 key 对应的 value 节点
func findMappingValue(node *yaml.Node, key string) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Kind == yaml.ScalarNode && node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// UpdateAPIURL 更新环境变量中的 API_URL
func UpdateAPIURL(installPath, apiURL string) error {
	envPath := getEnvPath(installPath)
	// #nosec G304 - 文件路径来自受信任的源（安装目录）
	envData, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("读取 .env 文件失败: %w", err)
	}

	envContent := string(envData)
	lines := strings.Split(envContent, "\n")
	found := false

	for i, line := range lines {
		if strings.HasPrefix(line, "API_URL=") {
			lines[i] = fmt.Sprintf("API_URL=%s", apiURL)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("API_URL=%s", apiURL))
	}

	// #nosec G703 - envPath 来自受信任的源（安装目录）
	return os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0600)
}
