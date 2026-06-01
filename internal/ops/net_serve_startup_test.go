package ops

import (
	"bytes"
	"context"
	"encoding/json"
	stdnet "net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeRawParty writes a raw (unsigned) org.Party to path.
func writeRawParty(t *testing.T, path string, party *org.Party) {
	t.Helper()
	data, err := json.Marshal(party)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, data, 0o644))
}

// writeKey writes a single public JWK to <dir>/<kid>.json, creating
// dir if needed.
func writeKey(t *testing.T, dir string, key *dsig.PrivateKey) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	b, err := json.Marshal(key.Public())
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, key.ID()+".json"), b, 0o644))
}

func writePrivate(t *testing.T, path string, key *dsig.PrivateKey) {
	t.Helper()
	b, err := json.Marshal(key)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, b, 0o600))
}

func dcFor(dir, domain string) domainConfig {
	return domainConfig{
		Domain:         domain,
		KeysDir:        filepath.Join(dir, "keys"),
		PrivateKeyFile: filepath.Join(dir, "private.jwk"),
		PartyFile:      filepath.Join(dir, "party.json"),
		InboxDir:       filepath.Join(dir, "inbox"),
		AllowFile:      filepath.Join(dir, "allow.json"),
	}
}

func TestEnsureKeysAutoGenerates(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "")

	before := time.Now().UTC()
	out := new(bytes.Buffer)
	keysByKID, err := ensureKeys(dc, out)
	require.NoError(t, err)
	after := time.Now().UTC()

	assert.DirExists(t, dc.KeysDir)
	assert.FileExists(t, dc.PrivateKeyFile)

	logged := out.String()
	assert.Contains(t, logged, "Generated new keypair")
	assert.Contains(t, logged, dc.PrivateKeyFile)

	privBytes, err := os.ReadFile(dc.PrivateKeyFile)
	require.NoError(t, err)
	priv := new(dsig.PrivateKey)
	require.NoError(t, json.Unmarshal(privBytes, priv))

	// The per-kid map served by the handler must include the freshly
	// generated key.
	require.Contains(t, keysByKID, priv.ID())

	// The published key is at <KeysDir>/<kid>.json on disk.
	keyFile := filepath.Join(dc.KeysDir, priv.ID()+".json")
	onDisk, err := os.ReadFile(keyFile)
	require.NoError(t, err)
	pk := new(dsig.PublicKey)
	require.NoError(t, json.Unmarshal(onDisk, pk))
	require.Equal(t, priv.ID(), pk.ID())

	// Freshly generated keys are stamped with valid_from = now and no
	// valid_until.
	require.NotNil(t, pk.ValidFrom, "valid_from must be stamped on a freshly generated key")
	assert.True(t,
		!pk.ValidFrom.Time.Before(before.Add(-time.Second)) &&
			!pk.ValidFrom.Time.After(after.Add(time.Second)),
		"valid_from %v should fall within [%v, %v]", pk.ValidFrom.Time, before, after,
	)
	assert.Nil(t, pk.ValidUntil)

	info, err := os.Stat(dc.PrivateKeyFile)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestEnsureKeysRejectsPartialState(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, dsig.NewES256Key()) // only the public side

	_, err := ensureKeys(dc, new(bytes.Buffer))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "inconsistent key setup")
}

func TestEnsureKeysRejectsMismatchedKid(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, dsig.NewES256Key())
	writePrivate(t, dc.PrivateKeyFile, dsig.NewES256Key())

	_, err := ensureKeys(dc, new(bytes.Buffer))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not published under")
}

func TestEnsureKeysRejectsMismatchedFilename(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "")
	// Write the public JWK under the wrong filename.
	require.NoError(t, os.MkdirAll(dc.KeysDir, 0o755))
	k := dsig.NewES256Key()
	b, err := json.Marshal(k.Public())
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dc.KeysDir, "wrong-name.json"), b, 0o644))
	writePrivate(t, dc.PrivateKeyFile, k)

	_, err = ensureKeys(dc, new(bytes.Buffer))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match JWK kid")
}

func TestNetServeHandlerPartyMissing(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)

	_, err := NetServeHandler(&NetServeOptions{
		PartyFile:      dc.PartyFile,
		KeysDir:        dc.KeysDir,
		PrivateKeyFile: dc.PrivateKeyFile,
		InboxDir:       dc.InboxDir,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "party file not found")
	assert.Contains(t, err.Error(), "gobl init")
}

