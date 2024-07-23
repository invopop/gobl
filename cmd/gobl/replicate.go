package main

import (
	"encoding/json"

	"github.com/invopop/gobl/internal/cli"
	"github.com/spf13/cobra"
)

type replicateOpts struct {
	*rootOpts
}

func replicate(root *rootOpts) *replicateOpts {
	return &replicateOpts{
		rootOpts: root,
	}
}

func (o *replicateOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.MaximumNArgs(2),
		RunE:  o.runE,
		Use:   "replicate [infile] [outfile]",
		Short: "Replicate a document from the provided input",
	}
	return cmd
}

func (o *replicateOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	out, err := o.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	rOpts := &cli.ReplicateOptions{
		ParseOptions: &cli.ParseOptions{
			Input: input,
		},
	}

	obj, err := cli.Replicate(ctx, rOpts)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	if o.indent {
		enc.SetIndent("", "\t")
	}

	return enc.Encode(obj)
}
