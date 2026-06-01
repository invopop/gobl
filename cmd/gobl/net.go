package main

import "github.com/spf13/cobra"

type netOpts struct {
	*rootOpts
}

func netCmd(root *rootOpts) *netOpts {
	return &netOpts{rootOpts: root}
}

func (n *netOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "net",
		Short: "GOBL Net operations (EXPERIMENTAL)",
		Long: "GOBL Net operations.\n\n" +
			"EXPERIMENTAL: GOBL Net is under active development. Commands, on-disk\n" +
			"layout, and the wire protocol may change without notice.",
	}
	cmd.AddCommand(netServe(n.rootOpts).cmd())
	cmd.AddCommand(netSend(n.rootOpts).cmd())
	cmd.AddCommand(netWho(n.rootOpts).cmd())
	return cmd
}
