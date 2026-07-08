package tools

import (
	"context"
	"fmt"
	"path"
	"strings"
	"unicode"
)

func normalizeArticleImageFilename(filename, extension string) string {
	filename = strings.ReplaceAll(strings.TrimSpace(filename), "\\", "/")
	filename = path.Base(filename)
	stem := strings.TrimSuffix(filename, path.Ext(filename))
	stem = strings.Trim(stem, " .")
	stem = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r == '/' || r == '\\' {
			return -1
		}
		return r
	}, stem)
	runes := []rune(stem)
	if len(runes) > 80 {
		runes = runes[:80]
	}
	stem = strings.TrimSpace(string(runes))
	if stem == "" || stem == "." || stem == ".." {
		stem = "article-image"
	}
	return stem + extension
}

func (w *ArticleImageWrapper) uploadValidatedArticleImage(
	ctx context.Context,
	filename string,
	imageInfo *validatedArticleImage,
) (string, error) {
	if w == nil || w.uploader == nil {
		return "", fmt.Errorf("文章图片上传服务未配置")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	fileURL, err := w.uploader.UploadArticleImage(
		ctx,
		filename,
		imageInfo.mimeType,
		imageInfo.data,
	)
	if err != nil {
		return "", fmt.Errorf("上传文章图片失败: %w", err)
	}
	return strings.TrimSpace(fileURL), nil
}
