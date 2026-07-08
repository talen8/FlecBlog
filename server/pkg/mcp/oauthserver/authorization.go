package oauthserver

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var pkceChallengePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{43}$`)

var authorizationPage = template.Must(template.New("authorize").Parse(`<!doctype html>
<html lang="zh-CN"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>FlecBlog MCP 授权</title><style>
body{font-family:system-ui,sans-serif;max-width:520px;margin:48px auto;padding:0 20px;line-height:1.5;color:#1f2937}
.card{border:1px solid #d1d5db;border-radius:14px;padding:24px}label{display:block;margin-top:14px;font-weight:600}
input{box-sizing:border-box;width:100%;padding:10px;margin-top:6px;border:1px solid #9ca3af;border-radius:8px}
button{margin-top:20px;width:100%;padding:11px;border:0;border-radius:8px;background:#111827;color:white;font-weight:700}
.scope{font-family:ui-monospace,monospace;font-size:13px;background:#f3f4f6;padding:8px;border-radius:8px;word-break:break-word}
.error{color:#b91c1c}.muted{color:#6b7280;font-size:14px}</style></head><body>
<div class="card"><h1>授权 FlecBlog MCP</h1>
<p class="muted">客户端：{{.ClientID}}</p><p class="muted">回调地址：</p><div class="scope">{{.RedirectURI}}</div><p>请求权限：</p><div class="scope">{{.Scopes}}</div>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/authorize">
<input type="hidden" name="request_id" value="{{.RequestID}}">
<label>管理员邮箱<input type="email" name="email" autocomplete="username" required autofocus></label>
<label>密码<input type="password" name="password" autocomplete="current-password" required></label>
<button type="submit">确认授权</button></form></div></body></html>`))

type authorizationPageData struct {
	RequestID   string
	ClientID    string
	RedirectURI string
	Scopes      string
	Error       string
}

func (s *Server) AuthorizationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleAuthorizationGet(w, r)
		case http.MethodPost:
			s.handleAuthorizationPost(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			writeOAuthError(w, http.StatusMethodNotAllowed, "invalid_request", "method not allowed")
		}
	})
}

func (s *Server) handleAuthorizationGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	clientID := query.Get("client_id")
	redirectURI := query.Get("redirect_uri")
	client, err := s.repo.GetClient(clientID)
	if err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_client", "unknown client")
		return
	}
	redirectURIs, err := decodeRedirectURIs(client.RedirectURIs)
	if err != nil || !containsString(redirectURIs, redirectURI) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "invalid redirect_uri")
		return
	}
	if query.Get("response_type") != "code" {
		writeOAuthError(w, http.StatusBadRequest, "unsupported_response_type", "response_type must be code")
		return
	}
	if query.Get("resource") != s.cfg.ResourceURI {
		writeOAuthError(w, http.StatusBadRequest, "invalid_target", "resource mismatch")
		return
	}
	challenge := query.Get("code_challenge")
	if query.Get("code_challenge_method") != "S256" || !pkceChallengePattern.MatchString(challenge) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "PKCE S256 is required")
		return
	}
	scopes, err := s.normalizeScopes(query.Get("scope"))
	if err != nil || !hasMCPAccessScope(scopes) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_scope", "invalid OAuth scope")
		return
	}
	now := time.Now()
	requestID, err := s.createPending(pendingAuthorization{
		ClientID:      clientID,
		PeerKey:       requestPeerKey(r),
		RedirectURI:   redirectURI,
		State:         query.Get("state"),
		Resource:      s.cfg.ResourceURI,
		Scopes:        scopes,
		CodeChallenge: challenge,
	}, now)
	if errors.Is(err, ErrStateLimitReached) {
		writeOAuthError(w, http.StatusServiceUnavailable, "temporarily_unavailable", "authorization request limit reached")
		return
	}
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to create authorization request")
		return
	}
	s.renderAuthorizationPage(w, http.StatusOK, authorizationPageData{
		RequestID:   requestID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scopes:      scopeString(scopes),
	})
}

func (s *Server) handleAuthorizationPost(w http.ResponseWriter, r *http.Request) {
	form, err := readLimitedForm(r)
	if err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "invalid authorization form")
		return
	}
	requestID := form.Get("request_id")
	now := time.Now()
	pending, ok := s.getPending(requestID, now)
	if !ok {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "authorization request expired")
		return
	}
	email := strings.TrimSpace(form.Get("email"))
	failureKey := authPeerKey(r, email)
	if allowed, retryAfter := s.beginAuthAttempt(failureKey, now); !allowed {
		w.Header().Set("Retry-After", strconv.Itoa(max(1, int(math.Ceil(retryAfter.Seconds())))))
		writeOAuthError(w, http.StatusTooManyRequests, "temporarily_unavailable", "too many authorization failures")
		return
	}
	identity, err := s.userService.AuthenticatePassword(email, form.Get("password"))
	success := err == nil && identity != nil && s.isAllowedRole(identity.Role)
	s.finishAuthAttempt(failureKey, time.Now(), success)
	if !success {
		s.renderAuthorizationPage(w, http.StatusUnauthorized, authorizationPageData{
			RequestID:   requestID,
			ClientID:    pending.ClientID,
			RedirectURI: pending.RedirectURI,
			Scopes:      scopeString(pending.Scopes),
			Error:       "邮箱、密码或管理员权限无效",
		})
		return
	}

	pending, ok = s.consumePending(requestID, now)
	if !ok {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "authorization request expired")
		return
	}

	code := randomToken(32)
	err = s.repo.StoreAuthorizationCodeBounded(&AuthorizationCode{
		CodeHash:      secretHash(code),
		ClientID:      pending.ClientID,
		UserID:        identity.ID,
		RedirectURI:   pending.RedirectURI,
		Resource:      pending.Resource,
		Scopes:        scopeString(pending.Scopes),
		CodeChallenge: pending.CodeChallenge,
		ExpiresAt:     now.Add(s.cfg.AuthCodeTTL),
		CreatedAt:     now,
	}, s.cfg.MaxAuthCodes, now)
	if errors.Is(err, ErrStateLimitReached) {
		writeOAuthError(w, http.StatusServiceUnavailable, "temporarily_unavailable", "authorization code limit reached")
		return
	}
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to issue authorization code")
		return
	}
	_ = s.repo.TouchClient(pending.ClientID, now)

	redirect, err := url.Parse(pending.RedirectURI)
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "invalid stored redirect URI")
		return
	}
	values := redirect.Query()
	values.Set("code", code)
	if pending.State != "" {
		values.Set("state", pending.State)
	}
	redirect.RawQuery = values.Encode()
	http.Redirect(w, r, redirect.String(), http.StatusFound)
}

func (s *Server) renderAuthorizationPage(w http.ResponseWriter, status int, data authorizationPageData) {
	secureAuthorizationHTMLHeaders(w, data.RedirectURI)
	w.WriteHeader(status)
	_ = authorizationPage.Execute(w, data)
}

func hasMCPAccessScope(scopes []string) bool {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, "mcp:") {
			return true
		}
	}
	return false
}

func pkceChallenge(verifier string) string {
	digest := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}
