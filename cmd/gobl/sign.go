package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/ops"
	goblnet "github.com/invopop/gobl/net"
)

type signOpts struct {
	*rootOpts
	set            map[string]string
	setFiles       map[string]string
	setStrings     map[string]string
	template       string
	privateKeyFile string
	domain         string
	audience       string
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
	f.StringVar(&opts.domain, "domain", "", "Sign with the key from ~/.config/gobl/<domain>/ and stamp iss=gobl:<domain>")
	f.StringVar(&opts.audience, "to", "", "GOBL Net address to bind the signature to (stamps aud=gobl:<to>)")
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

	keyFile := opts.privateKeyFile
	var iss, aud cbc.URI
	if opts.domain != "" {
		if cmd.Flags().Changed("key") {
			return errors.New("--domain and --key are mutually exclusive")
		}
		keyFile = filepath.Join(defaultConfigDir(), opts.domain, "private.jwk")
		iss = goblnet.Address(opts.domain).URI()
	}
	if opts.audience != "" {
		aud = goblnet.Address(opts.audience).URI()
	}

	key, err := loadPrivateKey(keyFile)
	if err != nil {
		return err
	}

	signOpts := &ops.SignOptions{
		ParseOptions: &ops.ParseOptions{
			Template:  template,
			Input:     input,
			SetFile:   opts.setFiles,
			SetYAML:   opts.set,
			SetString: opts.setStrings,
			DocType:   opts.docType,
		},
		PrivateKey: key,
		Iss:        iss,
		Aud:        aud,
	}

	env, err := ops.Sign(ctx, signOpts)
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
