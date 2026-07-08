package tools

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

const (
	maxArticleImageBytes     = 5 << 20
	maxArticleImageDimension = 16384
	maxArticleImagePixels    = 40_000_000
)

func decodeAndValidateArticleImage(raw string) (*validatedArticleImage, error) {
	payload, declaredMIME, err := splitArticleImageBase64(raw)
	if err != nil {
		return nil, err
	}

	data, err := decodeArticleImageBase64(payload)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("图片内容不能为空")
	}

	mimeType, extension, width, height, err := inspectArticleImage(data)
	if err != nil {
		return nil, err
	}
	if declaredMIME != "" && normalizeArticleImageMIME(declaredMIME) != mimeType {
		return nil, fmt.Errorf("data URL 声明类型 %q 与实际图片类型 %q 不一致", declaredMIME, mimeType)
	}
	if err := validateArticleImageDimensions(width, height); err != nil {
		return nil, err
	}

	digest := sha256.Sum256(data)
	return &validatedArticleImage{
		data:      data,
		mimeType:  mimeType,
		extension: extension,
		width:     width,
		height:    height,
		sha256:    hex.EncodeToString(digest[:]),
	}, nil
}

func splitArticleImageBase64(raw string) (payload, declaredMIME string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", fmt.Errorf("image_base64 不能为空")
	}
	if !strings.HasPrefix(strings.ToLower(raw), "data:") {
		return raw, "", nil
	}

	comma := strings.IndexByte(raw, ',')
	if comma < 0 {
		return "", "", fmt.Errorf("无效的 data URL")
	}
	meta := raw[len("data:"):comma]
	parts := strings.Split(meta, ";")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("data URL 必须使用 base64 编码")
	}
	declaredMIME = strings.ToLower(strings.TrimSpace(parts[0]))
	base64Mode := false
	for _, part := range parts[1:] {
		if strings.EqualFold(strings.TrimSpace(part), "base64") {
			base64Mode = true
			break
		}
	}
	if !base64Mode {
		return "", "", fmt.Errorf("data URL 必须使用 base64 编码")
	}
	return strings.TrimSpace(raw[comma+1:]), declaredMIME, nil
}

func decodeArticleImageBase64(payload string) ([]byte, error) {
	maxEncodedLen := base64.StdEncoding.EncodedLen(maxArticleImageBytes+1) + 1024
	if len(payload) > maxEncodedLen {
		return nil, fmt.Errorf("图片超过 MCP 最大允许大小 %d MiB", maxArticleImageBytes>>20)
	}

	encodings := []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding}
	var lastErr error
	for _, encoding := range encodings {
		decoder := base64.NewDecoder(encoding, strings.NewReader(payload))
		data, err := io.ReadAll(io.LimitReader(decoder, maxArticleImageBytes+1))
		if err != nil {
			lastErr = err
			continue
		}
		if len(data) > maxArticleImageBytes {
			return nil, fmt.Errorf("图片超过 MCP 最大允许大小 %d MiB", maxArticleImageBytes>>20)
		}
		return data, nil
	}
	return nil, fmt.Errorf("image_base64 解码失败: %w", lastErr)
}

func inspectArticleImage(data []byte) (mimeType, extension string, width, height int, err error) {
	switch {
	case len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}):
		cfg, decodeErr := png.DecodeConfig(bytes.NewReader(data))
		if decodeErr != nil {
			return "", "", 0, 0, fmt.Errorf("无效的 PNG 图片: %w", decodeErr)
		}
		return "image/png", ".png", cfg.Width, cfg.Height, nil
	case len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff:
		cfg, decodeErr := jpeg.DecodeConfig(bytes.NewReader(data))
		if decodeErr != nil {
			return "", "", 0, 0, fmt.Errorf("无效的 JPEG 图片: %w", decodeErr)
		}
		return "image/jpeg", ".jpg", cfg.Width, cfg.Height, nil
	case len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")):
		w, h, decodeErr := decodeWebPDimensions(data)
		if decodeErr != nil {
			return "", "", 0, 0, fmt.Errorf("无效的 WebP 图片: %w", decodeErr)
		}
		return "image/webp", ".webp", w, h, nil
	default:
		return "", "", 0, 0, fmt.Errorf("仅支持 PNG、JPEG、WebP 图片；SVG 等格式不允许上传")
	}
}

func validateArticleImageDimensions(width, height int) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("图片尺寸无效: %dx%d", width, height)
	}
	if width > maxArticleImageDimension || height > maxArticleImageDimension {
		return fmt.Errorf("图片单边尺寸超过限制 %d 像素: %dx%d", maxArticleImageDimension, width, height)
	}
	if int64(width)*int64(height) > maxArticleImagePixels {
		return fmt.Errorf("图片总像素超过限制 %d: %dx%d", maxArticleImagePixels, width, height)
	}
	return nil
}

