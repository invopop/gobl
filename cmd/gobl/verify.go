package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/cli"
)

type verifyOpts struct {
	publicKeyFile string
}

func verify() *verifyOpts {
	return &verifyOpts{}
}

func (v *verifyOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verify [infile]",
		Args: cobra.MaximumNArgs(1),
		RunE: v.runE,
	}

	f := cmd.Flags()

	f.StringVarP(&v.publicKeyFile, "key", "k", pubfileFromPriv(defaultKeyFilename), "Public key file for signature validation")

	return cmd
}

func (v *verifyOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	pbFilename, err := expandHome(v.publicKeyFile)
	if err != nil {
		return err
	}
	keyFile, err := os.Open(pbFilename)
	if err != nil {
		return err
	}
	defer keyFile.Close() // nolint:errcheck

	key := new(dsig.PublicKey)
	if err = json.NewDecoder(keyFile).Decode(key); err != nil {
		return err
	}

	return cli.Verify(ctx, input, key)
}
