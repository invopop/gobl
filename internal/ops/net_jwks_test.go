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
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// jwksKeysFromBody decodes the JWKS body to a slice of jose JWKs for
// easy field-level assertions.
func jwksKeysFromBody(t *testing.T, body []byte) []jose.JSONWebKey {
	t.Helper()
	out := struct {
		Keys []json.RawMessage `json:"keys"`
	}{}
	require.NoError(t, json.Unmarshal(body, &out))
	keys := make([]jose.JSONWebKey, 0, len(out.Keys))
	for _, raw := range out.Keys {
		var k jose.JSONWebKey
		require.NoError(t, json.Unmarshal(raw, &k))
		keys = append(keys, k)
	}
	return keys
}

func TestJWKSEndpointSingleKey(t *testing.T) {
	srv, _, _ := setupServerWithLog(t)

	resp, err := http.Get(srv.URL + net.JWKSPath)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	keys := jwksKeysFromBody(t, body)
	require.Len(t, keys, 1)
	assert.Equal(t, privateKey.ID(), keys[0].KeyID)
}

func TestJWKSEndpointNewestFirst(t *testing.T) {
	// Build a domain with two keys: one with an older valid_from and
	// one with a newer one. The JWKS response must put the newer
	// first.
	cfg := t.TempDir()
	dc := domainConfigFor(cfg, testServeDomain)
	require.NoError(t, os.MkdirAll(filepath.Join(cfg, testServeDomain), 0o700))
	writePrivate(t, dc.PrivateKeyFile, privateKey)
	require.NoError(t, os.MkdirAll(dc.KeysDir, 0o755))

	older := cal.TimestampOf(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	newer := cal.TimestampOf(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	olderKey := writeKeyWithValidFrom(t, dc.KeysDir, dsig.NewES256Key(), &older)
	newerKey := writeKeyWithValidFrom(t, dc.KeysDir, privateKey, &newer)
	writeRawParty(t, dc.PartyFile, &org.Party{Name: "Me"})

	client := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{}}))
	h, err := buildDomainHandler(dc, client, slog.New(slog.NewTextHandler(new(bytes.Buffer), nil)))
	require.NoError(t, err)
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL + net.JWKSPath)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	keys := jwksKeysFromBody(t, body)
	require.Len(t, keys, 2)
	assert.Equal(t, newerKey, keys[0].KeyID, "newest key (valid_from 2026) sorts first")
	assert.Equal(t, olderKey, keys[1].KeyID, "older key (valid_from 2024) sorts second")
}

func TestJWKSAccessLog(t *testing.T) {
	srv, buf, _ := setupServerWithLog(t)
	resp, err := http.Get(srv.URL + net.JWKSPath)
	require.NoError(t, err)
	_ = resp.Body.Close()
	out := buf.String()
	assert.Contains(t, out, "jwks.served")
	assert.Contains(t, out, "count=1")
	assert.Contains(t, out, "http_request")
	assert.Contains(t, out, "status=200")
}

func TestBuildJWKSSortFallback(t *testing.T) {
	// Two keys with no valid_from: fall back to kid descending. The
	// helper takes a kid-keyed map of raw JWK bytes — using
	// minimally-shaped EC JWKs keeps the test focused on ordering.
	keys := map[string][]byte{
		"aaa": jwkRawNoWindow(t, "aaa"),
		"bbb": jwkRawNoWindow(t, "bbb"),
		"ccc": jwkRawNoWindow(t, "ccc"),
	}
	body, count, err := buildJWKS(keys)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
	parsed := jwksKeysFromBody(t, body)
	ids := []string{parsed[0].KeyID, parsed[1].KeyID, parsed[2].KeyID}
	want := []string{"ccc", "bbb", "aaa"}
	assert.Equal(t, want, ids, "kid descending when no valid_from")
}

func TestBuildJWKSMixedValidFrom(t *testing.T) {
	// One key has valid_from, the other doesn't. The one with a
	// timestamp sorts before the one without, regardless of kid.
	ts := cal.TimestampOf(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	keys := map[string][]byte{
		"zzz-no-ts":   jwkRawNoWindow(t, "zzz-no-ts"),
		"aaa-with-ts": jwkRawWithValidFrom(t, "aaa-with-ts", &ts),
	}
	body, _, err := buildJWKS(keys)
	require.NoError(t, err)
	parsed := jwksKeysFromBody(t, body)
	assert.Equal(t, "aaa-with-ts", parsed[0].KeyID)
	assert.Equal(t, "zzz-no-ts", parsed[1].KeyID)
}

func TestBuildJWKSEmpty(t *testing.T) {
	body, count, err := buildJWKS(map[string][]byte{})
	require.NoError(t, err)
	assert.Equal(t, 0, count)
	assert.JSONEq(t, `{"keys":[]}`, string(body))
}

func TestBuildJWKSBadJSON(t *testing.T) {
	_, _, err := buildJWKS(map[string][]byte{
		"broken": []byte("not json"),
	})
	require.Error(t, err)
}

// --- helpers ---

// writeKeyWithValidFrom writes a single PublishedKey-shaped JWK with a
// custom valid_from to <dir>/<kid>.json and returns the kid. Used to
// pin chronological ordering across two keys.
func writeKeyWithValidFrom(t *testing.T, dir string, priv *dsig.PrivateKey, validFrom *cal.Timestamp) string {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	raw, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	pk := new(dsig.PublicKey)
	require.NoError(t, json.Unmarshal(raw, pk))
	pk.ValidFrom = validFrom
	out, err := json.Marshal(pk)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, priv.ID()+".json"), out, 0o644))
	return priv.ID()
}

func jwkRawNoWindow(t *testing.T, kid string) []byte {
	t.Helper()
	priv := dsig.NewES256Key()
	raw, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	// Rewrite the kid in the marshaled bytes via a parse+remarshal pass.
	m := map[string]any{}
	require.NoError(t, json.Unmarshal(raw, &m))
	m["kid"] = kid
	out, err := json.Marshal(m)
	require.NoError(t, err)
	return out
}

func jwkRawWithValidFrom(t *testing.T, kid string, ts *cal.Timestamp) []byte {
	t.Helper()
	priv := dsig.NewES256Key()
	raw, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	m := map[string]any{}
	require.NoError(t, json.Unmarshal(raw, &m))
	m["kid"] = kid
	m["valid_from"] = ts
	out, err := json.Marshal(m)
	require.NoError(t, err)
	return out
}

// Confirm equal-timestamp entries fall back deterministically to
// kid-descending order.
func TestBuildJWKSEqualValidFromTiebreak(t *testing.T) {
	ts := cal.TimestampOf(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	keys := map[string][]byte{
		"k-1": jwkRawWithValidFrom(t, "k-1", &ts),
		"k-2": jwkRawWithValidFrom(t, "k-2", &ts),
		"k-3": jwkRawWithValidFrom(t, "k-3", &ts),
	}
	body, _, err := buildJWKS(keys)
	require.NoError(t, err)
	parsed := jwksKeysFromBody(t, body)
	ids := []string{parsed[0].KeyID, parsed[1].KeyID, parsed[2].KeyID}
	assert.Equal(t, []string{"k-3", "k-2", "k-1"}, ids)
}
