package mcp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mcpauth "flec_blog/pkg/mcp/auth"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const testProtocolVersion = "2025-03-26"

func TestPublicHandlerSessionContracts(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	sessionA := initializeTestSession(t, server.URL)
	sessionB := initializeTestSession(t, server.URL)
	if sessionA == sessionB {
		t.Fatalf("two initialize requests returned the same session ID %q", sessionA)
	}

	postTestJSON(t, server.URL, sessionA, `{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}`, http.StatusAccepted)
	postTestJSON(t, server.URL, sessionB, `{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}`, http.StatusAccepted)
	postTestJSON(t, server.URL, sessionA, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`, http.StatusOK)
	postTestJSON(t, server.URL, sessionB, `{"jsonrpc":"2.0","id":3,"method":"tools/list","params":{}}`, http.StatusOK)

	postTestJSON(t, server.URL, "unknown-session", `{"jsonrpc":"2.0","id":4,"method":"tools/list","params":{}}`, http.StatusNotFound)

	deleteTestSession(t, server.URL, sessionA, http.StatusNoContent)
	postTestJSON(t, server.URL, sessionA, `{"jsonrpc":"2.0","id":5,"method":"tools/list","params":{}}`, http.StatusNotFound)

	// Terminating one session must not invalidate another active session.
	postTestJSON(t, server.URL, sessionB, `{"jsonrpc":"2.0","id":6,"method":"tools/list","params":{}}`, http.StatusOK)
	deleteTestSession(t, server.URL, sessionB, http.StatusNoContent)
}

type authRoundTripper func(*http.Request) (*http.Response, error)

func (f authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestPublicHandlerSessionAdmissionCapAndDeleteRecovery(t *testing.T) {
	handler := NewPublicHandlerWithOptions(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		PublicHandlerOptions{MaxSessions: 2},
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	sessionA := initializeTestSession(t, server.URL)
	sessionB := initializeTestSession(t, server.URL)
	if sessionA == sessionB {
		t.Fatalf("two initialize requests returned the same session ID %q", sessionA)
	}

	body := `{"jsonrpc":"2.0","id":30,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"capacity-test","version":"1.0"}}}`
	resp := doTestRequest(t, http.MethodPost, server.URL, "", body)
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("read capacity response: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable || !strings.Contains(string(data), "MCP_SESSION_LIMIT_REACHED") {
		t.Fatalf("capacity response status=%d body=%q", resp.StatusCode, data)
	}

	// Capacity pressure must not disturb already established sessions.
	postTestJSON(t, server.URL, sessionB, `{"jsonrpc":"2.0","id":31,"method":"tools/list","params":{}}`, http.StatusOK)

	deleteTestSession(t, server.URL, sessionA, http.StatusNoContent)
	sessionC := initializeTestSession(t, server.URL)
	if sessionC == "" || sessionC == sessionA {
		t.Fatalf("replacement session ID = %q", sessionC)
	}
	postTestJSON(t, server.URL, "unknown-session", `{"jsonrpc":"2.0","id":32,"method":"tools/list","params":{}}`, http.StatusNotFound)
}

func TestPublicHandlerFailedSessionlessPostDoesNotConsumeAdmission(t *testing.T) {
	handler := NewPublicHandlerWithOptions(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		PublicHandlerOptions{MaxSessions: 1},
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp := doTestRequest(t, http.MethodPost, server.URL, "", `{`)
	drainTestBody(t, resp)
	resp.Body.Close()
	if resp.Header.Get("Mcp-Session-Id") != "" {
		t.Fatalf("malformed request unexpectedly returned session ID %q", resp.Header.Get("Mcp-Session-Id"))
	}

	if sessionID := initializeTestSession(t, server.URL); sessionID == "" {
		t.Fatal("valid initialize was blocked after failed sessionless POST")
	}
}

func TestPublicHandlerSessionIdentityBinding(t *testing.T) {
	publicHandler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	bridgedHandler := sdkauth.RequireBearerToken(mcpauth.SDKTokenVerifierFromPrincipalContext, nil)(publicHandler)
	principals := map[string]*mcpauth.Principal{
		"user-a-high": testTransportOAuthPrincipal("user-a", mcpauth.ScopeAdmin),
		"user-a-low":  testTransportOAuthPrincipal("user-a", mcpauth.ScopeRead),
		"user-b":      testTransportOAuthPrincipal("user-b", mcpauth.ScopeRead),
		"static":      mcpauth.StaticPrincipal(),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := mcpauth.ExtractBearerToken(r.Header.Get("Authorization"))
		principal := principals[token]
		if !ok || principal == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := mcpauth.ContextWithPrincipal(r.Context(), principal)
		bridgedHandler.ServeHTTP(w, r.WithContext(ctx))
	}))
	defer server.Close()

	token := "user-a-high"
	httpClient := &http.Client{Transport: authRoundTripper(func(req *http.Request) (*http.Response, error) {
		clone := req.Clone(req.Context())
		clone.Header = req.Header.Clone()
		clone.Header.Set("Authorization", "Bearer "+token)
		return http.DefaultTransport.RoundTrip(clone)
	})}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "auth-binding-test", Version: "1.0"}, nil)
	session, err := client.Connect(context.Background(), &sdkmcp.StreamableClientTransport{Endpoint: server.URL, HTTPClient: httpClient}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer session.Close()

	token = "user-a-low"
	if _, err := session.ListTools(context.Background(), nil); err != nil {
		t.Fatalf("same identity with lower scopes rejected: %v", err)
	}
	publishResult, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{Name: "article_publish", Arguments: map[string]any{"id": 1}})
	if err != nil {
		t.Fatalf("call privileged tool with lower scopes: %v", err)
	}
	if !publishResult.IsError || len(publishResult.Content) == 0 {
		t.Fatalf("privileged tool result = %+v, want insufficient-scope error", publishResult)
	}
	text, ok := publishResult.Content[0].(*sdkmcp.TextContent)
	if !ok || !strings.Contains(text.Text, "insufficient MCP scope") {
		t.Fatalf("privileged tool error content = %#v", publishResult.Content)
	}
	token = "user-b"
	if _, err := session.ListTools(context.Background(), nil); err == nil {
		t.Fatal("different identity reused existing session")
	}

	token = "static"
	staticSession, err := client.Connect(context.Background(), &sdkmcp.StreamableClientTransport{Endpoint: server.URL, HTTPClient: httpClient}, nil)
	if err != nil {
		t.Fatalf("connect static session: %v", err)
	}
	defer staticSession.Close()
	token = "user-a-low"
	if _, err := staticSession.ListTools(context.Background(), nil); err == nil {
		t.Fatal("OAuth identity reused static session")
	}
}

func testTransportOAuthPrincipal(subject string, scopes ...string) *mcpauth.Principal {
	scopeSet := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		scopeSet[scope] = struct{}{}
	}
	return &mcpauth.Principal{
		Method:  "oauth",
		Subject: subject,
		Scopes:  scopeSet,
		Claims: map[string]any{
			"iss": "https://issuer.example",
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		},
	}
}

func initializeTestSession(t *testing.T, endpoint string) string {
	t.Helper()
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"flecblog-mcp-test","version":"1.0"}}}`
	resp := doTestRequest(t, http.MethodPost, endpoint, "", body)
	defer resp.Body.Close()
	drainTestBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("initialize status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	sessionID := resp.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		t.Fatal("initialize response missing Mcp-Session-Id")
	}
	return sessionID
}

func postTestJSON(t *testing.T, endpoint, sessionID, body string, wantStatus int) {
	t.Helper()
	resp := doTestRequest(t, http.MethodPost, endpoint, sessionID, body)
	defer resp.Body.Close()
	drainTestBody(t, resp)
	if resp.StatusCode != wantStatus {
		t.Fatalf("POST session %q status = %d, want %d", sessionID, resp.StatusCode, wantStatus)
	}
}

func deleteTestSession(t *testing.T, endpoint, sessionID string, wantStatus int) {
	t.Helper()
	resp := doTestRequest(t, http.MethodDelete, endpoint, sessionID, "")
	defer resp.Body.Close()
	drainTestBody(t, resp)
	if resp.StatusCode != wantStatus {
		t.Fatalf("DELETE session %q status = %d, want %d", sessionID, resp.StatusCode, wantStatus)
	}
}

func doTestRequest(t *testing.T, method, endpoint, sessionID, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, endpoint, strings.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Accept", "application/json, text/event-stream")
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
		req.Header.Set("Mcp-Protocol-Version", testProtocolVersion)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func drainTestBody(t *testing.T, resp *http.Response) {
	t.Helper()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		t.Fatalf("drain response body: %v", err)
	}
}
