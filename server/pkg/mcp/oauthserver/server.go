package oauthserver

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"flec_blog/internal/model"
	"flec_blog/internal/service"
	mcpauth "flec_blog/pkg/mcp/auth"

	"gorm.io/gorm"
)

const (
	maxRequestBodyBytes  = 64 * 1024
	maxRedirectURIs      = 10
	maxRedirectURIBytes  = 2048
	maxClientNameBytes   = 256
	maxPendingPerClient  = 8
	maxPendingPerPeer    = 32
	authFailureLimit     = 5
	authFailureWindow    = 15 * time.Minute
	maxAuthFailureKeys   = 1024
	dcrBurstLimitDivisor = 8
	minimumDCRBurstLimit = 4
)

type Server struct {
	cfg                       Config
	repo                      *Repository
	signer                    *signer
	userService               *service.UserService
	adminToolsEnabledProvider func() bool

	mu                sync.Mutex
	pending           map[string]pendingAuthorization
	authFailures      map[string]authFailureState
	registrationTimes []time.Time
}

type pendingAuthorization struct {
	ClientID      string
	PeerKey       string
	RedirectURI   string
	State         string
	Resource      string
	Scopes        []string
	CodeChallenge string
	ExpiresAt     time.Time
}

type authFailureState struct {
	Count    int
	InFlight int
	ResetAt  time.Time
	LastSeen time.Time
}

type ServerOptions struct {
	AdminToolsEnabledProvider func() bool
}

func New(cfg Config, db *gorm.DB, userService *service.UserService) (*Server, error) {
	return NewWithOptions(cfg, db, userService, ServerOptions{})
}

func NewWithOptions(cfg Config, db *gorm.DB, userService *service.UserService, options ServerOptions) (*Server, error) {
	if err := validateServerConfig(cfg); err != nil {
		return nil, err
	}
	if db == nil || userService == nil {
		return nil, errors.New("embedded OAuth server requires database and user service")
	}
	signer, err := loadSigner(cfg.SigningKeyFile)
	if err != nil {
		return nil, err
	}
	server := &Server{
		cfg:                       cfg,
		repo:                      NewRepository(db),
		signer:                    signer,
		userService:               userService,
		adminToolsEnabledProvider: options.AdminToolsEnabledProvider,
		pending:                   make(map[string]pendingAuthorization),
		authFailures:              make(map[string]authFailureState),
	}
	if err := server.repo.CleanupExpired(time.Now()); err != nil {
		return nil, fmt.Errorf("cleanup MCP OAuth state: %w", err)
	}
	return server, nil
}

func validateServerConfig(cfg Config) error {
	if !cfg.Enabled {
		return errors.New("embedded OAuth server is disabled")
	}
	if cfg.Issuer == "" || cfg.ResourceURI == "" || cfg.Audience == "" || cfg.SigningKeyFile == "" {
		return errors.New("embedded OAuth server configuration is incomplete")
	}
	if cfg.AuthorizationEndpoint == "" || cfg.TokenEndpoint == "" || cfg.RegistrationEndpoint == "" || cfg.JWKSURI == "" {
		return errors.New("embedded OAuth server endpoint configuration is incomplete")
	}
	if cfg.AccessTokenTTL <= 0 || cfg.RefreshTokenTTL <= 0 || cfg.ClientInactiveTTL <= 0 || cfg.ClientUnusedTTL <= 0 || cfg.AuthCodeTTL <= 0 || cfg.PendingRequestTTL <= 0 {
		return errors.New("embedded OAuth server TTL values must be positive")
	}
	if cfg.MaxClients <= 0 || cfg.MaxPending <= 0 || cfg.MaxAuthCodes <= 0 || cfg.MaxRefreshTokens <= 0 {
		return errors.New("embedded OAuth server state limits must be positive")
	}
	if len(cfg.SupportedScopes) == 0 || !hasMCPAccessScope(cfg.SupportedScopes) {
		return errors.New("embedded OAuth server must support at least one MCP scope")
	}
	supported := make(map[string]struct{}, len(cfg.SupportedScopes))
	for _, scope := range cfg.SupportedScopes {
		if scope == "" {
			return errors.New("embedded OAuth server scopes must not be empty")
		}
		supported[scope] = struct{}{}
	}
	for _, scope := range cfg.DefaultScopes {
		if _, ok := supported[scope]; !ok {
			return fmt.Errorf("embedded OAuth default scope %q is unsupported", scope)
		}
	}
	return nil
}

func (s *Server) Config() Config { return s.cfg }

