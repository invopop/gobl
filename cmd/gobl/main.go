// The gobl command provides a command-line interface to the GOBL library.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghodss/yaml"
	"github.com/invopop/gobl"
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
		RunE: build,
	})
	root.AddCommand(&cobra.Command{
		Use:  "verify",
		RunE: verify,
	})
	return root.ExecuteContext(ctx)
}

func build(*cobra.Command, []string) error {
	return nil
}

func verify(cmd *cobra.Command, args []string) error {
	input := cmd.InOrStdin()
	if len(args) > 0 && args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close() // nolint:errcheck
		input = f
	}
	in, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	env := new(gobl.Envelope)
	if err := yaml.Unmarshal(in, env); err != nil {
		return err
	}
	if err := env.Validate(); err != nil {
		return err
	}
	return env.Verify()
}
