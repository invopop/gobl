package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/cmd/gobl/internal"
)

type buildOpts struct {
	overwriteOutputFile bool
	inPlace             bool
	set                 map[string]string
	setFiles            map[string]string
	setStrings          map[string]string
	template            string
}

func build() *buildOpts {
	return &buildOpts{}
}

func (b *buildOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "build [infile] [outfile]",
		Args: cobra.MaximumNArgs(2),
		RunE: b.runE,
	}

	f := cmd.Flags()

	f.BoolVarP(&b.overwriteOutputFile, "force", "f", false, "force writing output file, even if it exists (only outputs JSON)")
	f.BoolVarP(&b.inPlace, "in-place", "w", false, "overwrite the input file in place")
	f.StringToStringVar(&b.set, "set", nil, "set value from the command line")
	f.StringToStringVar(&b.setFiles, "set-file", nil, "set value from the specified YAML or JSON file")
	f.StringToStringVar(&b.setStrings, "set-string", nil, "set STRING value from the command line")
	f.StringVarP(&b.template, "template", "T", "", "Template YAML/JSON file into which data is merged")

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

func cmdContext(cmd *cobra.Command) context.Context {
	if ctx := cmd.Context(); ctx != nil {
		return ctx
	}
	return context.Background()
}

func (b *buildOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := cmdContext(cmd)

	var template io.Reader
	if b.template != "" {
		f, err := os.Open(b.template)
		if err != nil {
			return err
		}
		defer f.Close() // nolint:errcheck
		template = f
	}

	input, err := openInput(cmd, args)
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
	defer input.Close() // nolint:errcheck

	env, err := internal.Build(ctx, internal.BuildOptions{
		Template:  template,
		Data:      input,
		SetFile:   b.setFiles,
		SetYAML:   b.set,
		SetString: b.setStrings,
	})
	if err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	return enc.Encode(env)
}