func decodeWebPDimensions(data []byte) (int, int, error) {
	if len(data) < 20 || !bytes.Equal(data[:4], []byte("RIFF")) || !bytes.Equal(data[8:12], []byte("WEBP")) {
		return 0, 0, fmt.Errorf("WebP RIFF 头无效")
	}
	riffSize := uint64(binary.LittleEndian.Uint32(data[4:8])) + 8
	if riffSize != uint64(len(data)) {
		return 0, 0, fmt.Errorf("WebP RIFF 长度不匹配")
	}

	offset := 12
	canvasWidth, canvasHeight := 0, 0
	imageWidth, imageHeight := 0, 0
	for offset < len(data) {
		if offset+8 > len(data) {
			return 0, 0, fmt.Errorf("WebP chunk 头被截断")
		}
		chunkType := string(data[offset : offset+4])
		chunkSize := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		payloadStart := offset + 8
		payloadEnd := payloadStart + chunkSize
		if chunkSize < 0 || payloadEnd < payloadStart || payloadEnd > len(data) {
			return 0, 0, fmt.Errorf("WebP chunk %q 长度无效", chunkType)
		}
		payload := data[payloadStart:payloadEnd]

		switch chunkType {
		case "VP8X":
			if len(payload) != 10 {
				return 0, 0, fmt.Errorf("VP8X chunk 长度无效")
			}
			if payload[0]&0x02 != 0 {
				return 0, 0, fmt.Errorf("不支持动画 WebP")
			}
			canvasWidth = 1 + int(uint32(payload[4])|uint32(payload[5])<<8|uint32(payload[6])<<16)
			canvasHeight = 1 + int(uint32(payload[7])|uint32(payload[8])<<8|uint32(payload[9])<<16)
		case "VP8L":
			w, h, parseErr := decodeVP8LDimensions(payload)
			if parseErr != nil {
				return 0, 0, parseErr
			}
			imageWidth, imageHeight = w, h
		case "VP8 ":
			w, h, parseErr := decodeVP8Dimensions(payload)
			if parseErr != nil {
				return 0, 0, parseErr
			}
			imageWidth, imageHeight = w, h
		}

		offset = payloadEnd
		if chunkSize%2 == 1 {
			offset++
		}
		if offset > len(data) {
			return 0, 0, fmt.Errorf("WebP chunk padding 被截断")
		}
	}
	if imageWidth == 0 || imageHeight == 0 {
		return 0, 0, fmt.Errorf("缺少 VP8/VP8L 图像数据")
	}
	if canvasWidth != 0 && (canvasWidth != imageWidth || canvasHeight != imageHeight) {
		return 0, 0, fmt.Errorf("VP8X 画布尺寸与图像尺寸不一致")
	}
	return imageWidth, imageHeight, nil
}

func decodeVP8LDimensions(payload []byte) (int, int, error) {
	if len(payload) < 5 || payload[0] != 0x2f {
		return 0, 0, fmt.Errorf("VP8L 头无效")
	}
	bits := uint32(payload[1]) | uint32(payload[2])<<8 | uint32(payload[3])<<16 | uint32(payload[4])<<24
	if bits>>29 != 0 {
		return 0, 0, fmt.Errorf("VP8L 版本无效")
	}
	width := 1 + int(bits&0x3fff)
	height := 1 + int((bits>>14)&0x3fff)
	return width, height, nil
}

func decodeVP8Dimensions(payload []byte) (int, int, error) {
	if len(payload) < 10 {
		return 0, 0, fmt.Errorf("VP8 头被截断")
	}
	frameTag := uint32(payload[0]) | uint32(payload[1])<<8 | uint32(payload[2])<<16
	if frameTag&1 != 0 {
		return 0, 0, fmt.Errorf("VP8 首帧不是关键帧")
	}
	if !bytes.Equal(payload[3:6], []byte{0x9d, 0x01, 0x2a}) {
		return 0, 0, fmt.Errorf("VP8 起始码无效")
	}
	firstPartitionSize := int((frameTag >> 5) & 0x7ffff)
	if 10+firstPartitionSize > len(payload) {
		return 0, 0, fmt.Errorf("VP8 首分区被截断")
	}
	width := int(binary.LittleEndian.Uint16(payload[6:8]) & 0x3fff)
	height := int(binary.LittleEndian.Uint16(payload[8:10]) & 0x3fff)
	return width, height, nil
}

func normalizeArticleImageMIME(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "image/jpg" {
		return "image/jpeg"
	}
	return value
}
