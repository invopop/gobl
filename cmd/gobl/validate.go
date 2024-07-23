package main

import (
	"github.com/invopop/gobl/internal/cli"
	"github.com/spf13/cobra"
)

type validateOpts struct {
	*rootOpts

	// Command options
	use   string
	short string
}

func validate(root *rootOpts) *validateOpts {
	return &validateOpts{
		rootOpts: root,
		use:      "validate [infile] [outfile]",
		short:    "Validate checks if the input is a valid GOBL document",
	}
}

func (opts *validateOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.MaximumNArgs(2),
		RunE:  opts.runE,
		Use:   opts.use,
		Short: opts.short,
	}

	return cmd
}

func (opts *validateOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

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

	return cli.Validate(ctx, input)
}
