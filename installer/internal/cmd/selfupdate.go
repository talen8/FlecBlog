package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/talen8/flecblog/installer/internal/core"
)

var selfUpdateCmd = &cobra.Command{
	Use:    "self-update",
	Short:  "更新安装器",
	Long:   "静默更新安装器到最新版本。",
	Hidden: true,
	RunE:   runSelfUpdate,
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	updater := core.NewSelfUpdater()
	_ = updater.Update()

	lockFile := filepath.Join(os.TempDir(), "flecb-update.lock")
	// #nosec G306 - lockFile 是临时文件，0644 权限足够
	_ = os.WriteFile(lockFile, []byte(time.Now().Format(time.RFC3339)), 0644)

	return nil
}
