package mcp

import (
	"context"
	"fmt"

	mcpauth "flec_blog/pkg/mcp/auth"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func readOnlyMCPToolAnnotations(title string) *sdkmcp.ToolAnnotations {
	closedWorld := false
	return &sdkmcp.ToolAnnotations{
		Title:         title,
		ReadOnlyHint:  true,
		OpenWorldHint: &closedWorld,
	}
}

func mutatingMCPToolAnnotations(title string, destructive, idempotent bool) *sdkmcp.ToolAnnotations {
	closedWorld := false
	return &sdkmcp.ToolAnnotations{
		Title:           title,
		DestructiveHint: &destructive,
		IdempotentHint:  idempotent,
		OpenWorldHint:   &closedWorld,
	}
}

func addScopedMCPTool[In, Out any](
	server *sdkmcp.Server,
	tool *sdkmcp.Tool,
	scope string,
	handler sdkmcp.ToolHandlerFor[In, Out],
) {
	sdkmcp.AddTool(server, tool, mcpauth.RequireScope(scope, handler))
}

func requireMCPAdminToolsEnabled[In, Out any](
	provider func() bool,
	next sdkmcp.ToolHandlerFor[In, Out],
) sdkmcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, request *sdkmcp.CallToolRequest, input In) (*sdkmcp.CallToolResult, Out, error) {
		var zero Out
		if provider == nil || !provider() {
			return nil, zero, fmt.Errorf("MCP 管理员工具未启用")
		}
		return next(ctx, request, input)
	}
}