func TestReadPartyEnvelopeRaw(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "d.example.com")
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Acme"})

	env, err := readPartyEnvelope(dc)
	require.NoError(t, err)
	require.False(t, env.Signed(), "party is returned unsigned; /who signs per request")
	party, ok := env.Extract().(*org.Party)
	require.True(t, ok)
	assert.Equal(t, "Acme", party.Name)
}

func TestReadPartyEnvelopeMissing(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "d.example.com")
	_, err := readPartyEnvelope(dc)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "party file not found")
}

func TestReadPartyEnvelopeReadError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("read-permission tests do not apply when running as root")
	}
	dir := t.TempDir()
	dc := dcFor(dir, "d.example.com")
	require.NoError(t, os.WriteFile(dc.PartyFile, []byte("{}"), 0o000))
	t.Cleanup(func() { _ = os.Chmod(dc.PartyFile, 0o644) })
	_, err := readPartyEnvelope(dc)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read party file")
}

func TestReadPartyEnvelopeSignedEnvelopeRoundTrip(t *testing.T) {
	// A pre-signed envelope on disk passes through unchanged.
	dir := t.TempDir()
	dc := dcFor(dir, "d.example.com")
	env, err := gobl.Envelop(&org.Party{Name: "Pre-signed"})
	require.NoError(t, err)
	require.NoError(t, env.Sign(privateKey,
		net.Address("d.example.com").URI(),
		net.Address("other.example").URI()))
	data, err := json.Marshal(env)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(dc.PartyFile, data, 0o644))

	got, err := readPartyEnvelope(dc)
	require.NoError(t, err)
	party, ok := got.Extract().(*org.Party)
	require.True(t, ok)
	assert.Equal(t, "Pre-signed", party.Name)
}

func TestReadPartyEnvelopeInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "d.example.com")
	require.NoError(t, os.WriteFile(dc.PartyFile, []byte("not json"), 0o644))
	_, err := readPartyEnvelope(dc)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestLoadAllowList(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "x.example")

	t.Run("absent file: present=false", func(t *testing.T) {
		set, present, err := loadAllowList(dc)
		require.NoError(t, err)
		assert.False(t, present)
		assert.Nil(t, set)
	})

	t.Run("empty AllowFile path: present=false", func(t *testing.T) {
		bare := domainConfig{} // AllowFile == ""
		set, present, err := loadAllowList(bare)
		require.NoError(t, err)
		assert.False(t, present)
		assert.Nil(t, set)
	})

	t.Run("valid list", func(t *testing.T) {
		require.NoError(t, os.WriteFile(dc.AllowFile, []byte(`["a.example","b.example"]`), 0o644))
		set, present, err := loadAllowList(dc)
		require.NoError(t, err)
		assert.True(t, present)
		assert.True(t, set["a.example"])
		assert.True(t, set["b.example"])
		assert.False(t, set["c.example"])
	})

	t.Run("invalid JSON", func(t *testing.T) {
		require.NoError(t, os.WriteFile(dc.AllowFile, []byte("not json"), 0o644))
		_, _, err := loadAllowList(dc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid allow list")
	})
}

func TestAllowed(t *testing.T) {
	t.Run("no list present: any caller accepted", func(t *testing.T) {
		assert.True(t, allowed(nil, false, "any.example"))
	})
	t.Run("list present: only listed caller accepted", func(t *testing.T) {
		set := map[net.Address]bool{"a.example": true}
		assert.True(t, allowed(set, true, "a.example"))
		assert.False(t, allowed(set, true, "b.example"))
	})
}

func TestDiscoverDomainsMissingConfigDir(t *testing.T) {
	// Non-existent config dir returns nil slice, no error.
	dcs, err := discoverDomains(filepath.Join(t.TempDir(), "does-not-exist"))
	require.NoError(t, err)
	assert.Empty(t, dcs)
}

func TestDomainNames(t *testing.T) {
	got := domainNames([]domainConfig{
		{Domain: "a.example"},
		{Domain: ""},
		{Domain: "b.example"},
	})
	assert.Equal(t, []string{"a.example", "b.example"}, got)
}

func TestStripPort(t *testing.T) {
	assert.Equal(t, "x.example", stripPort("x.example:8080"))
	assert.Equal(t, "x.example", stripPort("x.example"))
}

func TestFileExistsAndDirExists(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "f.txt")
	require.NoError(t, os.WriteFile(f, []byte("x"), 0o644))
	assert.True(t, fileExists(f))
	assert.False(t, fileExists(dir), "fileExists returns false for directories")
	assert.True(t, dirExists(dir))
	assert.False(t, dirExists(f), "dirExists returns false for files")
}

