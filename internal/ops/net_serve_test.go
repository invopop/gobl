package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mapFetcher struct {
	data map[string][]byte
}

func (m *mapFetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	body, ok := m.data[url]
	if !ok {
		return nil, net.ErrFetchFailed
	}
	return body, nil
}

// jwkBytes returns the single-JWK bytes served at the per-key endpoint
// for this key.
func jwkBytes(t *testing.T, key *dsig.PrivateKey) []byte {
	t.Helper()
	b, err := json.Marshal(key.Public())
	require.NoError(t, err)
	return b
}

const (
	testServeDomain = "me.example"
	testPeerDomain  = "peer.example"
)

var testPeerKey = dsig.NewES256Key()

// setupNetServer stands up a single-domain handler for testServeDomain
// (signed by the package privateKey) whose client can resolve both the
// served domain's and the peer's /keys. Returns the server and inbox dir.
func setupNetServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()):  jwkBytes(t, testPeerKey),
		net.Address(testServeDomain).KeyURL(privateKey.ID()): jwkBytes(t, privateKey),
	}}))

	h, err := buildDomainHandler(dc, client, io.Discard)
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv, dc.InboxDir
}

// signedRequest builds an envelope wrapping the peer's party, signed
// iss=peer, aud=<aud>.
func signedRequest(t *testing.T, aud net.Address) []byte {
	t.Helper()
	env, err := gobl.Envelop(&org.Party{Name: "Peer"})
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), aud.URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)
	return body
}

func TestNetServeKeys(t *testing.T) {
	srv, _ := setupNetServer(t)

	// Per-key endpoint: known kid returns the single JWK.
	resp, err := http.Get(srv.URL + net.KeyPath(privateKey.ID()))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	jwk := new(jose.JSONWebKey)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(jwk))
	assert.Equal(t, privateKey.ID(), jwk.KeyID)

	// Unknown kid returns 404 — no enumeration is exposed.
	resp404, err := http.Get(srv.URL + net.KeyPath("unknown-kid"))
	require.NoError(t, err)
	defer resp404.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusNotFound, resp404.StatusCode)

	// The bulk /keys endpoint no longer exists.
	respBulk, err := http.Get(srv.URL + net.KeysPath)
	require.NoError(t, err)
	defer respBulk.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusNotFound, respBulk.StatusCode)
}

func TestNetServeWhoExchange(t *testing.T) {
	srv, _ := setupNetServer(t)

	resp, err := http.Post(srv.URL+net.WhoPath, "application/json",
		bytes.NewReader(signedRequest(t, testServeDomain)))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	require.Equal(t, http.StatusOK, resp.StatusCode)

	env := new(gobl.Envelope)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(env))
	require.True(t, env.Signed())

	p, err := headSignedPayload(env)
	require.NoError(t, err)
	assert.Equal(t, net.Address(testServeDomain).URI(), p.Iss, "response signed by the served domain")
	assert.Equal(t, net.Address(testPeerDomain).URI(), p.Aud, "response bound to the caller")

	party, ok := env.Extract().(*org.Party)
	require.True(t, ok)
	assert.Equal(t, "Me", party.Name)
}

func TestNetServeWhoUnauthenticated(t *testing.T) {
	srv, _ := setupNetServer(t)
	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader([]byte("not json")))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestNetServeInboxAccepts(t *testing.T) {
	srv, inboxDir := setupNetServer(t)

	msg := &note.Message{Content: "hello inbox"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)

	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	files, err := os.ReadDir(inboxDir)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, env.Head.UUID.String()+".json", files[0].Name())
}

// callHandleWho drives the handleWho factory directly so we can craft
// corrupt internal state (bad partyEnvBytes, bad signing key) that the
// HTTP-level tests cannot reach via setupNetServer.
func callHandleWho(t *testing.T, partyEnvBytes []byte, priv *dsig.PrivateKey) *httptest.ResponseRecorder {
	t.Helper()
	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
	}}))
	self := net.Address(testServeDomain).URI()
	h := handleWho(client, partyEnvBytes, priv, self, nil, false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, net.WhoPath, bytes.NewReader(signedRequest(t, testServeDomain)))
	h(rec, req)
	return rec
}

