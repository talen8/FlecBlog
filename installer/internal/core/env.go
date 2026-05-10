package core

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/talen8/flecblog/installer/internal/ui"
)

type Environment struct {
	DockerVer  string
	ComposeVer string
}

type Ports struct {
	Server int
	Blog   int
	Admin  int
}

func DetectDocker() (*Environment, error) {
	env := &Environment{}

	out, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").Output()
	if err != nil {
		return nil, fmt.Errorf("docker 未安装或未启动")
	}
	env.DockerVer = strings.TrimSpace(string(out))

	if out, err := exec.Command("docker", "compose", "version", "--short").Output(); err == nil {
		env.ComposeVer = "v" + strings.TrimSpace(string(out))
	} else {
		if out, err := exec.Command("docker-compose", "version", "--short").Output(); err == nil {
			env.ComposeVer = "v" + strings.TrimSpace(string(out))
		} else {
			return nil, fmt.Errorf("docker compose 未安装")
		}
	}

	return env, nil
}

func GetConflictPorts(ports Ports) []int {
	var conflicts []int
	for _, port := range []int{ports.Server, ports.Blog, ports.Admin} {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			conflicts = append(conflicts, port)
		} else {
			if err := ln.Close(); err != nil {
				// 忽略关闭监听器的错误
				_ = err
			}
		}
	}
	return conflicts
}

func PromptDockerInstallation() error {
	fmt.Println("  Docker 未安装，请先安装 Docker")
	fmt.Println()
	fmt.Println("  使用以下命令安装 Docker:")
	fmt.Println()
	fmt.Printf("    %sbash <(curl -sSL https://linuxmirrors.cn/docker.sh)%s\n", ui.Green, ui.Reset)
	fmt.Println()
	fmt.Println("  安装完成后，执行以下命令继续安装:")
	fmt.Printf("    %sflecb install%s\n", ui.Green, ui.Reset)
	fmt.Println()

	return fmt.Errorf("docker 未安装")
}
