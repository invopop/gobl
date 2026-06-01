package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// hostRewrite routes every request to base, regardless of the request's
// host, so a test can use a real domain identity while talking to an
// httptest server.
type hostRewrite struct{ base string }

func (h hostRewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	u, _ := url.Parse(h.base)
	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	return http.DefaultTransport.RoundTrip(req)
}

func TestNetWho(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, "acme.example")
	dc := domainConfigFor(configDir, "acme.example")

	// The served domain's client resolves the caller's per-key endpoint
	// to verify the incoming request.
	serverClient := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		net.Address(testPeerDomain).KeyURL(testPeerKey.ID()): jwkBytes(t, testPeerKey),
	}}))
	handler, err := buildDomainHandler(dc, serverClient, new(bytes.Buffer))
	require.NoError(t, err)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// Read the freshly-generated private key for acme.example so the
	// test fetcher can serve its public counterpart at /keys/<kid>.
	privBytes, err := os.ReadFile(dc.PrivateKeyFile)
	require.NoError(t, err)
	targetKey := new(dsig.PrivateKey)
	require.NoError(t, json.Unmarshal(privBytes, targetKey))

	party, err := NetWho(context.Background(), &NetWhoOptions{
		Target:    "acme.example",
		From:      net.Address(testPeerDomain),
		FromKey:   testPeerKey,
		FromParty: &org.Party{Name: "Peer"},
		Insecure:  true,
		// POSTs to http://acme.example/... but routed to the test server.
		Client: &http.Client{Transport: hostRewrite{base: srv.URL}},
		// Resolves the target's per-key endpoint (the served domain's
		// published key).
		Fetcher: &mapFetcher{data: map[string][]byte{
			"http://acme.example" + net.KeyPath(targetKey.ID()): jwkBytes(t, targetKey),
		}},
	})
	require.NoError(t, err)
	require.NotNil(t, party)
	assert.Equal(t, "acme.example", party.Name)
	require.Len(t, party.Endpoints, 1)
	assert.Equal(t, "gobl:acme.example", party.Endpoints[0].URI.String())
}

func TestNetWhoMissingFrom(t *testing.T) {
	_, err := NetWho(context.Background(), &NetWhoOptions{Target: "acme.example"})
	require.Error(t, err)
}