func (s *Server) isAllowedRole(role model.UserRole) bool {
	return role == model.RoleAdmin || role == model.RoleSuperAdmin
}

func (s *Server) adminToolsEnabled() bool {
	return s.adminToolsEnabledProvider != nil && s.adminToolsEnabledProvider()
}

func (s *Server) normalizeScopes(raw string) ([]string, error) {
	requested := strings.Fields(raw)
	if len(requested) == 0 {
		requested = append([]string(nil), s.cfg.DefaultScopes...)
		if s.adminToolsEnabled() && !hasScope(requested, mcpauth.ScopeAdmin) {
			requested = append(requested, mcpauth.ScopeAdmin)
		}
	}
	supported := make(map[string]struct{}, len(s.cfg.SupportedScopes))
	for _, scope := range s.cfg.SupportedScopes {
		supported[scope] = struct{}{}
	}
	seen := make(map[string]struct{}, len(requested))
	normalized := make([]string, 0, len(requested))
	for _, scope := range requested {
		if scope == mcpauth.ScopeAdmin && !s.adminToolsEnabled() {
			return nil, fmt.Errorf("scope %q is disabled", scope)
		}
		if _, ok := supported[scope]; !ok {
			return nil, fmt.Errorf("unsupported scope %q", scope)
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		normalized = append(normalized, scope)
	}
	if len(normalized) == 0 {
		return nil, errors.New("at least one scope is required")
	}
	return normalized, nil
}

func scopeString(scopes []string) string {
	return strings.Join(scopes, " ")
}

func hasScope(scopes []string, wanted string) bool {
	for _, scope := range scopes {
		if scope == wanted {
			return true
		}
	}
	return false
}

func parseStoredScopes(raw string) []string {
	return strings.Fields(raw)
}

func validateRedirectURI(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		return errors.New("redirect URI must be absolute")
	}
	if parsed.User != nil || parsed.Fragment != "" {
		return errors.New("redirect URI has unsupported components")
	}
	if parsed.Scheme == "https" {
		return nil
	}
	if parsed.Scheme != "http" {
		return errors.New("redirect URI must use https or loopback http")
	}
	host := parsed.Hostname()
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	return errors.New("http redirect URI must be loopback")
}

func decodeRedirectURIs(raw string) ([]string, error) {
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil, err
	}
	return values, nil
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func constantTimeHashMatch(raw, expectedHash string) bool {
	actualHash := secretHash(raw)
	return subtle.ConstantTimeCompare([]byte(actualHash), []byte(expectedHash)) == 1
}

func (s *Server) validateClient(clientID, clientSecret string) (*OAuthClient, error) {
	client, err := s.repo.GetClient(clientID)
	if err != nil {
		return nil, err
	}
	if client.ClientSecretHash == "" {
		if clientSecret != "" {
			return nil, errors.New("invalid client credentials")
		}
		return client, nil
	}
	if clientSecret == "" || !constantTimeHashMatch(clientSecret, client.ClientSecretHash) {
		return nil, errors.New("invalid client credentials")
	}
	return client, nil
}

func (s *Server) beginClientRegistration(now time.Time) (bool, time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := now.Add(-s.cfg.ClientUnusedTTL)
	kept := s.registrationTimes[:0]
	for _, registeredAt := range s.registrationTimes {
		if registeredAt.After(cutoff) {
			kept = append(kept, registeredAt)
		}
	}
	s.registrationTimes = kept

	limit := s.dcrBurstLimit()
	if len(s.registrationTimes) >= limit {
		earliest := s.registrationTimes[0]
		for _, registeredAt := range s.registrationTimes[1:] {
			if registeredAt.Before(earliest) {
				earliest = registeredAt
			}
		}
		return false, positiveRetryAfter(earliest.Add(s.cfg.ClientUnusedTTL).Sub(now))
	}
	s.registrationTimes = append(s.registrationTimes, now)
	return true, 0
}

func (s *Server) dcrBurstLimit() int {
	limit := s.cfg.MaxClients / dcrBurstLimitDivisor
	if limit < minimumDCRBurstLimit {
		limit = minimumDCRBurstLimit
	}
	if limit >= s.cfg.MaxClients && s.cfg.MaxClients > 1 {
		limit = s.cfg.MaxClients - 1
	}
	if limit > s.cfg.MaxClients {
		limit = s.cfg.MaxClients
	}
	return max(1, limit)
}

