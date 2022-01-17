// The gobl command provides a command-line interface to the GOBL library.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
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

func extractDoc(env *gobl.Envelope) (gobl.Document, error) {
	switch env.Head.Type {
	case bill.InvoiceType:
		doc := new(bill.Invoice)
		err := env.Extract(doc)
		return doc, err
	default:
		doc := &genericDoc{
			typ: env.Head.Type,
		}
		err := env.Extract(doc)
		return doc, err
	}
}

func build(cmd *cobra.Command, args []string) error {
	env, err := readEnv(cmd, args)
	if err != nil {
		return err
	}
	if env.Document == nil {
		return errors.New("no document included")
	}
	doc, err := extractDoc(env)
	if err != nil {
		return err
	}
	if err := env.Insert(doc); err != nil {
		return err
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "\t")
	return enc.Encode(env)
}

type genericDoc struct {
	typ     string
	payload json.RawMessage
}

var _ gobl.Document = &genericDoc{}

func (d *genericDoc) Type() string { return d.typ }

func (d *genericDoc) MarshalJSON() ([]byte, error) {
	return d.payload, nil
}

func (d *genericDoc) UnmarshalJSON(p []byte) error {
	d.payload = p
	return nil
}
