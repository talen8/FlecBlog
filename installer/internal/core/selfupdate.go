package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	githubRepo = "talen8/flecblog"
)

type SelfUpdater struct{}

func NewSelfUpdater() *SelfUpdater {
	return &SelfUpdater{}
}

func (u *SelfUpdater) Update() error {
	tag, err := u.getLatestInstallerTag()
	if err != nil {
		return err
	}

	platform := GetPlatformSuffix()
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/flecb_%s", githubRepo, tag, platform)

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return err
	}

	tempFile := filepath.Join(os.TempDir(), "flecb-latest")
	if err := u.downloadFile(downloadURL, tempFile); err != nil {
		return err
	}

	// #nosec G302 - 可执行文件需要 0755 权限
	if err := os.Chmod(tempFile, 0755); err != nil {
		_ = os.Remove(tempFile)
		return err
	}

	backupPath := execPath + ".bak"
	if err := os.Rename(execPath, backupPath); err != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("备份可执行文件失败: %w", err)
	}

	if err := os.Rename(tempFile, execPath); err != nil {
		_ = os.Rename(backupPath, execPath)
		_ = os.Remove(tempFile)
		return err
	}

	_ = os.Remove(backupPath)
	return nil
}

func (u *SelfUpdater) downloadFile(url, targetPath string) error {
	// #nosec G107 - URL 来自 GitHub API，是受信任的源
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: %d", resp.StatusCode)
	}

	// #nosec G304 - targetPath 来自内部调用，是受信任的源
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	return err
}

func GetPlatformSuffix() string {
	return fmt.Sprintf("linux_%s", runtime.GOARCH)
}

// getLatestInstallerTag 获取最新的 installer release tag
// 从 /releases 中筛选 installer/v 前缀的 tag，与 worker.js 逻辑一致
func (u *SelfUpdater) getLatestInstallerTag() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", githubRepo)
	// #nosec G107 - URL 来自 GitHub API，是受信任的源
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取 releases 列表失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取 releases 列表失败: HTTP %d", resp.StatusCode)
	}

	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", fmt.Errorf("解析 releases 列表失败: %w", err)
	}

	for _, r := range releases {
		if strings.HasPrefix(r.TagName, "installer/") {
			return r.TagName, nil
		}
	}

	return "", fmt.Errorf("未找到 installer release")
}
