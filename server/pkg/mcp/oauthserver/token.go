package oauthserver

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"flec_blog/internal/service"
)

var pkceVerifierPattern = regexp.MustCompile(`^[A-Za-z0-9._~-]{43,128}$`)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

func (s *Server) TokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeOAuthError(w, http.StatusMethodNotAllowed, "invalid_request", "method not allowed")
			return
		}
		form, err := readLimitedForm(r)
		if err != nil {
			writeOAuthError(w, http.StatusBadRequest, "invalid_request", "invalid token request")
			return
		}
		client, err := s.validateClient(form.Get("client_id"), form.Get("client_secret"))
		if err != nil {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "invalid client credentials")
			return
		}
		switch form.Get("grant_type") {
		case "authorization_code":
			s.handleAuthorizationCodeGrant(w, form, client)
		case "refresh_token":
			s.handleRefreshTokenGrant(w, form, client)
		default:
			writeOAuthError(w, http.StatusBadRequest, "unsupported_grant_type", "unsupported grant_type")
		}
	})
}

func (s *Server) handleAuthorizationCodeGrant(w http.ResponseWriter, form url.Values, client *OAuthClient) {
	now := time.Now()
	verifier := form.Get("code_verifier")
	if !pkceVerifierPattern.MatchString(verifier) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "PKCE verification failed")
		return
	}
	resource := form.Get("resource")
	if resource != s.cfg.ResourceURI {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "authorization code binding mismatch")
		return
	}
	code, err := s.repo.ConsumeAuthorizationCodeBound(
		secretHash(form.Get("code")),
		client.ClientID,
		form.Get("redirect_uri"),
		resource,
		pkceChallenge(verifier),
		now,
	)
	if err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "invalid authorization code")
		return
	}
	identity, err := s.userService.GetEnabledIdentity(code.UserID)
	if err != nil || !s.isAllowedRole(identity.Role) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "local user is not authorized")
		return
	}
	scopes := parseStoredScopes(code.Scopes)
	response, err := s.issueTokens(identity, client.ClientID, scopes, now)
	if errors.Is(err, ErrStateLimitReached) {
		writeOAuthError(w, http.StatusServiceUnavailable, "temporarily_unavailable", "refresh token limit reached")
		return
	}
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to issue token")
		return
	}
	_ = s.repo.TouchClient(client.ClientID, now)
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleRefreshTokenGrant(w http.ResponseWriter, form url.Values, client *OAuthClient) {
	if resource := form.Get("resource"); resource != "" && resource != s.cfg.ResourceURI {
		writeOAuthError(w, http.StatusBadRequest, "invalid_target", "resource mismatch")
		return
	}
	now := time.Now()
	oldHash := secretHash(form.Get("refresh_token"))
	refresh, err := s.repo.GetRefreshToken(oldHash)
	if err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "invalid refresh token")
		return
	}
	if !refresh.ExpiresAt.After(now) || refresh.ClientID != client.ClientID {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "refresh token binding mismatch")
		return
	}
	identity, err := s.userService.GetEnabledIdentity(refresh.UserID)
	if err != nil || !s.isAllowedRole(identity.Role) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "local user is not authorized")
		return
	}
	if !refreshTokenGenerationMatches(refresh, identity.TokenVersion) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "refresh token security generation mismatch")
		return
	}
	scopes := parseStoredScopes(refresh.Scopes)
	response, err := s.rotateTokens(identity, client.ClientID, scopes, now, oldHash)
	if err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "failed to rotate refresh token")
		return
	}
	_ = s.repo.TouchClient(client.ClientID, now)
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) issueTokens(identity *service.PasswordIdentity, clientID string, scopes []string, now time.Time) (tokenResponse, error) {
	scope := scopeString(scopes)
	accessToken, err := s.signer.signAccessToken(s.cfg, identity.ID, identity.TokenVersion, scope, now)
	if err != nil {
		return tokenResponse{}, err
	}
	response := tokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.cfg.AccessTokenTTL.Seconds()),
		Scope:       scope,
	}
	refreshToken := randomToken(32)
	err = s.repo.StoreRefreshTokenBounded(&RefreshToken{
		TokenHash:    secretHash(refreshToken),
		ClientID:     clientID,
		UserID:       identity.ID,
		TokenVersion: tokenVersionPtr(identity.TokenVersion),
		Scopes:       scope,
		ExpiresAt:    now.Add(s.cfg.RefreshTokenTTL),
		CreatedAt:    now,
	}, s.cfg.MaxRefreshTokens, now)
	if err != nil {
		return tokenResponse{}, err
	}
	response.RefreshToken = refreshToken
	return response, nil
}

func (s *Server) rotateTokens(identity *service.PasswordIdentity, clientID string, scopes []string, now time.Time, oldHash string) (tokenResponse, error) {
	scope := scopeString(scopes)
	accessToken, err := s.signer.signAccessToken(s.cfg, identity.ID, identity.TokenVersion, scope, now)
	if err != nil {
		return tokenResponse{}, err
	}
	refreshToken := randomToken(32)
	replacement := &RefreshToken{
		TokenHash:    secretHash(refreshToken),
		ClientID:     clientID,
		UserID:       identity.ID,
		TokenVersion: tokenVersionPtr(identity.TokenVersion),
		Scopes:       scope,
		ExpiresAt:    now.Add(s.cfg.RefreshTokenTTL),
		CreatedAt:    now,
	}
	if err := s.repo.RotateRefreshToken(oldHash, clientID, identity.ID, replacement); err != nil {
		return tokenResponse{}, err
	}
	return tokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.AccessTokenTTL.Seconds()),
		RefreshToken: refreshToken,
		Scope:        scope,
	}, nil
}

func refreshTokenGenerationMatches(refresh *RefreshToken, current uint) bool {
	return refresh != nil && refresh.TokenVersion != nil && *refresh.TokenVersion == current
}

func tokenVersionPtr(version uint) *uint { return &version }
