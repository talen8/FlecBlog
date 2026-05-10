package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/talen8/flecblog/installer/internal/core"
	"github.com/talen8/flecblog/installer/internal/ui"
)

// flagUpgradeCheck 检查更新标志
var flagUpgradeCheck bool

// upgradeCmd 升级命令
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级 FlecBlog",
	Long:  "升级已安装的 FlecBlog 实例到最新版本。",
	RunE:  runUpgrade,
}

func init() {
	upgradeCmd.Flags().BoolVar(&flagUpgradeCheck, "check", false, "检查更新，不执行升级")
	rootCmd.AddCommand(upgradeCmd)
}

// runUpgrade 执行升级命令
func runUpgrade(cmd *cobra.Command, args []string) error {
	if flagUpgradeCheck {
		return checkUpgrade()
	}
	return doUpgrade()
}

// checkUpgrade 检查更新
func checkUpgrade() error {
	ui.PrintBanner()

	cfg, _ := LoadConfigOrExit()

	currentVersion := cfg.Version
	if currentVersion == "" {
		currentVersion = "unknown"
	}

	fmt.Printf("  当前官方组件版本: %s\n", currentVersion)
	fmt.Println()

	// 检查官方组件更新
	pc := core.NewPanelClient()
	latestVersion, err := pc.GetLatestVersion()
	switch {
	case err != nil:
		fmt.Printf("  %s检查官方组件更新失败: %v%s\n", ui.Yellow, err, ui.Reset)
	case latestVersion != currentVersion:
		fmt.Printf("  %s官方组件更新:%s\n", ui.Cyan, ui.Reset)
		fmt.Printf("    当前版本: %s\n", currentVersion)
		fmt.Printf("    最新版本: %s\n", latestVersion)
		fmt.Println()
	default:
		fmt.Printf("  %s官方组件: 已是最新版本 (%s)%s\n", ui.Green, currentVersion, ui.Reset)
		fmt.Println()
	}

	return nil
}

// doUpgrade 执行升级
func doUpgrade() error {
	ui.PrintBanner()

	cfg, configPath := LoadConfigOrExit()
	installPath := cfg.InstallPath
	fmt.Printf("  安装路径: %s (%s)\n", installPath, configPath)

	mode := core.DetectMode(installPath)
	if mode == "" {
		return fmt.Errorf("未在 %s 检测到 FlecBlog 安装", installPath)
	}

	fmt.Printf("  检测到 %s 模式安装\n", core.ModeLabel(mode))

	pc := core.NewPanelClient()
	latestVersion, err := pc.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("获取最新版本失败: %w", err)
	}

	currentVersion := cfg.Version
	if currentVersion == "" {
		currentVersion = "unknown"
	}

	fmt.Printf("  当前版本: %s\n", currentVersion)
	fmt.Printf("  最新版本: %s\n", latestVersion)

	if latestVersion == currentVersion {
		fmt.Printf("\n  %s已是最新版本，无需升级%s\n", ui.Green, ui.Reset)
		return nil
	}

	fmt.Print("\n  确认升级到最新版本? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "" && confirm != "y" && confirm != "Y" {
		fmt.Println("  已取消升级")
		return nil
	}

	cfg.Version = latestVersion
	if _, err := cfg.Save(); err != nil {
		fmt.Printf("  %s警告: 保存配置失败: %v%s\n", ui.Yellow, err, ui.Reset)
	}

	if err := core.WriteComposeFile(installPath, cfg); err != nil {
		return fmt.Errorf("更新 docker-compose.yml 失败: %w", err)
	}

	d := core.NewDeployer(cfg)
	if err := d.Upgrade(); err != nil {
		return fmt.Errorf("升级失败: %w", err)
	}

	fmt.Println("\n  升级完成")
	return nil
}
