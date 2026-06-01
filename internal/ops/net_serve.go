package ops

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	stdnet "net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
)

const (
	netServeShutdownTimeout = 10 * time.Second
	netInboxMaxBody         = 1 << 20 // 1 MiB

	defaultHTTPPort  = 80
	defaultHTTPSPort = 443

	acmeStagingDirectoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// NetServeOptions configures the GOBL Net HTTP server.
type NetServeOptions struct {
	// ConfigDir is the base directory whose <domain>/ subdirectories are
	// auto-discovered when no explicit single identity is provided.
	ConfigDir string

	// Explicit single-identity ("manual") mode: when PartyFile or KeysDir
	// is set, exactly one identity is served from these paths.
	PartyFile      string
	KeysDir        string // directory of <kid>.json public JWK files
	PrivateKeyFile string
	InboxDir       string

	Client *net.Client // optional; defaults to net.NewClient()
	Out    io.Writer   // optional; defaults to os.Stdout

	// Port overrides (zero means use the default — 80 / 443).
	HTTPPort  int
	HTTPSPort int

	// ACME options. ACMELive and ACMETest are mutually exclusive.
	ACMELive  bool
	ACMETest  bool
	Domain    string // restricts multi-domain discovery to one, or names the manual identity
	ACMEEmail string
	CertDir   string

	// File-based TLS. CertFile and KeyFile must be supplied together.
	CertFile string
	KeyFile  string
}

// domainConfig groups the on-disk paths that make up one GOBL Net
// identity. The directory name is the domain.
type domainConfig struct {
	Domain         string
	KeysDir        string // directory of <kid>.json public JWK files
	PrivateKeyFile string
	PartyFile      string
	InboxDir       string
	AllowFile      string
}

// domainConfigFor builds the standard paths for a domain inside configDir.
func domainConfigFor(configDir, domain string) domainConfig {
	dir := filepath.Join(configDir, domain)
	return domainConfig{
		Domain:         domain,
		KeysDir:        filepath.Join(dir, "keys"),
		PrivateKeyFile: filepath.Join(dir, "private.jwk"),
		PartyFile:      filepath.Join(dir, "party.json"),
		InboxDir:       filepath.Join(dir, "inbox"),
		AllowFile:      filepath.Join(dir, "allow.json"),
	}
}

// loadAllowList reads <domain>/allow.json (a JSON array of GOBL Net
// addresses). It returns the set of accepted addresses and whether a
// list is configured at all. An absent file means "accept any verified
// caller" (present == false).
func loadAllowList(dc domainConfig) (map[net.Address]bool, bool, error) {
	if dc.AllowFile == "" || !fileExists(dc.AllowFile) {
		return nil, false, nil
	}
	data, err := os.ReadFile(dc.AllowFile)
	if err != nil {
		return nil, false, fmt.Errorf("net serve: read allow list: %w", err)
	}
	var addrs []net.Address
	if err := json.Unmarshal(data, &addrs); err != nil {
		return nil, false, fmt.Errorf("net serve: invalid allow list: %w", err)
	}
	set := make(map[net.Address]bool, len(addrs))
	for _, a := range addrs {
		set[a] = true
	}
	return set, true, nil
}

// allowed reports whether addr may call a protected endpoint: any
// verified caller when no list is configured, otherwise only listed ones.
func allowed(set map[net.Address]bool, present bool, addr net.Address) bool {
	return !present || set[addr]
}

// discoverDomains lists the immediate subdirectories of configDir (skipping
// "certs") that look like a domain identity (containing a keys/ dir
// and/or a party.json), returning a domainConfig for each.
func discoverDomains(configDir string) ([]domainConfig, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("net serve: read config dir: %w", err)
	}
	var out []domainConfig
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "certs" {
			continue
		}
		dc := domainConfigFor(configDir, e.Name())
		if dirExists(dc.KeysDir) || fileExists(dc.PartyFile) {
			out = append(out, dc)
		}
	}
	return out, nil
}

