package mcp

import (
	"context"
	"net/http/httptest"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestPublicServerExposesWorkflowInstructions(t *testing.T) {
	session := connectPublicMCPTestClient(t)
	initResult := session.InitializeResult()
	if initResult == nil {
		t.Fatal("InitializeResult() = nil")
	}
	if initResult.Instructions != publicServerInstructions {
		t.Fatalf("instructions = %q, want %q", initResult.Instructions, publicServerInstructions)
	}
}

func TestPublicArticleToolAnnotationsMatchRiskSemantics(t *testing.T) {
	session := connectPublicMCPTestClient(t)
	ctx := context.Background()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}

	toolByName := make(map[string]*sdkmcp.Tool, len(result.Tools))
	for _, tool := range result.Tools {
		toolByName[tool.Name] = tool
	}

	for _, name := range []string{"article_list", "article_get", "article_search", "stats_query", "user_list", "user_get"} {
		annotations := requireToolAnnotations(t, toolByName, name)
		if !annotations.ReadOnlyHint {
			t.Errorf("%s readOnlyHint = false, want true", name)
		}
		assertClosedWorldHint(t, name, annotations)
	}

	for _, name := range []string{"article_create_draft", "article_image_upload", "user_create"} {
		annotations := requireToolAnnotations(t, toolByName, name)
		if annotations.DestructiveHint == nil || *annotations.DestructiveHint {
			t.Errorf("%s destructiveHint = %v, want false", name, annotations.DestructiveHint)
		}
		assertClosedWorldHint(t, name, annotations)
	}

	for _, name := range []string{"article_update_draft", "article_update_published", "article_delete", "user_update", "user_delete"} {
		annotations := requireToolAnnotations(t, toolByName, name)
		if annotations.DestructiveHint == nil || !*annotations.DestructiveHint {
			t.Errorf("%s destructiveHint = %v, want true", name, annotations.DestructiveHint)
		}
		assertClosedWorldHint(t, name, annotations)
	}

	publish := requireToolAnnotations(t, toolByName, "article_publish")
	if publish.DestructiveHint == nil || !*publish.DestructiveHint {
		t.Errorf("article_publish destructiveHint = %v, want true", publish.DestructiveHint)
	}
	if !publish.IdempotentHint {
		t.Error("article_publish idempotentHint = false, want true")
	}
	assertClosedWorldHint(t, "article_publish", publish)

	unpublish := requireToolAnnotations(t, toolByName, "article_unpublish")
	if unpublish.DestructiveHint == nil || !*unpublish.DestructiveHint {
		t.Errorf("article_unpublish destructiveHint = %v, want true", unpublish.DestructiveHint)
	}
	if !unpublish.IdempotentHint {
		t.Error("article_unpublish idempotentHint = false, want true")
	}
	assertClosedWorldHint(t, "article_unpublish", unpublish)
}

func connectPublicMCPTestClient(t *testing.T) *sdkmcp.ClientSession {
	t.Helper()
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-metadata-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })
	return session
}

func requireToolAnnotations(t *testing.T, tools map[string]*sdkmcp.Tool, name string) *sdkmcp.ToolAnnotations {
	t.Helper()
	tool := tools[name]
	if tool == nil {
		t.Fatalf("tool %q not found", name)
	}
	if tool.Annotations == nil {
		t.Fatalf("tool %q annotations = nil", name)
	}
	return tool.Annotations
}

func assertClosedWorldHint(t *testing.T, name string, annotations *sdkmcp.ToolAnnotations) {
	t.Helper()
	if annotations.OpenWorldHint == nil || *annotations.OpenWorldHint {
		t.Errorf("%s openWorldHint = %v, want false", name, annotations.OpenWorldHint)
	}
}
