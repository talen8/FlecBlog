package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/talen8/flecblog/installer/internal/core"
	"github.com/talen8/flecblog/installer/internal/ui"
)

// installCmd 安装命令
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装 FlecBlog",
	Long:  "交互式安装 FlecBlog 博客系统，通过问答引导完成配置和部署。",
	RunE:  runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

// runInstall 执行安装
func runInstall(cmd *cobra.Command, args []string) error {
	ui.PrintBanner()

	// 检测是否已安装
	if cfg, configPath, err := core.LoadConfig(); err == nil && cfg != nil {
		fmt.Printf("  检测到已安装的 FlecBlog\n")
		fmt.Printf("  配置文件: %s\n", configPath)
		fmt.Printf("  安装路径: %s\n", cfg.InstallPath)
		fmt.Printf("  部署模式: %s\n", core.ModeLabel(cfg.DeployMode))
		fmt.Println("\n  如需重新安装，请先执行 'flecb uninstall' 卸载现有实例")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)

	// 选择部署方式
	deployMode := core.PromptDeployMode(reader)

	// 环境检测
	env, err := core.DetectDocker()
	if err != nil {
		_ = core.PromptDockerInstallation()
		return fmt.Errorf("环境检测失败: %w", err)
	}
	ui.PrintDockerEnv(env.DockerVer, env.ComposeVer)

	// 收集配置
	cfg, err := core.PromptConfig(reader, deployMode)
	if err != nil {
		return fmt.Errorf("配置错误: %w", err)
	}

	// 端口冲突检测
	for {
		ports := core.Ports{
			Server: cfg.ServerPort,
			Blog:   cfg.BlogPort,
			Admin:  cfg.AdminPort,
		}
		conflicts := core.GetConflictPorts(ports)
		if len(conflicts) == 0 {
			break
		}
		fmt.Println("\n  以下端口已被占用:")
		for _, port := range conflicts {
			fmt.Printf("    - %d\n", port)
		}
		fmt.Println("\n  请重新设置端口")
		core.PromptPorts(reader, cfg)
	}

	// 确认安装
	fmt.Println("\n  配置确认:")
	fmt.Printf("    部署方式: %s\n", core.ModeLabel(cfg.DeployMode))
	fmt.Printf("    安装路径: %s\n", cfg.InstallPath)
	fmt.Printf("    端口: Server=%d, Blog=%d, Admin=%d\n", cfg.ServerPort, cfg.BlogPort, cfg.AdminPort)
	fmt.Print("\n  确认安装? [Y/n]: ")
	var confirm string
	_, _ = fmt.Scanln(&confirm)
	if confirm != "" && confirm != "y" && confirm != "Y" {
		fmt.Println("  已取消安装")
		return nil
	}

	// 获取最新系统版本
	pc := core.NewPanelClient()
	systemVersion, _ := pc.GetLatestVersion()
	if systemVersion == "" {
		systemVersion = "latest"
	}
	cfg.Version = systemVersion

	// 执行部署
	d := core.NewDeployer(cfg)
	if err := d.Deploy(); err != nil {
		fmt.Fprintf(os.Stderr, "\n  %s安装失败: %v%s\n", ui.Red, err, ui.Reset)
		fmt.Println("  正在清理...")
		d.Rollback(true)
		return err
	}

	// 保存配置
	if configPath, err := cfg.Save(); err == nil {
		fmt.Printf("  配置已保存: %s\n", configPath)
	} else {
		fmt.Fprintf(os.Stderr, "\n  %s注意: 配置保存失败: %v%s\n", ui.Yellow, err, ui.Reset)
	}

	ui.PrintSuccess(cfg.Domain, cfg.BlogPort, cfg.AdminPort, cfg.ServerPort, cfg.InstallPath)
	return nil
}
