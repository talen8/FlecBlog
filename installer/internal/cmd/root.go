package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/talen8/flecblog/installer/internal/ui"
)

var (
	version   string
	buildTime string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "flecb",
	Short: "FlecBlog 一键安装器",
	Long:  "FlecBlog 博客系统交互安装工具。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "flecb" || cmd.Name() == "help" || cmd.Name() == "self-update" {
			return nil
		}
		if os.Geteuid() != 0 {
			fmt.Printf("  %s错误: FlecBlog 安装器需要 root 权限，请使用 sudo 运行%s\n", ui.Red, ui.Reset)
			os.Exit(1)
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "self-update" || cmd.Name() == "help" {
			return
		}
		startBackgroundUpdate()
	},
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s (built %s)", version, buildTime)
	rootCmd.SetVersionTemplate("flecb {{.Version}}\n")
}

// SetVersionInfo 设置版本信息
func SetVersionInfo(v, t string) {
	version = v
	buildTime = t
	ui.SetVersion(v)
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func startBackgroundUpdate() {
	if !shouldUpdate() {
		return
	}

	execPath, err := os.Executable()
	if err != nil {
		return
	}

	// #nosec G204 - execPath 来自 os.Executable()，是受信任的源
	cmd := exec.Command(execPath, "self-update")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	_ = cmd.Start()
}

func shouldUpdate() bool {
	lockFile := filepath.Join(os.TempDir(), "flecb-update.lock")
	info, err := os.Stat(lockFile)
	if err != nil {
		return true
	}

	return time.Since(info.ModTime()) > 24*time.Hour
}
