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
		Use:  "verify [infile]",
		Args: cobra.MaximumNArgs(1),
		RunE: verify,
	})
	root.AddCommand(&cobra.Command{
		Use:  "build [infile] [outfile]",
		Args: cobra.MaximumNArgs(2),
		RunE: build,
	})
	return root.ExecuteContext(ctx)
}

func readEnv(cmd *cobra.Command, args []string) (*gobl.Envelope, error) {
	input := cmd.InOrStdin()
	if len(args) > 0 && args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}
		defer f.Close() // nolint:errcheck
		input = f
	}
	in, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	env := new(gobl.Envelope)
	if err := yaml.Unmarshal(in, env); err != nil {
		return nil, err
	}
	if err := env.Validate(); err != nil {
		return nil, err
	}
	return env, nil
}

func verify(cmd *cobra.Command, args []string) error {
	env, err := readEnv(cmd, args)
	if err != nil {
		return err
	}

	return env.Verify()
}

func build(cmd *cobra.Command, args []string) error {
	env, err := readEnv(cmd, args)
	if err != nil {
		return err
	}
	doc := new(gobl.Document)
	if err := env.Extract(doc); err != nil {
		return err
	}
	return nil
}