// NetServeHandler builds a single-identity HTTP handler from explicit
// options (manual mode). Multi-domain serving uses buildRouter. It is
// exported so tests can drive the resulting handler via httptest.
func NetServeHandler(opts *NetServeOptions) (http.Handler, error) {
	out := opts.Out
	if out == nil {
		out = io.Discard
	}
	client := opts.Client
	if client == nil {
		client = net.NewClient()
	}
	dc := domainConfig{
		Domain:         opts.Domain,
		KeysDir:        opts.KeysDir,
		PrivateKeyFile: opts.PrivateKeyFile,
		PartyFile:      opts.PartyFile,
		InboxDir:       opts.InboxDir,
	}
	return buildDomainHandler(dc, client, out)
}

// buildDomainHandler prepares one domain's on-disk state (keys, party,
// inbox, allow-list) and returns its mux.
//
//   - GET  /keys  — open, serves the public JWKS.
//   - POST /who   — authenticated party exchange (see handleWho).
//   - POST /inbox — authenticated envelope delivery (see handleInbox).
func buildDomainHandler(dc domainConfig, client *net.Client, out io.Writer) (http.Handler, error) {
	keysByKID, err := ensureKeys(dc, out)
	if err != nil {
		return nil, err
	}
	priv, err := loadPrivateKeyFile(dc.PrivateKeyFile)
	if err != nil {
		return nil, err
	}
	partyEnv, err := readPartyEnvelope(dc)
	if err != nil {
		return nil, err
	}
	partyEnvBytes, err := json.Marshal(partyEnv) // canonical, unsigned, stable UUID
	if err != nil {
		return nil, fmt.Errorf("net serve: marshal party: %w", err)
	}
	if err := os.MkdirAll(dc.InboxDir, 0o755); err != nil {
		return nil, fmt.Errorf("net serve: create inbox dir: %w", err)
	}
	allow, present, err := loadAllowList(dc)
	if err != nil {
		return nil, err
	}
	var self cbc.URI
	if dc.Domain != "" {
		self = net.Address(dc.Domain).URI()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET "+net.KeysPath+"/{kid}", handleKey(keysByKID))
	mux.HandleFunc("POST "+net.WhoPath, handleWho(client, partyEnvBytes, priv, self, allow, present))
	mux.HandleFunc("POST "+net.InboxPath, handleInbox(client, dc.InboxDir, self, allow, present))
	return mux, nil
}

// handleKey serves a single published JWK by its kid path value, or 404
// if the kid is not in the domain's published set.
func handleKey(keysByKID map[string][]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kid := r.PathValue("kid")
		body, ok := keysByKID[kid]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, _ = w.Write(body)
	}
}

// buildRouter returns an HTTP handler dispatching by the request Host
// header to the matching domain's handler. A single unnamed identity
// (manual mode without a domain) is served for all hosts.
func buildRouter(domains []domainConfig, client *net.Client, out io.Writer) (http.Handler, error) {
	if len(domains) == 1 && domains[0].Domain == "" {
		return buildDomainHandler(domains[0], client, out)
	}
	handlers := make(map[string]http.Handler, len(domains))
	for _, dc := range domains {
		h, err := buildDomainHandler(dc, client, out)
		if err != nil {
			return nil, err
		}
		handlers[dc.Domain] = h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := stripPort(r.Host)
		if h, ok := handlers[host]; ok {
			h.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	}), nil
}

