package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadConfigFromEnvDefaultsToStatic(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "")
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv: %v", err)
	}
	if cfg.Mode != ModeStatic {
		t.Fatalf("mode = %q, want static", cfg.Mode)
	}
	if cfg.ResourceURI != "" || cfg.MetadataURL != "" {
		t.Fatalf("static config unexpectedly contains OAuth metadata: %+v", cfg)
	}
}

func TestLoadConfigFromEnvOAuth(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "hybrid")
	t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://mcp.example.test/public/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "https://issuer.example.test/tenant/")
	t.Setenv("MCP_OAUTH_JWKS_URL", "https://issuer.example.test/tenant/jwks")
	t.Setenv("MCP_OAUTH_AUDIENCE", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv: %v", err)
	}
	if cfg.Mode != ModeHybrid {
		t.Fatalf("mode = %q, want hybrid", cfg.Mode)
	}
	if cfg.Audience != cfg.ResourceURI {
		t.Fatalf("audience = %q, want resource %q", cfg.Audience, cfg.ResourceURI)
	}
	if cfg.Issuer != "https://issuer.example.test/tenant/" {
		t.Fatalf("issuer = %q, want exact configured issuer", cfg.Issuer)
	}
	if cfg.MetadataPath != "/.well-known/oauth-protected-resource/public/mcp" {
		t.Fatalf("metadata path = %q", cfg.MetadataPath)
	}
	if cfg.MetadataURL != "https://mcp.example.test/.well-known/oauth-protected-resource/public/mcp" {
		t.Fatalf("metadata URL = %q", cfg.MetadataURL)
	}
}

func TestLoadConfigFromEnvRejectsUnsafeRemoteHTTP(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "oauth")
	t.Setenv("MCP_OAUTH_RESOURCE_URI", "http://mcp.example.test/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "https://issuer.example.test")
	t.Setenv("MCP_OAUTH_JWKS_URL", "https://issuer.example.test/jwks")
	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatal("unsafe remote HTTP resource error = nil")
	}

	t.Setenv("MCP_OAUTH_RESOURCE_URI", "http://127.0.0.1:8080/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "http://127.0.0.1:9090")
	t.Setenv("MCP_OAUTH_JWKS_URL", "http://127.0.0.1:9090/jwks")
	if _, err := LoadConfigFromEnv(); err != nil {
		t.Fatalf("loopback HTTP should be allowed for local testing: %v", err)
	}
}

func TestMetadataHandlerAndChallenge(t *testing.T) {
	cfg := testOAuthConfig("http://127.0.0.1:9/jwks")
	authenticator, err := NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, cfg.MetadataPath, nil)
	recorder := httptest.NewRecorder()
	authenticator.MetadataHandler().ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("metadata status = %d", recorder.Code)
	}
	var metadata ProtectedResourceMetadata
	if err := json.Unmarshal(recorder.Body.Bytes(), &metadata); err != nil {
		t.Fatalf("decode metadata: %v", err)
	}
	if metadata.Resource != cfg.ResourceURI {
		t.Fatalf("resource = %q", metadata.Resource)
	}
	if len(metadata.AuthorizationServers) != 1 || metadata.AuthorizationServers[0] != cfg.Issuer {
		t.Fatalf("authorization_servers = %v", metadata.AuthorizationServers)
	}
	if len(metadata.ScopesSupported) != len(defaultSupportedScopes) {
		t.Fatalf("scopes_supported = %v", metadata.ScopesSupported)
	}

	challenge := authenticator.UnauthorizedChallenge()
	if !strings.Contains(challenge, `resource_metadata="`+cfg.MetadataURL+`"`) {
		t.Fatalf("challenge missing resource metadata: %q", challenge)
	}
	if !strings.Contains(challenge, `scope="mcp:read mcp:draft mcp:publish mcp:manage"`) {
		t.Fatalf("challenge missing default scopes: %q", challenge)
	}
}

