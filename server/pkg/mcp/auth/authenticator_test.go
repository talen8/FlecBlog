package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestAuthenticatorOAuthAndHybrid(t *testing.T) {
	key := generateRSAKey(t)
	jwks := newRotatingJWKSServer(t, &key.PublicKey, "key-1")
	cfg := testOAuthConfig(jwks.URL + "/jwks")

	oauth, err := NewAuthenticator(cfg, func() string { return "legacy-secret" })
	if err != nil {
		t.Fatalf("NewAuthenticator oauth: %v", err)
	}
	token := signRS256Token(t, key, "key-1", jwt.MapClaims{
		"iss":   cfg.Issuer,
		"aud":   cfg.Audience,
		"sub":   "user-123",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"scope": "mcp:read mcp:draft",
	})
	principal, err := oauth.AuthenticateBearer(context.Background(), token)
	if err != nil {
		t.Fatalf("AuthenticateBearer oauth: %v", err)
	}
	if principal.Method != "oauth" || principal.Subject != "user-123" {
		t.Fatalf("oauth principal = %+v", principal)
	}
	if !principal.HasScope(ScopeRead) || !principal.HasScope(ScopeDraft) || principal.HasScope(ScopePublish) {
		t.Fatalf("oauth scopes = %v", principal.Scopes)
	}
	hybridCfg := cfg
	hybridCfg.Mode = ModeHybrid
	hybrid, err := NewAuthenticator(hybridCfg, func() string { return "legacy-secret" })
	if err != nil {
		t.Fatalf("NewAuthenticator hybrid: %v", err)
	}
	staticPrincipal, err := hybrid.AuthenticateBearer(context.Background(), "legacy-secret")
	if err != nil {
		t.Fatalf("AuthenticateBearer hybrid static: %v", err)
	}
	if staticPrincipal.Method != "static" || !staticPrincipal.HasScope(ScopeAdmin) {
		t.Fatalf("hybrid static principal = %+v", staticPrincipal)
	}
	oauthPrincipal, err := hybrid.AuthenticateBearer(context.Background(), token)
	if err != nil {
		t.Fatalf("AuthenticateBearer hybrid oauth: %v", err)
	}
	if oauthPrincipal.Method != "oauth" || oauthPrincipal.Subject != "user-123" {
		t.Fatalf("hybrid oauth principal = %+v", oauthPrincipal)
	}
}

func TestAuthenticatorExternalOAuthRawAdminClaimsDoNotBindLocalUser(t *testing.T) {
	key := generateRSAKey(t)
	jwks := newRotatingJWKSServer(t, &key.PublicKey, "key-1")
	cfg := testOAuthConfig(jwks.URL + "/jwks")

	authenticator, err := NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	token := signRS256Token(t, key, "key-1", jwt.MapClaims{
		"iss":         cfg.Issuer,
		"aud":         cfg.Audience,
		"sub":         "external-admin-subject",
		"exp":         time.Now().Add(time.Hour).Unix(),
		"scope":       ScopeAdmin,
		"email":       "local-admin@example.test",
		"mcp_user_id": float64(42),
	})

	principal, err := authenticator.AuthenticateBearer(context.Background(), token)
	if err != nil {
		t.Fatalf("AuthenticateBearer: %v", err)
	}
	if !principal.HasScope(ScopeAdmin) {
		t.Fatal("external OAuth admin scope was not preserved")
	}
	if userID, ok := principal.VerifiedLocalUserID(); ok {
		t.Fatalf("external OAuth raw claims created verified local binding user_id=%d", userID)
	}
}

