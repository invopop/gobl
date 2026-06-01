package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/internal/ops"
)

type netServeOpts struct {
	*rootOpts
	configDir  string
	partyFile  string
	keysDir    string
	privateKey string
	inboxDir   string

	httpPort  int
	httpsPort int

	acmeLive  bool
	acmeTest  bool
	domain    string
	acmeEmail string
	certDir   string

	tlsCert string
	tlsKey  string
}

func netServe(root *rootOpts) *netServeOpts {
	return &netServeOpts{rootOpts: root}
}

func (s *netServeOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the GOBL Net well-known endpoints (EXPERIMENTAL)",
		Long: "Serve the GOBL Net well-known endpoints.\n\n" +
			"EXPERIMENTAL: GOBL Net is under active development and may change without notice.",
		RunE: s.runE,
	}
	configDir := defaultConfigDir()

	f := cmd.Flags()
	f.StringVar(&s.configDir, "config-dir", configDir, "Base directory; its <domain>/ subdirectories are auto-discovered and served, routed by Host")
	f.StringVar(&s.partyFile, "party", "", "Manual single-identity mode: party.json (raw org.Party or signed envelope) served at /.well-known/gobl/who")
	f.StringVarP(&s.keysDir, "keys-dir", "k", "", "Manual single-identity mode: directory of <kid>.json public keys published at /.well-known/gobl/keys/<kid>")
	f.StringVar(&s.privateKey, "private-key", "", "Manual single-identity mode: private key paired with the JWKS")
	f.StringVar(&s.inboxDir, "inbox", "", "Manual single-identity mode: directory to write accepted envelopes into")

	f.IntVar(&s.httpPort, "http-port", 80, "HTTP listen port")
	f.IntVar(&s.httpsPort, "https-port", 443, "HTTPS listen port (used only when a TLS source is configured)")

	f.BoolVar(&s.acmeLive, "acme-live", false, "Activate HTTPS via Let's Encrypt production directory")
	f.BoolVar(&s.acmeTest, "acme-test", false, "Activate HTTPS via Let's Encrypt staging directory (for testing)")
	f.StringVar(&s.domain, "domain", "", "Hostname the ACME client is allowed to issue for; MUST match the participant's GOBL Net address")
	f.StringVar(&s.acmeEmail, "acme-email", "", "Account email for ACME registration")
	f.StringVar(&s.certDir, "cert-dir", "", "Directory to cache ACME-issued certificates (default <config-dir>/certs)")
	f.StringVar(&s.tlsCert, "tls-cert", "", "PEM-encoded TLS certificate; activates HTTPS with file-based TLS")
	f.StringVar(&s.tlsKey, "tls-key", "", "PEM-encoded TLS private key paired with --tls-cert")

	return cmd
}

func (s *netServeOpts) runE(cmd *cobra.Command, _ []string) error {
	if s.configDir == "" {
		s.configDir = defaultConfigDir()
	}
	if s.certDir == "" {
		s.certDir = filepath.Join(s.configDir, "certs")
	}
	if err := s.validate(); err != nil {
		return err
	}

	opts := &ops.NetServeOptions{
		ConfigDir: s.configDir,
		Out:       cmd.OutOrStdout(),

		HTTPPort:  s.httpPort,
		HTTPSPort: s.httpsPort,

		ACMELive:  s.acmeLive,
		ACMETest:  s.acmeTest,
		Domain:    s.domain,
		ACMEEmail: s.acmeEmail,
		CertDir:   s.certDir,

		CertFile: s.tlsCert,
		KeyFile:  s.tlsKey,
	}

	// Manual single-identity mode: triggered by an explicit --party or
	// --keys-dir. Unset companion paths default to the flat config-dir
	// layout.
	if cmd.Flags().Changed("party") || cmd.Flags().Changed("keys-dir") {
		opts.PartyFile = orDefault(s.partyFile, filepath.Join(s.configDir, "party.json"))
		opts.KeysDir = orDefault(s.keysDir, filepath.Join(s.configDir, "keys"))
		opts.PrivateKeyFile = orDefault(s.privateKey, filepath.Join(s.configDir, "private.jwk"))
		opts.InboxDir = orDefault(s.inboxDir, filepath.Join(s.configDir, "inbox"))
	}

	ctx := commandContext(cmd)
	return ops.NetServe(ctx, opts)
}

func orDefault(v, def string) string {
	if v != "" {
		return v
	}
	return def
}

func (s *netServeOpts) validate() error {
	if s.acmeLive && s.acmeTest {
		return errors.New("--acme-live and --acme-test are mutually exclusive")
	}
	acme := s.acmeLive || s.acmeTest
	fileTLS := s.tlsCert != "" || s.tlsKey != ""
	if acme && fileTLS {
		return errors.New("--acme-* and --tls-cert/--tls-key are mutually exclusive")
	}
	// Note: --domain is optional here. When absent and ACME is active,
	// the domain is derived from the party's GOBL inbox at startup; the
	// ops layer fails clearly if neither source provides one.
	if (s.tlsCert == "") != (s.tlsKey == "") {
		return errors.New("--tls-cert and --tls-key must be provided together")
	}
	return nil
}

func defaultConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "gobl"
	}
	return filepath.Join(home, ".config", "gobl")
}
