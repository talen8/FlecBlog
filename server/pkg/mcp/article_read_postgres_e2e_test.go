package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/internal/testutil"
	mcpauth "flec_blog/pkg/mcp/auth"
	"flec_blog/pkg/mcp/tools"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGranularArticleReadToolsOverMCPPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	draft := model.Article{
		Title:     "MCP Read E2E Draft",
		Slug:      "mcp-read-e2e-draft",
		Content:   "draftneedle0706 only appears in this draft body",
		Summary:   "draft summary",
		IsPublish: false,
	}
	published := model.Article{
		Title:     "MCP Read E2E Published",
		Slug:      "mcp-read-e2e-published",
		Content:   "published body",
		Summary:   "published summary",
		IsPublish: true,
	}
	if err := db.Create(&draft).Error; err != nil {
		t.Fatalf("create draft article: %v", err)
	}
	if err := db.Create(&published).Error; err != nil {
		t.Fatalf("create published article: %v", err)
	}

	articleRepo := repository.NewArticleRepository(db)
	tagRepo := repository.NewTagRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	articleService := service.NewArticleService(articleRepo, tagRepo, categoryRepo, commentRepo, nil, db)

	handler := NewPublicHandler(articleService, nil, nil, nil, nil, nil, nil, nil, nil)
	server := newMCPTestServer(t, handler)

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-mcp-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	listResult := callMCPTool(t, ctx, session, "article_list", map[string]any{
		"is_publish": false,
		"page_size":  10,
	})
	var listOutput tools.ArticleCollectionOutput
	decodeStructuredContent(t, listResult.StructuredContent, &listOutput)
	if !containsArticleID(listOutput.Items, draft.ID) {
		t.Fatalf("article_list draft items = %v, want draft ID %d", articleIDs(listOutput.Items), draft.ID)
	}
	if containsArticleID(listOutput.Items, published.ID) {
		t.Fatalf("article_list draft items unexpectedly contain published ID %d", published.ID)
	}

	searchResult := callMCPTool(t, ctx, session, "article_search", map[string]any{
		"keyword":   "draftneedle0706",
		"page_size": 10,
	})
	var searchOutput tools.ArticleCollectionOutput
	decodeStructuredContent(t, searchResult.StructuredContent, &searchOutput)
	if !containsArticleID(searchOutput.Items, draft.ID) {
		t.Fatalf("article_search items = %v, want draft ID %d", articleIDs(searchOutput.Items), draft.ID)
	}

	getResult := callMCPTool(t, ctx, session, "article_get", map[string]any{"id": draft.ID})
	var getOutput tools.ArticleGetOutput
	decodeStructuredContent(t, getResult.StructuredContent, &getOutput)
	if getOutput.Item.ID != draft.ID {
		t.Fatalf("article_get ID = %d, want %d", getOutput.Item.ID, draft.ID)
	}
	if getOutput.Item.Content != draft.Content {
		t.Fatalf("article_get content = %q, want %q", getOutput.Item.Content, draft.Content)
	}
}

func newMCPTestServer(t *testing.T, handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}) *httptest.Server {
	t.Helper()
	bridgedHandler := sdkauth.RequireBearerToken(mcpauth.SDKTokenVerifierFromPrincipalContext, nil)(handler)
	staticAuthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := mcpauth.ContextWithPrincipal(r.Context(), mcpauth.StaticPrincipal())
		r.Header.Set("Authorization", "Bearer test-static")
		bridgedHandler.ServeHTTP(w, r.WithContext(ctx))
	})
	server := httptest.NewServer(staticAuthorized)
	t.Cleanup(server.Close)
	return server
}

func callMCPTool(t *testing.T, ctx context.Context, session *sdkmcp.ClientSession, name string, arguments any) *sdkmcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(ctx, &sdkmcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	if result.IsError {
		t.Fatalf("call %s returned isError=true: %v", name, result.Content)
	}
	return result
}

func decodeStructuredContent(t *testing.T, content any, out any) {
	t.Helper()
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
}

func containsArticleID(items []tools.ArticleItem, id uint) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func articleIDs(items []tools.ArticleItem) []uint {
	ids := make([]uint, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}
