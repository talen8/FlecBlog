package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"flec_blog/config"
	mcpauth "flec_blog/pkg/mcp/auth"

	"github.com/gin-gonic/gin"
)

func TestMCPAuthPreservesExactStaticSecretAndInjectsPrincipal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conf := &config.Config{AI: config.AIConfig{MCPSecret: "exact-secret"}}
	router := gin.New()
	router.GET("/mcp", MCPAuth(conf), func(c *gin.Context) {
		principal, ok := mcpauth.PrincipalFromContext(c.Request.Context())
		if !ok || principal.Method != "static" || !principal.HasScope(mcpauth.ScopeAdmin) {
			c.String(http.StatusInternalServerError, "missing static principal")
			return
		}
		c.String(http.StatusOK, principal.Subject)
	})

	accepted := httptest.NewRecorder()
	acceptedReq := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	acceptedReq.Header.Set("Authorization", "Bearer exact-secret")
	router.ServeHTTP(accepted, acceptedReq)
	if accepted.Code != http.StatusOK {
		t.Fatalf("exact secret status = %d body=%s", accepted.Code, accepted.Body.String())
	}

	rejected := httptest.NewRecorder()
	rejectedReq := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	rejectedReq.Header.Set("Authorization", "Bearer exact-secret ")
	router.ServeHTTP(rejected, rejectedReq)
	if rejected.Code != http.StatusUnauthorized {
		t.Fatalf("trailing-space secret status = %d, want 401", rejected.Code)
	}
}

func TestMCPResourceAuthChallengesMissingBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := mcpauth.Config{
		Mode:            mcpauth.ModeOAuth,
		ResourceURI:     "https://mcp.example.test/mcp",
		Issuer:          "https://issuer.example.test/",
		JWKSURL:         "https://issuer.example.test/jwks",
		Audience:        "https://mcp.example.test/mcp",
		MetadataURL:     "https://mcp.example.test/.well-known/oauth-protected-resource/mcp",
		MetadataPath:    "/.well-known/oauth-protected-resource/mcp",
		ChallengeScopes: []string{mcpauth.ScopeRead, mcpauth.ScopeDraft},
	}
	authenticator, err := mcpauth.NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	router := gin.New()
	router.GET("/mcp", MCPResourceAuth(authenticator), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/mcp", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("missing bearer status = %d, want 401", recorder.Code)
	}
	challenge := recorder.Header().Get("WWW-Authenticate")
	if challenge == "" || challenge != authenticator.UnauthorizedChallenge() {
		t.Fatalf("WWW-Authenticate = %q, want %q", challenge, authenticator.UnauthorizedChallenge())
	}
}

func TestMCPResourceAuthHybridAcceptsLegacySecretWithoutJWKS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := mcpauth.Config{
		Mode:        mcpauth.ModeHybrid,
		ResourceURI: "https://mcp.example.test/mcp",
		Issuer:      "https://issuer.example.test/",
		JWKSURL:     "http://127.0.0.1:1/jwks",
		Audience:    "https://mcp.example.test/mcp",
	}
	authenticator, err := mcpauth.NewAuthenticator(cfg, func() string { return "legacy-secret" })
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	router := gin.New()
	router.GET("/mcp", MCPResourceAuth(authenticator), func(c *gin.Context) {
		principal, ok := mcpauth.PrincipalFromContext(c.Request.Context())
		if !ok || principal.Method != "static" {
			c.String(http.StatusInternalServerError, "missing static principal")
			return
		}
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer legacy-secret")
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("hybrid static status = %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestMCPAuthWithSecretProviderReflectsRotation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "first-secret"
	router := gin.New()
	router.GET("/mcp", MCPAuthWithSecretProvider(func() string { return secret }), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := func(token string) int {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(recorder, req)
		return recorder.Code
	}

	if got := request("first-secret"); got != http.StatusNoContent {
		t.Fatalf("initial secret status = %d, want %d", got, http.StatusNoContent)
	}
	secret = "second-secret"
	if got := request("first-secret"); got != http.StatusUnauthorized {
		t.Fatalf("rotated old secret status = %d, want %d", got, http.StatusUnauthorized)
	}
	if got := request("second-secret"); got != http.StatusNoContent {
		t.Fatalf("rotated new secret status = %d, want %d", got, http.StatusNoContent)
	}
}