func (s *Server) createPending(p pendingAuthorization, now time.Time) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupInMemoryLocked(now)
	if len(s.pending) >= s.cfg.MaxPending {
		return "", ErrStateLimitReached
	}
	clientPending := 0
	peerPending := 0
	for _, existing := range s.pending {
		if existing.ClientID == p.ClientID {
			clientPending++
		}
		if p.PeerKey != "" && existing.PeerKey == p.PeerKey {
			peerPending++
		}
	}
	if clientPending >= maxPendingPerClient || (p.PeerKey != "" && peerPending >= maxPendingPerPeer) {
		return "", ErrStateLimitReached
	}
	id := randomToken(24)
	p.ExpiresAt = now.Add(s.cfg.PendingRequestTTL)
	s.pending[id] = p
	return id, nil
}

func (s *Server) getPending(id string, now time.Time) (pendingAuthorization, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupInMemoryLocked(now)
	p, ok := s.pending[id]
	return p, ok
}

func (s *Server) consumePending(id string, now time.Time) (pendingAuthorization, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupInMemoryLocked(now)
	pending, ok := s.pending[id]
	if !ok {
		return pendingAuthorization{}, false
	}
	delete(s.pending, id)
	return pending, true
}

func (s *Server) cleanupInMemoryLocked(now time.Time) {
	for id, p := range s.pending {
		if !p.ExpiresAt.After(now) {
			delete(s.pending, id)
		}
	}
	for key, failure := range s.authFailures {
		if failure.InFlight == 0 && !failure.ResetAt.After(now) {
			delete(s.authFailures, key)
		}
	}
}

func requestPeerHost(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return strings.TrimSpace(host)
}

func requestPeerKey(r *http.Request) string {
	// Deliberately use the transport peer, not untrusted forwarded headers.
	return secretHash(requestPeerHost(r))
}

func authPeerKey(r *http.Request, email string) string {
	return secretHash(strings.ToLower(strings.TrimSpace(email)) + "|" + requestPeerHost(r))
}

func (s *Server) beginAuthAttempt(key string, now time.Time) (bool, time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupInMemoryLocked(now)

	failure, ok := s.authFailures[key]
	if ok && failure.Count+failure.InFlight >= authFailureLimit {
		return false, positiveRetryAfter(failure.ResetAt.Sub(now))
	}
	if !ok {
		if len(s.authFailures) >= maxAuthFailureKeys && !s.evictLowRiskAuthFailureLocked() {
			return false, s.nextAuthFailureRetryLocked(now)
		}
		failure = authFailureState{ResetAt: now.Add(authFailureWindow)}
	}
	failure.InFlight++
	failure.LastSeen = now
	s.authFailures[key] = failure
	return true, 0
}

func (s *Server) finishAuthAttempt(key string, now time.Time, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupInMemoryLocked(now)

	failure, ok := s.authFailures[key]
	if !ok {
		return
	}
	if failure.InFlight > 0 {
		failure.InFlight--
	}
	failure.LastSeen = now
	if success {
		failure.Count = 0
		if failure.InFlight == 0 {
			delete(s.authFailures, key)
			return
		}
		failure.ResetAt = now.Add(authFailureWindow)
		s.authFailures[key] = failure
		return
	}
	if !failure.ResetAt.After(now) {
		failure.Count = 0
		failure.ResetAt = now.Add(authFailureWindow)
	}
	failure.Count++
	s.authFailures[key] = failure
}

func (s *Server) evictLowRiskAuthFailureLocked() bool {
	var candidateKey string
	var candidate authFailureState
	found := false
	for key, failure := range s.authFailures {
		if failure.InFlight != 0 || failure.Count >= authFailureLimit {
			continue
		}
		if !found || failure.Count < candidate.Count ||
			(failure.Count == candidate.Count && failure.LastSeen.Before(candidate.LastSeen)) ||
			(failure.Count == candidate.Count && failure.LastSeen.Equal(candidate.LastSeen) && key < candidateKey) {
			candidateKey = key
			candidate = failure
			found = true
		}
	}
	if !found {
		return false
	}
	delete(s.authFailures, candidateKey)
	return true
}

func (s *Server) nextAuthFailureRetryLocked(now time.Time) time.Duration {
	var retryAfter time.Duration
	for _, failure := range s.authFailures {
		remaining := failure.ResetAt.Sub(now)
		if remaining <= 0 {
			continue
		}
		if retryAfter == 0 || remaining < retryAfter {
			retryAfter = remaining
		}
	}
	if retryAfter == 0 {
		retryAfter = authFailureWindow
	}
	return positiveRetryAfter(retryAfter)
}

func positiveRetryAfter(retryAfter time.Duration) time.Duration {
	if retryAfter <= 0 {
		return time.Second
	}
	return retryAfter
}
