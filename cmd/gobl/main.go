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
	root.AddCommand(build())
	root.AddCommand(version())
	return root
}

func inputFilename(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func readEnv(cmd *cobra.Command, args []string) (*gobl.Envelope, error) {
	input := cmd.InOrStdin()
	if inFile := inputFilename(args); inFile != "" {
		f, err := os.Open(inFile)
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
		return nil, fmt.Errorf("unrecognized document type: %s", env.Head.Type)
	}
}

type buildOpts struct {
	overwriteOutputFile bool
	inPlace             bool
}

func build() *cobra.Command {
	opts := &buildOpts{}
	cmd := &cobra.Command{
		Use:  "build [infile] [outfile]",
		Args: cobra.MaximumNArgs(2),
		RunE: opts.RunE,
	}

	f := cmd.Flags()

	f.BoolVarP(&opts.overwriteOutputFile, "force", "f", false, "force writing output file, even if it exists")
	f.BoolVarP(&opts.inPlace, "in-place", "w", false, "overwrite the input file in place")

	return cmd
}

func (b *buildOpts) outputFilename(args []string) string {
	if b.inPlace {
		return inputFilename(args)
	}
	if len(args) >= 2 && args[1] != "-" {
		return args[1]
	}
	return ""
}

func (b *buildOpts) RunE(cmd *cobra.Command, args []string) error {
	env, err := readEnv(cmd, args)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if outFile := b.outputFilename(args); outFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if !b.overwriteOutputFile && !b.inPlace {
			flags |= os.O_EXCL
		}
		f, err := os.OpenFile(outFile, flags, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close() // nolint:errcheck
		out = f
	} else if b.inPlace {
		return errors.New("cannot overwrite STDIN")
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
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	return enc.Encode(env)
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
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "GOBL version %s", gobl.VERSION)
		},
	}
}
