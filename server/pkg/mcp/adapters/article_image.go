package adapters

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	"flec_blog/internal/service"
	"flec_blog/pkg/upload"
)

const articleImageUploadType upload.Type = "文章图片"

// FileServiceArticleImageUploader persists MCP-validated images through the existing file service.
type FileServiceArticleImageUploader struct {
	fileService *service.FileService
	publicHost  string
}

func NewFileServiceArticleImageUploader(
	fileService *service.FileService,
	publicAPIURL string,
) (*FileServiceArticleImageUploader, error) {
	if fileService == nil {
		return nil, fmt.Errorf("file service is required")
	}
	publicHost, err := publicHostFromAPIURL(publicAPIURL)
	if err != nil {
		return nil, err
	}
	return &FileServiceArticleImageUploader{
		fileService: fileService,
		publicHost:  publicHost,
	}, nil
}

func (u *FileServiceArticleImageUploader) UploadArticleImage(
	ctx context.Context,
	filename, mimeType string,
	data []byte,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("image data is empty")
	}
	return u.fileService.UploadFromReader(
		bytes.NewReader(data),
		filename,
		mimeType,
		articleImageUploadType,
		0,
		u.publicHost,
	)
}

func publicHostFromAPIURL(rawAPIURL string) (string, error) {
	rawAPIURL = strings.TrimSpace(rawAPIURL)
	if rawAPIURL == "" {
		return "", fmt.Errorf("API_URL is required")
	}
	parsed, err := url.Parse(rawAPIURL)
	if err != nil {
		return "", fmt.Errorf("invalid API_URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("API_URL must use http or https")
	}
	if parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("unsupported API_URL format")
	}
	return (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String(), nil
}
