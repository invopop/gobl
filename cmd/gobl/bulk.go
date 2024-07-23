package main

import (
	"encoding/json"

	"github.com/invopop/gobl/internal/cli"
	"github.com/spf13/cobra"
)

type bulkOpts struct {
	*rootOpts
}

func (o *bulkOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	in, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer in.Close() // nolint:errcheck

	out, err := o.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	enc := json.NewEncoder(out)
	if o.indent {
		enc.SetIndent("", "\t")
	}
	opts := &cli.BulkOptions{
		In: in,
	}
	for result := range cli.Bulk(ctx, opts) {
		if err := enc.Encode(result); err != nil {
			return err
		}
	}
	return nil
}
