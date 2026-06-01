package ops

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// initTestDomain scaffolds a domain under configDir via InitDomain.
func initTestDomain(t *testing.T, configDir, domain string) {
	t.Helper()
	require.NoError(t, InitDomain(&InitOptions{
		ConfigDir: configDir,
		Domain:    domain,
		Name:      domain,
		Out:       new(bytes.Buffer),
	}))
}

func TestInitDomain(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "billing.invopop.com")

	dir := filepath.Join(configDir, "billing.invopop.com")
	assert.DirExists(t, filepath.Join(dir, "keys"))
	assert.FileExists(t, filepath.Join(dir, "private.jwk"))
	assert.FileExists(t, filepath.Join(dir, "party.json"))
	assert.DirExists(t, filepath.Join(dir, "inbox"))

	// Exactly one published key file matching the private key's kid.
	entries, err := os.ReadDir(filepath.Join(dir, "keys"))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.True(t, strings.HasSuffix(entries[0].Name(), ".json"))

	info, err := os.Stat(filepath.Join(dir, "private.jwk"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	pb, err := os.ReadFile(filepath.Join(dir, "party.json"))
	require.NoError(t, err)
	party := new(org.Party)
	require.NoError(t, json.Unmarshal(pb, party))
	require.Len(t, party.Endpoints, 1)
	assert.Equal(t, "gobl:billing.invopop.com", party.Endpoints[0].URI.String())

	err = InitDomain(&InitOptions{ConfigDir: configDir, Domain: "billing.invopop.com", Out: new(bytes.Buffer)})
	require.Error(t, err)
}

func TestInitDomainMissingDomain(t *testing.T) {
	err := InitDomain(&InitOptions{ConfigDir: t.TempDir()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "domain is required")
}

func TestInitDomainExistingNotEmpty(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "a.example")
	// Re-init without --force fails.
	err := InitDomain(&InitOptions{
		ConfigDir: configDir,
		Domain:    "a.example",
		Out:       new(bytes.Buffer),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestInitDomainPartyWriteFails(t *testing.T) {
	// Pre-stage <configDir>/<domain>/party.json as a directory so
	// os.WriteFile fails after generateKeypair succeeds.
	configDir := t.TempDir()
	dir := filepath.Join(configDir, "x.example")
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "party.json"), 0o755))

	err := InitDomain(&InitOptions{
		ConfigDir: configDir,
		Domain:    "x.example",
		Force:     true,
		Out:       new(bytes.Buffer),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "write party")
}

func TestInitDomainInboxIsFile(t *testing.T) {
	// Pre-stage <configDir>/<domain>/inbox as a regular file so
	// os.MkdirAll(InboxDir) fails after the key + party writes succeed.
	configDir := t.TempDir()
	dir := filepath.Join(configDir, "x.example")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "inbox"), []byte("file-not-dir"), 0o644))

	err := InitDomain(&InitOptions{
		ConfigDir: configDir,
		Domain:    "x.example",
		Force:     true,
		Out:       new(bytes.Buffer),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create inbox dir")
}

func TestInitDomainMkdirError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("write-permission tests do not apply when running as root")
	}
	// Place ConfigDir inside an unwritable parent so MkdirAll fails.
	parent := t.TempDir()
	ro := filepath.Join(parent, "ro")
	require.NoError(t, os.MkdirAll(ro, 0o500))
	t.Cleanup(func() { _ = os.Chmod(ro, 0o755) })

	err := InitDomain(&InitOptions{
		ConfigDir: filepath.Join(ro, "sub"),
		Domain:    "x.example",
		Out:       new(bytes.Buffer),
	})
	require.Error(t, err)
}

func TestInitDomainForceOverwrite(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "a.example")
	// --force allows re-init over a populated directory.
	err := InitDomain(&InitOptions{
		ConfigDir: configDir,
		Domain:    "a.example",
		Force:     true,
		Out:       new(bytes.Buffer),
	})
	require.NoError(t, err, "Force re-init should succeed; new key adds alongside the old one")
}

func TestInitDomainDefaultsStdout(t *testing.T) {
	// Smoke: omitting Out routes to os.Stdout without panicking. We
	// can't intercept os.Stdout here cleanly, so just verify the call
	// returns success.
	configDir := t.TempDir()
	err := InitDomain(&InitOptions{ConfigDir: configDir, Domain: "default-out.example"})
	require.NoError(t, err)
}

func TestDiscoverDomainsSkipsCerts(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "a.example")
	initTestDomain(t, configDir, "b.example")
	require.NoError(t, os.MkdirAll(filepath.Join(configDir, "certs"), 0o755))

	domains, err := discoverDomains(configDir)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"a.example", "b.example"}, domainNames(domains))
}

func TestMultiDomainRouter(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "a.example")
	initTestDomain(t, configDir, "b.example")
	domains, err := discoverDomains(configDir)
	require.NoError(t, err)

	peerKey := dsig.NewES256Key()
	const peer = "peer.example"
	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(peer).KeyURL(peerKey.ID()): jwkBytes(t, peerKey),
	}}))

	router, err := buildRouter(domains, client, discardLog())
	require.NoError(t, err)
	srv := httptest.NewServer(router)
	defer srv.Close()

	// POST /who on each host returns a party signed by that host, bound to peer.
	for _, host := range []string{"a.example", "b.example"} {
		reqEnv, err := gobl.Envelop(&org.Party{Name: "Peer"})
		require.NoError(t, err)
		require.NoError(t, reqEnv.Sign(peerKey, net.Address(peer).URI(), net.Address(host).URI()))
		body, err := json.Marshal(reqEnv)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, srv.URL+net.WhoPath, bytes.NewReader(body))
		req.Host = host
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		env := new(gobl.Envelope)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(env))
		_ = resp.Body.Close()
		p, err := headSignedPayload(env)
		require.NoError(t, err)
		assert.Equal(t, net.Address(host).URI(), p.Iss)
		assert.Equal(t, net.Address(peer).URI(), p.Aud)
	}

	// Unknown host -> 404.
	req, _ := http.NewRequest(http.MethodPost, srv.URL+net.WhoPath, bytes.NewReader(signedRequest(t, "zzz.example")))
	req.Host = "zzz.example"
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	_ = resp.Body.Close()

	// Inbox POST with Host a.example lands in a.example/inbox.
	msg := &note.Message{Content: "routed"}
	msg.SetUUID(uuid.V7())
	denv, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, denv.Sign(peerKey, net.Address(peer).URI(), net.Address("a.example").URI()))
	body, err := json.Marshal(denv)
	require.NoError(t, err)

	ireq, _ := http.NewRequest(http.MethodPost, srv.URL+net.InboxPath, bytes.NewReader(body))
	ireq.Host = "a.example"
	ireq.Header.Set("Content-Type", "application/json")
	iresp, err := http.DefaultClient.Do(ireq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, iresp.StatusCode)
	_ = iresp.Body.Close()

	aFiles, err := os.ReadDir(filepath.Join(configDir, "a.example", "inbox"))
	require.NoError(t, err)
	assert.Len(t, aFiles, 1)
	bFiles, err := os.ReadDir(filepath.Join(configDir, "b.example", "inbox"))
	require.NoError(t, err)
	assert.Empty(t, bFiles)
}
