package main

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	goblmcp "github.com/invopop/gobl/internal/mcp"
)

type mcpOpts struct{}

func mcpServe() *mcpOpts {
	return &mcpOpts{}
}

func (o *mcpOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Launch an MCP server over stdio",
		RunE:  o.runE,
	}
	return cmd
}

func (o *mcpOpts) runE(_ *cobra.Command, _ []string) error {
	srv := goblmcp.NewServer()
	return server.ServeStdio(srv)
}