func TestUnauthorizedChallengeTracksAdminToolsCapability(t *testing.T) {
	enabled := false
	authenticator := &Authenticator{
		config: Config{
			Mode:            ModeOAuth,
			MetadataURL:     "https://mcp.example/.well-known/oauth-protected-resource/mcp",
			ChallengeScopes: []string{ScopeRead, ScopeDraft, ScopePublish, ScopeManage},
		},
		adminToolsEnabledProvider: func() bool { return enabled },
	}

	challenge := authenticator.UnauthorizedChallenge()
	if strings.Contains(challenge, ScopeAdmin) {
		t.Fatalf("disabled challenge unexpectedly contains %s: %q", ScopeAdmin, challenge)
	}

	enabled = true
	challenge = authenticator.UnauthorizedChallenge()
	if !strings.Contains(challenge, `scope="mcp:read mcp:draft mcp:publish mcp:manage mcp:admin"`) {
		t.Fatalf("enabled challenge missing admin scope: %q", challenge)
	}

	enabled = false
	challenge = authenticator.UnauthorizedChallenge()
	if strings.Contains(challenge, ScopeAdmin) {
		t.Fatalf("disabled challenge retained admin scope: %q", challenge)
	}
}

func TestExtractBearerToken(t *testing.T) {
	for _, header := range []string{"Bearer abc", "bearer abc"} {
		token, ok := ExtractBearerToken(header)
		if !ok || token != "abc" {
			t.Errorf("ExtractBearerToken(%q) = %q, %t", header, token, ok)
		}
	}
	for _, header := range []string{"", "abc", "Bearer", "Bearer a b", "Basic abc", "Bearer abc ", " Bearer abc", "Bearer  abc"} {
		if _, ok := ExtractBearerToken(header); ok {
			t.Errorf("ExtractBearerToken(%q) unexpectedly accepted", header)
		}
	}
}

func TestLoadConfigFromEnvEmbeddedDerivesJWKSURL(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "hybrid")
	t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://aa.example.test/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "https://aa.example.test")
	t.Setenv("MCP_OAUTH_JWKS_URL", "")
	t.Setenv("MCP_OAUTH_EMBEDDED_SERVER", "true")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv: %v", err)
	}
	if cfg.JWKSURL != "https://aa.example.test/.well-known/jwks.json" {
		t.Fatalf("JWKSURL = %q", cfg.JWKSURL)
	}
}

func TestLoadConfigFromEnvExternalOAuthStillRequiresJWKSURL(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "oauth")
	t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://aa.example.test/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "https://issuer.example.test")
	t.Setenv("MCP_OAUTH_JWKS_URL", "")
	t.Setenv("MCP_OAUTH_EMBEDDED_SERVER", "false")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatal("missing external JWKS URL error = nil")
	}
}

func TestLoadConfigFromEnvCustomChallengeScopes(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "hybrid")
	t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://aa.example.test/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "https://aa.example.test")
	t.Setenv("MCP_OAUTH_JWKS_URL", "https://aa.example.test/.well-known/jwks.json")
	t.Setenv("MCP_OAUTH_CHALLENGE_SCOPES", "mcp:read,mcp:admin mcp:read")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv: %v", err)
	}
	want := []string{ScopeRead, ScopeAdmin}
	if len(cfg.ChallengeScopes) != len(want) {
		t.Fatalf("ChallengeScopes = %v, want %v", cfg.ChallengeScopes, want)
	}
	for i := range want {
		if cfg.ChallengeScopes[i] != want[i] {
			t.Fatalf("ChallengeScopes = %v, want %v", cfg.ChallengeScopes, want)
		}
	}
}

func TestLoadConfigFromEnvRejectsUnsupportedChallengeScope(t *testing.T) {
	t.Setenv("MCP_AUTH_MODE", "hybrid")
	t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://aa.example.test/mcp")
	t.Setenv("MCP_OAUTH_ISSUER", "https://aa.example.test")
	t.Setenv("MCP_OAUTH_JWKS_URL", "https://aa.example.test/.well-known/jwks.json")
	t.Setenv("MCP_OAUTH_CHALLENGE_SCOPES", "mcp:read mcp:root")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatal("unsupported challenge scope error = nil")
	}
}
