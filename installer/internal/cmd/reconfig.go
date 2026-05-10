package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/talen8/flecblog/installer/internal/core"
	"github.com/talen8/flecblog/installer/internal/ui"
)

var (
	flagServerPort        int
	flagBlogPort          int
	flagAdminPort         int
	flagAPIURL            string
	flagReconfigNoConfirm bool
)

var reconfigCmd = &cobra.Command{
	Use:   "reconfig",
	Short: "修改 FlecBlog 配置",
	Long:  "修改已安装的 FlecBlog 配置，包括端口、API地址等。不指定参数则进入交互模式。",
	RunE:  runReconfig,
}

func init() {
	reconfigCmd.Flags().IntVar(&flagServerPort, "server-port", 0, "后端服务端口")
	reconfigCmd.Flags().IntVar(&flagBlogPort, "blog-port", 0, "博客前台端口")
	reconfigCmd.Flags().IntVar(&flagAdminPort, "admin-port", 0, "管理后台端口")
	reconfigCmd.Flags().StringVar(&flagAPIURL, "api-url", "", "后端 API 地址")
	reconfigCmd.Flags().BoolVar(&flagReconfigNoConfirm, "yes", false, "跳过确认提示")
	rootCmd.AddCommand(reconfigCmd)
}

func runReconfig(cmd *cobra.Command, args []string) error {
	ui.PrintBanner()

	cfg, configPath := LoadConfigOrExit()
	installPath := cfg.InstallPath

	fmt.Printf("  安装路径: %s (%s)\n", installPath, configPath)

	mode := core.DetectMode(installPath)
	if mode == "" {
		return fmt.Errorf("未在 %s 检测到 FlecBlog 安装", installPath)
	}
	fmt.Printf("  部署模式: %s\n", core.ModeLabel(mode))

	reader := bufio.NewReader(os.Stdin)

	currentServerPort, currentBlogPort, currentAdminPort, currentAPIURL, err := core.ParseComposePorts(installPath)
	if err != nil {
		fmt.Printf("  %s警告: 无法解析当前端口配置: %v%s\n", ui.Yellow, err, ui.Reset)
		currentServerPort, currentBlogPort, currentAdminPort = cfg.ServerPort, cfg.BlogPort, cfg.AdminPort
	}

	fmt.Println("\n  当前配置:")
	fmt.Printf("    后端服务端口: %d\n", currentServerPort)
	fmt.Printf("    博客前台端口: %d\n", currentBlogPort)
	fmt.Printf("    管理后台端口: %d\n", currentAdminPort)
	fmt.Printf("    API 地址: %s\n", currentAPIURL)

	newServerPort, newBlogPort, newAdminPort, newAPIURL := currentServerPort, currentBlogPort, currentAdminPort, currentAPIURL

	if err := applyFlags(&newServerPort, &newBlogPort, &newAdminPort, &newAPIURL); err != nil {
		return err
	}

	if !flagReconfigNoConfirm && !hasChanges(currentServerPort, currentBlogPort, currentAdminPort, currentAPIURL, newServerPort, newBlogPort, newAdminPort, newAPIURL) {
		if !ConfirmAction(reader, "\n  进入交互模式，按 Enter 保持当前值\n\n  继续? [Y/n]: ", true) {
			fmt.Println("  已取消")
			return nil
		}
		runInteractiveConfig(reader, &newServerPort, &newBlogPort, &newAdminPort, &newAPIURL, currentServerPort, currentBlogPort, currentAdminPort, currentAPIURL)
	}

	if !hasChanges(currentServerPort, currentBlogPort, currentAdminPort, currentAPIURL, newServerPort, newBlogPort, newAdminPort, newAPIURL) {
		fmt.Println("\n  配置未发生变化")
		return nil
	}

	showChanges(currentServerPort, currentBlogPort, currentAdminPort, currentAPIURL, newServerPort, newBlogPort, newAdminPort, newAPIURL)

	if !flagReconfigNoConfirm {
		if !ConfirmAction(reader, "\n  确认修改配置? [y/N]: ", false) {
			fmt.Println("  已取消修改")
			return nil
		}
	}

	if err := core.UpdateComposePorts(installPath, newServerPort, newBlogPort, newAdminPort); err != nil {
		return fmt.Errorf("更新端口配置失败: %w", err)
	}

	if newAPIURL != currentAPIURL {
		if err := core.UpdateAPIURL(installPath, newAPIURL); err != nil {
			return fmt.Errorf("更新 API 地址失败: %w", err)
		}
	}

	cfg.ServerPort, cfg.BlogPort, cfg.AdminPort = newServerPort, newBlogPort, newAdminPort
	if _, err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Printf("\n  %s配置已更新%s\n", ui.Green, ui.Reset)
	fmt.Println("  如果服务正在运行，请重启服务使配置生效")

	return nil
}

func applyFlags(serverPort, blogPort, adminPort *int, apiURL *string) error {
	if flagServerPort != 0 {
		if err := ValidatePort(flagServerPort); err != nil {
			return err
		}
		*serverPort = flagServerPort
	}
	if flagBlogPort != 0 {
		if err := ValidatePort(flagBlogPort); err != nil {
			return err
		}
		*blogPort = flagBlogPort
	}
	if flagAdminPort != 0 {
		if err := ValidatePort(flagAdminPort); err != nil {
			return err
		}
		*adminPort = flagAdminPort
	}
	if flagAPIURL != "" {
		*apiURL = flagAPIURL
	}
	return nil
}

func hasChanges(oldServer, oldBlog, oldAdmin int, oldAPI string, newServer, newBlog, newAdmin int, newAPI string) bool {
	return oldServer != newServer || oldBlog != newBlog || oldAdmin != newAdmin || oldAPI != newAPI
}

func runInteractiveConfig(reader *bufio.Reader, serverPort, blogPort, adminPort *int, apiURL *string, currentServer, currentBlog, currentAdmin int, currentAPI string) {
	if input := promptInput(reader, "  后端服务端口", currentServer); input != "" {
		if p, err := strconv.Atoi(input); err == nil && p > 0 && p < 65536 {
			*serverPort = p
		}
	}

	if input := promptInput(reader, "  博客前台端口", currentBlog); input != "" {
		if p, err := strconv.Atoi(input); err == nil && p > 0 && p < 65536 {
			*blogPort = p
		}
	}

	if input := promptInput(reader, "  管理后台端口", currentAdmin); input != "" {
		if p, err := strconv.Atoi(input); err == nil && p > 0 && p < 65536 {
			*adminPort = p
		}
	}

	if input := promptInput(reader, "  API 地址", currentAPI); input != "" {
		*apiURL = input
	}
}

func promptInput(reader *bufio.Reader, label string, defaultVal interface{}) string {
	fmt.Printf("%s [%v]: ", label, defaultVal)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func showChanges(oldServer, oldBlog, oldAdmin int, oldAPI string, newServer, newBlog, newAdmin int, newAPI string) {
	fmt.Println("\n  新配置:")
	if newServer != oldServer {
		fmt.Printf("    后端服务端口: %d -> %d\n", oldServer, newServer)
	}
	if newBlog != oldBlog {
		fmt.Printf("    博客前台端口: %d -> %d\n", oldBlog, newBlog)
	}
	if newAdmin != oldAdmin {
		fmt.Printf("    管理后台端口: %d -> %d\n", oldAdmin, newAdmin)
	}
	if newAPI != oldAPI {
		fmt.Printf("    API 地址: %s -> %s\n", oldAPI, newAPI)
	}
}
