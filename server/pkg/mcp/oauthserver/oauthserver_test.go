package oauthserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	mcpauth "flec_blog/pkg/mcp/auth"

	"github.com/golang-jwt/jwt/v4"
)

func TestRefreshTokenGenerationMatchesFailsClosed(t *testing.T) {
	current := uint(7)
	if refreshTokenGenerationMatches(nil, current) {
		t.Fatal("nil refresh token generation accepted")
	}
	if refreshTokenGenerationMatches(&RefreshToken{}, current) {
		t.Fatal("legacy refresh token without generation accepted")
	}
	if refreshTokenGenerationMatches(&RefreshToken{TokenVersion: tokenVersionPtr(6)}, current) {
		t.Fatal("stale refresh token generation accepted")
	}
	if !refreshTokenGenerationMatches(&RefreshToken{TokenVersion: tokenVersionPtr(current)}, current) {
		t.Fatal("matching refresh token generation rejected")
	}
}

func TestEmbeddedPrincipalValidatorRequiresTokenVersion(t *testing.T) {
	server := &Server{}
	principal := &mcpauth.Principal{
		Method:  "oauth",
		Subject: "user:42",
		Claims:  map[string]any{"mcp_user_id": float64(42)},
	}
	if err := server.OAuthPrincipalValidator()(context.Background(), principal); err == nil {
		t.Fatal("embedded principal without token_version accepted")
	}
}

func TestLoadConfigFromEnvEmbeddedOAuth(t *testing.T) {
	resource := mcpauth.Config{
		Mode:            mcpauth.ModeHybrid,
		ResourceURI:     "https://aa.example.test/mcp",
		Issuer:          "https://aa.example.test",
		Audience:        "https://aa.example.test/mcp",
		SupportedScopes: []string{mcpauth.ScopeRead, mcpauth.ScopeDraft},
		ChallengeScopes: []string{mcpauth.ScopeRead},
	}

	t.Setenv("MCP_OAUTH_EMBEDDED_SERVER", "")
	cfg, err := LoadConfigFromEnv(resource)
	if err != nil || cfg.Enabled {
		t.Fatalf("disabled config = %+v, err=%v", cfg, err)
	}

	t.Setenv("MCP_OAUTH_EMBEDDED_SERVER", "true")
	t.Setenv("MCP_OAUTH_SIGNING_KEY_FILE", "/tmp/test-key")
	cfg, err = LoadConfigFromEnv(resource)
	if err != nil {
		t.Fatalf("LoadConfigFromEnv: %v", err)
	}
	if !cfg.Enabled || cfg.AuthorizationEndpoint != "https://aa.example.test/authorize" || cfg.JWKSURI != "https://aa.example.test/.well-known/jwks.json" {
		t.Fatalf("embedded config = %+v", cfg)
	}
	if !containsString(cfg.SupportedScopes, ScopeOfflineAccess) || !containsString(cfg.DefaultScopes, mcpauth.ScopeRead) {
		t.Fatalf("scope config = supported %v default %v", cfg.SupportedScopes, cfg.DefaultScopes)
	}
	if cfg.ClientUnusedTTL != defaultClientUnusedTTL {
		t.Fatalf("ClientUnusedTTL = %s, want %s", cfg.ClientUnusedTTL, defaultClientUnusedTTL)
	}

	resource.Issuer = "https://aa.example.test/"
	cfg, err = LoadConfigFromEnv(resource)
	if err != nil {
		t.Fatalf("LoadConfigFromEnv trailing slash issuer: %v", err)
	}
	if cfg.Issuer != resource.Issuer {
		t.Fatalf("embedded issuer = %q, want exact resource issuer %q", cfg.Issuer, resource.Issuer)
	}
	if cfg.AuthorizationEndpoint != "https://aa.example.test/authorize" || cfg.TokenEndpoint != "https://aa.example.test/token" || cfg.JWKSURI != "https://aa.example.test/.well-known/jwks.json" {
		t.Fatalf("trailing slash endpoints = %+v", cfg)
	}

	resource.Issuer = "https://aa.example.test/tenant"
	if _, err := LoadConfigFromEnv(resource); err == nil {
		t.Fatal("issuer path error = nil")
	}
}