func TestAuthenticatorRejectsInvalidOAuthTokens(t *testing.T) {
	key := generateRSAKey(t)
	jwks := newRotatingJWKSServer(t, &key.PublicKey, "key-1")
	cfg := testOAuthConfig(jwks.URL + "/jwks")
	authenticator, err := NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}

	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "wrong audience",
			claims: jwt.MapClaims{
				"iss": cfg.Issuer, "aud": "https://other.example/mcp", "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "wrong issuer",
			claims: jwt.MapClaims{
				"iss": "https://other.example", "aud": cfg.Audience, "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "expired",
			claims: jwt.MapClaims{
				"iss": cfg.Issuer, "aud": cfg.Audience, "sub": "user-1", "exp": time.Now().Add(-time.Minute).Unix(),
			},
		},
		{
			name: "missing exp",
			claims: jwt.MapClaims{
				"iss": cfg.Issuer, "aud": cfg.Audience, "sub": "user-1",
			},
		},
		{
			name: "missing subject",
			claims: jwt.MapClaims{
				"iss": cfg.Issuer, "aud": cfg.Audience, "exp": time.Now().Add(time.Hour).Unix(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := signRS256Token(t, key, "key-1", tc.claims)
			if _, err := authenticator.AuthenticateBearer(context.Background(), token); err == nil {
				t.Fatal("AuthenticateBearer error = nil, want rejection")
			}
		})
	}

	hsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": cfg.Issuer, "aud": cfg.Audience, "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
	})
	hsToken.Header["kid"] = "key-1"
	hsSigned, err := hsToken.SignedString([]byte("not-a-public-key"))
	if err != nil {
		t.Fatalf("sign HS256 token: %v", err)
	}
	if _, err := authenticator.AuthenticateBearer(context.Background(), hsSigned); err == nil {
		t.Fatal("HS256 token error = nil, want rejection")
	}
}

func TestAuthenticatorRefreshesJWKSOnUnknownKid(t *testing.T) {
	key1 := generateRSAKey(t)
	key2 := generateRSAKey(t)
	var mu sync.RWMutex
	currentKey := &key1.PublicKey
	currentKid := "key-1"
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		key := currentKey
		kid := currentKid
		mu.Unlock()
		writeJWKS(t, w, key, kid)
	}))
	defer server.Close()

	cfg := testOAuthConfig(server.URL)
	authenticator, err := NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	claims := func() jwt.MapClaims {
		return jwt.MapClaims{"iss": cfg.Issuer, "aud": cfg.Audience, "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix()}
	}
	if _, err := authenticator.AuthenticateBearer(context.Background(), signRS256Token(t, key1, "key-1", claims())); err != nil {
		t.Fatalf("authenticate first key: %v", err)
	}

	mu.Lock()
	currentKey = &key2.PublicKey
	currentKid = "key-2"
	mu.Unlock()
	if _, err := authenticator.AuthenticateBearer(context.Background(), signRS256Token(t, key2, "key-2", claims())); err != nil {
		t.Fatalf("authenticate rotated key: %v", err)
	}

	key3 := generateRSAKey(t)
	if _, err := authenticator.AuthenticateBearer(context.Background(), signRS256Token(t, key3, "random-kid", claims())); err == nil {
		t.Fatal("random unknown kid error = nil, want rejection")
	}
	mu.RLock()
	gotRequests := requestCount
	mu.RUnlock()
	if gotRequests != 2 {
		t.Fatalf("JWKS request count = %d, want 2 after rate-limited unknown kid", gotRequests)
	}
}

