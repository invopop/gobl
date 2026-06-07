package net

import (
	"context"
	"encoding/json"
	"errors"
	stdnet "net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFetcher struct {
	data []byte
	err  error
	url  string // records the URL that was fetched
}

func (m *mockFetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	m.url = url
	return m.data, m.err
}

func TestFetchPublicKey(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	pub := key.Public()

	pubData, err := json.Marshal(pub)
	require.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		mock := &mockFetcher{data: pubData}
		c := NewClient(WithFetcher(mock))

		pk, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), key.ID())
		require.NoError(t, err)
		assert.Equal(t, key.ID(), pk.ID())
		assert.Equal(t,
			"https://billing.invopop.com/.well-known/gobl/keys/"+key.ID(),
			mock.url, "client hits the per-key endpoint",
		)
	})

	t.Run("kid mismatch", func(t *testing.T) {
		// Fetcher returns a JWK whose kid does not match the requested kid.
		mock := &mockFetcher{data: pubData}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), "other-kid")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mock := &mockFetcher{data: []byte("not json")}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), key.ID())
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("fetch error", func(t *testing.T) {
		mock := &mockFetcher{err: ErrFetchFailed}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), key.ID())
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("invalid address", func(t *testing.T) {
		mock := &mockFetcher{data: pubData}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchPublicKey(ctx, Address(""), key.ID())
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressEmpty))
	})

	t.Run("empty kid", func(t *testing.T) {
		mock := &mockFetcher{data: pubData}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})
}

func TestNewHTTPFetcher(t *testing.T) {
	f := NewHTTPFetcher()
	require.NotNil(t, f)
	require.NotNil(t, f.Client)
	assert.Equal(t, defaultTimeout, f.Client.Timeout)
}

func TestHTTPFetcherFetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		defer srv.Close()

		body, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL+"/x")
		require.NoError(t, err)
		assert.Equal(t, `{"ok":true}`, string(body))
	})

	t.Run("non-200", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "nope", http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
		assert.Contains(t, err.Error(), "HTTP 404")
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := NewHTTPFetcher().Fetch(context.Background(), "://broken")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("nil context", func(t *testing.T) {
		// NewRequestWithContext panics on a nil context, so the func
		// must error rather than crash.
		var ctx context.Context //nolint:staticcheck
		_, err := NewHTTPFetcher().Fetch(ctx, "http://example.invalid")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("transport error", func(t *testing.T) {
		// Unreachable address (port 1 is privileged + usually closed).
		f := &HTTPFetcher{Client: &http.Client{Timeout: 100 * time.Millisecond}}
		_, err := f.Fetch(context.Background(), "http://127.0.0.1:1/x")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("body too large is truncated", func(t *testing.T) {
		// Send a body larger than maxBodySize. The fetcher uses
		// io.LimitReader, so it returns exactly maxBodySize bytes and
		// no error.
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Stream more than 1 MiB.
			w.Header().Set("Content-Type", "application/json")
			big := strings.Repeat("a", maxBodySize+1024)
			_, _ = w.Write([]byte(big))
		}))
		defer srv.Close()

		body, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL)
		require.NoError(t, err)
		assert.Equal(t, maxBodySize, len(body))
	})
}

func TestHTTPFetcherRejectsNonPublicAddresses(t *testing.T) {
	// The default HTTPFetcher must refuse to dial loopback / private /
	// link-local / unspecified addresses to prevent SSRF via a signed
	// `iss` URI.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	_, err := NewHTTPFetcher().Fetch(context.Background(), srv.URL+"/x")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrFetchFailed))
	assert.Contains(t, err.Error(), "refusing to dial non-public address")
}

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"8.8.8.8", true},
		{"1.1.1.1", true},
		{"2606:4700:4700::1111", true},
		{"127.0.0.1", false},       // loopback
		{"::1", false},             // loopback
		{"10.0.0.1", false},        // private
		{"192.168.1.1", false},     // private
		{"172.16.0.1", false},      // private
		{"169.254.169.254", false}, // link-local (AWS metadata)
		{"fe80::1", false},         // link-local
		{"0.0.0.0", false},         // unspecified
		{"::", false},              // unspecified
		{"224.0.0.1", false},       // multicast
	}
	for _, tc := range tests {
		t.Run(tc.ip, func(t *testing.T) {
			ip := stdnet.ParseIP(tc.ip)
			require.NotNil(t, ip, "parse")
			assert.Equal(t, tc.want, isPublicIP(ip))
		})
	}

	t.Run("nil IP", func(t *testing.T) {
		assert.False(t, isPublicIP(nil))
	})
}
