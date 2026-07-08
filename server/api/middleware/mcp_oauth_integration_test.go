package middleware

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"flec_blog/config"
	blogmcp "flec_blog/pkg/mcp"
	mcpauth "flec_blog/pkg/mcp/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPResourceAuthOAuthOverStreamableHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeMiddlewareTestJWKS(t, w, &key.PublicKey, "oauth-key")
	}))
	defer jwksServer.Close()

	cfg := mcpauth.Config{
		Mode:            mcpauth.ModeOAuth,
		ResourceURI:     "https://mcp.example.test/mcp",
		Issuer:          "https://issuer.example.test/",
		JWKSURL:         jwksServer.URL,
		Audience:        "https://mcp.example.test/mcp",
		MetadataURL:     "https://mcp.example.test/.well-known/oauth-protected-resource/mcp",
		MetadataPath:    "/.well-known/oauth-protected-resource/mcp",
		SupportedScopes: []string{mcpauth.ScopeRead, mcpauth.ScopeDraft, mcpauth.ScopePublish, mcpauth.ScopeManage},
		ChallengeScopes: []string{mcpauth.ScopeRead, mcpauth.ScopeDraft},
	}
	authenticator, err := mcpauth.NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}

	claims := jwt.MapClaims{
		"iss":   cfg.Issuer,
		"aud":   cfg.Audience,
		"sub":   "oauth-user-1",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"scope": "mcp:read mcp:draft",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessToken.Header["kid"] = "oauth-key"
	signedToken, err := accessToken.SignedString(key)
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}

	router := gin.New()
	mcpHandler := blogmcp.NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	mcpHandler = sdkauth.RequireBearerToken(mcpauth.SDKTokenVerifierFromPrincipalContext, nil)(mcpHandler)
	router.Any("/mcp", MCPResourceAuth(authenticator), gin.WrapH(mcpHandler))
	server := httptest.NewServer(router)
	defer server.Close()

	httpClient := &http.Client{Transport: bearerRoundTripper{
		base:  http.DefaultTransport,
		token: signedToken,
	}}
	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "oauth-integration-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{
		Endpoint:   server.URL + "/mcp",
		HTTPClient: httpClient,
	}, nil)
	if err != nil {
		t.Fatalf("connect OAuth MCP client: %v", err)
	}
	defer session.Close()

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools over OAuth MCP: %v", err)
	}
	if len(result.Tools) == 0 {
		t.Fatal("tools/list returned no tools")
	}

	publishResult, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "article_publish",
		Arguments: map[string]any{"id": 0},
	})
	if err != nil {
		t.Fatalf("call article_publish with read/draft token: %v", err)
	}
	if !publishResult.IsError {
		t.Fatal("article_publish insufficient-scope result isError = false")
	}
	publishError := toolResultText(t, publishResult)
	if !strings.Contains(publishError, "insufficient MCP scope") || !strings.Contains(publishError, mcpauth.ScopePublish) {
		t.Fatalf("article_publish scope error = %q", publishError)
	}

	readResult, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "article_search",
		Arguments: map[string]any{"keyword": "   "},
	})
	if err != nil {
		t.Fatalf("call article_search with read token: %v", err)
	}
	if !readResult.IsError {
		t.Fatal("article_search blank-keyword result isError = false")
	}
	if strings.Contains(toolResultText(t, readResult), "insufficient MCP scope") {
		t.Fatal("article_search read token was blocked by scope authorization")
	}
}

type bearerRoundTripper struct {
	base  http.RoundTripper
	token string
}

func (rt bearerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+rt.token)
	return rt.base.RoundTrip(clone)
}

func writeMiddlewareTestJWKS(t *testing.T, w http.ResponseWriter, key *rsa.PublicKey, kid string) {
	t.Helper()
	document := map[string]any{
		"keys": []map[string]any{{
			"kty": "RSA",
			"kid": kid,
			"use": "sig",
			"alg": "RS256",
			"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
		}},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(document); err != nil {
		t.Fatalf("encode JWKS: %v", err)
	}
}

func TestMCPAuthStaticOverStreamableHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conf := &config.Config{AI: config.AIConfig{MCPSecret: "legacy-static-secret"}}

	router := gin.New()
	handler := gin.WrapH(blogmcp.NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil))
	router.Any("/mcp", MCPAuth(conf), handler)
	server := httptest.NewServer(router)
	defer server.Close()

	httpClient := &http.Client{Transport: bearerRoundTripper{
		base:  http.DefaultTransport,
		token: "legacy-static-secret",
	}}
	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "static-integration-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{
		Endpoint:   server.URL + "/mcp",
		HTTPClient: httpClient,
	}, nil)
	if err != nil {
		t.Fatalf("connect static MCP client: %v", err)
	}
	defer session.Close()

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools over static MCP: %v", err)
	}
	if len(result.Tools) == 0 {
		t.Fatal("tools/list returned no tools")
	}

	deleteResult, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "article_delete",
		Arguments: map[string]any{"id": 0},
	})
	if err != nil {
		t.Fatalf("call article_delete with static secret: %v", err)
	}
	if !deleteResult.IsError {
		t.Fatal("article_delete id=0 result isError = false")
	}
	if strings.Contains(toolResultText(t, deleteResult), "insufficient MCP scope") {
		t.Fatal("static principal lost admin capability")
	}
}

func toolResultText(t *testing.T, result *sdkmcp.CallToolResult) string {
	t.Helper()
	data, err := json.Marshal(result.Content)
	if err != nil {
		t.Fatalf("marshal tool result content: %v", err)
	}
	return string(data)
}
