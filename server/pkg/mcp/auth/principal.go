package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
)

type SecretProvider func() string

const (
	staticSubject         = "mcp-static-secret"
	tokenInfoPrincipalKey = "flecblog.mcp.principal" //nolint:gosec // context key name, not a credential
	bridgeTokenInfoTTL    = time.Minute
)

type principalContextKey struct{}

type Principal struct {
	Method  string
	Subject string
	Scopes  map[string]struct{}
	Claims  map[string]any
}

func StaticPrincipal() *Principal {
	return &Principal{
		Method:  "static",
		Subject: staticSubject,
		Scopes:  map[string]struct{}{"*": {}},
	}
}

func ContextWithPrincipal(ctx context.Context, principal *Principal) context.Context {
	if principal == nil {
		return ctx
	}
	return context.WithValue(ctx, principalContextKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*Principal)
	return principal, ok && principal != nil
}

func SDKTokenVerifierFromPrincipalContext(ctx context.Context, _ string, _ *http.Request) (*sdkauth.TokenInfo, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w: MCP authorization principal is missing", sdkauth.ErrInvalidToken)
	}
	return sdkTokenInfoFromPrincipal(principal, time.Now())
}

func PrincipalFromTokenInfo(info *sdkauth.TokenInfo) (*Principal, bool) {
	if info == nil || info.Extra == nil {
		return nil, false
	}
	principal, ok := info.Extra[tokenInfoPrincipalKey].(*Principal)
	return principal, ok && principal != nil
}

func sdkTokenInfoFromPrincipal(principal *Principal, now time.Time) (*sdkauth.TokenInfo, error) {
	identity, err := sessionIdentityKey(principal)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", sdkauth.ErrInvalidToken, err)
	}
	scopes := make([]string, 0, len(principal.Scopes))
	for scope := range principal.Scopes {
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)

	return &sdkauth.TokenInfo{
		Scopes:     scopes,
		Expiration: now.Add(bridgeTokenInfoTTL),
		UserID:     identity,
		Extra: map[string]any{
			tokenInfoPrincipalKey: principal,
		},
	}, nil
}

func sessionIdentityKey(principal *Principal) (string, error) {
	if principal == nil {
		return "", fmt.Errorf("principal is nil")
	}
	method := strings.TrimSpace(principal.Method)
	subject := strings.TrimSpace(principal.Subject)
	if subject == "" {
		return "", fmt.Errorf("principal subject is missing")
	}

	var material string
	switch method {
	case "static":
		material = "static\x00" + subject
	default:
		return "", fmt.Errorf("unsupported principal method %q", method)
	}

	digest := sha256.Sum256([]byte(material))
	return "flecblog-mcp-v1:" + hex.EncodeToString(digest[:]), nil
}

func (p *Principal) HasScope(scope string) bool {
	if p == nil {
		return false
	}
	if _, ok := p.Scopes["*"]; ok {
		return true
	}
	if _, ok := p.Scopes[ScopeAdmin]; ok {
		return true
	}
	_, ok := p.Scopes[scope]
	return ok
}