func stripPort(host string) string {
	if h, _, err := stdnet.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

// ensureKeys returns a map of kid → single-JWK JSON bytes for every
// public key published by this domain. Each key lives in its own file
// at <keysDir>/<kid>.json. If neither the keys directory nor the
// private key file exist, a fresh ECDSA P-256 keypair is generated and
// persisted (single-key bootstrap). If only one of the two exists the
// setup is inconsistent.
func ensureKeys(dc domainConfig, out io.Writer) (map[string][]byte, error) {
	keysExists := dirExists(dc.KeysDir)
	privExists := fileExists(dc.PrivateKeyFile)

	switch {
	case keysExists && privExists:
		keysByKID, err := readKeysDir(dc.KeysDir)
		if err != nil {
			return nil, err
		}
		if len(keysByKID) == 0 {
			return nil, fmt.Errorf("net serve: keys directory %s contains no JWKs", dc.KeysDir)
		}
		priv, err := loadPrivateKeyFile(dc.PrivateKeyFile)
		if err != nil {
			return nil, err
		}
		if _, ok := keysByKID[priv.ID()]; !ok {
			return nil, fmt.Errorf("net serve: private key kid %q is not published under %s", priv.ID(), dc.KeysDir)
		}
		return keysByKID, nil

	case !keysExists && !privExists:
		return generateKeypair(dc.KeysDir, dc.PrivateKeyFile, out)

	default:
		present, missing := dc.KeysDir, dc.PrivateKeyFile
		if !keysExists {
			present, missing = dc.PrivateKeyFile, dc.KeysDir
		}
		return nil, fmt.Errorf(
			"net serve: inconsistent key setup — %s exists but %s does not. "+
				"Remove both to auto-generate, or supply both.",
			present, missing,
		)
	}
}

// readKeysDir reads each <kid>.json file in dir, validates that the
// JWK's kid matches the filename stem, and returns the raw file bytes
// keyed by kid. Non-JSON entries and subdirectories are ignored so the
// operator can drop sidecar files alongside their keys.
func readKeysDir(dir string) (map[string][]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("net serve: read keys dir: %w", err)
	}
	keysByKID := make(map[string][]byte, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		kid := strings.TrimSuffix(name, ".json")
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("net serve: read %s: %w", name, err)
		}
		pk := new(dsig.PublicKey)
		if err := json.Unmarshal(data, pk); err != nil {
			return nil, fmt.Errorf("net serve: %s: invalid JWK: %w", name, err)
		}
		if pk.ID() != kid {
			return nil, fmt.Errorf("net serve: %s: filename kid %q does not match JWK kid %q", name, kid, pk.ID())
		}
		keysByKID[kid] = data
	}
	return keysByKID, nil
}

// generateKeypair creates an ECDSA P-256 keypair, writes the private
// key to privFile (0600) and the public key to keysDir/<kid>.json
// (stamping valid_from = now), logs the action, and returns the
// per-kid JWK map.
func generateKeypair(keysDir, privFile string, out io.Writer) (map[string][]byte, error) {
	if err := os.MkdirAll(filepath.Dir(privFile), 0o700); err != nil {
		return nil, fmt.Errorf("net serve: create config dir: %w", err)
	}
	priv := dsig.NewES256Key()
	privBytes, err := json.MarshalIndent(priv, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("net serve: marshal private key: %w", err)
	}
	if err := os.WriteFile(privFile, privBytes, 0o600); err != nil {
		return nil, fmt.Errorf("net serve: write private key: %w", err)
	}
	pubBytes, err := publishedKeyBytes(priv)
	if err != nil {
		return nil, fmt.Errorf("net serve: marshal public key: %w", err)
	}
	if err := os.MkdirAll(keysDir, 0o755); err != nil {
		return nil, fmt.Errorf("net serve: create keys dir: %w", err)
	}
	keyFile := filepath.Join(keysDir, priv.ID()+".json")
	if err := os.WriteFile(keyFile, pubBytes, 0o644); err != nil {
		return nil, fmt.Errorf("net serve: write key file: %w", err)
	}
	fmt.Fprintf(out, "Generated new keypair (kid=%s): %s, %s\n", priv.ID(), privFile, keyFile) //nolint:errcheck
	return map[string][]byte{priv.ID(): pubBytes}, nil
}

