package main

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/internal/ops"
)

type initCmdOpts struct {
	*rootOpts
	configDir string
	name      string
	force     bool
}

func initCmd(root *rootOpts) *initCmdOpts {
	return &initCmdOpts{rootOpts: root}
}

func (o *initCmdOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <domain>",
		Short: "Initialise a new GOBL Net domain identity (EXPERIMENTAL)",
		Long: "Initialise a new GOBL Net domain identity.\n\n" +
			"EXPERIMENTAL: GOBL Net is under active development and may change without notice.",
		Args: cobra.ExactArgs(1),
		RunE: o.runE,
	}
	f := cmd.Flags()
	f.StringVar(&o.configDir, "config-dir", defaultConfigDir(), "Base directory for domain identities")
	f.StringVar(&o.name, "name", "", "Party name to seed into the generated party.json")
	f.BoolVarP(&o.force, "force", "f", false, "Overwrite an existing non-empty domain directory")
	return cmd
}

func (o *initCmdOpts) runE(cmd *cobra.Command, args []string) error {
	domain := args[0]
	if domain == "" {
		return errors.New("a domain is required")
	}
	return ops.InitDomain(&ops.InitOptions{
		ConfigDir: o.configDir,
		Domain:    domain,
		Name:      o.name,
		Force:     o.force,
		Out:       cmd.OutOrStdout(),
	})
}