func TestLoadPrivateKeyFileErrors(t *testing.T) {
	dir := t.TempDir()
	t.Run("missing", func(t *testing.T) {
		_, err := loadPrivateKeyFile(filepath.Join(dir, "nope.jwk"))
		require.Error(t, err)
	})
	t.Run("bad JSON", func(t *testing.T) {
		p := filepath.Join(dir, "bad.jwk")
		require.NoError(t, os.WriteFile(p, []byte("not json"), 0o600))
		_, err := loadPrivateKeyFile(p)
		require.Error(t, err)
	})
}

func TestResolveDomains(t *testing.T) {
	t.Run("manual mode via KeysDir", func(t *testing.T) {
		dcs, err := resolveDomains(&NetServeOptions{KeysDir: "/keys", Domain: "x"})
		require.NoError(t, err)
		require.Len(t, dcs, 1)
		assert.Equal(t, "x", dcs[0].Domain)
		assert.Equal(t, "/keys", dcs[0].KeysDir)
	})
	t.Run("manual mode via PartyFile", func(t *testing.T) {
		dcs, err := resolveDomains(&NetServeOptions{PartyFile: "/p"})
		require.NoError(t, err)
		require.Len(t, dcs, 1)
	})
	t.Run("no config dir", func(t *testing.T) {
		_, err := resolveDomains(&NetServeOptions{})
		require.Error(t, err)
	})
	t.Run("discovered domains", func(t *testing.T) {
		configDir := t.TempDir()
		initTestDomain(t, configDir, "a.example")
		initTestDomain(t, configDir, "b.example")
		dcs, err := resolveDomains(&NetServeOptions{ConfigDir: configDir})
		require.NoError(t, err)
		assert.Len(t, dcs, 2)
	})
	t.Run("named domain found", func(t *testing.T) {
		configDir := t.TempDir()
		initTestDomain(t, configDir, "a.example")
		dcs, err := resolveDomains(&NetServeOptions{ConfigDir: configDir, Domain: "a.example"})
		require.NoError(t, err)
		require.Len(t, dcs, 1)
		assert.Equal(t, "a.example", dcs[0].Domain)
	})
	t.Run("named domain bootstrapped", func(t *testing.T) {
		configDir := t.TempDir()
		// No domains on disk — resolveDomains returns the constructed config.
		dcs, err := resolveDomains(&NetServeOptions{ConfigDir: configDir, Domain: "fresh.example"})
		require.NoError(t, err)
		require.Len(t, dcs, 1)
		assert.Equal(t, "fresh.example", dcs[0].Domain)
	})
}

func TestEnsureKeysOnlyKeysDirExists(t *testing.T) {
	// keys dir present but no private.jwk -> inconsistent setup with
	// keys as `present` (covers the other arm of the default branch).
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, dsig.NewES256Key())
	_, err := ensureKeys(dc, new(bytes.Buffer))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "inconsistent key setup")
}

func TestEnsureKeysBadPrivateKey(t *testing.T) {
	// keys dir + private.jwk both exist but private.jwk is unparseable.
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, dsig.NewES256Key())
	require.NoError(t, os.WriteFile(dc.PrivateKeyFile, []byte("not json"), 0o600))
	_, err := ensureKeys(dc, new(bytes.Buffer))
	require.Error(t, err)
}

func TestEnsureKeysEmptyKeysDir(t *testing.T) {
	dir := t.TempDir()
	dc := dcFor(dir, "")
	// keys dir exists but is empty (no <kid>.json files).
	require.NoError(t, os.MkdirAll(dc.KeysDir, 0o755))
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	_, err := ensureKeys(dc, new(bytes.Buffer))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contains no JWKs")
}

func TestReadKeysDirInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "bad.json"), []byte("not json"), 0o644))
	_, err := readKeysDir(dir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JWK")
}

func TestReadKeysDirIgnoresNonJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	// Non-JSON file: ignored.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello"), 0o644))
	// Subdirectory: ignored.
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "subdir"), 0o755))
	got, err := readKeysDir(dir)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestNetServeHandlerDefaultsClient(t *testing.T) {
	// Omitting Out + Client routes to defaults without panic.
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "X"})

	h, err := NetServeHandler(&NetServeOptions{
		PartyFile:      dc.PartyFile,
		KeysDir:        dc.KeysDir,
		PrivateKeyFile: dc.PrivateKeyFile,
		InboxDir:       dc.InboxDir,
	})
	require.NoError(t, err)
	require.NotNil(t, h)
}

func TestBuildRouterSingleUnnamed(t *testing.T) {
	// One unnamed identity: router shortcircuits to a single handler.
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Solo"})

	h, err := buildRouter([]domainConfig{dc}, nil, new(bytes.Buffer))
	require.NoError(t, err)
	require.NotNil(t, h)
}

