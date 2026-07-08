package tools

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"
)

type recordingArticleImageUploader struct {
	filename string
	mimeType string
	data     []byte
	url      string
	err      error
}

func (u *recordingArticleImageUploader) UploadArticleImage(
	_ context.Context,
	filename, mimeType string,
	data []byte,
) (string, error) {
	u.filename = filename
	u.mimeType = mimeType
	u.data = append([]byte(nil), data...)
	return u.url, u.err
}

func TestDecodeAndValidateArticleImageFormats(t *testing.T) {
	pngData := encodeArticleImageTestPNG(t, 3, 2)
	pngInfo, err := decodeAndValidateArticleImage(base64.StdEncoding.EncodeToString(pngData))
	if err != nil {
		t.Fatalf("decode PNG: %v", err)
	}
	if pngInfo.mimeType != "image/png" || pngInfo.extension != ".png" || pngInfo.width != 3 || pngInfo.height != 2 {
		t.Fatalf("PNG info = %+v", pngInfo)
	}

	jpegData := encodeArticleImageTestJPEG(t, 4, 5)
	jpegInfo, err := decodeAndValidateArticleImage("data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(jpegData))
	if err != nil {
		t.Fatalf("decode JPEG: %v", err)
	}
	if jpegInfo.mimeType != "image/jpeg" || jpegInfo.width != 4 || jpegInfo.height != 5 {
		t.Fatalf("JPEG info = %+v", jpegInfo)
	}

	webpInfo, err := decodeAndValidateArticleImage(base64.StdEncoding.EncodeToString(encodeArticleImageTestWebP(7, 9)))
	if err != nil {
		t.Fatalf("decode WebP: %v", err)
	}
	if webpInfo.mimeType != "image/webp" || webpInfo.width != 7 || webpInfo.height != 9 {
		t.Fatalf("WebP info = %+v", webpInfo)
	}
}

func TestDecodeAndValidateArticleImageRejectsUnsafeInput(t *testing.T) {
	svg := base64.StdEncoding.EncodeToString([]byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`))
	if _, err := decodeAndValidateArticleImage(svg); err == nil {
		t.Fatal("SVG error = nil, want rejection")
	}

	pngData := encodeArticleImageTestPNG(t, 2, 2)
	mismatch := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(pngData)
	if _, err := decodeAndValidateArticleImage(mismatch); err == nil {
		t.Fatal("MIME mismatch error = nil, want rejection")
	}

	maxEncodedLen := base64.StdEncoding.EncodedLen(maxArticleImageBytes+1) + 1024
	if _, err := decodeArticleImageBase64(strings.Repeat("A", maxEncodedLen+1)); err == nil {
		t.Fatal("oversized base64 error = nil, want rejection")
	}
}

func TestValidateArticleImageDimensions(t *testing.T) {
	if err := validateArticleImageDimensions(maxArticleImageDimension+1, 1); err == nil {
		t.Fatal("oversized dimension error = nil")
	}
	if err := validateArticleImageDimensions(10000, 5000); err == nil {
		t.Fatal("oversized pixel count error = nil")
	}
	if err := validateArticleImageDimensions(4000, 3000); err != nil {
		t.Fatalf("bounded dimensions rejected: %v", err)
	}
}

func TestUploadArticleImageToolUsesUploaderPort(t *testing.T) {
	pngData := encodeArticleImageTestPNG(t, 6, 4)
	uploader := &recordingArticleImageUploader{url: "https://cdn.example/cover.png"}
	wrapper := NewArticleImageWrapper(uploader)

	_, output, err := wrapper.UploadArticleImageTool(context.Background(), nil, ArticleImageUploadInput{
		ImageBase64: base64.StdEncoding.EncodeToString(pngData),
		Filename:    "../../cover.svg",
	})
	if err != nil {
		t.Fatalf("UploadArticleImageTool: %v", err)
	}
	if uploader.filename != "cover.png" {
		t.Fatalf("filename = %q, want cover.png", uploader.filename)
	}
	if uploader.mimeType != "image/png" {
		t.Fatalf("mime type = %q, want image/png", uploader.mimeType)
	}
	if !bytes.Equal(uploader.data, pngData) {
		t.Fatal("uploaded bytes differ")
	}
	if output.URL != "https://cdn.example/cover.png" || output.MIMEType != "image/png" || output.Width != 6 || output.Height != 4 {
		t.Fatalf("output = %+v", output)
	}
	if output.SizeBytes != len(pngData) || len(output.SHA256) != 64 {
		t.Fatalf("output metadata = %+v", output)
	}
}

func TestUploadArticleImageToolRequiresUploader(t *testing.T) {
	wrapper := NewArticleImageWrapper(nil)
	_, _, err := wrapper.UploadArticleImageTool(context.Background(), nil, ArticleImageUploadInput{
		ImageBase64: base64.StdEncoding.EncodeToString(encodeArticleImageTestPNG(t, 2, 2)),
	})
	if err == nil || !strings.Contains(err.Error(), "上传服务未配置") {
		t.Fatalf("error = %v, want missing uploader", err)
	}
}

func TestUploadArticleImageToolPropagatesUploaderError(t *testing.T) {
	wantErr := errors.New("storage unavailable")
	wrapper := NewArticleImageWrapper(&recordingArticleImageUploader{err: wantErr})
	_, _, err := wrapper.UploadArticleImageTool(context.Background(), nil, ArticleImageUploadInput{
		ImageBase64: base64.StdEncoding.EncodeToString(encodeArticleImageTestPNG(t, 2, 2)),
	})
	if err == nil || !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapped uploader error", err)
	}
}

func encodeArticleImageTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}
	return buf.Bytes()
}

func encodeArticleImageTestJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{G: 255, A: 255})
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("encode JPEG: %v", err)
	}
	return buf.Bytes()
}

func encodeArticleImageTestWebP(width, height int) []byte {
	bits := uint32(width-1) | uint32(height-1)<<14
	payload := []byte{0x2f, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24)}
	fileSize := 12 + 8 + len(payload) + 1
	data := make([]byte, fileSize)
	copy(data[:4], "RIFF")
	binary.LittleEndian.PutUint32(data[4:8], uint32(fileSize-8))
	copy(data[8:12], "WEBP")
	copy(data[12:16], "VP8L")
	binary.LittleEndian.PutUint32(data[16:20], uint32(len(payload)))
	copy(data[20:25], payload)
	return data
}
