package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/internal/ops"
	goblnet "github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
)

type netWhoOpts struct {
	*rootOpts
	from     string
	insecure bool
}

func netWho(root *rootOpts) *netWhoOpts {
	return &netWhoOpts{rootOpts: root}
}

func (w *netWhoOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "who <address>",
		Short: "Look up the party a GOBL Net domain belongs to (EXPERIMENTAL)",
		Long: "Fetch and verify the org.Party published at a GOBL Net address's\n" +
			"/.well-known/gobl/who endpoint. The request is authenticated as the\n" +
			"--from domain, so /who is a mutual party exchange.\n\n" +
			"EXPERIMENTAL: GOBL Net is under active development and may change without notice.",
		Args: cobra.ExactArgs(1),
		RunE: w.runE,
	}
	f := cmd.Flags()
	f.StringVar(&w.from, "from", "", "Local domain identity (~/.config/gobl/<from>/) used to sign the request")
	f.BoolVar(&w.insecure, "insecure", false, "Query over plain HTTP and permit host:port (development)")
	_ = cmd.MarkFlagRequired("from")
	return cmd
}

func (w *netWhoOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := commandContext(cmd)

	if w.from == "" {
		return errors.New("--from is required to authenticate the request")
	}
	dir := filepath.Join(defaultConfigDir(), w.from)

	key, err := loadPrivateKey(filepath.Join(dir, "private.jwk"))
	if err != nil {
		return err
	}
	partyData, err := os.ReadFile(filepath.Join(dir, "party.json"))
	if err != nil {
		return err
	}
	party := new(org.Party)
	if err := json.Unmarshal(partyData, party); err != nil {
		return err
	}

	result, err := ops.NetWho(ctx, &ops.NetWhoOptions{
		Target:    goblnet.Address(args[0]),
		From:      goblnet.Address(w.from),
		FromKey:   key,
		FromParty: party,
		Insecure:  w.insecure,
	})
	if err != nil {
		return err
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	if w.indent {
		enc.SetIndent("", "\t")
	}
	return enc.Encode(result)
}
