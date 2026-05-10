package core

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateBackup 创建备份
func CreateBackup(installPath string, w io.Writer) error {
	gzWriter := gzip.NewWriter(w)
	defer func() {
		_ = gzWriter.Close()
	}()

	tw := tar.NewWriter(gzWriter)
	defer func() {
		_ = tw.Close()
	}()

	dataDir := getDataPath(installPath)
	postgresDir := getPostgresDataPath(installPath)

	if _, err := os.Stat(dataDir); err == nil {
		if err := addDirToTar(tw, dataDir, "data"); err != nil {
			return fmt.Errorf("备份 data 目录失败: %w", err)
		}
	}
	if _, err := os.Stat(postgresDir); err == nil {
		if err := addDirToTar(tw, postgresDir, "postgres_data"); err != nil {
			return fmt.Errorf("备份 postgres_data 目录失败: %w", err)
		}
	}

	return nil
}

// RestoreBackup 恢复备份
func RestoreBackup(installPath string, r io.Reader) error {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("解压备份文件失败: %w", err)
	}
	defer func() {
		_ = gzReader.Close()
	}()

	return extractTar(gzReader, installPath)
}

func addDirToTar(tw *tar.Writer, srcPath, archivePath string) error {
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			header.Name = filepath.ToSlash(archivePath + "/")
		} else {
			header.Name = filepath.ToSlash(filepath.Join(archivePath, relPath))
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// #nosec G304 G122 - 文件路径来自受信任的源（安装目录）
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		_, err = io.Copy(tw, file)
		return err
	})
}

func extractTar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取备份文件失败: %w", err)
		}

		// 清理路径并防止目录遍历攻击
		cleanName := filepath.Clean(header.Name)
		if strings.Contains(cleanName, "..") {
			continue // 跳过包含 .. 的路径
		}
		targetPath := filepath.Join(dest, cleanName)

		// 确保目标路径在目标目录内
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(dest)) {
			continue // 跳过越界路径
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// #nosec G301 - 目录权限使用 0750
			if err := os.MkdirAll(targetPath, 0750); err != nil {
				return fmt.Errorf("创建目录失败: %w", err)
			}
		case tar.TypeReg:
			// #nosec G301 - 父目录权限使用 0750
			if err := os.MkdirAll(filepath.Dir(targetPath), 0750); err != nil {
				return fmt.Errorf("创建父目录失败: %w", err)
			}

			// #nosec G304 - 文件路径已经过清理
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
			if err != nil {
				return fmt.Errorf("创建文件失败: %w", err)
			}

			// #nosec G110 - 限制复制的数据量以防止解压炸弹
			const maxSize = 10 * 1024 * 1024 * 1024 // 10GB 限制
			if _, err := io.CopyN(outFile, tr, maxSize); err != nil && err != io.EOF {
				_ = outFile.Close()
				return fmt.Errorf("写入文件失败: %w", err)
			}
			_ = outFile.Close()
		}
	}

	return nil
}
