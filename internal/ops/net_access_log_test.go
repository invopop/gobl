package ops

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
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

// setupServerWithLog stands up the test domain handler chain with a
// captured logger so individual test cases can assert on log lines.
func setupServerWithLog(t *testing.T) (*httptest.Server, *bytes.Buffer, string) {
	t.Helper()
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
		net.Address(testServeDomain).KeyURL(privateKey.ID()): jwkBytes(t, privateKey),
	}}))

	buf := new(bytes.Buffer)
	log := slog.New(slog.NewTextHandler(buf, nil))
	h, err := buildDomainHandler(dc, client, log)
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv, buf, dc.InboxDir
}

func TestAccessLogKeysLookup(t *testing.T) {
	srv, buf, _ := setupServerWithLog(t)

	// Known kid -> 200 + keys.lookup found=true.
	resp, err := http.Get(srv.URL + net.KeyPath(privateKey.ID()))
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	out := buf.String()
	assert.Contains(t, out, "keys.lookup")
	assert.Contains(t, out, "kid="+privateKey.ID())
	assert.Contains(t, out, "found=true")
	assert.Contains(t, out, "http_request")
	assert.Contains(t, out, "status=200")

	buf.Reset()
	// Unknown kid -> 404 + keys.lookup found=false.
	resp404, err := http.Get(srv.URL + net.KeyPath("ghost"))
	require.NoError(t, err)
	_ = resp404.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp404.StatusCode)
	out = buf.String()
	assert.Contains(t, out, "keys.lookup")
	assert.Contains(t, out, "found=false")
	assert.Contains(t, out, "status=404")
}

func TestAccessLogWhoRejectsBadBody(t *testing.T) {
	srv, buf, _ := setupServerWithLog(t)
	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", strings.NewReader("not json"))
	require.NoError(t, err)
	_ = resp.Body.Close()
	out := buf.String()
	assert.Contains(t, out, "who.rejected")
	assert.Contains(t, out, "reason=bad_body")
	assert.Contains(t, out, "status=400")
}

func TestAccessLogWhoRejectsVerifyFailed(t *testing.T) {
	srv, buf, _ := setupServerWithLog(t)
	// Signed by an iss we don't have keys for.
	other := dsig.NewES256Key()
	env, err := gobl.Envelop(&org.Party{Name: "Stranger"})
	require.NoError(t, err)
	require.NoError(t, env.Sign(other, net.Address("unknown.example").URI(), net.Address(testServeDomain).URI()))
	body, _ := json.Marshal(env)
	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	_ = resp.Body.Close()
	out := buf.String()
	assert.Contains(t, out, "who.rejected")
	assert.Contains(t, out, "reason=verify_failed")
	assert.Contains(t, out, "status=401")
}

func TestAccessLogWhoNotAllowed(t *testing.T) {
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writeKey(t, dc.KeysDir, privateKey)
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})
	// allow-list excludes the peer.
	require.NoError(t, os.WriteFile(dc.AllowFile, []byte(`["other.example"]`), 0o644))

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
	}}))
	buf := new(bytes.Buffer)
	log := slog.New(slog.NewTextHandler(buf, nil))
	h, err := buildDomainHandler(dc, client, log)
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader(signedRequest(t, testServeDomain)))
	require.NoError(t, err)
	_ = resp.Body.Close()
	out := buf.String()
	assert.Contains(t, out, "who.rejected")
	assert.Contains(t, out, "reason=not_allowed")
	assert.Contains(t, out, "status=403")
}

func TestAccessLogWhoExchangeSuccess(t *testing.T) {
	srv, buf, _ := setupServerWithLog(t)
	resp, err := http.Post(srv.URL+net.WhoPath, "application/json", bytes.NewReader(signedRequest(t, testServeDomain)))
	require.NoError(t, err)
	_ = resp.Body.Close()
	out := buf.String()
	assert.Contains(t, out, "who.exchange")
	assert.Contains(t, out, "caller="+testPeerDomain)
	assert.Contains(t, out, "status=200")
}

func TestAccessLogInboxAccepted(t *testing.T) {
	srv, buf, inboxDir := setupServerWithLog(t)

	msg := &note.Message{Content: "logged"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address(testServeDomain).URI()))
	body, _ := json.Marshal(env)

	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	out := buf.String()
	assert.Contains(t, out, "inbox.accepted")
	assert.Contains(t, out, "envelope="+env.Head.UUID.String())
	assert.Contains(t, out, "status=202")

	// Sanity: the envelope was persisted.
	files, _ := os.ReadDir(inboxDir)
	require.Len(t, files, 1)
}

func TestAccessLogInboxAudMismatch(t *testing.T) {
	srv, buf, _ := setupServerWithLog(t)
	msg := &note.Message{Content: "wrong aud"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address("other.example").URI()))
	body, _ := json.Marshal(env)
	resp, err := http.Post(srv.URL+net.InboxPath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	_ = resp.Body.Close()
	out := buf.String()
	assert.Contains(t, out, "inbox.rejected")
	assert.Contains(t, out, "reason=aud_mismatch")
	assert.Contains(t, out, "status=401")
}

func TestStatusRecorderImplicit200(t *testing.T) {
	// When the inner handler writes a body without an explicit
	// WriteHeader, the recorder treats it as 200.
	rec := &statusRecorder{ResponseWriter: httptest.NewRecorder()}
	_, _ = rec.Write([]byte("hello"))
	assert.Equal(t, http.StatusOK, rec.status)
}

func TestStatusRecorderRespectsFirst(t *testing.T) {
	rec := &statusRecorder{ResponseWriter: httptest.NewRecorder()}
	rec.WriteHeader(http.StatusCreated)
	rec.WriteHeader(http.StatusBadRequest) // second call must not overwrite
	assert.Equal(t, http.StatusCreated, rec.status)
}

func TestAccessLogMiddlewareDirect(t *testing.T) {
	// Drive the middleware directly to confirm field shape.
	buf := new(bytes.Buffer)
	log := slog.New(slog.NewTextHandler(buf, nil))
	h := accessLog(log, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Host = "acme.example:8080"
	h.ServeHTTP(rec, req)

	out := buf.String()
	assert.Contains(t, out, "http_request")
	assert.Contains(t, out, "method=GET")
	assert.Contains(t, out, "path=/x")
	assert.Contains(t, out, "host=acme.example")
	assert.Contains(t, out, "remote=10.0.0.1:1234")
	assert.Contains(t, out, "status=200")
	assert.Contains(t, out, "duration_ms=")
}
