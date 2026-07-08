package oauthserver

import (
	"context"
	"fmt"

	mcpauth "flec_blog/pkg/mcp/auth"

	"github.com/golang-jwt/jwt/v4"
)

// OAuthKeyResolver returns an in-process resolver for access tokens signed by this embedded server.
func (s *Server) OAuthKeyResolver() mcpauth.OAuthKeyResolver {
	return func(_ context.Context, token *jwt.Token) (any, error) {
		if token == nil || token.Method == nil || token.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
			return nil, fmt.Errorf("embedded OAuth token must use EdDSA")
		}
		kid, _ := token.Header["kid"].(string)
		if kid == "" || kid != s.signer.kid {
			return nil, fmt.Errorf("embedded OAuth token kid is unknown")
		}
		return s.signer.publicKey, nil
	}
}