// dirExists reports whether path exists and is a directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// publishedKeyBytes marshals the public counterpart of priv as a
// dsig.PublicKey with valid_from stamped to the current UTC time.
func publishedKeyBytes(priv *dsig.PrivateKey) ([]byte, error) {
	pubJSON, err := json.Marshal(priv.Public())
	if err != nil {
		return nil, err
	}
	pk := new(dsig.PublicKey)
	if err := json.Unmarshal(pubJSON, pk); err != nil {
		return nil, err
	}
	now := cal.TimestampNow()
	pk.ValidFrom = &now
	return json.Marshal(pk)
}

func loadPrivateKeyFile(path string) (*dsig.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("net serve: read private key: %w", err)
	}
	k := new(dsig.PrivateKey)
	if err := json.Unmarshal(b, k); err != nil {
		return nil, fmt.Errorf("net serve: invalid private key: %w", err)
	}
	return k, nil
}

// readPartyEnvelope reads the domain's party.json (a raw org.Party or an
// envelope, possibly already signed by an external authority) and returns
// it as an unsigned *gobl.Envelope. The /who handler signs a fresh copy
// per request with iss=self, aud=requester.
func readPartyEnvelope(dc domainConfig) (*gobl.Envelope, error) {
	data, err := os.ReadFile(dc.PartyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"net serve: party file not found at %s — create one with `gobl init %s` "+
					"or supply a raw org.Party / signed envelope",
				dc.PartyFile, dc.Domain,
			)
		}
		return nil, fmt.Errorf("net serve: read party file: %w", err)
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err == nil && env.Document != nil && !env.Document.IsEmpty() {
		return env, nil
	}
	// Not an envelope — parse as a raw org.Party and wrap it.
	party := new(org.Party)
	if err := json.Unmarshal(data, party); err != nil {
		return nil, fmt.Errorf("net serve: party file: invalid JSON: %w", err)
	}
	env, err = gobl.Envelop(party)
	if err != nil {
		return nil, fmt.Errorf("net serve: party file: %w", err)
	}
	return env, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// resolveDomains determines which identities to serve: an explicit single
// identity (manual mode) when PartyFile/KeysDir are set, otherwise the
// domains discovered under ConfigDir (optionally filtered by Domain).
func resolveDomains(opts *NetServeOptions) ([]domainConfig, error) {
	if opts.PartyFile != "" || opts.KeysDir != "" {
		return []domainConfig{{
			Domain:         opts.Domain,
			KeysDir:        opts.KeysDir,
			PrivateKeyFile: opts.PrivateKeyFile,
			PartyFile:      opts.PartyFile,
			InboxDir:       opts.InboxDir,
		}}, nil
	}
	if opts.ConfigDir == "" {
		return nil, errors.New("net serve: no config dir configured")
	}
	all, err := discoverDomains(opts.ConfigDir)
	if err != nil {
		return nil, err
	}
	if opts.Domain != "" {
		for _, dc := range all {
			if dc.Domain == opts.Domain {
				return []domainConfig{dc}, nil
			}
		}
		// Not yet on disk — construct it (keys auto-generate; party required).
		return []domainConfig{domainConfigFor(opts.ConfigDir, opts.Domain)}, nil
	}
	return all, nil
}

func domainNames(domains []domainConfig) []string {
	var names []string
	for _, dc := range domains {
		if dc.Domain != "" {
			names = append(names, dc.Domain)
		}
	}
	return names
}

