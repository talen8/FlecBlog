package oauthserver

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"
)

type registrationRequest struct {
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types,omitempty"`
	ResponseTypes           []string `json:"response_types,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
	ClientName              string   `json:"client_name,omitempty"`
}

type registrationResponse struct {
	ClientID                string   `json:"client_id"`
	ClientSecret            string   `json:"client_secret,omitempty"`
	ClientSecretExpiresAt   *int64   `json:"client_secret_expires_at,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
}

func (s *Server) RegistrationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeOAuthError(w, http.StatusMethodNotAllowed, "invalid_request", "method not allowed")
			return
		}
		var req registrationRequest
		if err := readLimitedJSON(r, &req); err != nil {
			writeOAuthError(w, http.StatusBadRequest, "invalid_client_metadata", "invalid registration request")
			return
		}
		redirectURIs, authMethod, err := validateRegistrationRequest(req)
		if err != nil {
			writeOAuthError(w, http.StatusBadRequest, "invalid_client_metadata", err.Error())
			return
		}
		now := time.Now()
		if allowed, retryAfter := s.beginClientRegistration(now); !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(max(1, int(math.Ceil(retryAfter.Seconds())))))
			writeOAuthError(w, http.StatusTooManyRequests, "temporarily_unavailable", "OAuth client registration rate limit reached")
			return
		}
		encodedRedirects, err := json.Marshal(redirectURIs)
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to encode client metadata")
			return
		}
		clientID := "mcp-" + randomToken(18)
		clientSecret := ""
		clientSecretHash := ""
		if authMethod == "client_secret_post" {
			clientSecret = randomToken(32)
			clientSecretHash = secretHash(clientSecret)
		}
		err = s.repo.CreateClientBounded(&OAuthClient{
			ClientID:         clientID,
			ClientSecretHash: clientSecretHash,
			RedirectURIs:     string(encodedRedirects),
			CreatedAt:        now,
		}, s.cfg.MaxClients, now.Add(-s.cfg.ClientInactiveTTL), now.Add(-s.cfg.ClientUnusedTTL))
		if errors.Is(err, ErrStateLimitReached) {
			writeOAuthError(w, http.StatusServiceUnavailable, "temporarily_unavailable", "OAuth client limit reached")
			return
		}
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to register client")
			return
		}
		response := registrationResponse{
			ClientID:                clientID,
			ClientSecret:            clientSecret,
			RedirectURIs:            redirectURIs,
			GrantTypes:              []string{"authorization_code", "refresh_token"},
			ResponseTypes:           []string{"code"},
			TokenEndpointAuthMethod: authMethod,
		}
		if authMethod == "client_secret_post" {
			neverExpires := int64(0)
			response.ClientSecretExpiresAt = &neverExpires
		}
		writeJSONStatus(w, http.StatusCreated, response)
	})
}

func validateRegistrationRequest(req registrationRequest) ([]string, string, error) {
	if len(req.RedirectURIs) == 0 || len(req.RedirectURIs) > maxRedirectURIs {
		return nil, "", errors.New("redirect_uris must contain 1 to 10 entries")
	}
	if len(req.ClientName) > maxClientNameBytes {
		return nil, "", errors.New("client_name is too long")
	}
	if len(req.GrantTypes) > 2 {
		return nil, "", errors.New("too many grant_types")
	}
	if len(req.ResponseTypes) > 1 {
		return nil, "", errors.New("too many response_types")
	}
	if len(req.GrantTypes) > 0 && !stringSubset(req.GrantTypes, []string{"authorization_code", "refresh_token"}) {
		return nil, "", errors.New("unsupported grant_types")
	}
	if len(req.ResponseTypes) > 0 && !stringSubset(req.ResponseTypes, []string{"code"}) {
		return nil, "", errors.New("unsupported response_types")
	}
	authMethod := req.TokenEndpointAuthMethod
	if authMethod == "" {
		authMethod = "client_secret_post"
	}
	if authMethod != "client_secret_post" && authMethod != "none" {
		return nil, "", errors.New("unsupported token_endpoint_auth_method")
	}
	seen := make(map[string]struct{}, len(req.RedirectURIs))
	redirectURIs := make([]string, 0, len(req.RedirectURIs))
	for _, redirectURI := range req.RedirectURIs {
		if len(redirectURI) > maxRedirectURIBytes {
			return nil, "", errors.New("redirect_uri is too long")
		}
		if err := validateRedirectURI(redirectURI); err != nil {
			return nil, "", errors.New("invalid redirect_uri")
		}
		if _, ok := seen[redirectURI]; ok {
			continue
		}
		seen[redirectURI] = struct{}{}
		redirectURIs = append(redirectURIs, redirectURI)
	}
	return redirectURIs, authMethod, nil
}

func stringSubset(actual, supported []string) bool {
	allowed := make(map[string]struct{}, len(supported))
	for _, value := range supported {
		allowed[value] = struct{}{}
	}
	for _, value := range actual {
		if _, ok := allowed[value]; !ok {
			return false
		}
	}
	return true
}
