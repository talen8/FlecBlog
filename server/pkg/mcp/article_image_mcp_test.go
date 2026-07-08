package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const testOnePixelPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

type mcpArticleImageUploaderStub struct {
	filename string
	mimeType string
	data     []byte
}

func (u *mcpArticleImageUploaderStub) UploadArticleImage(
	_ context.Context,
	filename, mimeType string,
	data []byte,
) (string, error) {
	u.filename = filename
	u.mimeType = mimeType
	u.data = append([]byte(nil), data...)
	return "https://cdn.example/article/cover.png", nil
}

func TestArticleImageUploadOverMCP(t *testing.T) {
	uploader := &mcpArticleImageUploaderStub{}
	handler := NewPublicHandlerWithOptions(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		PublicHandlerOptions{ArticleImageUploader: uploader},
	)
	server := newMCPTestServer(t, handler)
	defer server.Close()

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-mcp-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	result, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name: "article_image_upload",
		Arguments: map[string]any{
			"image_base64": testOnePixelPNGBase64,
			"filename":     "../../cover.svg",
		},
	})
	if err != nil {
		t.Fatalf("call article_image_upload: %v", err)
	}
	if result.IsError {
		t.Fatalf("article_image_upload returned isError=true: %v", result.Content)
	}

	wantData, err := base64.StdEncoding.DecodeString(testOnePixelPNGBase64)
	if err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	if uploader.filename != "cover.png" || uploader.mimeType != "image/png" || !bytes.Equal(uploader.data, wantData) {
		t.Fatalf("uploader call = filename %q mime %q bytes_equal=%v", uploader.filename, uploader.mimeType, bytes.Equal(uploader.data, wantData))
	}

	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var output tools.ArticleImageUploadOutput
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
	if output.URL != "https://cdn.example/article/cover.png" || output.MIMEType != "image/png" || output.Width != 1 || output.Height != 1 {
		t.Fatalf("output = %+v", output)
	}
}