func TestJWKSClientRejectsRedirect(t *testing.T) {
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"keys":[]}`))
	}))
	defer redirectTarget.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget.URL, http.StatusFound)
	}))
	defer redirector.Close()

	key := generateRSAKey(t)
	cfg := testOAuthConfig(redirector.URL)
	authenticator, err := NewAuthenticator(cfg, nil)
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	token := signRS256Token(t, key, "key-1", jwt.MapClaims{
		"iss": cfg.Issuer, "aud": cfg.Audience, "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
	})
	if _, err := authenticator.AuthenticateBearer(context.Background(), token); err == nil {
		t.Fatal("redirecting JWKS error = nil, want rejection")
	}
}

func testOAuthConfig(jwksURL string) Config {
	return Config{
		Mode:            ModeOAuth,
		ResourceURI:     "https://mcp.example.test/mcp",
		Issuer:          "https://issuer.example.test/",
		JWKSURL:         jwksURL,
		Audience:        "https://mcp.example.test/mcp",
		MetadataURL:     "https://mcp.example.test/.well-known/oauth-protected-resource/mcp",
		MetadataPath:    "/.well-known/oauth-protected-resource/mcp",
		SupportedScopes: append([]string(nil), defaultSupportedScopes...),
		ChallengeScopes: append([]string(nil), defaultChallengeScopes...),
	}
}

func generateRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	return key
}

func signRS256Token(t *testing.T, key *rsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign RS256 token: %v", err)
	}
	return signed
}

func newRotatingJWKSServer(t *testing.T, key *rsa.PublicKey, kid string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jwks" {
			http.NotFound(w, r)
			return
		}
		writeJWKS(t, w, key, kid)
	}))
}

func writeJWKS(t *testing.T, w http.ResponseWriter, key *rsa.PublicKey, kid string) {
	t.Helper()
	exponent := big.NewInt(int64(key.E)).Bytes()
	document := map[string]any{
		"keys": []map[string]any{{
			"kty": "RSA",
			"kid": kid,
			"use": "sig",
			"alg": "RS256",
			"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(exponent),
		}},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(document); err != nil {
		t.Fatalf("encode JWKS: %v", err)
	}
}

func TestExtractScopesSupportsScopeAndScpClaims(t *testing.T) {
	claims := jwt.MapClaims{
		"scope": "mcp:read mcp:draft",
		"scp":   []any{"mcp:publish", "mcp:manage extra"},
	}
	scopes := extractScopes(claims)
	for _, scope := range []string{"mcp:read", "mcp:draft", "mcp:publish", "mcp:manage", "extra"} {
		if _, ok := scopes[scope]; !ok {
			t.Errorf("scope %q missing from %v", scope, scopes)
		}
	}
}

func TestStaticSecretComparisonIsExact(t *testing.T) {
	authenticator, err := NewAuthenticator(Config{Mode: ModeStatic}, func() string { return "exact-secret" })
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	if _, err := authenticator.AuthenticateBearer(context.Background(), "exact-secret"); err != nil {
		t.Fatalf("exact static secret rejected: %v", err)
	}
	for _, token := range []string{"", "exact-secret ", "EXACT-secret", "exact-secret-x"} {
		if _, err := authenticator.AuthenticateBearer(context.Background(), token); err == nil {
			t.Errorf("static token %q unexpectedly accepted", token)
		}
	}
}

func TestNoSensitiveTokenMaterialInErrors(t *testing.T) {
	authenticator, err := NewAuthenticator(Config{Mode: ModeStatic}, func() string { return "secret" })
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	_, err = authenticator.AuthenticateBearer(context.Background(), "super-sensitive-token")
	if err == nil {
		t.Fatal("invalid token error = nil")
	}
	if strings.Contains(err.Error(), "super-sensitive-token") {
		t.Fatalf("error leaks token material: %v", err)
	}
}

func TestAuthenticatorRejectsOversizedBearerToken(t *testing.T) {
	authenticator, err := NewAuthenticator(Config{Mode: ModeStatic}, func() string { return "secret" })
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	if _, err := authenticator.AuthenticateBearer(context.Background(), strings.Repeat("x", maxBearerTokenBytes+1)); err == nil {
		t.Fatal("oversized bearer token error = nil, want rejection")
	}
}

func TestParseRSAJWKRejectsWeakModulus(t *testing.T) {
	weakKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("generate weak RSA key: %v", err)
	}
	raw := rawJWK{
		Kty: "RSA",
		N:   base64.RawURLEncoding.EncodeToString(weakKey.PublicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(weakKey.PublicKey.E)).Bytes()),
	}
	if _, err := parseRSAJWK(raw); err == nil {
		t.Fatal("weak RSA JWK error = nil, want rejection")
	}
}
