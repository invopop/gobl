package main

import (
	"context"

	"github.com/spf13/cobra"
)

func commandContext(cmd *cobra.Command) context.Context {
	ctx := cmd.Context()
	// When tests manually invoke our `runE` methods, the command
	// context is nil, even though the doc block of `cmd.Context` hints at
	// different behavior.
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
