package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/dsig"
)

type keygenOpts struct {
	filename  string
	overwrite bool
}

func keygen() *keygenOpts {
	return &keygenOpts{}
}

func (k *keygenOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "keygen [flags] [outfile]",
		Args: cobra.MaximumNArgs(1),
		RunE: k.runE,
	}

	f := cmd.Flags()

	f.BoolVarP(&k.overwrite, "force", "f", false, "force writing output file, even if it exists")

	return cmd
}

func outputKeyfile(args []string) string {
	if len(args) == 0 {
		return "~/.gobl/id_es256"
	}
	return args[0]
}

func (k *keygenOpts) runE(cmd *cobra.Command, args []string) error {
	key := dsig.NewES256Key()
	priv, err := json.Marshal(key)
	if err != nil {
		return err
	}
	outfile := outputKeyfile(args)
	if outfile == "-" {
		fmt.Fprintln(cmd.OutOrStdout(), string(priv))
		return nil
	}
	return nil
}
