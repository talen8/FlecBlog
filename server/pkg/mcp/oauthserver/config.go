package oauthserver

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	mcpauth "flec_blog/pkg/mcp/auth"
)

const ScopeOfflineAccess = "offline_access"

const (
	defaultAccessTokenTTL    = time.Hour
	defaultRefreshTokenTTL   = 30 * 24 * time.Hour
	defaultClientInactiveTTL = 30 * 24 * time.Hour
	defaultClientUnusedTTL   = 15 * time.Minute
	defaultAuthCodeTTL       = 5 * time.Minute
	defaultPendingRequestTTL = 10 * time.Minute
	defaultMaxClients        = 128
	defaultMaxPending        = 256
	defaultMaxAuthCodes      = 512
	defaultMaxRefreshTokens  = 512
)

type Config struct {
	Enabled bool

	Issuer      string
	ResourceURI string
	Audience    string

	AuthorizationEndpoint string
	TokenEndpoint         string
	RegistrationEndpoint  string
	JWKSURI               string
	SigningKeyFile        string

	SupportedScopes []string
	DefaultScopes   []string

	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
	ClientInactiveTTL time.Duration
	ClientUnusedTTL   time.Duration
	AuthCodeTTL       time.Duration
	PendingRequestTTL time.Duration

	MaxClients       int
	MaxPending       int
	MaxAuthCodes     int
	MaxRefreshTokens int
}

func LoadConfigFromEnv(resourceConfig mcpauth.Config) (Config, error) {
	enabled, err := strconv.ParseBool(strings.TrimSpace(os.Getenv("MCP_OAUTH_EMBEDDED_SERVER")))
	if err != nil && strings.TrimSpace(os.Getenv("MCP_OAUTH_EMBEDDED_SERVER")) != "" {
		return Config{}, fmt.Errorf("MCP_OAUTH_EMBEDDED_SERVER must be a boolean")
	}
	if !enabled {
		return Config{}, nil
	}
	if resourceConfig.Mode == mcpauth.ModeStatic {
		return Config{}, fmt.Errorf("embedded OAuth server requires oauth or hybrid MCP auth mode")
	}
	if resourceConfig.Issuer == "" || resourceConfig.ResourceURI == "" || resourceConfig.Audience == "" {
		return Config{}, fmt.Errorf("embedded OAuth server requires issuer, resource URI, and audience")
	}
	issuerURL, err := url.Parse(resourceConfig.Issuer)
	if err != nil || issuerURL.Scheme == "" || issuerURL.Host == "" {
		return Config{}, fmt.Errorf("embedded OAuth issuer is invalid")
	}
	if issuerURL.Path != "" && issuerURL.Path != "/" {
		return Config{}, fmt.Errorf("embedded OAuth issuer must use the origin root path")
	}
	issuer := resourceConfig.Issuer
	endpointRoot := strings.TrimSuffix(resourceConfig.Issuer, "/")
	signingKeyFile := strings.TrimSpace(os.Getenv("MCP_OAUTH_SIGNING_KEY_FILE"))
	if signingKeyFile == "" {
		return Config{}, fmt.Errorf("MCP_OAUTH_SIGNING_KEY_FILE is required when embedded OAuth server is enabled")
	}

	scopes := append([]string(nil), resourceConfig.SupportedScopes...)
	scopes = append(scopes, ScopeOfflineAccess)

	return Config{
		Enabled:               true,
		Issuer:                issuer,
		ResourceURI:           resourceConfig.ResourceURI,
		Audience:              resourceConfig.Audience,
		AuthorizationEndpoint: endpointRoot + "/authorize",
		TokenEndpoint:         endpointRoot + "/token",
		RegistrationEndpoint:  endpointRoot + "/register",
		JWKSURI:               endpointRoot + "/.well-known/jwks.json",
		SigningKeyFile:        signingKeyFile,
		SupportedScopes:       scopes,
		DefaultScopes:         append([]string(nil), resourceConfig.ChallengeScopes...),
		AccessTokenTTL:        defaultAccessTokenTTL,
		RefreshTokenTTL:       defaultRefreshTokenTTL,
		ClientInactiveTTL:     defaultClientInactiveTTL,
		ClientUnusedTTL:       defaultClientUnusedTTL,
		AuthCodeTTL:           defaultAuthCodeTTL,
		PendingRequestTTL:     defaultPendingRequestTTL,
		MaxClients:            defaultMaxClients,
		MaxPending:            defaultMaxPending,
		MaxAuthCodes:          defaultMaxAuthCodes,
		MaxRefreshTokens:      defaultMaxRefreshTokens,
	}, nil
}
