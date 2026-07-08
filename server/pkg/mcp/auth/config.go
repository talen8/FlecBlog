package auth

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
)

// Mode MCP HTTP authentication mode.
type Mode string

const (
	ModeStatic Mode = "static"
	ModeOAuth  Mode = "oauth"
	ModeHybrid Mode = "hybrid"
)

const (
	ScopeRead    = "mcp:read"
	ScopeDraft   = "mcp:draft"
	ScopePublish = "mcp:publish"
	ScopeManage  = "mcp:manage"
	ScopeAdmin   = "mcp:admin"
)

var defaultSupportedScopes = []string{
	ScopeRead,
	ScopeDraft,
	ScopePublish,
	ScopeManage,
	ScopeAdmin,
}

var defaultChallengeScopes = []string{
	ScopeRead,
	ScopeDraft,
	ScopePublish,
	ScopeManage,
}

// Config configures MCP resource-server authentication.
type Config struct {
	Mode            Mode
	ResourceURI     string
	Issuer          string
	JWKSURL         string
	Audience        string
	MetadataURL     string
	MetadataPath    string
	SupportedScopes []string
	ChallengeScopes []string
}

// LoadConfigFromEnv loads MCP auth configuration without network access.
// Static mode is the default and requires no additional variables.
func LoadConfigFromEnv() (Config, error) {
	mode := Mode(strings.ToLower(strings.TrimSpace(os.Getenv("MCP_AUTH_MODE"))))
	if mode == "" {
		mode = ModeStatic
	}
	if mode != ModeStatic && mode != ModeOAuth && mode != ModeHybrid {
		return Config{}, fmt.Errorf("MCP_AUTH_MODE must be static, oauth, or hybrid")
	}

	cfg := Config{
		Mode:            mode,
		SupportedScopes: append([]string(nil), defaultSupportedScopes...),
		ChallengeScopes: append([]string(nil), defaultChallengeScopes...),
	}
	if mode == ModeStatic {
		return cfg, nil
	}

	challengeScopes, err := loadChallengeScopesFromEnv(cfg.SupportedScopes)
	if err != nil {
		return Config{}, err
	}
	cfg.ChallengeScopes = challengeScopes

	resourceURI, err := validateAbsoluteHTTPURL("MCP_OAUTH_RESOURCE_URI", os.Getenv("MCP_OAUTH_RESOURCE_URI"))
	if err != nil {
		return Config{}, err
	}
	issuer, err := validateAbsoluteHTTPURL("MCP_OAUTH_ISSUER", os.Getenv("MCP_OAUTH_ISSUER"))
	if err != nil {
		return Config{}, err
	}
	jwksRaw := strings.TrimSpace(os.Getenv("MCP_OAUTH_JWKS_URL"))
	if jwksRaw == "" && strings.EqualFold(strings.TrimSpace(os.Getenv("MCP_OAUTH_EMBEDDED_SERVER")), "true") {
		issuerRoot := strings.TrimSuffix(issuer.String(), "/")
		jwksRaw = issuerRoot + "/.well-known/jwks.json"
	}
	jwksURL, err := validateAbsoluteHTTPURL("MCP_OAUTH_JWKS_URL", jwksRaw)
	if err != nil {
		return Config{}, err
	}

	if err := requireSecureRemoteURL("MCP_OAUTH_RESOURCE_URI", resourceURI); err != nil {
		return Config{}, err
	}
	if err := requireSecureRemoteURL("MCP_OAUTH_ISSUER", issuer); err != nil {
		return Config{}, err
	}
	if err := requireSecureRemoteURL("MCP_OAUTH_JWKS_URL", jwksURL); err != nil {
		return Config{}, err
	}

	audience := strings.TrimSpace(os.Getenv("MCP_OAUTH_AUDIENCE"))
	if audience == "" {
		audience = resourceURI.String()
	}

	metadataPath := protectedResourceMetadataPath(resourceURI)
	metadataURL := (&url.URL{
		Scheme: resourceURI.Scheme,
		Host:   resourceURI.Host,
		Path:   metadataPath,
	}).String()

	cfg.ResourceURI = resourceURI.String()
	cfg.Issuer = issuer.String()
	cfg.JWKSURL = jwksURL.String()
	cfg.Audience = audience
	cfg.MetadataPath = metadataPath
	cfg.MetadataURL = metadataURL
	return cfg, nil
}

func loadChallengeScopesFromEnv(supported []string) ([]string, error) {
	raw := strings.TrimSpace(os.Getenv("MCP_OAUTH_CHALLENGE_SCOPES"))
	if raw == "" {
		return append([]string(nil), defaultChallengeScopes...), nil
	}
	requested := strings.Fields(strings.ReplaceAll(raw, ",", " "))
	if len(requested) == 0 {
		return nil, fmt.Errorf("MCP_OAUTH_CHALLENGE_SCOPES must contain at least one scope")
	}
	allowed := make(map[string]struct{}, len(supported))
	for _, scope := range supported {
		allowed[scope] = struct{}{}
	}
	seen := make(map[string]struct{}, len(requested))
	result := make([]string, 0, len(requested))
	for _, scope := range requested {
		if _, ok := allowed[scope]; !ok {
			return nil, fmt.Errorf("MCP_OAUTH_CHALLENGE_SCOPES contains unsupported scope %q", scope)
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		result = append(result, scope)
	}
	return result, nil
}

func validateAbsoluteHTTPURL(name, raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("%s is required for oauth/hybrid mode", name)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%s is invalid: %w", name, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("%s must use http or https", name)
	}
	if parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("%s has unsupported URL components", name)
	}
	return parsed, nil
}

func requireSecureRemoteURL(name string, parsed *url.URL) error {
	if parsed.Scheme == "https" {
		return nil
	}
	host := parsed.Hostname()
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	return fmt.Errorf("%s must use https unless the host is loopback", name)
}

func protectedResourceMetadataPath(resourceURI *url.URL) string {
	resourcePath := path.Clean("/" + strings.TrimSpace(resourceURI.EscapedPath()))
	if resourcePath == "/" || resourcePath == "/." {
		return "/.well-known/oauth-protected-resource"
	}
	return "/.well-known/oauth-protected-resource" + resourcePath
}
