package tools

import (
	"strings"
	"testing"

	mcpauth "flec_blog/pkg/mcp/auth"
)

func TestResolveOAuthOperatorRejectsUnverifiedRawClaims(t *testing.T) {
	tests := []struct {
		name   string
		claims map[string]any
	}{
		{name: "raw user id", claims: map[string]any{"mcp_user_id": float64(42)}},
		{name: "raw email", claims: map[string]any{"email": "local-admin@example.test"}},
		{name: "both raw claims", claims: map[string]any{
			"mcp_user_id": float64(42),
			"email":       "local-admin@example.test",
		}},
	}

	resolver := &userOperatorResolver{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			principal := &mcpauth.Principal{
				Method:  "oauth",
				Subject: "external-user",
				Scopes:  map[string]struct{}{mcpauth.ScopeAdmin: {}},
				Claims:  tc.claims,
			}
			if _, err := resolver.resolveOAuthOperator(principal); err == nil {
				t.Fatal("unverified external OAuth claims resolved a local operator")
			} else if !strings.Contains(err.Error(), "已验证的本地用户绑定") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
