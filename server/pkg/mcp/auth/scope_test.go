package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type scopeTestInput struct {
	Value string `json:"value"`
}

type scopeTestOutput struct {
	Value string `json:"value"`
}

func TestRequireScopeFailsClosedWithoutCurrentRequestTokenInfo(t *testing.T) {
	wrapped := RequireScope(ScopeRead, func(context.Context, *sdkmcp.CallToolRequest, scopeTestInput) (*sdkmcp.CallToolResult, scopeTestOutput, error) {
		return nil, scopeTestOutput{Value: "called"}, nil
	})
	_, _, err := wrapped(context.Background(), nil, scopeTestInput{})
	if err == nil {
		t.Fatal("missing principal error = nil")
	}
}

func TestRequireScopeEnforcesOAuthScopes(t *testing.T) {
	called := false
	wrapped := RequireScope(ScopePublish, func(context.Context, *sdkmcp.CallToolRequest, scopeTestInput) (*sdkmcp.CallToolResult, scopeTestOutput, error) {
		called = true
		return nil, scopeTestOutput{Value: "ok"}, nil
	})
	staleHigh := oauthTestPrincipal("user-1", ScopeAdmin)
	currentLow := oauthTestPrincipal("user-1", ScopeRead)
	ctx := ContextWithPrincipal(context.Background(), staleHigh)
	_, _, err := wrapped(ctx, currentRequestForPrincipal(t, currentLow), scopeTestInput{})
	var scopeErr *InsufficientScopeError
	if !errors.As(err, &scopeErr) || scopeErr.Required != ScopePublish {
		t.Fatalf("error = %v, want publish scope error", err)
	}
	if called {
		t.Fatal("handler called despite insufficient scope")
	}
}

func TestRequireScopePreservesAllStaticCapabilities(t *testing.T) {
	wrapped := RequireScope(ScopeAdmin, func(context.Context, *sdkmcp.CallToolRequest, scopeTestInput) (*sdkmcp.CallToolResult, scopeTestOutput, error) {
		return nil, scopeTestOutput{Value: "ok"}, nil
	})
	_, output, err := wrapped(context.Background(), currentRequestForPrincipal(t, StaticPrincipal()), scopeTestInput{})
	if err != nil {
		t.Fatalf("static principal rejected: %v", err)
	}
	if output.Value != "ok" {
		t.Fatalf("output = %+v", output)
	}
}

func TestAdminScopeCoversAllOAuthScopes(t *testing.T) {
	principal := oauthTestPrincipal("admin-1", ScopeAdmin)
	for _, scope := range []string{ScopeRead, ScopeDraft, ScopePublish, ScopeManage, ScopeAdmin} {
		if !principal.HasScope(scope) {
			t.Errorf("admin principal missing implied scope %q", scope)
		}
	}
}

func currentRequestForPrincipal(t *testing.T, principal *Principal) *sdkmcp.CallToolRequest {
	t.Helper()
	info, err := sdkTokenInfoFromPrincipal(principal, time.Now())
	if err != nil {
		t.Fatalf("sdkTokenInfoFromPrincipal: %v", err)
	}
	return &sdkmcp.CallToolRequest{Extra: &sdkmcp.RequestExtra{TokenInfo: info}}
}

func oauthTestPrincipal(subject string, scopes ...string) *Principal {
	scopeSet := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		scopeSet[scope] = struct{}{}
	}
	return &Principal{
		Method:  "oauth",
		Subject: subject,
		Scopes:  scopeSet,
		Claims: map[string]any{
			"iss": "https://issuer.example",
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		},
	}
}
