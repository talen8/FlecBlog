package tools

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ArticleImageUploadInput article_image_upload 输入。
// ImageBase64 可为原始 base64 或 data:image/...;base64,...。
type ArticleImageUploadInput struct {
	ImageBase64 string `json:"image_base64"`
	Filename    string `json:"filename,omitempty"`
}

// ArticleImageUploadOutput article_image_upload 输出。
type ArticleImageUploadOutput struct {
	URL       string `json:"url"`
	MIMEType  string `json:"mime_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SizeBytes int    `json:"size_bytes"`
	SHA256    string `json:"sha256"`
}

// ArticleImageUploader isolates MCP input handling from the storage implementation.
// The MCP tool owns base64 decoding and image validation; the uploader owns persistence.
type ArticleImageUploader interface {
	UploadArticleImage(ctx context.Context, filename, mimeType string, data []byte) (string, error)
}

type validatedArticleImage struct {
	data      []byte
	mimeType  string
	extension string
	width     int
	height    int
	sha256    string
}

// ArticleImageWrapper 文章图片上传包装器。
type ArticleImageWrapper struct {
	uploader ArticleImageUploader
}

func NewArticleImageWrapper(uploader ArticleImageUploader) *ArticleImageWrapper {
	return &ArticleImageWrapper{uploader: uploader}
}

// UploadArticleImageTool 安全校验 base64 图片后，通过注入的 uploader 持久化。
func (w *ArticleImageWrapper) UploadArticleImageTool(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	input ArticleImageUploadInput,
) (*sdkmcp.CallToolResult, ArticleImageUploadOutput, error) {
	imageInfo, err := decodeAndValidateArticleImage(input.ImageBase64)
	if err != nil {
		return nil, ArticleImageUploadOutput{}, err
	}

	filename := normalizeArticleImageFilename(input.Filename, imageInfo.extension)
	fileURL, err := w.uploadValidatedArticleImage(ctx, filename, imageInfo)
	if err != nil {
		return nil, ArticleImageUploadOutput{}, err
	}
	if fileURL == "" {
		return nil, ArticleImageUploadOutput{}, fmt.Errorf("文章图片上传未返回文件 URL")
	}

	return nil, ArticleImageUploadOutput{
		URL:       fileURL,
		MIMEType:  imageInfo.mimeType,
		Width:     imageInfo.width,
		Height:    imageInfo.height,
		SizeBytes: len(imageInfo.data),
		SHA256:    imageInfo.sha256,
	}, nil
}
