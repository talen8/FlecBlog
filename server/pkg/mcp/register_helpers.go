package mcp

import (
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
