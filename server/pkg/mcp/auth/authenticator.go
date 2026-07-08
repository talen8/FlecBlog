package auth

import (
	"context"
	"crypto/subtle"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const maxBearerTokenBytes = 16 << 10

// SecretProvider returns the current legacy MCP secret.
type SecretProvider func() string

// OAuthKeyResolver resolves the public key for an OAuth JWT.
type OAuthKeyResolver func(context.Context, *jwt.Token) (any, error)

// OAuthPrincipalValidator applies optional issuer-specific identity checks after JWT validation.
type OAuthPrincipalValidator func(context.Context, *Principal) error

// AdminToolsEnabledProvider returns the current runtime MCP admin-tools capability state.
type AdminToolsEnabledProvider func() bool

// AuthenticatorOptions configures optional local OAuth integrations.
type AuthenticatorOptions struct {
	OAuthKeyResolver          OAuthKeyResolver
	OAuthPrincipalValidator   OAuthPrincipalValidator
	AdminToolsEnabledProvider AdminToolsEnabledProvider
}

// Authenticator validates legacy static secrets and OAuth access tokens.
type Authenticator struct {
	config                    Config
	secret                    SecretProvider
	jwks                      *jwksCache
	keyResolver               OAuthKeyResolver
	principalValidator        OAuthPrincipalValidator
	adminToolsEnabledProvider AdminToolsEnabledProvider
}

// NewAuthenticator creates a dual-mode MCP authenticator.
func NewAuthenticator(config Config, secret SecretProvider) (*Authenticator, error) {
	return NewAuthenticatorWithOptions(config, secret, AuthenticatorOptions{})
}

// NewAuthenticatorWithOptions creates an authenticator with optional local OAuth key resolution.
func NewAuthenticatorWithOptions(config Config, secret SecretProvider, options AuthenticatorOptions) (*Authenticator, error) {
	if config.Mode != ModeStatic && config.Mode != ModeOAuth && config.Mode != ModeHybrid {
		return nil, fmt.Errorf("unsupported MCP auth mode %q", config.Mode)
	}
	if secret == nil {
		secret = func() string { return "" }
	}
	authenticator := &Authenticator{
		config:                    config,
		secret:                    secret,
		keyResolver:               options.OAuthKeyResolver,
		principalValidator:        options.OAuthPrincipalValidator,
		adminToolsEnabledProvider: options.AdminToolsEnabledProvider,
	}
	if config.Mode == ModeOAuth || config.Mode == ModeHybrid {
		if config.ResourceURI == "" || config.Issuer == "" || config.Audience == "" {
			return nil, fmt.Errorf("oauth/hybrid mode requires resource URI, issuer, and audience")
		}
		if authenticator.keyResolver == nil {
			if config.JWKSURL == "" {
				return nil, fmt.Errorf("oauth/hybrid mode requires JWKS URL without a local key resolver")
			}
			authenticator.jwks = newJWKSCache(config.JWKSURL)
		}
	}
	return authenticator, nil
}

// Config returns a defensive copy of the authenticator configuration.
func (a *Authenticator) Config() Config {
	cfg := a.config
	cfg.SupportedScopes = append([]string(nil), a.config.SupportedScopes...)
	cfg.ChallengeScopes = append([]string(nil), a.config.ChallengeScopes...)
	return cfg
}

// AuthenticateBearer validates a bearer token according to the configured mode.
func (a *Authenticator) AuthenticateBearer(ctx context.Context, token string) (*Principal, error) {
	if token == "" {
		return nil, fmt.Errorf("bearer token is empty")
	}
	if len(token) > maxBearerTokenBytes {
		return nil, fmt.Errorf("bearer token exceeds size limit")
	}

	switch a.config.Mode {
	case ModeStatic:
		return a.authenticateStatic(token)
	case ModeOAuth:
		return a.authenticateOAuth(ctx, token)
	case ModeHybrid:
		if principal, err := a.authenticateStatic(token); err == nil {
			return principal, nil
		}
		return a.authenticateOAuth(ctx, token)
	default:
		return nil, fmt.Errorf("unsupported MCP auth mode %q", a.config.Mode)
	}
}

func (a *Authenticator) authenticateStatic(token string) (*Principal, error) {
	secret := a.secret()
	if secret == "" || subtle.ConstantTimeCompare([]byte(token), []byte(secret)) != 1 {
		return nil, fmt.Errorf("static MCP secret is invalid")
	}
	return StaticPrincipal(), nil
}

func (a *Authenticator) authenticateOAuth(ctx context.Context, rawToken string) (*Principal, error) {
	parser := &jwt.Parser{
		ValidMethods:         mapKeys(allowedJWTAlgorithms),
		SkipClaimsValidation: true,
	}
	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if a.keyResolver != nil {
			return a.keyResolver(ctx, token)
		}
		return a.jwks.keyForToken(ctx, token)
	})
	if err != nil || token == nil || !token.Valid {
		if err == nil {
			err = fmt.Errorf("JWT is invalid")
		}
		return nil, fmt.Errorf("validate OAuth access token: %w", err)
	}

	now := time.Now().Unix()
	if !claims.VerifyExpiresAt(now, true) {
		return nil, fmt.Errorf("OAuth access token is expired or missing exp")
	}
	if !claims.VerifyNotBefore(now, false) {
		return nil, fmt.Errorf("OAuth access token is not active yet")
	}
	if !claims.VerifyIssuer(a.config.Issuer, true) {
		return nil, fmt.Errorf("OAuth access token issuer mismatch")
	}
	if !claims.VerifyAudience(a.config.Audience, true) {
		return nil, fmt.Errorf("OAuth access token audience mismatch")
	}

	subject, _ := claims["sub"].(string)
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return nil, fmt.Errorf("OAuth access token subject is missing")
	}

	scopes := extractScopes(claims)
	principal := &Principal{
		Method:  "oauth",
		Subject: subject,
		Scopes:  scopes,
		Claims:  cloneClaims(claims),
	}
	if a.principalValidator != nil {
		if err := a.principalValidator(ctx, principal); err != nil {
			return nil, fmt.Errorf("validate OAuth principal: %w", err)
		}
	}
	return principal, nil
}

func extractScopes(claims jwt.MapClaims) map[string]struct{} {
	scopes := make(map[string]struct{})
	addScopeString := func(value string) {
		for _, scope := range strings.Fields(value) {
			if scope != "" {
				scopes[scope] = struct{}{}
			}
		}
	}

	switch value := claims["scope"].(type) {
	case string:
		addScopeString(value)
	case []any:
		for _, item := range value {
			if scope, ok := item.(string); ok {
				addScopeString(scope)
			}
		}
	}
	switch value := claims["scp"].(type) {
	case string:
		addScopeString(value)
	case []any:
		for _, item := range value {
			if scope, ok := item.(string); ok {
				addScopeString(scope)
			}
		}
	}
	return scopes
}

func cloneClaims(claims jwt.MapClaims) map[string]any {
	copyClaims := make(map[string]any, len(claims))
	for key, value := range claims {
		copyClaims[key] = value
	}
	return copyClaims
}

func mapKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
