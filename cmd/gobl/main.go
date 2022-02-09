// The gobl command provides a command-line interface to the GOBL library.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/iotools"
)

func main() {
	if err := run(); err != nil {
		echoErr := new(echo.HTTPError)
		if errors.As(err, &echoErr) {
			msg := echoErr.Message
			int := echoErr.Internal
			switch {
			case msg != "" && int != nil:
				err = fmt.Errorf("%v: %w", msg, int)
			case int != nil:
				err = int
			default:
				err = fmt.Errorf("%v", msg)
			}
		}
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return root().ExecuteContext(ctx)
}

func root() *cobra.Command {
	root := &cobra.Command{
		Use:           "gobl",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(&cobra.Command{
		Use:  "verify [infile]",
		Args: cobra.MaximumNArgs(1),
		RunE: verify,
	})
	root.AddCommand(build().cmd())
	root.AddCommand(version())
	root.AddCommand(serve().cmd())
	return root
}

func inputFilename(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func openInput(cmd *cobra.Command, args []string) (io.ReadCloser, error) {
	if inFile := inputFilename(args); inFile != "" {
		return os.Open(inFile)
	}
	return ioutil.NopCloser(cmd.InOrStdin()), nil
}

func readEnv(cmd *cobra.Command, args []string) (*gobl.Envelope, error) {
	input, err := openInput(cmd, args)
	if err != nil {
		return nil, err
	}
	defer input.Close() // nolint:errcheck
	in, err := ioutil.ReadAll(iotools.CancelableReader(cmdContext(cmd), input))
	if err != nil {
		return nil, err
	}
	env := new(gobl.Envelope)
	if err := yaml.Unmarshal(in, env); err != nil {
		return nil, err
	}
	return env, nil
}

func verify(cmd *cobra.Command, args []string) error {
	env, err := readEnv(cmd, args)
	if err != nil {
		return err
	}
	if err := env.Validate(); err != nil {
		return err
	}

	return env.Verify()
}

type genericDoc struct {
	typ     string
	payload json.RawMessage
}

var _ gobl.Document = &genericDoc{}

func (d *genericDoc) Type() string { return d.typ }

// MarshalJSON satisfies the json.Marshaler interface.
func (d *genericDoc) MarshalJSON() ([]byte, error) { // nolint:unparam
	return d.payload, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (d *genericDoc) UnmarshalJSON(p []byte) error { // nolint:unparam
	d.payload = p
	return nil
}

func version() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "GOBL version %s\n", gobl.VERSION)
		},
	}
}
