package ui

import (
	"fmt"
	"net"

	"github.com/schollz/progressbar/v3"
)

// 颜色常量
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Cyan   = "\033[36m"
)

// Version 版本号
var Version = "dev"

// SetVersion 设置版本号
func SetVersion(v string) {
	Version = v
}

// PrintBanner 打印 Banner
func PrintBanner() {
	fmt.Printf(`
%s    ███████╗██╗     ███████╗ ██████╗
    ██╔════╝██║     ██╔════╝██╔════╝
    █████╗  ██║     █████╗  ██║
    ██╔══╝  ██║     ██╔══╝  ██║
    ██║     ███████╗███████╗╚██████╗
    ╚═╝     ╚══════╝╚══════╝ ╚═════╝%s

%s    FlecBlog 博客系统安装器 v%s%s

`, Cyan, Reset, Bold, Version, Reset)
}

// PrintDockerEnv 打印 Docker 环境信息
func PrintDockerEnv(dockerVer, composeVer string) {
	fmt.Printf("  %s系统检测%s\n", Bold, Reset)
	fmt.Printf("  Docker %s\n", dockerVer)
	fmt.Printf("  Docker Compose %s\n", composeVer)
}

// Progress 进度条
type Progress struct {
	bar *progressbar.ProgressBar
}

// NewProgress 创建进度条
func NewProgress(total int, description string) *Progress {
	bar := progressbar.NewOptions(total,
		progressbar.OptionSetDescription(fmt.Sprintf("  %s", description)),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)
	return &Progress{bar: bar}
}

// Describe 更新进度条描述
func (p *Progress) Describe(msg string) {
	p.bar.Describe(fmt.Sprintf("  %s", msg))
}

// Add 增加进度
func (p *Progress) Add(n int) {
	_ = p.bar.Add(n)
}

// Finish 完成进度
func (p *Progress) Finish() {
	_ = p.bar.Finish()
}

// Fail 标记失败
func (p *Progress) Fail() {
	_ = p.bar.Clear()
	fmt.Printf("  %s部署中断%s\n", Red, Reset)
}

// PrintSuccess 打印安装成功信息
func PrintSuccess(domain string, blogPort, adminPort, serverPort int, installPath string) {
	var blogURL, adminURL, serverURL string
	if net.ParseIP(domain) != nil || domain == "localhost" {
		blogURL = fmt.Sprintf("http://%s:%d", domain, blogPort)
		adminURL = fmt.Sprintf("http://%s:%d", domain, adminPort)
		serverURL = fmt.Sprintf("http://%s:%d", domain, serverPort)
	} else {
		blogURL = fmt.Sprintf("https://%s", domain)
		adminURL = fmt.Sprintf("https://%s", domain)
		serverURL = fmt.Sprintf("https://%s", domain)
	}

	fmt.Printf(`
%s  安装完成%s

  -------------------------------------

  博客前台:  %s
  管理后台:  %s
  后端服务:  %s

  默认账号:  admin@example.com
  默认密码:  123456

  配置目录:  %s

  -------------------------------------

  %s请登录管理后台修改默认密码%s

`, Green, Reset,
		blogURL, adminURL, serverURL,
		installPath,
		Yellow, Reset,
	)
}