func TestNetServeRunCancel(t *testing.T) {
	// End-to-end smoke of NetServe via a temp config dir. Picks a free
	// unprivileged port (race-safe enough for tests), then cancels the
	// context to drive the graceful-shutdown branch.
	configDir := t.TempDir()
	initTestDomain(t, configDir, "x.example")

	port := freePort(t)
	ctx, cancel := context.WithCancel(context.Background())
	out := new(bytes.Buffer)
	doneCh := make(chan error, 1)
	go func() {
		doneCh <- NetServe(ctx, &NetServeOptions{
			ConfigDir: configDir,
			HTTPPort:  port,
			Out:       out,
		})
	}()
	// Let the goroutine reach Serve before shutting it down.
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case err := <-doneCh:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("NetServe did not return after cancel")
	}
	// Confirm it logged the listening address.
	assert.Contains(t, out.String(), "GOBL Net listening on HTTP")
	assert.Contains(t, out.String(), "Shutting down")
}

// freePort returns a TCP port that's currently free on 127.0.0.1.
// There's a small race window between close and reuse, but it is good
// enough for short-lived test binds.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := stdnet.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*stdnet.TCPAddr).Port
	_ = ln.Close()
	return port
}

func TestNetServeNoDomains(t *testing.T) {
	configDir := t.TempDir()
	// Empty config dir → discoverDomains returns nothing → error.
	err := NetServe(context.Background(), &NetServeOptions{ConfigDir: configDir})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no domains configured")
}

func TestNetServeWithACMETest(t *testing.T) {
	// ACMETest with a named domain: drives the ACME case in NetServe
	// (sets up the autocert manager). We use HTTPPort=freePort to
	// avoid privilege issues, and cancel quickly to shut down.
	configDir := t.TempDir()
	initTestDomain(t, configDir, "x.example")

	ctx, cancel := context.WithCancel(context.Background())
	out := new(bytes.Buffer)
	doneCh := make(chan error, 1)
	go func() {
		doneCh <- NetServe(ctx, &NetServeOptions{
			ConfigDir: configDir,
			ACMETest:  true,
			Domain:    "x.example",
			HTTPPort:  freePort(t),
			HTTPSPort: freePort(t),
			CertDir:   filepath.Join(configDir, "certs"),
			Out:       out,
		})
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		t.Fatal("NetServe did not return after cancel")
	}
	assert.Contains(t, out.String(), "ACME enabled for domains")
}

func TestNetServeACMEManualMode(t *testing.T) {
	// Manual mode (no Domain) + ACME requires named domains -> error.
	dir := t.TempDir()
	dc := dcFor(dir, "")
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Solo"})

	err := NetServe(context.Background(), &NetServeOptions{
		KeysDir:        dc.KeysDir,
		PrivateKeyFile: dc.PrivateKeyFile,
		PartyFile:      dc.PartyFile,
		InboxDir:       dc.InboxDir,
		ACMETest:       true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ACME requires named domains")
}

func TestNetServeCertFileMissing(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "x.example")
	err := NetServe(context.Background(), &NetServeOptions{
		ConfigDir: configDir,
		CertFile:  "/no/such/cert.pem",
		KeyFile:   "/no/such/key.pem",
		HTTPPort:  freePort(t),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "load TLS keypair")
}

func TestNetServeBuildRouterError(t *testing.T) {
	// Domain exists on disk but has an inconsistent key state: only
	// the keys/ dir is populated, no private.jwk. buildRouter fails.
	configDir := t.TempDir()
	domain := "broken.example"
	dc := domainConfigFor(configDir, domain)
	require.NoError(t, os.MkdirAll(filepath.Join(configDir, domain), 0o755))
	writeKey(t, dc.KeysDir, dsig.NewES256Key())
	err := NetServe(context.Background(), &NetServeOptions{ConfigDir: configDir})
	require.Error(t, err)
}

func TestGenerateKeypairWriteErrors(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("write-permission tests do not apply when running as root")
	}
	t.Run("private dir not writable", func(t *testing.T) {
		dir := t.TempDir()
		// Make the parent dir read-only so MkdirAll(parent, ...) for
		// the private key file fails.
		ro := filepath.Join(dir, "ro")
		require.NoError(t, os.MkdirAll(ro, 0o500))
		t.Cleanup(func() { _ = os.Chmod(ro, 0o755) })

		_, err := generateKeypair(filepath.Join(dir, "keys"), filepath.Join(ro, "sub", "private.jwk"), new(bytes.Buffer))
		require.Error(t, err)
	})
	t.Run("keys dir not writable", func(t *testing.T) {
		dir := t.TempDir()
		ro := filepath.Join(dir, "ro")
		require.NoError(t, os.MkdirAll(ro, 0o500))
		t.Cleanup(func() { _ = os.Chmod(ro, 0o755) })

		// Private file path is fine; keysDir path is inside a non-writable parent.
		_, err := generateKeypair(filepath.Join(ro, "keys"), filepath.Join(dir, "private.jwk"), new(bytes.Buffer))
		require.Error(t, err)
	})
}

func TestReadKeysDirNonExistent(t *testing.T) {
	_, err := readKeysDir(filepath.Join(t.TempDir(), "missing"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read keys dir")
}

func TestLoadAllowListReadError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("read-permission tests do not apply when running as root")
	}
	dir := t.TempDir()
	dc := dcFor(dir, "")
	require.NoError(t, os.WriteFile(dc.AllowFile, []byte("[]"), 0o000))
	t.Cleanup(func() { _ = os.Chmod(dc.AllowFile, 0o644) })
	_, _, err := loadAllowList(dc)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read allow list")
}

func TestDiscoverDomainsReadError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("read-permission tests do not apply when running as root")
	}
	dir := t.TempDir()
	require.NoError(t, os.Chmod(dir, 0o000))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	_, err := discoverDomains(dir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read config dir")
}