func TestSignerLoadsRestrictedEd25519KeyAndSignsJWT(t *testing.T) {
	keyPath, publicKey := writeTestSigningKey(t, 0o600)
	loaded, err := loadSigner(keyPath)
	if err != nil {
		t.Fatalf("loadSigner: %v", err)
	}
	cfg := Config{Issuer: "https://issuer.example.test", Audience: "https://resource.example.test/mcp", AccessTokenTTL: time.Hour}
	now := time.Now().UTC().Truncate(time.Second)
	raw, err := loaded.signAccessToken(cfg, 42, 7, "mcp:read offline_access", now)
	if err != nil {
		t.Fatalf("signAccessToken: %v", err)
	}
	claims := &accessTokenClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) { return publicKey, nil }, jwt.WithValidMethods([]string{"EdDSA"}))
	if err != nil || !token.Valid {
		t.Fatalf("parse signed token: valid=%v err=%v", token != nil && token.Valid, err)
	}
	if claims.Issuer != cfg.Issuer || claims.Subject != "user:42" || claims.MCPUserID != 42 || claims.TokenVersion != 7 || claims.Scope != "mcp:read offline_access" {
		t.Fatalf("claims = %+v", claims)
	}
	if len(claims.Audience) != 1 || claims.Audience[0] != cfg.Audience {
		t.Fatalf("audience = %v", claims.Audience)
	}
	if token.Header["kid"] != loaded.kid {
		t.Fatalf("kid = %v, want %s", token.Header["kid"], loaded.kid)
	}
}

func TestSignerRejectsLoosePermissions(t *testing.T) {
	keyPath, _ := writeTestSigningKey(t, 0o644)
	if _, err := loadSigner(keyPath); err == nil {
		t.Fatal("loadSigner loose permissions error = nil")
	}
}

func writeTestSigningKey(t *testing.T, mode os.FileMode) (string, ed25519.PublicKey) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("MarshalPKCS8PrivateKey: %v", err)
	}
	path := filepath.Join(t.TempDir(), "oauth-ed25519.pem")
	data := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(path, data, mode); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path, publicKey
}

func TestRegistrationValidation(t *testing.T) {
	good := registrationRequest{
		RedirectURIs:  []string{"https://chatgpt.example.test/callback", "http://127.0.0.1:8080/callback"},
		GrantTypes:    []string{"authorization_code"},
		ResponseTypes: []string{"code"},
	}
	redirects, authMethod, err := validateRegistrationRequest(good)
	if err != nil || len(redirects) != 2 || authMethod != "client_secret_post" {
		t.Fatalf("good registration redirects=%v auth_method=%q err=%v", redirects, authMethod, err)
	}
	bad := good
	bad.RedirectURIs = []string{"http://evil.example.test/callback"}
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("remote http redirect error = nil")
	}
	bad = good
	bad.GrantTypes = []string{"client_credentials"}
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("unsupported grant error = nil")
	}
	public := good
	public.TokenEndpointAuthMethod = "none"
	_, publicMethod, err := validateRegistrationRequest(public)
	if err != nil || publicMethod != "none" {
		t.Fatalf("public registration auth_method=%q err=%v", publicMethod, err)
	}
	bad = good
	bad.TokenEndpointAuthMethod = "private_key_jwt"
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("unsupported token endpoint auth method error = nil")
	}
	bad = good
	bad.ClientName = strings.Repeat("x", maxClientNameBytes+1)
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("oversized client_name error = nil")
	}
	bad = good
	bad.RedirectURIs = []string{"https://chatgpt.example.test/" + strings.Repeat("x", maxRedirectURIBytes)}
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("oversized redirect_uri error = nil")
	}
	bad = good
	bad.GrantTypes = []string{"authorization_code", "refresh_token", "authorization_code"}
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("excess grant_types error = nil")
	}
	bad = good
	bad.ResponseTypes = []string{"code", "code"}
	if _, _, err := validateRegistrationRequest(bad); err == nil {
		t.Fatal("excess response_types error = nil")
	}
}

func TestNormalizeScopes(t *testing.T) {
	s := &Server{cfg: Config{
		SupportedScopes: []string{"mcp:read", "mcp:draft", ScopeOfflineAccess},
		DefaultScopes:   []string{"mcp:read"},
	}}
	got, err := s.normalizeScopes("mcp:read mcp:read offline_access")
	if err != nil || len(got) != 2 || got[0] != "mcp:read" || got[1] != ScopeOfflineAccess {
		t.Fatalf("normalize scopes = %v err=%v", got, err)
	}
	if _, err := s.normalizeScopes("mcp:admin"); err == nil {
		t.Fatal("unknown scope error = nil")
	}
	got, err = s.normalizeScopes("")
	if err != nil || len(got) != 1 || got[0] != "mcp:read" {
		t.Fatalf("default scopes = %v err=%v", got, err)
	}
}

