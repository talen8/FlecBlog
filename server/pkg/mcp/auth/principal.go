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

const (
	staticSubject         = "mcp-static-secret"
	tokenInfoPrincipalKey = "flecblog.mcp.principal"
	bridgeTokenInfoTTL    = time.Minute
)

type principalContextKey struct{}

// Principal is the authenticated MCP caller identity.
type Principal struct {
	Method  string
	Subject string
	Scopes  map[string]struct{}
	Claims  map[string]any

	verifiedLocalUserID uint
}

// BindVerifiedLocalUser records an application-validated local user binding.
// Raw OAuth claims must never be treated as an equivalent binding.
func (p *Principal) BindVerifiedLocalUser(userID uint) error {
	if p == nil {
		return fmt.Errorf("principal is nil")
	}
	if userID == 0 {
		return fmt.Errorf("local user ID is empty")
	}
	p.verifiedLocalUserID = userID
	return nil
}

// VerifiedLocalUserID returns the application-validated local user binding.
func (p *Principal) VerifiedLocalUserID() (uint, bool) {
	if p == nil || p.verifiedLocalUserID == 0 {
		return 0, false
	}
	return p.verifiedLocalUserID, true
}

// StaticPrincipal returns the full-capability identity used by the legacy MCP secret.
func StaticPrincipal() *Principal {
	return &Principal{
		Method:  "static",
		Subject: staticSubject,
		Scopes:  map[string]struct{}{"*": {}},
	}
}

// ContextWithPrincipal attaches an authenticated principal to a request context.
func ContextWithPrincipal(ctx context.Context, principal *Principal) context.Context {
	if principal == nil {
		return ctx
	}
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext returns the authenticated MCP principal, if present.
func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*Principal)
	return principal, ok && principal != nil
}

// SDKTokenVerifierFromPrincipalContext bridges an already-authenticated FlecBlog
// principal into go-sdk TokenInfo. The outer FlecBlog middleware remains the
// credential verifier; this adapter enables SDK session binding and per-request
// transport metadata without re-validating the bearer token.
func SDKTokenVerifierFromPrincipalContext(ctx context.Context, _ string, _ *http.Request) (*sdkauth.TokenInfo, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w: MCP authorization principal is missing", sdkauth.ErrInvalidToken)
	}
	return sdkTokenInfoFromPrincipal(principal, time.Now())
}

// PrincipalFromTokenInfo returns the current-request FlecBlog principal carried
// through go-sdk RequestExtra metadata.
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
	case "oauth":
		issuer, _ := principal.Claims["iss"].(string)
		issuer = strings.TrimSpace(issuer)
		if issuer == "" {
			return "", fmt.Errorf("OAuth principal issuer is missing")
		}
		material = "oauth\x00" + issuer + "\x00" + subject
	default:
		return "", fmt.Errorf("unsupported principal method %q", method)
	}

	digest := sha256.Sum256([]byte(material))
	return "flecblog-mcp-v1:" + hex.EncodeToString(digest[:]), nil
}

// HasScope reports whether the principal can use the requested scope.
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
