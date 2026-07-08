package oauthserver

import (
	"encoding/json"
	"net/http"
)

type AuthorizationServerMetadata struct {
	Issuer                        string   `json:"issuer"`
	AuthorizationEndpoint         string   `json:"authorization_endpoint"`
	TokenEndpoint                 string   `json:"token_endpoint"`
	RegistrationEndpoint          string   `json:"registration_endpoint"`
	JWKSURI                       string   `json:"jwks_uri"`
	ResponseTypesSupported        []string `json:"response_types_supported"`
	GrantTypesSupported           []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
	ScopesSupported               []string `json:"scopes_supported"`
	TokenEndpointAuthMethods      []string `json:"token_endpoint_auth_methods_supported"`
	ResourceParameterSupported    bool     `json:"resource_parameter_supported"`
}

func (s *Server) AuthorizationServerMetadataHandler() http.Handler {
	metadata := AuthorizationServerMetadata{
		Issuer:                        s.cfg.Issuer,
		AuthorizationEndpoint:         s.cfg.AuthorizationEndpoint,
		TokenEndpoint:                 s.cfg.TokenEndpoint,
		RegistrationEndpoint:          s.cfg.RegistrationEndpoint,
		JWKSURI:                       s.cfg.JWKSURI,
		ResponseTypesSupported:        []string{"code"},
		GrantTypesSupported:           []string{"authorization_code", "refresh_token"},
		CodeChallengeMethodsSupported: []string{"S256"},
		ScopesSupported:               append([]string(nil), s.cfg.SupportedScopes...),
		TokenEndpointAuthMethods:      []string{"client_secret_post", "none"},
		ResourceParameterSupported:    true,
	}
	return jsonHandler(metadata, "public, max-age=300")
}

func (s *Server) JWKSHandler() http.Handler {
	return jsonHandler(s.signer.jwks(), "public, max-age=300")
}

func jsonHandler(payload any, cacheControl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", cacheControl)
		_ = json.NewEncoder(w).Encode(payload)
	})
}
