package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/talen8/flecblog/installer/internal/core"
)

func LoadConfigOrExit() (*core.Config, string) {
	cfg, path, err := core.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
	return cfg, path
}

func ConfirmAction(reader *bufio.Reader, prompt string, defaultYes bool) bool {
	fmt.Print(prompt)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)

	if defaultYes {
		return confirm == "" || confirm == "y" || confirm == "Y"
	}
	return confirm == "y" || confirm == "Y"
}

func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("无效的端口号: %d", port)
	}
	return nil
}
