package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestPublicHandlerListsGranularArticleTools(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-mcp-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}

	want := map[string]bool{
		"article_list":             false,
		"article_get":              false,
		"article_search":           false,
		"article_create_draft":     false,
		"article_update_draft":     false,
		"article_publish":          false,
		"article_update_published": false,
		"article_unpublish":        false,
		"article_delete":           false,
		"article_image_upload":     false,
		"user_list":                false,
		"user_get":                 false,
		"user_create":              false,
		"user_update":              false,
		"user_delete":              false,
	}
	denied := map[string]bool{
		"article_manage": false,
		"user_manage":    false,
	}
	for _, tool := range result.Tools {
		for _, forbidden := range []string{
			"作者原有",
			"真实本地 operator",
			"继续遵守作者",
			"原权限规则",
			"必须映射真实本地管理员",
		} {
			if strings.Contains(tool.Description, forbidden) {
				t.Errorf("tool %q description contains internal wording %q: %q", tool.Name, forbidden, tool.Description)
			}
		}
		if _, blocked := denied[tool.Name]; blocked {
			denied[tool.Name] = true
		}
		if _, ok := want[tool.Name]; !ok {
			continue
		}
		want[tool.Name] = true
		if tool.InputSchema == nil {
			t.Errorf("tool %q missing input schema", tool.Name)
		}
		if tool.OutputSchema == nil {
			t.Errorf("tool %q missing output schema", tool.Name)
		}
		if tool.Name == "article_create_draft" {
			schemaJSON, err := json.Marshal(tool.InputSchema)
			if err != nil {
				t.Fatalf("marshal %s input schema: %v", tool.Name, err)
			}
			if bytes.Contains(schemaJSON, []byte("is_publish")) {
				t.Errorf("tool %q unexpectedly exposes is_publish in input schema", tool.Name)
			}
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("tool %q not found in tools/list", name)
		}
	}
	for name, found := range denied {
		if found {
			t.Errorf("tool %q unexpectedly exposed in tools/list", name)
		}
	}
}

func TestArticleSearchToolErrorUsesMCPIsError(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
		Name:      "article_search",
		Arguments: map[string]any{"keyword": "   "},
	})
	if err != nil {
		t.Fatalf("call article_search: %v", err)
	}
	if !result.IsError {
		t.Fatal("article_search blank keyword result isError = false, want true")
	}
}

func TestArticleCreateDraftValidationUsesMCPIsError(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
		Name: "article_create_draft",
		Arguments: map[string]any{
			"title":   "   ",
			"content": "draft body",
		},
	})
	if err != nil {
		t.Fatalf("call article_create_draft: %v", err)
	}
	if !result.IsError {
		t.Fatal("article_create_draft blank title result isError = false, want true")
	}
}

func TestArticlePublishValidationUsesMCPIsError(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
		Name:      "article_publish",
		Arguments: map[string]any{"id": 0},
	})
	if err != nil {
		t.Fatalf("call article_publish: %v", err)
	}
	if !result.IsError {
		t.Fatal("article_publish id=0 result isError = false, want true")
	}
}

func TestArticleImageUploadRejectsSVGWithMCPIsError(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	server := newMCPTestServer(t, handler)
	defer server.Close()

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-mcp-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	svg := base64.StdEncoding.EncodeToString([]byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`))
	result, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "article_image_upload",
		Arguments: map[string]any{"image_base64": svg},
	})
	if err != nil {
		t.Fatalf("call article_image_upload: %v", err)
	}
	if !result.IsError {
		t.Fatal("article_image_upload SVG result isError = false, want true")
	}
}

func TestPublicHandlerToolCallsFailClosedWithoutPrincipal(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-no-principal-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	result, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "article_publish",
		Arguments: map[string]any{"id": 0},
	})
	if err != nil {
		t.Fatalf("call article_publish: %v", err)
	}
	if !result.IsError {
		t.Fatal("tool call without principal isError = false, want true")
	}
	content, err := json.Marshal(result.Content)
	if err != nil {
		t.Fatalf("marshal tool result: %v", err)
	}
	if !bytes.Contains(content, []byte("current MCP request authorization metadata is missing")) {
		t.Fatalf("missing-principal error content = %s", content)
	}
}

func TestRequireMCPAdminToolsEnabledTracksRuntimeProvider(t *testing.T) {
	enabled := false
	called := 0
	wrapped := requireMCPAdminToolsEnabled(
		func() bool { return enabled },
		func(context.Context, *sdkmcp.CallToolRequest, struct{}) (*sdkmcp.CallToolResult, string, error) {
			called++
			return nil, "ok", nil
		},
	)

	_, _, err := wrapped(context.Background(), nil, struct{}{})
	if err == nil || !strings.Contains(err.Error(), "管理员工具未启用") {
		t.Fatalf("disabled gate error = %v", err)
	}
	if called != 0 {
		t.Fatalf("disabled gate called handler %d times", called)
	}

	enabled = true
	_, output, err := wrapped(context.Background(), nil, struct{}{})
	if err != nil || output != "ok" {
		t.Fatalf("enabled gate output=%q err=%v", output, err)
	}
	if called != 1 {
		t.Fatalf("enabled gate called handler %d times", called)
	}

	enabled = false
	_, _, err = wrapped(context.Background(), nil, struct{}{})
	if err == nil {
		t.Fatal("disabled gate remained open after runtime toggle")
	}
	if called != 1 {
		t.Fatalf("re-disabled gate called handler %d times", called)
	}
}
