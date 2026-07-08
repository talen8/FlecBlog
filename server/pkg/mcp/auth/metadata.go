package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ProtectedResourceMetadata is the RFC 9728 metadata advertised by the MCP resource server.
type ProtectedResourceMetadata struct {
	Resource             string   `json:"resource"`
	AuthorizationServers []string `json:"authorization_servers"`
	ScopesSupported      []string `json:"scopes_supported,omitempty"`
	BearerMethods        []string `json:"bearer_methods_supported,omitempty"`
}

// MetadataHandler serves OAuth protected-resource metadata.
func (a *Authenticator) MetadataHandler() http.Handler {
	metadata := ProtectedResourceMetadata{
		Resource:             a.config.ResourceURI,
		AuthorizationServers: []string{a.config.Issuer},
		ScopesSupported:      append([]string(nil), a.config.SupportedScopes...),
		BearerMethods:        []string{"header"},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		if err := json.NewEncoder(w).Encode(metadata); err != nil {
			return
		}
	})
}

// UnauthorizedChallenge returns the OAuth Bearer challenge for this MCP resource.
func (a *Authenticator) UnauthorizedChallenge() string {
	if a.config.Mode == ModeStatic || a.config.MetadataURL == "" {
		return "Bearer"
	}
	parts := []string{fmt.Sprintf(`resource_metadata="%s"`, a.config.MetadataURL)}
	challengeScopes := a.currentChallengeScopes()
	if len(challengeScopes) > 0 {
		parts = append(parts, fmt.Sprintf(`scope="%s"`, strings.Join(challengeScopes, " ")))
	}
	return "Bearer " + strings.Join(parts, ", ")
}

// ExtractBearerToken parses an Authorization header using the Bearer scheme.
func ExtractBearerToken(header string) (string, bool) {
	space := strings.IndexByte(header, ' ')
	if space <= 0 || !strings.EqualFold(header[:space], "Bearer") {
		return "", false
	}
	token := header[space+1:]
	if token == "" || strings.TrimSpace(token) != token || strings.ContainsAny(token, " \t\r\n") {
		return "", false
	}
	return token, true
}

func (a *Authenticator) currentChallengeScopes() []string {
	scopes := append([]string(nil), a.config.ChallengeScopes...)
	if a.adminToolsEnabledProvider == nil || !a.adminToolsEnabledProvider() {
		return scopes
	}
	for _, scope := range scopes {
		if scope == ScopeAdmin {
			return scopes
		}
	}
	return append(scopes, ScopeAdmin)
}