func TestReadKeysDirReadError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("read-permission tests do not apply when running as root")
	}
	dir := t.TempDir()
	keysDir := filepath.Join(dir, "keys")
	require.NoError(t, os.MkdirAll(keysDir, 0o755))
	// A file inside that we cannot read.
	bad := filepath.Join(keysDir, "abc.json")
	require.NoError(t, os.WriteFile(bad, []byte(`{}`), 0o000))
	t.Cleanup(func() { _ = os.Chmod(bad, 0o644) })

	_, err := readKeysDir(keysDir)
	require.Error(t, err)
}

func TestBuildRouterPropagatesDomainError(t *testing.T) {
	// Domain with no keys/private key and no party -> ensureKeys fails.
	dir := t.TempDir()
	dc := dcFor(dir, "broken.example")
	// Write keys/ dir without private.jwk to trigger the "inconsistent" path.
	writeKey(t, dc.KeysDir, dsig.NewES256Key())
	_, err := buildRouter([]domainConfig{dc}, nil, new(bytes.Buffer))
	require.Error(t, err)
}

// TestBuildDomainHandlerErrors covers the buildDomainHandler error
// branches that come after ensureKeys: bad private key, missing party,
// malformed allow-list, and inbox-mkdir-fail.
func TestBuildDomainHandlerErrors(t *testing.T) {
	t.Run("bad private key", func(t *testing.T) {
		dir := t.TempDir()
		dc := dcFor(dir, "")
		writeKey(t, dc.KeysDir, privateKey)
		require.NoError(t, os.WriteFile(dc.PrivateKeyFile, []byte("not json"), 0o600))
		_, err := buildDomainHandler(dc, nil, new(bytes.Buffer))
		require.Error(t, err)
	})

	t.Run("missing party", func(t *testing.T) {
		dir := t.TempDir()
		dc := dcFor(dir, "")
		writeKey(t, dc.KeysDir, privateKey)
		writePrivate(t, dc.PrivateKeyFile, privateKey)
		_, err := buildDomainHandler(dc, nil, new(bytes.Buffer))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "party file not found")
	})

	t.Run("bad allow list", func(t *testing.T) {
		dir := t.TempDir()
		dc := dcFor(dir, "")
		writeKey(t, dc.KeysDir, privateKey)
		writePrivate(t, dc.PrivateKeyFile, privateKey)
		writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})
		require.NoError(t, os.WriteFile(dc.AllowFile, []byte("not json"), 0o644))
		_, err := buildDomainHandler(dc, nil, new(bytes.Buffer))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid allow list")
	})

	t.Run("inbox is a file", func(t *testing.T) {
		dir := t.TempDir()
		dc := dcFor(dir, "")
		writeKey(t, dc.KeysDir, privateKey)
		writePrivate(t, dc.PrivateKeyFile, privateKey)
		writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})
		// Pre-create dc.InboxDir as a regular file so MkdirAll fails.
		require.NoError(t, os.WriteFile(dc.InboxDir, []byte("x"), 0o644))
		_, err := buildDomainHandler(dc, nil, new(bytes.Buffer))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create inbox dir")
	})
}
