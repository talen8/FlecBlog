package mcp

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

// registerTools composes domain-specific MCP tool registrations.
func (s *publicServer) registerTools(server *sdkmcp.Server) {
	s.registerArticleTools(server)
	s.registerContentTools(server)
	s.registerAdminTools(server)
}
