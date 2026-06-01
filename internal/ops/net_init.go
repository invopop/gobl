package ops

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
)

// InitOptions configures InitDomain.
type InitOptions struct {
	ConfigDir string
	Domain    string
	Name      string // optional party name seed
	Force     bool   // overwrite a non-empty existing directory
	Out       io.Writer
	Log       *slog.Logger // optional; defaults to slog.Default()
}

// InitDomain scaffolds a new GOBL Net domain identity under
// <ConfigDir>/<Domain>/: a single private key (private.jwk), the
// matching public key as keys/<kid>.json (stamped with valid_from), a
// raw org.Party template with a pre-filled gobl: endpoint, and an
// inbox/ directory. The party is intentionally left unsigned — serve
// signs it on demand.
func InitDomain(opts *InitOptions) error {
	log := logger(opts.Log)
	if opts.Domain == "" {
		return fmt.Errorf("init: domain is required")
	}

	dc := domainConfigFor(opts.ConfigDir, opts.Domain)
	dir := filepath.Join(opts.ConfigDir, opts.Domain)

	if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 && !opts.Force {
		return fmt.Errorf("init: %s already exists and is not empty (use --force to overwrite)", dir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("init: create domain dir: %w", err)
	}

	if _, err := generateKeypair(dc.KeysDir, dc.PrivateKeyFile, log); err != nil {
		return err
	}

	party := &org.Party{
		Name: opts.Name,
		Endpoints: []*org.Endpoint{
			{URI: cbc.URI(net.Scheme + ":" + opts.Domain)},
		},
	}
	partyBytes, err := json.MarshalIndent(party, "", "  ")
	if err != nil {
		return fmt.Errorf("init: marshal party: %w", err)
	}
	if err := os.WriteFile(dc.PartyFile, partyBytes, 0o644); err != nil {
		return fmt.Errorf("init: write party: %w", err)
	}
	if err := os.MkdirAll(dc.InboxDir, 0o755); err != nil {
		return fmt.Errorf("init: create inbox dir: %w", err)
	}

	log.Info("initialised domain", "domain", opts.Domain, "party", dc.PartyFile, "inbox", dc.InboxDir)
	return nil
}
