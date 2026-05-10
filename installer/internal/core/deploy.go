package core

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/talen8/flecblog/installer/internal/ui"
)

// Deployer 部署器
type Deployer struct {
	cfg *Config
}

// NewDeployer 创建部署器实例
func NewDeployer(cfg *Config) *Deployer {
	return &Deployer{cfg: cfg}
}

// Deploy 执行 Docker 部署
func (d *Deployer) Deploy() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"创建目录", d.createDirs},
		{"生成配置文件", d.generateDockerConfig},
		{"拉取镜像", d.pullImages},
		{"启动服务", d.startCompose},
		{"等待服务就绪", d.waitForReady},
	}

	bar := ui.NewProgress(len(steps), "Docker 部署")
	for _, step := range steps {
		bar.Describe(step.name)
		if err := step.fn(); err != nil {
			bar.Fail()
			return fmt.Errorf("%s: %w", step.name, err)
		}
		bar.Add(1)
	}
	bar.Finish()
	return nil
}

// Upgrade 执行升级
func (d *Deployer) Upgrade() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"拉取最新镜像", d.pullImages},
		{"重启服务", d.startCompose},
		{"等待服务就绪", d.waitForReady},
	}

	bar := ui.NewProgress(len(steps), "Docker 升级")
	for _, step := range steps {
		bar.Describe(step.name)
		if err := step.fn(); err != nil {
			bar.Fail()
			return fmt.Errorf("%s: %w", step.name, err)
		}
		bar.Add(1)
	}
	bar.Finish()
	return nil
}

// Rollback 执行回滚
func (d *Deployer) Rollback(purge bool) {
	args := []string{"compose", "down"}
	if purge {
		args = append(args, "-v")
	}
	cmd := DockerCmd(args...)
	cmd.Dir = d.cfg.InstallPath
	_ = cmd.Run()

	if purge {
		_ = os.Remove(getEnvPath(d.cfg.InstallPath))
		_ = os.Remove(getComposePath(d.cfg.InstallPath))

		if configPath, err := GetConfigPath(); err == nil {
			_ = os.Remove(configPath)
		}

		if err := os.RemoveAll(getDataPath(d.cfg.InstallPath)); err != nil {
			fmt.Printf("  %s警告: %s%s\n", ui.Yellow, err, ui.Reset)
		}
		if err := os.RemoveAll(getPostgresDataPath(d.cfg.InstallPath)); err != nil {
			fmt.Printf("  %s警告: %s%s\n", ui.Yellow, err, ui.Reset)
		}
	}
}

func (d *Deployer) createDirs() error {
	return os.MkdirAll(getDataPath(d.cfg.InstallPath), 0750)
}

func (d *Deployer) generateDockerConfig() error {
	if err := d.writeEnvFile(); err != nil {
		return fmt.Errorf("生成 .env: %w", err)
	}
	if err := d.writeComposeFile(); err != nil {
		return fmt.Errorf("生成 docker-compose.yml: %w", err)
	}
	return nil
}

func (d *Deployer) pullImages() error {
	images := []string{
		"postgres:15-alpine",
		d.cfg.GetServerImage(),
		d.cfg.GetBlogImage(),
		d.cfg.GetAdminImage(),
	}
	for _, img := range images {
		if err := d.pullImageWithRetry(img, 10); err != nil {
			return fmt.Errorf("拉取 %s 失败: %w", img, err)
		}
	}
	return nil
}

func (d *Deployer) pullImageWithRetry(image string, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(5 * time.Second)
		}
		cmd := DockerCmd("pull", image)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return lastErr
}

func (d *Deployer) startCompose() error {
	cmd := DockerCmd("compose", "up", "-d")
	cmd.Dir = d.cfg.InstallPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *Deployer) waitForReady() error {
	url := fmt.Sprintf("http://localhost:%d/", d.cfg.ServerPort)
	timeout := 120 * time.Second
	deadline := time.Now().Add(timeout)

	fmt.Println("  等待服务初始化（首次启动可能需要几分钟）...")

	attempts := 0
	for time.Now().Before(deadline) {
		attempts++
		// #nosec G107 - URL 是本地服务地址，来自配置
		resp, err := http.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		if attempts%5 == 0 {
			fmt.Printf("  服务启动中... (已等待 %.0f 秒)\n", time.Since(deadline.Add(-timeout)).Seconds())
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("\n  诊断信息:")
	fmt.Printf("    检查 URL: %s\n", url)
	fmt.Println("    查看日志: docker logs flec_server")
	fmt.Println("    查看状态: docker ps -a")

	return fmt.Errorf("服务在 %s 内未就绪，请运行 'docker logs flec_server' 查看日志", timeout)
}

func (d *Deployer) writeEnvFile() error {
	return WriteEnvFile(d.cfg.InstallPath, d.cfg.DBPassword, d.cfg.JWTSecret, d.cfg.APIURL)
}

func (d *Deployer) writeComposeFile() error {
	return WriteComposeFile(d.cfg.InstallPath, d.cfg)
}

func DockerCmd(args ...string) *exec.Cmd {
	// #nosec G204 - docker 命令是硬编码的，args 来自受信任的源
	return exec.Command("docker", args...)
}
