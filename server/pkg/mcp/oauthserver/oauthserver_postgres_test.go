package oauthserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"flec_blog/api/middleware"
	"flec_blog/config"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/internal/testutil"
	blogmcp "flec_blog/pkg/mcp"
	mcpauth "flec_blog/pkg/mcp/auth"

	"github.com/gin-gonic/gin"
	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var requestIDPattern = regexp.MustCompile(`name="request_id" value="([^"]+)"`)

type oauthBearerTransport struct {
	base  http.RoundTripper
	token string
}

func (t oauthBearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(clone)
}

func TestEmbeddedOAuthFlowOverPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	password := "oauth-e2e-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &model.User{
		Email:       fmt.Sprintf("mcp-oauth-%d@example.test", time.Now().UnixNano()),
		Password:    string(hash),
		HasPassword: true,
		Nickname:    "MCP OAuth E2E",
		Role:        model.RoleAdmin,
		IsEnabled:   true,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create OAuth test user: %v", err)
	}
	t.Cleanup(func() { _ = db.Unscoped().Delete(&model.User{}, user.ID).Error })

	userService := service.NewUserService(repository.NewUserRepository(db), nil, &config.Config{})
	keyPath, _ := writeTestSigningKey(t, 0o600)

	var embedded *Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/oauth-authorization-server", func(w http.ResponseWriter, r *http.Request) {
		embedded.AuthorizationServerMetadataHandler().ServeHTTP(w, r)
	})
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		embedded.JWKSHandler().ServeHTTP(w, r)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		embedded.RegistrationHandler().ServeHTTP(w, r)
	})
	mux.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		embedded.AuthorizationHandler().ServeHTTP(w, r)
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		embedded.TokenHandler().ServeHTTP(w, r)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resourceURI := ts.URL + "/mcp"
	cfg := Config{
		Enabled:               true,
		Issuer:                ts.URL,
		ResourceURI:           resourceURI,
		Audience:              resourceURI,
		AuthorizationEndpoint: ts.URL + "/authorize",
		TokenEndpoint:         ts.URL + "/token",
		RegistrationEndpoint:  ts.URL + "/register",
		JWKSURI:               ts.URL + "/.well-known/jwks.json",
		SigningKeyFile:        keyPath,
		SupportedScopes:       []string{mcpauth.ScopeRead, mcpauth.ScopeDraft, mcpauth.ScopePublish, mcpauth.ScopeManage, mcpauth.ScopeAdmin, ScopeOfflineAccess},
		DefaultScopes:         []string{mcpauth.ScopeRead},
		AccessTokenTTL:        time.Hour,
		RefreshTokenTTL:       24 * time.Hour,
		ClientInactiveTTL:     24 * time.Hour,
		ClientUnusedTTL:       15 * time.Minute,
		AuthCodeTTL:           5 * time.Minute,
		PendingRequestTTL:     10 * time.Minute,
		MaxClients:            16,
		MaxPending:            16,
		MaxAuthCodes:          32,
		MaxRefreshTokens:      32,
	}
	embedded, err = NewWithOptions(
		cfg,
		db,
		userService,
		ServerOptions{AdminToolsEnabledProvider: func() bool { return true }},
	)
	if err != nil {
		t.Fatalf("New embedded OAuth: %v", err)
	}

	client := registerTestClient(t, ts.URL+"/register")
	verifier := strings.Repeat("a", 64)
	scopes := strings.Join([]string{mcpauth.ScopeRead, mcpauth.ScopeAdmin}, " ")
	authURL := ts.URL + "/authorize?" + url.Values{
		"response_type":         {"code"},
		"client_id":             {client.ClientID},
		"redirect_uri":          {"https://client.example.test/callback"},
		"resource":              {resourceURI},
		"scope":                 {scopes},
		"state":                 {"state-123"},
		"code_challenge":        {pkceChallenge(verifier)},
		"code_challenge_method": {"S256"},
	}.Encode()

	requestID := fetchAuthorizationRequestID(t, authURL, "https://client.example.test/callback")
	code := submitAuthorizationLogin(t, ts.URL+"/authorize", requestID, user.Email, password)
	assertTokenError(t, ts.URL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://client.example.test/callback"},
		"resource":      {resourceURI},
		"code_verifier": {strings.Repeat("b", 64)},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	}, "invalid_grant")
	first := exchangeAuthorizationCode(t, ts.URL+"/token", client, code, verifier, resourceURI)
	if strings.Contains(first.Scope, ScopeOfflineAccess) {
		t.Fatalf("token scope unexpectedly contains %q: %q", ScopeOfflineAccess, first.Scope)
	}
	if first.AccessToken == "" || first.RefreshToken == "" {
		t.Fatalf("token response without offline_access must include access and refresh tokens: %+v", first)
	}

	resourceAuth, err := mcpauth.NewAuthenticatorWithOptions(mcpauth.Config{
		Mode:        mcpauth.ModeOAuth,
		ResourceURI: resourceURI,
		Issuer:      ts.URL,
		Audience:    resourceURI,
	}, nil, mcpauth.AuthenticatorOptions{
		OAuthKeyResolver:        embedded.OAuthKeyResolver(),
		OAuthPrincipalValidator: embedded.OAuthPrincipalValidator(),
	})
	if err != nil {
		t.Fatalf("NewAuthenticator: %v", err)
	}
	gin.SetMode(gin.TestMode)
	mcpRouter := gin.New()
	mcpHTTPHandler := blogmcp.NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	mcpHTTPHandler = sdkauth.RequireBearerToken(mcpauth.SDKTokenVerifierFromPrincipalContext, nil)(mcpHTTPHandler)
	mcpRouter.Any("/mcp", middleware.MCPResourceAuth(resourceAuth), gin.WrapH(mcpHTTPHandler))
	mux.Handle("/mcp", mcpRouter)

	principal, err := resourceAuth.AuthenticateBearer(context.Background(), first.AccessToken)
	if err != nil {
		t.Fatalf("AuthenticateBearer: %v", err)
	}
	if principal.Method != "oauth" || principal.Subject != fmt.Sprintf("user:%d", user.ID) || !principal.HasScope(mcpauth.ScopeRead) || !principal.HasScope(mcpauth.ScopeAdmin) {
		t.Fatalf("principal = %+v", principal)
	}
	if got := principal.Claims["mcp_user_id"]; got != float64(user.ID) {
		t.Fatalf("mcp_user_id claim = %#v", got)
	}
	if boundUserID, ok := principal.VerifiedLocalUserID(); !ok || boundUserID != user.ID {
		t.Fatalf("verified local binding = (%d, %v), want (%d, true)", boundUserID, ok, user.ID)
	}
	probeMCPBearer(t, ts.URL+"/mcp", first.AccessToken)

	assertTokenError(t, ts.URL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://client.example.test/callback"},
		"resource":      {resourceURI},
		"code_verifier": {verifier},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	}, "invalid_grant")

	otherClient := registerTestClient(t, ts.URL+"/register")
	assertTokenError(t, ts.URL+"/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {first.RefreshToken},
		"resource":      {resourceURI},
		"client_id":     {otherClient.ClientID},
		"client_secret": {otherClient.ClientSecret},
	}, "invalid_grant")

	second := refreshAccessToken(t, ts.URL+"/token", client, first.RefreshToken, resourceURI)
	if second.RefreshToken == "" || second.RefreshToken == first.RefreshToken {
		t.Fatalf("refresh rotation failed: first=%q second=%q", first.RefreshToken, second.RefreshToken)
	}
	assertTokenError(t, ts.URL+"/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {first.RefreshToken},
		"resource":      {resourceURI},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	}, "invalid_grant")

	publicClient := registerTestClientWithMethod(t, ts.URL+"/register", "none")
	publicVerifier := strings.Repeat("c", 64)
	publicAuthURL := ts.URL + "/authorize?" + url.Values{
		"response_type":         {"code"},
		"client_id":             {publicClient.ClientID},
		"redirect_uri":          {"https://client.example.test/callback"},
		"resource":              {resourceURI},
		"scope":                 {mcpauth.ScopeRead},
		"state":                 {"state-123"},
		"code_challenge":        {pkceChallenge(publicVerifier)},
		"code_challenge_method": {"S256"},
	}.Encode()
	publicRequestID := fetchAuthorizationRequestID(t, publicAuthURL, "https://client.example.test/callback")
	publicCode := submitAuthorizationLogin(t, ts.URL+"/authorize", publicRequestID, user.Email, password)
	publicFirst := exchangePublicAuthorizationCode(t, ts.URL+"/token", publicClient, publicCode, publicVerifier, resourceURI)
	if publicFirst.AccessToken == "" || publicFirst.RefreshToken == "" {
		t.Fatalf("public client token response missing tokens: %+v", publicFirst)
	}
	publicSecond := refreshPublicAccessToken(t, ts.URL+"/token", publicClient, publicFirst.RefreshToken, resourceURI)
	if publicSecond.RefreshToken == "" || publicSecond.RefreshToken == publicFirst.RefreshToken {
		t.Fatalf("public refresh rotation failed: first=%q second=%q", publicFirst.RefreshToken, publicSecond.RefreshToken)
	}
	assertTokenError(t, ts.URL+"/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {publicFirst.RefreshToken},
		"resource":      {resourceURI},
		"client_id":     {publicClient.ClientID},
	}, "invalid_grant")

	if err := db.Model(&model.User{}).Where("id = ?", user.ID).Update("is_enabled", false).Error; err != nil {
		t.Fatalf("disable test user: %v", err)
	}
	if _, err := resourceAuth.AuthenticateBearer(context.Background(), second.AccessToken); err == nil {
		t.Fatal("disabled local user access token was accepted")
	}
	assertTokenError(t, ts.URL+"/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {second.RefreshToken},
		"resource":      {resourceURI},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	}, "invalid_grant")
}

