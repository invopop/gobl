package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/cli"
)

type signOpts struct {
	*rootOpts
	set            map[string]string
	setFiles       map[string]string
	setStrings     map[string]string
	template       string
	privateKeyFile string
	docType        string

	// Command options
	use   string
	short string
}

func sign(root *rootOpts) *signOpts {
	return &signOpts{
		rootOpts: root,
		use:      "sign [infile] [outfile]",
		short:    "Signs an envelope using a private key",
	}
}

func (opts *signOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.MaximumNArgs(2),
		RunE:  opts.runE,
		Use:   opts.use,
		Short: opts.short,
	}

	f := cmd.Flags()
	f.StringToStringVar(&opts.set, "set", nil, "Set value from the command line")
	f.StringToStringVar(&opts.setFiles, "set-file", nil, "Set value from the specified YAML or JSON file")
	f.StringToStringVar(&opts.setStrings, "set-string", nil, "Set STRING value from the command line")
	f.StringVarP(&opts.template, "template", "T", "", "Template YAML/JSON file into which data is merged")
	f.StringVarP(&opts.privateKeyFile, "key", "k", defaultKeyFilename, "Private key file for signing")
	f.StringVarP(&opts.docType, "type", "t", "", "Specify the document type")

	return cmd
}

func (opts *signOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	var template io.Reader
	if opts.template != "" {
		f, err := os.Open(opts.template)
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

	out, err := opts.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	key, err := loadPrivateKey(opts.privateKeyFile)
	if err != nil {
		return err
	}

	signOpts := &cli.SignOptions{
		ParseOptions: &cli.ParseOptions{
			Template:  template,
			Input:     input,
			SetFile:   opts.setFiles,
			SetYAML:   opts.set,
			SetString: opts.setStrings,
			DocType:   opts.docType,
		},
		PrivateKey: key,
	}

	env, err := cli.Sign(ctx, signOpts)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	if opts.indent {
		enc.SetIndent("", "\t")
	}

	return enc.Encode(env)
}

func loadPrivateKey(file string) (*dsig.PrivateKey, error) {
	pkFilename, err := expandHome(file)
	if err != nil {
		return nil, err
	}
	keyFile, err := os.Open(pkFilename)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close() // nolint:errcheck

	key := new(dsig.PrivateKey)
	if err = json.NewDecoder(keyFile).Decode(key); err != nil {
		return nil, err
	}

	return key, nil
}
