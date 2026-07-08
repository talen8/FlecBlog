package oauthserver

import (
	"context"
	"fmt"
	"math"

	mcpauth "flec_blog/pkg/mcp/auth"
)

func (s *Server) OAuthPrincipalValidator() mcpauth.OAuthPrincipalValidator {
	return func(_ context.Context, principal *mcpauth.Principal) error {
		if principal == nil || principal.Method != "oauth" {
			return fmt.Errorf("embedded OAuth principal is invalid")
		}
		userID, err := embeddedUintClaim(principal.Claims["mcp_user_id"], false)
		if err != nil {
			return fmt.Errorf("embedded OAuth mcp_user_id claim is invalid: %w", err)
		}
		tokenVersion, err := embeddedUintClaim(principal.Claims["token_version"], true)
		if err != nil {
			return fmt.Errorf("embedded OAuth token_version claim is invalid: %w", err)
		}
		if principal.Subject != fmt.Sprintf("user:%d", userID) {
			return fmt.Errorf("embedded OAuth subject binding mismatch")
		}

		identity, err := s.userService.GetEnabledIdentity(userID)
		if err != nil {
			return fmt.Errorf("embedded OAuth local user is unavailable: %w", err)
		}
		if !s.isAllowedRole(identity.Role) {
			return fmt.Errorf("embedded OAuth local user role is not authorized")
		}
		if identity.TokenVersion != tokenVersion {
			return fmt.Errorf("embedded OAuth token security generation mismatch")
		}
		if err := principal.BindVerifiedLocalUser(userID); err != nil {
			return fmt.Errorf("embedded OAuth local user binding failed: %w", err)
		}
		return nil
	}
}

func embeddedUintClaim(value any, allowZero bool) (uint, error) {
	number, ok := value.(float64)
	if !ok || math.IsNaN(number) || math.IsInf(number, 0) || math.Trunc(number) != number {
		return 0, fmt.Errorf("must be an integer")
	}
	if number < 0 || (!allowZero && number == 0) || number > float64(^uint(0)) {
		return 0, fmt.Errorf("out of range")
	}
	return uint(number), nil
}
