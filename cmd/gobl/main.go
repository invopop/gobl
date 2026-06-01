// Package main provides a command-line interface to the GOBL library.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
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

// newLogger builds the slog.Logger used for all operator-facing log
// output. The result writes one entry per line to stderr; result
// output (signed envelopes, /who party JSON, version) lives on stdout
// and is not affected by this flag.
func newLogger(jsonMode bool) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var h slog.Handler
	if jsonMode {
		h = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		h = slog.NewTextHandler(os.Stderr, opts)
	}
	return slog.New(h)
}

func printError(err error) {
	// Normalise to a *gobl.Error so every report carries a "key" and
	// (when present) a "message" + structured faults.
	var ge *gobl.Error
	if !errors.As(err, &ge) {
		ge = gobl.ErrInternal.WithCause(err)
	}
	attrs := []any{"key", ge.Key().String()}
	if msg := ge.Message(); msg != "" {
		attrs = append(attrs, "message", msg)
	}
	if faults := ge.Faults(); faults != nil {
		attrs = append(attrs, "faults", faults)
	}
	slog.Error("command failed", attrs...)
}