func TestNetWhoMissingTarget(t *testing.T) {
	_, err := NetWho(context.Background(), &NetWhoOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target address is required")
}

// staticHandler returns a fixed status+body — useful for error-path tests.
func staticHandler(status int, body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
}

func newWhoOpts(t *testing.T, target string, srvURL string) *NetWhoOptions {
	t.Helper()
	return &NetWhoOptions{
		Target:    net.Address(target),
		From:      net.Address(testPeerDomain),
		FromKey:   testPeerKey,
		FromParty: &org.Party{Name: "Peer"},
		Insecure:  true,
		Client:    &http.Client{Transport: hostRewrite{base: srvURL}},
		Fetcher:   &mapFetcher{data: map[string][]byte{}},
	}
}

func TestNetWhoNon200(t *testing.T) {
	srv := httptest.NewServer(staticHandler(http.StatusForbidden, "no"))
	defer srv.Close()
	_, err := NetWho(context.Background(), newWhoOpts(t, "acme.example", srv.URL))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 403")
}

func TestNetWhoInvalidResponseJSON(t *testing.T) {
	srv := httptest.NewServer(staticHandler(http.StatusOK, "not json"))
	defer srv.Close()
	_, err := NetWho(context.Background(), newWhoOpts(t, "acme.example", srv.URL))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid /who response")
}

func TestNetWhoUnsignedResponse(t *testing.T) {
	srv := httptest.NewServer(staticHandler(http.StatusOK, `{"doc":{}}`))
	defer srv.Close()
	_, err := NetWho(context.Background(), newWhoOpts(t, "acme.example", srv.URL))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not signed")
}

// TestNetWhoResponseWrongIssuer: response is signed but by the peer
// (i.e., not by the target). Verification loop finds no matching iss.
func TestNetWhoResponseWrongIssuer(t *testing.T) {
	// Build a signed envelope where iss/aud are reversed from what NetWho
	// expects to find on a /who response.
	env, err := gobl.Envelop(&org.Party{Name: "Wrong"})
	require.NoError(t, err)
	// iss = peer (caller) — but NetWho expects iss=target.
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	_, err = NetWho(context.Background(), newWhoOpts(t, testServeDomain, srv.URL))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not signed by")
}

// TestNetWhoTransportError exercises the default http.Client + default
// Fetcher branches plus the client.Do error path. Uses port 1 which is
// closed on a non-root host.
func TestNetWhoTransportError(t *testing.T) {
	_, err := NetWho(context.Background(), &NetWhoOptions{
		Target:    net.Address("127.0.0.1:1"),
		From:      net.Address(testPeerDomain),
		FromKey:   testPeerKey,
		FromParty: &org.Party{Name: "Peer"},
		Insecure:  true,
		// Client + Fetcher omitted to exercise the default branches.
	})
	require.Error(t, err)
}

// TestNetWhoResponseAudMismatch: response is correctly signed by the
// target but the aud names someone else.
func TestNetWhoResponseAudMismatch(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, testServeDomain)
	dc := domainConfigFor(configDir, testServeDomain)
	privBytes, err := os.ReadFile(dc.PrivateKeyFile)
	require.NoError(t, err)
	targetKey := new(dsig.PrivateKey)
	require.NoError(t, json.Unmarshal(privBytes, targetKey))

	env, err := gobl.Envelop(&org.Party{Name: "X"})
	require.NoError(t, err)
	require.NoError(t, env.Sign(targetKey,
		net.Address(testServeDomain).URI(),
		net.Address("other.example").URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	opts := newWhoOpts(t, testServeDomain, srv.URL)
	opts.Fetcher = &mapFetcher{data: map[string][]byte{
		"http://" + testServeDomain + net.KeyPath(targetKey.ID()): jwkBytes(t, targetKey),
	}}
	_, err = NetWho(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "audience mismatch")
}

// TestNetWhoResponseDocNotParty: response is correctly signed by the
// target, but the document is not an org.Party.
func TestNetWhoResponseDocNotParty(t *testing.T) {
	configDir := t.TempDir()
	initTestDomain(t, configDir, testServeDomain)
	dc := domainConfigFor(configDir, testServeDomain)
	privBytes, err := os.ReadFile(dc.PrivateKeyFile)
	require.NoError(t, err)
	targetKey := new(dsig.PrivateKey)
	require.NoError(t, json.Unmarshal(privBytes, targetKey))

	// Wrap a non-party document.
	wrap, err := gobl.Envelop(&org.Endpoint{URI: "gobl:x.example"})
	require.NoError(t, err)
	require.NoError(t, wrap.Sign(targetKey, net.Address(testServeDomain).URI(), net.Address(testPeerDomain).URI()))
	body, err := json.Marshal(wrap)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	opts := newWhoOpts(t, testServeDomain, srv.URL)
	opts.Fetcher = &mapFetcher{data: map[string][]byte{
		"http://" + testServeDomain + net.KeyPath(targetKey.ID()): jwkBytes(t, targetKey),
	}}
	_, err = NetWho(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not an org.Party")
}

// TestSchemeRewriteFetcher confirms scheme/host rewriting plus error
// passthrough for malformed input.
func TestSchemeRewriteFetcher(t *testing.T) {
	t.Run("rewrites scheme + host", func(t *testing.T) {
		var seen string
		inner := stubFetcher(func(_ context.Context, u string) ([]byte, error) {
			seen = u
			return []byte("ok"), nil
		})
		f := &schemeRewriteFetcher{base: "http://localhost:1234", inner: inner}
		body, err := f.Fetch(context.Background(), "https://acme.example/.well-known/gobl/keys/abc")
		require.NoError(t, err)
		assert.Equal(t, "ok", string(body))
		assert.Equal(t, "http://localhost:1234/.well-known/gobl/keys/abc", seen)
	})

	t.Run("invalid raw URL falls through unchanged", func(t *testing.T) {
		var seen string
		inner := stubFetcher(func(_ context.Context, u string) ([]byte, error) {
			seen = u
			return []byte("x"), nil
		})
		f := &schemeRewriteFetcher{base: "http://localhost:1234", inner: inner}
		_, err := f.Fetch(context.Background(), "://broken")
		require.NoError(t, err)
		assert.Equal(t, "://broken", seen)
	})
}

type stubFetcher func(context.Context, string) ([]byte, error)

func (s stubFetcher) Fetch(ctx context.Context, u string) ([]byte, error) { return s(ctx, u) }