func TestAuthorizationFailureRateLimitThresholdAndExpiry(t *testing.T) {
	s := &Server{
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState),
	}
	now := time.Now().UTC()
	key := "rate-limit-test"

	for i := 0; i < authFailureLimit; i++ {
		allowed, retryAfter := s.beginAuthAttempt(key, now)
		if !allowed || retryAfter != 0 {
			t.Fatalf("attempt %d allowed=%v retryAfter=%s", i+1, allowed, retryAfter)
		}
		s.finishAuthAttempt(key, now, false)
	}
	allowed, retryAfter := s.beginAuthAttempt(key, now)
	if allowed || retryAfter <= 0 {
		t.Fatalf("rate limit allowed=%v retryAfter=%s, want blocked with positive retry", allowed, retryAfter)
	}

	successKey := "success-clears-test"
	allowed, retryAfter = s.beginAuthAttempt(successKey, now)
	if !allowed || retryAfter != 0 {
		t.Fatalf("success attempt allowed=%v retryAfter=%s", allowed, retryAfter)
	}
	s.finishAuthAttempt(successKey, now, true)
	s.mu.Lock()
	_, stillTracked := s.authFailures[successKey]
	s.mu.Unlock()
	if stillTracked {
		t.Fatal("successful authorization attempt remained tracked")
	}

	afterExpiry := now.Add(authFailureWindow + time.Second)
	allowed, retryAfter = s.beginAuthAttempt(key, afterExpiry)
	if !allowed || retryAfter != 0 {
		t.Fatalf("attempt after expiry allowed=%v retryAfter=%s", allowed, retryAfter)
	}
	s.finishAuthAttempt(key, afterExpiry, false)
}

func TestAuthorizationFailureRateLimitBoundsConcurrentAttempts(t *testing.T) {
	s := &Server{
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState),
	}
	now := time.Now().UTC()
	const attempts = authFailureLimit * 4
	results := make(chan bool, attempts)
	var wg sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, _ := s.beginAuthAttempt("concurrent-key", now)
			results <- allowed
		}()
	}
	wg.Wait()
	close(results)

	allowedCount := 0
	for allowed := range results {
		if allowed {
			allowedCount++
		}
	}
	if allowedCount != authFailureLimit {
		t.Fatalf("concurrent allowed attempts = %d, want %d", allowedCount, authFailureLimit)
	}
}

func TestAuthorizationFailureEvictionProtectsBlockedAndInflightKeys(t *testing.T) {
	s := &Server{
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState),
	}
	now := time.Now().UTC()
	blockedKey := "blocked-key"
	inflightKey := "inflight-key"
	s.authFailures[blockedKey] = authFailureState{
		Count:    authFailureLimit,
		ResetAt:  now.Add(authFailureWindow),
		LastSeen: now.Add(-time.Hour),
	}
	s.authFailures[inflightKey] = authFailureState{
		InFlight: 1,
		ResetAt:  now.Add(authFailureWindow),
		LastSeen: now.Add(-time.Hour),
	}
	for i := len(s.authFailures); i < maxAuthFailureKeys; i++ {
		key := fmt.Sprintf("low-risk-%04d", i)
		s.authFailures[key] = authFailureState{
			Count:    1,
			ResetAt:  now.Add(authFailureWindow),
			LastSeen: now.Add(time.Duration(i) * time.Second),
		}
	}

	allowed, retryAfter := s.beginAuthAttempt("new-key", now)
	if !allowed || retryAfter != 0 {
		t.Fatalf("new key allowed=%v retryAfter=%s", allowed, retryAfter)
	}
	if _, ok := s.authFailures[blockedKey]; !ok {
		t.Fatal("blocked key was evicted")
	}
	if _, ok := s.authFailures[inflightKey]; !ok {
		t.Fatal("in-flight key was evicted")
	}
	if len(s.authFailures) != maxAuthFailureKeys {
		t.Fatalf("auth failure map size = %d, want %d", len(s.authFailures), maxAuthFailureKeys)
	}
}

func TestAuthorizationFailureCapacityFailsClosedWithoutSafeEviction(t *testing.T) {
	s := &Server{
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState, maxAuthFailureKeys),
	}
	now := time.Now().UTC()
	for i := 0; i < maxAuthFailureKeys; i++ {
		s.authFailures[fmt.Sprintf("blocked-%04d", i)] = authFailureState{
			Count:    authFailureLimit,
			ResetAt:  now.Add(authFailureWindow),
			LastSeen: now,
		}
	}

	allowed, retryAfter := s.beginAuthAttempt("new-key", now)
	if allowed || retryAfter <= 0 {
		t.Fatalf("capacity saturation allowed=%v retryAfter=%s, want fail-closed", allowed, retryAfter)
	}
	if len(s.authFailures) != maxAuthFailureKeys {
		t.Fatalf("auth failure map size = %d, want %d", len(s.authFailures), maxAuthFailureKeys)
	}
}

