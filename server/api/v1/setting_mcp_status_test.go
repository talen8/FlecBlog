package v1

import "testing"

func TestLoadMCPAuthStatusFromEnv(t *testing.T) {
	t.Run("static", func(t *testing.T) {
		t.Setenv("MCP_AUTH_MODE", "static")
		status, err := loadMCPAuthStatusFromEnv()
		if err != nil {
			t.Fatalf("load status: %v", err)
		}
		if status.Mode != "static" || status.OAuth != "disabled" {
			t.Fatalf("status = %+v", status)
		}
	})

	t.Run("hybrid embedded", func(t *testing.T) {
		t.Setenv("MCP_AUTH_MODE", "hybrid")
		t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://mcp.example.test/mcp")
		t.Setenv("MCP_OAUTH_ISSUER", "https://mcp.example.test")
		t.Setenv("MCP_OAUTH_JWKS_URL", "")
		t.Setenv("MCP_OAUTH_EMBEDDED_SERVER", "true")
		t.Setenv("MCP_OAUTH_SIGNING_KEY_FILE", "/tmp/test-ed25519.pem")
		status, err := loadMCPAuthStatusFromEnv()
		if err != nil {
			t.Fatalf("load status: %v", err)
		}
		if status.Mode != "hybrid" || status.OAuth != "embedded" {
			t.Fatalf("status = %+v", status)
		}
	})

	t.Run("oauth external", func(t *testing.T) {
		t.Setenv("MCP_AUTH_MODE", "oauth")
		t.Setenv("MCP_OAUTH_RESOURCE_URI", "https://mcp.example.test/mcp")
		t.Setenv("MCP_OAUTH_ISSUER", "https://issuer.example.test")
		t.Setenv("MCP_OAUTH_JWKS_URL", "https://issuer.example.test/jwks")
		t.Setenv("MCP_OAUTH_EMBEDDED_SERVER", "false")
		status, err := loadMCPAuthStatusFromEnv()
		if err != nil {
			t.Fatalf("load status: %v", err)
		}
		if status.Mode != "oauth" || status.OAuth != "external" {
			t.Fatalf("status = %+v", status)
		}
	})
}
