package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/cli"
	goblnet "github.com/invopop/gobl/net"
)

type verifyOpts struct {
	publicKeyFile string
	address       string
	remote        bool
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
	f.StringVarP(&v.address, "address", "a", "", "GOBL Net address (FQDN) for remote key discovery")
	f.BoolVarP(&v.remote, "remote", "r", false, "Auto-discover keys from the signature's gn header")

	return cmd
}

func (v *verifyOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	if v.address != "" || v.remote {
		var addr goblnet.Address
		if v.address != "" {
			addr, err = goblnet.ParseAddress(v.address)
			if err != nil {
				return err
			}
		}
		return cli.VerifyRemote(ctx, input, goblnet.NewClient(), addr)
	}

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