func TestClientRegistrationBurstLimitIsAtomicAndExpires(t *testing.T) {
	s := &Server{
		cfg: Config{
			MaxClients:      defaultMaxClients,
			ClientUnusedTTL: defaultClientUnusedTTL,
		},
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState),
	}
	now := time.Now().UTC()
	limit := s.dcrBurstLimit()
	if limit != defaultMaxClients/dcrBurstLimitDivisor {
		t.Fatalf("DCR burst limit = %d", limit)
	}
	small := &Server{cfg: Config{MaxClients: 4, ClientUnusedTTL: time.Minute}}
	if got := small.dcrBurstLimit(); got != 3 {
		t.Fatalf("small DCR burst limit = %d, want 3", got)
	}

	const multiplier = 4
	results := make(chan bool, defaultMaxClients*multiplier)
	var wg sync.WaitGroup
	for i := 0; i < defaultMaxClients*multiplier; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, _ := s.beginClientRegistration(now)
			results <- allowed
		}()
	}
	wg.Wait()
	close(results)

	allowedCount := 0
	for allowed := range results {
		if allowed {
			allowedCount++
		}
	}
	if allowedCount != limit {
		t.Fatalf("concurrent DCR registrations allowed = %d, want %d", allowedCount, limit)
	}
	if allowed, retryAfter := s.beginClientRegistration(now); allowed || retryAfter <= 0 {
		t.Fatalf("DCR burst saturation allowed=%v retryAfter=%s", allowed, retryAfter)
	}

	afterWindow := now.Add(defaultClientUnusedTTL + time.Second)
	if allowed, retryAfter := s.beginClientRegistration(afterWindow); !allowed || retryAfter != 0 {
		t.Fatalf("DCR after window allowed=%v retryAfter=%s", allowed, retryAfter)
	}
}

func TestPendingAuthorizationLimitAndExpiryCleanup(t *testing.T) {
	s := &Server{
		cfg: Config{
			MaxPending:        2,
			PendingRequestTTL: time.Minute,
		},
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState),
	}
	now := time.Now().UTC()
	pending := pendingAuthorization{ClientID: "client", RedirectURI: "https://client.example.test/callback"}

	if _, err := s.createPending(pending, now); err != nil {
		t.Fatalf("create first pending: %v", err)
	}
	if _, err := s.createPending(pending, now); err != nil {
		t.Fatalf("create second pending: %v", err)
	}
	if _, err := s.createPending(pending, now); !errors.Is(err, ErrStateLimitReached) {
		t.Fatalf("third pending error = %v, want ErrStateLimitReached", err)
	}

	if _, err := s.createPending(pending, now.Add(time.Minute+time.Second)); err != nil {
		t.Fatalf("create pending after expiry cleanup: %v", err)
	}
	if len(s.pending) != 1 {
		t.Fatalf("pending count after expiry cleanup = %d, want 1", len(s.pending))
	}
}

func TestPendingAuthorizationPerClientAndPeerCaps(t *testing.T) {
	newServer := func() *Server {
		return &Server{
			cfg: Config{
				MaxPending:        defaultMaxPending,
				PendingRequestTTL: time.Minute,
			},
			pending:      make(map[string]pendingAuthorization),
			authFailures: make(map[string]authFailureState),
		}
	}
	now := time.Now().UTC()

	clientBound := newServer()
	for i := 0; i < maxPendingPerClient; i++ {
		_, err := clientBound.createPending(pendingAuthorization{ClientID: "client-a", PeerKey: "peer-a"}, now)
		if err != nil {
			t.Fatalf("create per-client pending %d: %v", i+1, err)
		}
	}
	if _, err := clientBound.createPending(pendingAuthorization{ClientID: "client-a", PeerKey: "peer-b"}, now); !errors.Is(err, ErrStateLimitReached) {
		t.Fatalf("per-client overflow error = %v", err)
	}

	peerBound := newServer()
	for i := 0; i < maxPendingPerPeer; i++ {
		_, err := peerBound.createPending(pendingAuthorization{
			ClientID: fmt.Sprintf("client-%03d", i),
			PeerKey:  "shared-peer",
		}, now)
		if err != nil {
			t.Fatalf("create per-peer pending %d: %v", i+1, err)
		}
	}
	if _, err := peerBound.createPending(pendingAuthorization{ClientID: "client-overflow", PeerKey: "shared-peer"}, now); !errors.Is(err, ErrStateLimitReached) {
		t.Fatalf("per-peer overflow error = %v", err)
	}

	afterExpiry := now.Add(time.Minute + time.Second)
	if _, err := peerBound.createPending(pendingAuthorization{ClientID: "client-new", PeerKey: "shared-peer"}, afterExpiry); err != nil {
		t.Fatalf("per-peer cap did not recover after expiry: %v", err)
	}
}

