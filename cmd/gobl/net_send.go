package main

import (
	"github.com/spf13/cobra"

	"github.com/invopop/gobl/internal/ops"
	goblnet "github.com/invopop/gobl/net"
)

type netSendOpts struct {
	*rootOpts
	to       string
	insecure bool
}

func netSend(root *rootOpts) *netSendOpts {
	return &netSendOpts{rootOpts: root}
}

func (s *netSendOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [infile]",
		Short: "Send a signed GOBL envelope to a GOBL Net inbox (EXPERIMENTAL)",
		Long: "Send a signed GOBL envelope to a GOBL Net inbox.\n\n" +
			"EXPERIMENTAL: GOBL Net is under active development and may change without notice.",
		Args: cobra.MaximumNArgs(1),
		RunE: s.runE,
	}
	f := cmd.Flags()
	f.StringVarP(&s.to, "to", "t", "", "Destination GOBL Net address (FQDN, or host:port with --insecure)")
	f.BoolVar(&s.insecure, "insecure", false, "Use plain HTTP and permit host:port form in --to (development)")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func (s *netSendOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	return ops.NetSend(ctx, &ops.NetSendOptions{
		Input:    input,
		To:       goblnet.Address(s.to),
		Insecure: s.insecure,
	})
}
