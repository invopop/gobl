// Package main provides a command-line interface to the GOBL library.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl"
)

// build data provided by goreleaser and mage setup
var (
	version = "dev"
	date    = ""
)

var versionOutput = struct {
	Version string `json:"version"`
	Date    string `json:"date,omitempty"`
}{
	Version: version,
	Date:    date,
}

func main() {
	if err := run(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return root().cmd().ExecuteContext(ctx)
}

func inputFilename(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "\t") // always indent version
			return enc.Encode(versionOutput)
		},
	}
}

func encode(in any, out io.WriteCloser, indent bool) error {
	enc := json.NewEncoder(out)
	if indent {
		enc.SetIndent("", "\t")
	}
	return enc.Encode(in)
}

func printError(err error) {
	// Normalise to a *gobl.Error so the JSON body always carries a
	// "key" and "message"; a plain error would otherwise marshal to "{}".
	var ge *gobl.Error
	if errors.As(err, &ge) {
		err = ge
	} else {
		err = gobl.ErrInternal.WithCause(err)
	}
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "\t") // always indent errors
	if encErr := enc.Encode(err); encErr != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