type testClientRegistration struct {
	ClientID                string `json:"client_id"`
	ClientSecret            string `json:"client_secret"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
}

func probeMCPBearer(t *testing.T, endpoint, token string) {
	t.Helper()
	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "embedded-oauth-postgres-e2e", Version: "1.0"}, nil)
	httpClient := &http.Client{Transport: oauthBearerTransport{base: http.DefaultTransport, token: token}}
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: endpoint, HTTPClient: httpClient}, nil)
	if err != nil {
		t.Fatalf("connect MCP with embedded OAuth token: %v", err)
	}
	defer session.Close()
	listed, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list MCP tools with embedded OAuth token: %v", err)
	}
	if len(listed.Tools) == 0 {
		t.Fatal("embedded OAuth MCP tools/list returned no tools")
	}
}

func registerTestClient(t *testing.T, endpoint string) testClientRegistration {
	t.Helper()
	return registerTestClientWithMethod(t, endpoint, "client_secret_post")
}

func registerTestClientWithMethod(t *testing.T, endpoint, authMethod string) testClientRegistration {
	t.Helper()
	body := fmt.Sprintf(`{"redirect_uris":["https://client.example.test/callback"],"grant_types":["authorization_code","refresh_token"],"response_types":["code"],"token_endpoint_auth_method":%q}`, authMethod)
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("register client: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d", resp.StatusCode)
	}
	var result testClientRegistration
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode registration: %v", err)
	}
	if result.TokenEndpointAuthMethod != authMethod {
		t.Fatalf("registered auth method = %q, want %q", result.TokenEndpointAuthMethod, authMethod)
	}
	if authMethod == "none" && result.ClientSecret != "" {
		t.Fatal("public client registration unexpectedly returned client_secret")
	}
	if authMethod == "client_secret_post" && result.ClientSecret == "" {
		t.Fatal("confidential client registration missing client_secret")
	}
	return result
}

func fetchAuthorizationRequestID(t *testing.T, endpoint, expectedRedirectURI string) string {
	t.Helper()
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("authorization GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("authorization GET status = %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read authorization page: %v", err)
	}
	if !strings.Contains(string(body), expectedRedirectURI) {
		t.Fatalf("authorization page missing redirect URI %q", expectedRedirectURI)
	}
	match := requestIDPattern.FindStringSubmatch(string(body))
	if len(match) != 2 || match[1] == "" {
		t.Fatalf("authorization page missing request_id")
	}
	return match[1]
}

func submitAuthorizationLogin(t *testing.T, endpoint, requestID, email, password string) string {
	t.Helper()
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.PostForm(endpoint, url.Values{
		"request_id": {requestID},
		"email":      {email},
		"password":   {password},
	})
	if err != nil {
		t.Fatalf("authorization POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("authorization POST status = %d", resp.StatusCode)
	}
	location, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		t.Fatalf("parse redirect location: %v", err)
	}
	if location.Query().Get("state") != "state-123" {
		t.Fatalf("redirect state = %q", location.Query().Get("state"))
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatal("redirect missing authorization code")
	}
	return code
}

func exchangeAuthorizationCode(t *testing.T, endpoint string, client testClientRegistration, code, verifier, resource string) tokenResponse {
	t.Helper()
	return postTokenRequest(t, endpoint, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://client.example.test/callback"},
		"resource":      {resource},
		"code_verifier": {verifier},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	})
}

func refreshAccessToken(t *testing.T, endpoint string, client testClientRegistration, refreshToken, resource string) tokenResponse {
	t.Helper()
	return postTokenRequest(t, endpoint, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"resource":      {resource},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	})
}

func exchangePublicAuthorizationCode(t *testing.T, endpoint string, client testClientRegistration, code, verifier, resource string) tokenResponse {
	t.Helper()
	return postTokenRequest(t, endpoint, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://client.example.test/callback"},
		"resource":      {resource},
		"code_verifier": {verifier},
		"client_id":     {client.ClientID},
	})
}

func refreshPublicAccessToken(t *testing.T, endpoint string, client testClientRegistration, refreshToken, resource string) tokenResponse {
	t.Helper()
	return postTokenRequest(t, endpoint, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"resource":      {resource},
		"client_id":     {client.ClientID},
	})
}

func postTokenRequest(t *testing.T, endpoint string, form url.Values) tokenResponse {
	t.Helper()
	resp, err := http.PostForm(endpoint, form)
	if err != nil {
		t.Fatalf("token request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var problem oauthErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&problem)
		t.Fatalf("token status = %d problem=%+v", resp.StatusCode, problem)
	}
	var result tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	return result
}

func assertTokenError(t *testing.T, endpoint string, form url.Values, wantCode string) {
	t.Helper()
	resp, err := http.PostForm(endpoint, form)
	if err != nil {
		t.Fatalf("token error request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		t.Fatalf("token error status = %d, want failure", resp.StatusCode)
	}
	var problem oauthErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		t.Fatalf("decode token error: %v", err)
	}
	if problem.Error != wantCode {
		t.Fatalf("token error = %q, want %q", problem.Error, wantCode)
	}
}

func TestOAuthRepositoryReclaimsUnusedClientAfterGraceOverPostgres(t *testing.T) {
	db := openOAuthTestPostgres(t)
	repo := NewRepository(db)
	now := time.Now().UTC()

	var existingCount int64
	if err := db.Model(&OAuthClient{}).Count(&existingCount).Error; err != nil {
		t.Fatalf("count existing clients: %v", err)
	}
	staleID := "mcp-stale-" + randomToken(8)
	newID := "mcp-new-" + randomToken(8)
	stale := &OAuthClient{
		ClientID:         staleID,
		ClientSecretHash: secretHash(randomToken(16)),
		RedirectURIs:     `[]`,
		CreatedAt:        now.Add(-90 * 24 * time.Hour),
	}
	if err := db.Create(stale).Error; err != nil {
		t.Fatalf("create stale client: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Where("client_id IN ?", []string{staleID, newID}).Delete(&OAuthClient{}).Error
	})

	fresh := &OAuthClient{
		ClientID:         newID,
		ClientSecretHash: secretHash(randomToken(16)),
		RedirectURIs:     `[]`,
		CreatedAt:        now,
	}
	if err := repo.CreateClientBounded(fresh, int(existingCount)+1, now.Add(-30*24*time.Hour), now.Add(-15*time.Minute)); err != nil {
		t.Fatalf("CreateClientBounded with reclaim: %v", err)
	}
	var staleCount int64
	if err := db.Model(&OAuthClient{}).Where("client_id = ?", staleID).Count(&staleCount).Error; err != nil {
		t.Fatalf("count stale client: %v", err)
	}
	if staleCount != 0 {
		t.Fatalf("stale client count = %d, want 0", staleCount)
	}
	var freshCount int64
	if err := db.Model(&OAuthClient{}).Where("client_id = ?", newID).Count(&freshCount).Error; err != nil {
		t.Fatalf("count fresh client: %v", err)
	}
	if freshCount != 1 {
		t.Fatalf("fresh client count = %d, want 1", freshCount)
	}
}

func TestOAuthRepositoryProtectsFreshUnusedClientOverPostgres(t *testing.T) {
	db := openOAuthTestPostgres(t)
	repo := NewRepository(db)
	now := time.Now().UTC()

	var existingCount int64
	if err := db.Model(&OAuthClient{}).Count(&existingCount).Error; err != nil {
		t.Fatalf("count existing clients: %v", err)
	}
	protectedID := "mcp-protected-" + randomToken(8)
	newID := "mcp-overflow-" + randomToken(8)
	protected := &OAuthClient{
		ClientID:         protectedID,
		ClientSecretHash: secretHash(randomToken(16)),
		RedirectURIs:     `[]`,
		CreatedAt:        now,
	}
	if err := db.Create(protected).Error; err != nil {
		t.Fatalf("create protected client: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Where("client_id IN ?", []string{protectedID, newID}).Delete(&OAuthClient{}).Error
	})

	overflow := &OAuthClient{
		ClientID:         newID,
		ClientSecretHash: secretHash(randomToken(16)),
		RedirectURIs:     `[]`,
		CreatedAt:        now,
	}
	err := repo.CreateClientBounded(overflow, int(existingCount)+1, time.Time{}, time.Time{})
	if !errors.Is(err, ErrStateLimitReached) {
		t.Fatalf("CreateClientBounded error = %v, want ErrStateLimitReached", err)
	}
	var protectedCount int64
	if err := db.Model(&OAuthClient{}).Where("client_id = ?", protectedID).Count(&protectedCount).Error; err != nil {
		t.Fatalf("count protected client: %v", err)
	}
	if protectedCount != 1 {
		t.Fatalf("protected client count = %d, want 1", protectedCount)
	}
	var overflowCount int64
	if err := db.Model(&OAuthClient{}).Where("client_id = ?", newID).Count(&overflowCount).Error; err != nil {
		t.Fatalf("count overflow client: %v", err)
	}
	if overflowCount != 0 {
		t.Fatalf("overflow client count = %d, want 0", overflowCount)
	}
}

func TestRefreshRotationRollbackPreservesOldTokenOverPostgres(t *testing.T) {
	db := openOAuthTestPostgres(t)
	now := time.Now().UTC()

	user := &model.User{
		Email:       fmt.Sprintf("mcp-rotation-%d@example.test", time.Now().UnixNano()),
		Password:    "unused",
		HasPassword: true,
		Nickname:    "MCP Rotation E2E",
		Role:        model.RoleAdmin,
		IsEnabled:   true,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create rotation user: %v", err)
	}
	clientID := "mcp-rotation-" + randomToken(8)
	client := &OAuthClient{
		ClientID:         clientID,
		ClientSecretHash: secretHash(randomToken(16)),
		RedirectURIs:     `[]`,
		CreatedAt:        now,
	}
	if err := db.Create(client).Error; err != nil {
		t.Fatalf("create rotation client: %v", err)
	}
	oldHash := secretHash("old-" + randomToken(16))
	conflictHash := secretHash("conflict-" + randomToken(16))
	oldToken := &RefreshToken{TokenHash: oldHash, ClientID: clientID, UserID: user.ID, Scopes: mcpauth.ScopeRead, ExpiresAt: now.Add(time.Hour), CreatedAt: now}
	conflict := &RefreshToken{TokenHash: conflictHash, ClientID: clientID, UserID: user.ID, Scopes: mcpauth.ScopeRead, ExpiresAt: now.Add(time.Hour), CreatedAt: now}
	if err := db.Create(oldToken).Error; err != nil {
		t.Fatalf("create old refresh token: %v", err)
	}
	if err := db.Create(conflict).Error; err != nil {
		t.Fatalf("create conflicting refresh token: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Where("token_hash IN ?", []string{oldHash, conflictHash}).Delete(&RefreshToken{}).Error
		_ = db.Where("client_id = ?", clientID).Delete(&OAuthClient{}).Error
		_ = db.Unscoped().Delete(&model.User{}, user.ID).Error
	})

	replacement := &RefreshToken{TokenHash: conflictHash, ClientID: clientID, UserID: user.ID, Scopes: mcpauth.ScopeRead, ExpiresAt: now.Add(2 * time.Hour), CreatedAt: now}
	if err := NewRepository(db).RotateRefreshToken(oldHash, clientID, user.ID, replacement); err == nil {
		t.Fatal("RotateRefreshToken conflict error = nil")
	}
	var oldCount int64
	if err := db.Model(&RefreshToken{}).Where("token_hash = ?", oldHash).Count(&oldCount).Error; err != nil {
		t.Fatalf("count old refresh token: %v", err)
	}
	if oldCount != 1 {
		t.Fatalf("old refresh token count = %d, want 1 after rollback", oldCount)
	}
}

func openOAuthTestPostgres(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.OpenMCPTestPostgres(t)
	return db
}