// NetServe runs the GOBL Net HTTP server. It always serves over plain
// HTTP and, when a TLS source is configured, additionally over HTTPS with
// identical content (no HTTP→HTTPS redirect). In the default mode it
// discovers every <domain>/ directory under ConfigDir and routes requests
// by the HTTP Host header. The server shuts down gracefully on ctx cancel.
func NetServe(ctx context.Context, opts *NetServeOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.Client == nil {
		opts.Client = net.NewClient()
	}

	domains, err := resolveDomains(opts)
	if err != nil {
		return err
	}
	if len(domains) == 0 {
		return gobl.ErrInput.WithReason("net serve: no domains configured — run `gobl init <domain>` or pass --party/--keys")
	}

	router, err := buildRouter(domains, opts.Client, opts.Out)
	if err != nil {
		return err
	}

	httpHandler := router
	var tlsConfig *tls.Config

	switch {
	case opts.ACMELive || opts.ACMETest:
		names := domainNames(domains)
		if len(names) == 0 {
			return gobl.ErrInput.WithReason("net serve: ACME requires named domains — use --domain or per-domain config directories")
		}
		m := newAutocertManager(opts, names)
		httpHandler = m.HTTPHandler(router)
		tlsConfig = m.TLSConfig()
		fmt.Fprintf(opts.Out, "ACME enabled for domains: %s\n", strings.Join(names, ", ")) //nolint:errcheck
	case opts.CertFile != "" && opts.KeyFile != "":
		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return fmt.Errorf("net serve: load TLS keypair: %w", err)
		}
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	httpPort := opts.HTTPPort
	if httpPort == 0 {
		httpPort = defaultHTTPPort
	}
	httpsPort := opts.HTTPSPort
	if httpsPort == 0 {
		httpsPort = defaultHTTPSPort
	}

	httpLn, err := listenTCP(httpPort)
	if err != nil {
		return err
	}

	var httpsLn stdnet.Listener
	if tlsConfig != nil {
		httpsLn, err = listenTCP(httpsPort)
		if err != nil {
			_ = httpLn.Close()
			return err
		}
	}

	return serveOnListeners(ctx, opts, httpHandler, router, tlsConfig, httpLn, httpsLn)
}

// serveOnListeners runs the HTTP (and optionally HTTPS) servers on the
// provided listeners. Both listeners are closed by the http.Server lifecycle.
func serveOnListeners(
	ctx context.Context,
	opts *NetServeOptions,
	httpHandler http.Handler,
	httpsHandler http.Handler,
	tlsConfig *tls.Config,
	httpLn stdnet.Listener,
	httpsLn stdnet.Listener,
) error {
	httpSrv := &http.Server{
		Handler:           httpHandler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	var httpsSrv *http.Server
	if httpsLn != nil {
		httpsSrv = &http.Server{
			Handler:           httpsHandler,
			TLSConfig:         tlsConfig,
			ReadHeaderTimeout: 10 * time.Second,
		}
	}

	srvCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Fprintf(opts.Out, "GOBL Net listening on HTTP %s\n", httpLn.Addr()) //nolint:errcheck
	if httpsLn != nil {
		fmt.Fprintf(opts.Out, "GOBL Net listening on HTTPS %s\n", httpsLn.Addr()) //nolint:errcheck
	}

	errCh := make(chan error, 2)
	go func() {
		err := httpSrv.Serve(httpLn)
		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http: %w", err)
			cancel()
			return
		}
		errCh <- nil
	}()
	if httpsSrv != nil {
		go func() {
			err := httpsSrv.ServeTLS(httpsLn, "", "")
			if !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("https: %w", err)
				cancel()
				return
			}
			errCh <- nil
		}()
	}

	<-srvCtx.Done()
	fmt.Fprintln(opts.Out, "Shutting down...") //nolint:errcheck

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), netServeShutdownTimeout)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	if httpsSrv != nil {
		_ = httpsSrv.Shutdown(shutdownCtx)
	}

	expected := 1
	if httpsSrv != nil {
		expected = 2
	}
	var firstErr error
	for i := 0; i < expected; i++ {
		if err := <-errCh; err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func newAutocertManager(opts *NetServeOptions, domains []string) *autocert.Manager {
	certDir := opts.CertDir
	if certDir == "" {
		certDir = "certs"
	}
	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
		Email:      opts.ACMEEmail,
	}
	if opts.ACMETest {
		m.Client = &acme.Client{DirectoryURL: acmeStagingDirectoryURL}
	}
	return m
}

