package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/talen8/flecblog/installer/internal/core"
	"github.com/talen8/flecblog/installer/internal/ui"
)

var (
	flagBackupOutput string
	flagRestoreInput string
	flagNoConfirm    bool
)

// backupCmd 备份命令
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份和恢复 FlecBlog 数据",
	Long:  "备份 FlecBlog 的数据目录到压缩文件，或从备份文件恢复数据。",
	RunE:  runBackup,
}

func init() {
	backupCmd.Flags().StringVar(&flagBackupOutput, "output", "", "备份文件输出路径 (留空则使用默认路径)")
	backupCmd.Flags().StringVar(&flagRestoreInput, "restore", "", "从指定备份文件恢复数据")
	backupCmd.Flags().BoolVar(&flagNoConfirm, "yes", false, "跳过确认提示")
	rootCmd.AddCommand(backupCmd)
}

// runBackup 执行备份或恢复
func runBackup(cmd *cobra.Command, args []string) error {
	if flagRestoreInput != "" {
		return runRestore()
	}
	return runBackupCreate()
}

// runBackupCreate 执行备份
func runBackupCreate() error {
	ui.PrintBanner()

	cfg, configPath := LoadConfigOrExit()
	installPath := cfg.InstallPath
	fmt.Printf("  安装路径: %s (%s)\n", installPath, configPath)

	dataDir := filepath.Join(installPath, "data")
	postgresDir := filepath.Join(installPath, "postgres_data")

	if !dirExists(dataDir) && !dirExists(postgresDir) {
		return fmt.Errorf("未找到数据目录，请确认 FlecBlog 已正确安装")
	}

	outputPath := flagBackupOutput
	if outputPath == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputPath = filepath.Join(installPath, fmt.Sprintf("flecblog-backup-%s.tar.gz", timestamp))
	}

	fmt.Printf("  备份文件: %s\n", outputPath)

	if dirExists(dataDir) {
		fmt.Printf("  备份目录: data/\n")
	}
	if dirExists(postgresDir) {
		fmt.Printf("  备份目录: postgres_data/\n")
	}

	fmt.Println("\n  正在备份...")

	// #nosec G304 - outputPath 来自用户输入或生成的默认路径
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %w", err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	if err := core.CreateBackup(installPath, outFile); err != nil {
		return fmt.Errorf("备份失败: %w", err)
	}

	fileInfo, _ := os.Stat(outputPath)
	fmt.Printf("\n  %s备份完成%s\n", ui.Green, ui.Reset)
	fmt.Printf("  文件大小: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	fmt.Printf("  文件路径: %s\n", outputPath)

	return nil
}

// runRestore 执行恢复
func runRestore() error {
	ui.PrintBanner()

	cfg, configPath := LoadConfigOrExit()
	installPath := cfg.InstallPath
	fmt.Printf("  安装路径: %s (%s)\n", installPath, configPath)

	if _, err := os.Stat(flagRestoreInput); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", flagRestoreInput)
	}

	fmt.Printf("  备份文件: %s\n", flagRestoreInput)

	if !flagNoConfirm {
		fmt.Printf("\n  %s警告: 此操作将覆盖现有的 data/ 和 postgres_data/ 目录%s\n", ui.Yellow, ui.Reset)
		fmt.Print("  确认恢复? [y/N]: ")
		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("  已取消恢复")
			return nil
		}
	}

	fmt.Println("\n  正在恢复...")

	// #nosec G304 - flagRestoreInput 来自用户输入的命令行参数
	inFile, err := os.Open(flagRestoreInput)
	if err != nil {
		return fmt.Errorf("打开备份文件失败: %w", err)
	}
	defer func() {
		_ = inFile.Close()
	}()

	if err := core.RestoreBackup(installPath, inFile); err != nil {
		return fmt.Errorf("恢复失败: %w", err)
	}

	fmt.Printf("\n  %s恢复完成%s\n", ui.Green, ui.Reset)
	fmt.Printf("  %s建议: 如果服务正在运行，请重启服务以使数据生效%s\n", ui.Yellow, ui.Reset)

	return nil
}

// dirExists 检查目录是否存在
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
