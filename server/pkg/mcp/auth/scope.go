package auth

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// InsufficientScopeError reports a missing OAuth scope without exposing token details.
type InsufficientScopeError struct {
	Required string
}

func (e *InsufficientScopeError) Error() string {
	return fmt.Sprintf("insufficient MCP scope: requires %s", e.Required)
}

// RequireScope wraps a typed MCP tool handler with fail-closed scope authorization.
// Authorization is derived from the current HTTP request's go-sdk TokenInfo, not
// from the long-lived Streamable session connection context.
func RequireScope[In, Out any](required string, next sdkmcp.ToolHandlerFor[In, Out]) sdkmcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, request *sdkmcp.CallToolRequest, input In) (*sdkmcp.CallToolResult, Out, error) {
		var zero Out
		if request == nil || request.Extra == nil || request.Extra.TokenInfo == nil {
			return nil, zero, fmt.Errorf("current MCP request authorization metadata is missing")
		}
		principal, ok := PrincipalFromTokenInfo(request.Extra.TokenInfo)
		if !ok {
			return nil, zero, fmt.Errorf("current MCP request authorization principal is missing")
		}
		if !principal.HasScope(required) {
			return nil, zero, &InsufficientScopeError{Required: required}
		}

		// Rebind the current-request principal for downstream handlers such as the
		// local operator resolver. This deliberately overwrites any principal value
		// inherited from the session's initialize connection context.
		ctx = ContextWithPrincipal(ctx, principal)
		return next(ctx, request, input)
	}
}