// listenTCP binds to the requested port on all interfaces. On EACCES it
// returns a wrapped error that guides the operator to a fix.
func listenTCP(port int) (stdnet.Listener, error) {
	addr := ":" + strconv.Itoa(port)
	ln, err := stdnet.Listen("tcp", addr)
	if err == nil {
		return ln, nil
	}
	if errors.Is(err, syscall.EACCES) {
		return nil, fmt.Errorf(
			"net serve: cannot bind %s — permission denied. "+
				"Use --http-port / --https-port to pick an unprivileged port, "+
				"grant the binary CAP_NET_BIND_SERVICE "+
				"(setcap 'cap_net_bind_service=+ep' <binary>), or run with sudo / "+
				"inside a container that maps the host port externally",
			addr,
		)
	}
	return nil, fmt.Errorf("net serve: listen %s: %w", addr, err)
}

func serveBytes(body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, _ = w.Write(body)
	}
}

// handleWho answers an authenticated party-exchange request. The caller
// POSTs a signed envelope (iss=gobl:caller, aud=gobl:self); the server
// verifies it, allow-lists the caller, and responds with its own party
// envelope signed with iss/aud reversed (iss=gobl:self, aud=gobl:caller).
func handleWho(client *net.Client, partyEnvBytes []byte, priv *dsig.PrivateKey, self cbc.URI, allow map[net.Address]bool, present bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, netInboxMaxBody))
		if err != nil {
			http.Error(w, "could not read body", http.StatusBadRequest)
			return
		}
		req := new(gobl.Envelope)
		if err := json.Unmarshal(body, req); err != nil {
			http.Error(w, "invalid envelope JSON", http.StatusBadRequest)
			return
		}
		caller, err := client.VerifyEnvelope(r.Context(), req, self)
		if err != nil {
			http.Error(w, "request verification failed: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if !allowed(allow, present, caller) {
			http.Error(w, "caller not accepted", http.StatusForbidden)
			return
		}

		resp := new(gobl.Envelope)
		if err := json.Unmarshal(partyEnvBytes, resp); err != nil {
			http.Error(w, "could not load party", http.StatusInternalServerError)
			return
		}
		if err := resp.Sign(priv, self, caller.URI()); err != nil {
			http.Error(w, "could not sign party: "+err.Error(), http.StatusInternalServerError)
			return
		}
		out, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "could not encode party", http.StatusInternalServerError)
			return
		}
		serveBytes(out)(w, r)
	}
}

func handleInbox(client *net.Client, dir string, self cbc.URI, allow map[net.Address]bool, present bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, netInboxMaxBody))
		if err != nil {
			http.Error(w, "could not read body", http.StatusBadRequest)
			return
		}

		env := new(gobl.Envelope)
		if err := json.Unmarshal(body, env); err != nil {
			http.Error(w, "invalid envelope JSON", http.StatusBadRequest)
			return
		}

		if err := env.Validate(); err != nil {
			http.Error(w, "envelope failed validation: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}

		// Audience binding is optional for deliveries: verify the signer
		// (iss) without enforcing aud, then reject only if a present aud
		// names a different recipient.
		sender, err := client.VerifyEnvelope(r.Context(), env, "")
		if err != nil {
			http.Error(w, "signature verification failed: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if p, perr := head.SignedPayload(env.Signatures[0]); perr == nil && self != "" && p.Aud != "" && p.Aud != self {
			http.Error(w, "envelope audience does not match this inbox", http.StatusUnauthorized)
			return
		}
		if !allowed(allow, present, sender) {
			http.Error(w, "sender not accepted", http.StatusForbidden)
			return
		}

		filename := filepath.Join(dir, env.Head.UUID.String()+".json")
		f, err := os.Create(filename)
		if err != nil {
			http.Error(w, "could not write inbox file", http.StatusInternalServerError)
			return
		}
		defer f.Close() //nolint:errcheck
		if _, err := f.Write(body); err != nil {
			http.Error(w, "could not write inbox file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
