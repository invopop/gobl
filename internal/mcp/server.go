// Package mcp provides a Model Context Protocol server for GOBL.
package mcp

import (
	"github.com/invopop/gobl"
	"github.com/mark3labs/mcp-go/server"
)

// NewServer creates and configures a new MCP server with all GOBL tools
// and resources registered.
func NewServer() *server.MCPServer {
	srv := server.NewMCPServer(
		"gobl",
		string(gobl.VERSION),
	)

	registerTools(srv)
	registerResources(srv)

	return srv
}
