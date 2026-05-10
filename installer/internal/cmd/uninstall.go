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

// flagPurge 彻底删除标志
var flagPurge bool

// uninstallCmd 卸载命令
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "卸载 FlecBlog",
	Long:  "卸载已安装的 FlecBlog 实例，停止服务并清理配置文件。",
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVar(&flagPurge, "purge", false, "删除所有内容包括配置文件和数据")
	rootCmd.AddCommand(uninstallCmd)
}

// runUninstall 执行卸载
func runUninstall(cmd *cobra.Command, args []string) error {
	ui.PrintBanner()

	cfg, configPath := LoadConfigOrExit()
	installPath := cfg.InstallPath
	fmt.Printf("  安装路径: %s (%s)\n", installPath, configPath)

	mode := core.DetectMode(installPath)
	if mode == "" {
		return fmt.Errorf("未在 %s 检测到 FlecBlog 安装", installPath)
	}

	fmt.Printf("  检测到 %s 模式安装\n", core.ModeLabel(mode))

	fmt.Printf("\n  %s警告: 此操作将停止所有服务%s\n", ui.Yellow, ui.Reset)
	if flagPurge {
		fmt.Printf("  %s警告: --purge 将删除所有内容包括配置文件和数据%s\n", ui.Yellow, ui.Reset)
	}
	fmt.Print("  确认卸载? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("  已取消卸载")
		return nil
	}

	d := core.NewDeployer(cfg)
	d.Rollback(flagPurge)

	fmt.Println("\n  卸载完成")
	if !flagPurge {
		fmt.Printf("  %s注意: 配置文件和数据已保留，如需彻底删除请使用 --purge%s\n", ui.Yellow, ui.Reset)
	}
	return nil
}