func TestHandleWhoBadPartyBytes(t *testing.T) {
	rec := callHandleWho(t, []byte("not json"), privateKey)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "could not load party")
}

func TestHandleWhoSignFails(t *testing.T) {
	// Valid party bytes but a zero-value PrivateKey so resp.Sign errors.
	env, err := gobl.Envelop(&org.Party{Name: "Me"})
	require.NoError(t, err)
	partyBytes, err := json.Marshal(env)
	require.NoError(t, err)
	rec := callHandleWho(t, partyBytes, &dsig.PrivateKey{})
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "could not sign party")
}

func TestNetServeInboxValidationFails(t *testing.T) {
	srv, _ := setupNetServer(t)
	// Envelope JSON that parses but lacks required fields (digest, etc.)
	// so env.Validate fails with 422.
	resp, err := http.Post(srv.URL+net.InboxPath, "application/json",
		bytes.NewReader([]byte(`{"$schema":"https://gobl.org/draft-0/envelope","head":{"uuid":"01906c00-0000-7000-0000-000000000000","dig":{"alg":"sha256","val":"x"}},"doc":null}`)))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	// 422 (validation), or other 4xx — anything that isn't 202.
	assert.NotEqual(t, http.StatusAccepted, resp.StatusCode)
}

func TestNetServeInboxWriteFails(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("write-permission tests do not apply when running as root")
	}
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
	}}))
	h, err := buildDomainHandler(dc, client, io.Discard)
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	defer srv.Close()

	// Make the inbox directory read-only so os.Create fails.
	require.NoError(t, os.Chmod(dc.InboxDir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dc.InboxDir, 0o755) })

	msg := &note.Message{Content: "fail to write"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)

	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestNetServeInboxRejectsBadJSON(t *testing.T) {
	srv, _ := setupNetServer(t)
	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader([]byte("not json")))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestNetServeWhoUnauthorizedSignature(t *testing.T) {
	srv, _ := setupNetServer(t)
	// Build a request signed by an iss whose /keys the server can't resolve.
	other := dsig.NewES256Key()
	env, err := gobl.Envelop(&org.Party{Name: "Stranger"})
	require.NoError(t, err)
	require.NoError(t, env.Sign(other, net.Address("unknown.example").URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)

	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestNetServeWhoForbidden(t *testing.T) {
	// Allow list rejecting the peer triggers 403.
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})
	// Allow-list contains a different caller.
	require.NoError(t, os.WriteFile(dc.AllowFile, []byte(`["other.example"]`), 0o644))

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
	}}))
	h, err := buildDomainHandler(dc, client, io.Discard)
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader(signedRequest(t, testServeDomain)))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestNetServeWhoUnauthorized(t *testing.T) {
	srv, _ := setupNetServer(t)
	// Signed but with aud != self -> /who server rejects with 401.
	env, err := gobl.Envelop(&org.Party{Name: "Peer"})
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address("other.example").URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)
	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestNetServeInboxAudMismatch(t *testing.T) {
	srv, _ := setupNetServer(t)
	// /inbox tolerates missing aud, but rejects an aud naming a different recipient.
	msg := &note.Message{Content: "wrong aud"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address("other.example").URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)
	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestNetServeInboxForbidden(t *testing.T) {
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})
	require.NoError(t, os.WriteFile(dc.AllowFile, []byte(`["other.example"]`), 0o644))

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
	}}))
	h, err := buildDomainHandler(dc, client, io.Discard)
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	defer srv.Close()

	msg := &note.Message{Content: "rejected"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)
	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestNetServeInboxRejectsBadSignature(t *testing.T) {
	srv, inboxDir := setupNetServer(t)

	// Signed by a key whose /keys the server cannot resolve for the iss.
	other := dsig.NewES256Key()
	msg := &note.Message{Content: "bad sig"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(other, net.Address("unknown.example").URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)

	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	files, err := os.ReadDir(inboxDir)
	require.NoError(t, err)
	assert.Empty(t, files)
}

func headSignedPayload(env *gobl.Envelope) (*head.SigningPayload, error) {
	return head.SignedPayload(env.Signatures[0])
}

func readDirNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}