func TestAuthorizationCSPAllowsValidatedRedirectOrigin(t *testing.T) {
	tests := []struct {
		name        string
		redirectURI string
		wantOrigin  string
	}{
		{
			name:        "https callback",
			redirectURI: "https://chatgpt.com/connector/oauth/example",
			wantOrigin:  "https://chatgpt.com",
		},
		{
			name:        "loopback callback",
			redirectURI: "http://127.0.0.1:19195/callback",
			wantOrigin:  "http://127.0.0.1:19195",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			secureAuthorizationHTMLHeaders(recorder, tt.redirectURI)
			policy := recorder.Header().Get("Content-Security-Policy")
			want := "form-action 'self' " + tt.wantOrigin + ";"
			if !strings.Contains(policy, want) {
				t.Fatalf("CSP = %q, want %q", policy, want)
			}
			if strings.Contains(policy, "/connector/oauth/example") || strings.Contains(policy, "/callback") {
				t.Fatalf("CSP leaked callback path: %q", policy)
			}
		})
	}
}

func TestAuthorizationCSPRejectsUnsafeOriginSerialization(t *testing.T) {
	recorder := httptest.NewRecorder()
	secureAuthorizationHTMLHeaders(recorder, "https://example.test;script-src-elem/callback")
	policy := recorder.Header().Get("Content-Security-Policy")
	if policy != "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'" {
		t.Fatalf("unsafe redirect CSP = %q", policy)
	}
}

func TestPendingAuthorizationConsumeOnceConcurrently(t *testing.T) {
	s := &Server{
		cfg: Config{
			MaxPending:        4,
			PendingRequestTTL: time.Minute,
		},
		pending:      make(map[string]pendingAuthorization),
		authFailures: make(map[string]authFailureState),
	}
	now := time.Now().UTC()
	requestID, err := s.createPending(pendingAuthorization{
		ClientID:    "client",
		RedirectURI: "https://client.example.test/callback",
	}, now)
	if err != nil {
		t.Fatalf("create pending: %v", err)
	}

	const workers = 16
	start := make(chan struct{})
	results := make(chan bool, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, ok := s.consumePending(requestID, now)
			results <- ok
		}()
	}
	close(start)
	wg.Wait()
	close(results)

	wins := 0
	for ok := range results {
		if ok {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("successful concurrent consumes = %d, want 1", wins)
	}
}

func TestNormalizeScopesTracksAdminToolsCapability(t *testing.T) {
	enabled := false
	server := &Server{
		cfg: Config{
			SupportedScopes: []string{mcpauth.ScopeRead, mcpauth.ScopeAdmin},
			DefaultScopes:   []string{mcpauth.ScopeRead},
		},
		adminToolsEnabledProvider: func() bool { return enabled },
	}

	scopes, err := server.normalizeScopes("")
	if err != nil {
		t.Fatalf("disabled default scopes: %v", err)
	}
	if len(scopes) != 1 || scopes[0] != mcpauth.ScopeRead {
		t.Fatalf("disabled default scopes = %v", scopes)
	}
	if _, err := server.normalizeScopes(mcpauth.ScopeRead + " " + mcpauth.ScopeAdmin); err == nil {
		t.Fatal("disabled explicit admin scope error = nil")
	}

	enabled = true
	scopes, err = server.normalizeScopes("")
	if err != nil {
		t.Fatalf("enabled default scopes: %v", err)
	}
	if len(scopes) != 2 || scopes[0] != mcpauth.ScopeRead || scopes[1] != mcpauth.ScopeAdmin {
		t.Fatalf("enabled default scopes = %v", scopes)
	}
	if _, err := server.normalizeScopes(mcpauth.ScopeRead + " " + mcpauth.ScopeAdmin); err != nil {
		t.Fatalf("enabled explicit admin scope rejected: %v", err)
	}

	enabled = false
	if _, err := server.normalizeScopes(mcpauth.ScopeAdmin); err == nil {
		t.Fatal("disabled admin scope remained usable after runtime toggle")
	}
}
