// The gobl command provides a command-line interface to the GOBL library.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	root := &cobra.Command{
		Use:           "gobl",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(&cobra.Command{
		Use:  "build",
		Args: cobra.ExactArgs(1),
		RunE: build,
	})
	root.AddCommand(&cobra.Command{
		Use:  "verify",
		Args: cobra.ExactArgs(1),
		RunE: verify,
	})
	return root.ExecuteContext(ctx)
}

func build(*cobra.Command, []string) error {
	return nil
}

func verify(*cobra.Command, []string) error {
	return nil
}
