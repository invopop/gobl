package main

import (
	"io"
	"os"

	"github.com/invopop/gobl/internal/cli"
	"github.com/spf13/cobra"
)

type buildOpts struct {
	*rootOpts
	set        map[string]string
	setFiles   map[string]string
	setStrings map[string]string
	template   string
	docType    string
	envelop    bool

	// Command options
	use   string
	short string
}

func build(root *rootOpts) *buildOpts {
	return &buildOpts{
		rootOpts: root,
		use:      "build [infile] [outfile]",
		short:    "Calculate and validate a document, wrapping it in envelope if needed",
	}
}

func (b *buildOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.MaximumNArgs(2),
		RunE:  b.runE,
		Use:   b.use,
		Short: b.short,
	}

	f := cmd.Flags()
	f.StringToStringVar(&b.set, "set", nil, "set value from the command line")
	f.StringToStringVar(&b.setFiles, "set-file", nil, "set value from the specified YAML or JSON file")
	f.StringToStringVar(&b.setStrings, "set-string", nil, "set STRING value from the command line")
	f.StringVarP(&b.template, "template", "T", "", "template YAML/JSON file into which data is merged")
	f.StringVarP(&b.docType, "type", "t", "", "specify the document type")
	f.BoolVarP(&b.envelop, "envelop", "e", false, "insert document into an envelope")

	return cmd
}

func (b *buildOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

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
	defer input.Close() // nolint:errcheck

	out, err := b.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	buildOpts := &cli.BuildOptions{
		ParseOptions: &cli.ParseOptions{
			Template:  template,
			Input:     input,
			SetFile:   b.setFiles,
			SetYAML:   b.set,
			SetString: b.setStrings,
			DocType:   b.docType,
			Envelop:   b.envelop,
		},
	}

	res, err := cli.Build(ctx, buildOpts)
	if err != nil {
		return err
	}
	return encode(res, out, b.indent)
}
