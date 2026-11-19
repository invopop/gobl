package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (m *mockFetcher) Fetch(_ context.Context, url string, _ http.Header) ([]byte, error) {
	m.url = url
	return m.data, m.err
}

func (m *mockFetcher) Post(_ context.Context, _ string, _ []byte) error {
	return ErrFetchFailed
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

// countingFetcher wraps a Fetcher and counts fetches per URL.
type countingFetcher struct {
	inner Fetcher
	n     map[string]int
}

func (f *countingFetcher) Fetch(ctx context.Context, url string, header http.Header) ([]byte, error) {
	f.n[url]++
	return f.inner.Fetch(ctx, url, header)
}

func (f *countingFetcher) Post(ctx context.Context, url string, body []byte) error {
	return f.inner.Post(ctx, url, body)
}

func TestFetchKeyCache(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	addr := Address("billing.invopop.com")
	url := addr.KeyURL(key.ID())

	pubData, err := json.Marshal(key.Public())
	require.NoError(t, err)

	newCounting := func() *countingFetcher {
		return &countingFetcher{
			inner: &mapFetcher{data: map[string][]byte{url: pubData}},
			n:     map[string]int{},
		}
	}

	t.Run("second fetch is served from cache", func(t *testing.T) {
		f := newCounting()
		c := NewClient(WithFetcher(f))
		for range 3 {
			pk, err := c.FetchKey(ctx, addr, key.ID())
			require.NoError(t, err)
			assert.Equal(t, key.ID(), pk.ID())
		}
		assert.Equal(t, 1, f.n[url])
	})

	t.Run("zero TTL disables the cache", func(t *testing.T) {
		f := newCounting()
		c := NewClient(WithFetcher(f), WithKeyCacheTTL(0))
		_, err := c.FetchKey(ctx, addr, key.ID())
		require.NoError(t, err)
		_, err = c.FetchKey(ctx, addr, key.ID())
		require.NoError(t, err)
		assert.Equal(t, 2, f.n[url])
	})

	t.Run("entries expire after the TTL", func(t *testing.T) {
		f := newCounting()
		c := NewClient(WithFetcher(f), WithKeyCacheTTL(10*time.Millisecond))
		_, err := c.FetchKey(ctx, addr, key.ID())
		require.NoError(t, err)
		time.Sleep(20 * time.Millisecond)
		_, err = c.FetchKey(ctx, addr, key.ID())
		require.NoError(t, err)
		assert.Equal(t, 2, f.n[url])
	})

	t.Run("failed fetches are not cached", func(t *testing.T) {
		f := &countingFetcher{inner: &mapFetcher{}, n: map[string]int{}}
		c := NewClient(WithFetcher(f))
		_, err := c.FetchKey(ctx, addr, key.ID())
		require.Error(t, err)
		_, err = c.FetchKey(ctx, addr, key.ID())
		require.Error(t, err)
		assert.Equal(t, 2, f.n[url])
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

		body, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL+"/x", nil)
		require.NoError(t, err)
		assert.Equal(t, `{"ok":true}`, string(body))
	})

	t.Run("204 returns ErrNoContent", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		_, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL, nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoContent))
	})

	t.Run("non-200", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "nope", http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL, nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
		assert.Contains(t, err.Error(), "HTTP 404")
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := NewHTTPFetcher().Fetch(context.Background(), "://broken", nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("nil context", func(t *testing.T) {
		// NewRequestWithContext panics on a nil context, so the func
		// must error rather than crash.
		var ctx context.Context //nolint:staticcheck
		_, err := NewHTTPFetcher().Fetch(ctx, "http://example.invalid", nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("transport error is retryable", func(t *testing.T) {
		// Unreachable address (port 1 is privileged + usually closed).
		f := &HTTPFetcher{Client: &http.Client{Timeout: 100 * time.Millisecond}}
		_, err := f.Fetch(context.Background(), "http://127.0.0.1:1/x", nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
	})

	t.Run("5xx is retryable", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "boom", http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL, nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
		assert.False(t, errors.Is(err, ErrFetchFailed))
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

		body, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL, nil)
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

	_, err := NewHTTPFetcher().Fetch(context.Background(), srv.URL+"/x", nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrFetchFailed))
	assert.Contains(t, err.Error(), "refusing to dial non-public address")
}

func TestSafeDialContext(t *testing.T) {
	ctx := context.Background()

	t.Run("address without a port", func(t *testing.T) {
		_, err := safeDialContext(ctx, "tcp", "example.com")
		require.Error(t, err)
	})

	t.Run("unresolvable host is retryable", func(t *testing.T) {
		// .invalid is reserved (RFC 2606) and never resolves; DNS
		// failures are transient conditions.
		_, err := safeDialContext(ctx, "tcp", "does-not-exist.invalid:443")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
	})

	t.Run("non-public address", func(t *testing.T) {
		_, err := safeDialContext(ctx, "tcp", "127.0.0.1:443")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
		assert.Contains(t, err.Error(), "refusing to dial non-public address")
	})
}

func TestDialAddrs(t *testing.T) {
	ctx := context.Background()
	ln, err := stdnet.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close() //nolint:errcheck
	_, port, err := stdnet.SplitHostPort(ln.Addr().String())
	require.NoError(t, err)
	loopback := stdnet.ParseIP("127.0.0.1")

	t.Run("connects to the first address accepting", func(t *testing.T) {
		closed, err := stdnet.Listen("tcp", "127.0.0.2:0")
		if err != nil {
			t.Skip("127.0.0.2 not available")
		}
		require.NoError(t, closed.Close())
		conn, err := dialAddrs(ctx, "tcp", "example.com", []stdnet.IP{stdnet.ParseIP("127.0.0.2"), loopback}, port)
		require.NoError(t, err)
		require.NoError(t, conn.Close())
	})

	t.Run("wraps the last dial error as unavailable", func(t *testing.T) {
		cctx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := dialAddrs(cctx, "tcp", "example.com", []stdnet.IP{loopback}, port)
		require.ErrorIs(t, err, ErrUnavailable)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("no addresses", func(t *testing.T) {
		_, err := dialAddrs(ctx, "tcp", "example.com", nil, port)
		require.ErrorIs(t, err, ErrUnavailable)
		assert.ErrorContains(t, err, "no addresses for example.com")
	})
}

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"8.8.8.8", true},
		{"1.1.1.1", true},
		{"2606:4700:4700::1111", true},
		{"127.0.0.1", false},           // loopback
		{"::1", false},                 // loopback
		{"10.0.0.1", false},            // private
		{"192.168.1.1", false},         // private
		{"172.16.0.1", false},          // private
		{"169.254.169.254", false},     // link-local (AWS metadata)
		{"fe80::1", false},             // link-local
		{"0.0.0.0", false},             // unspecified
		{"::", false},                  // unspecified
		{"224.0.0.1", false},           // multicast
		{"100.64.0.1", false},          // shared (RFC 6598)
		{"100.127.255.254", false},     // shared (RFC 6598)
		{"100.128.0.1", true},          // just past shared
		{"0.1.2.3", false},             // this network
		{"192.0.0.8", false},           // IETF protocol assignments
		{"198.18.0.1", false},          // benchmarking
		{"240.0.0.1", false},           // reserved
		{"255.255.255.255", false},     // broadcast
		{"fd00::1", false},             // unique local (RFC 4193)
		{"64:ff9b::a00:1", false},      // NAT64 to 10.0.0.1
		{"64:ff9b:1::1", false},        // local-use NAT64
		{"::ffff:10.0.0.1", false},     // IPv4-mapped private
		{"192.0.2.1", false},           // TEST-NET-1
		{"198.51.100.1", false},        // TEST-NET-2
		{"203.0.113.1", false},         // TEST-NET-3
		{"192.88.99.1", false},         // 6to4 relay anycast
		{"2001:db8::1", false},         // documentation
		{"3fff::1", false},             // documentation
		{"2001:0:4136:e378::1", false}, // Teredo
		{"2002:a00:1::1", false},       // 6to4 for 10.0.0.1
		{"100::1", false},              // discard-only
		{"5f00::1", false},             // segment routing
		{"2a00:1450:4001::1", true},    // global unicast
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

func TestHTTPFetcherFetch202(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	_, err := newHTTPFetcher(true).Fetch(context.Background(), srv.URL, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrPending))
}

func TestHTTPFetcherPostNilContext(t *testing.T) {
	// NewRequestWithContext panics on a nil context, so the func must
	// error rather than crash.
	var ctx context.Context //nolint:staticcheck
	err := NewHTTPFetcher().Post(ctx, "http://example.invalid", []byte(`{}`))
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrFetchFailed))
}

func TestAuthHeaderMintFailure(t *testing.T) {
	// An identity whose address cannot be canonicalized fails token
	// minting, which must surface from the calling operation.
	key := dsig.NewES256Key()
	c := NewClient(
		WithFetcher(&mapFetcher{}),
		WithIdentity("not valid!", key),
	)
	_, err := c.Who(context.Background(), "receiver.example.com")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrAddressInvalid))

}

func TestKeyCacheEviction(t *testing.T) {
	key := dsig.NewES256Key()
	pk := key.Public()

	t.Run("purges expired entries when full", func(t *testing.T) {
		c := NewClient(WithFetcher(&mapFetcher{}))
		past := time.Now().Add(-time.Minute)
		for i := range maxKeyCacheEntries {
			c.keyCache[fmt.Sprintf("https://k%d.example/key", i)] = keyCacheEntry{key: pk, exp: past}
		}
		c.storeKey("https://fresh.example/key", pk)
		assert.LessOrEqual(t, len(c.keyCache), maxKeyCacheEntries)
		assert.NotNil(t, c.cachedKey("https://fresh.example/key"))
	})

	t.Run("evicts arbitrarily when full of live entries", func(t *testing.T) {
		c := NewClient(WithFetcher(&mapFetcher{}))
		future := time.Now().Add(time.Hour)
		for i := range maxKeyCacheEntries {
			c.keyCache[fmt.Sprintf("https://k%d.example/key", i)] = keyCacheEntry{key: pk, exp: future}
		}
		c.storeKey("https://fresh.example/key", pk)
		assert.LessOrEqual(t, len(c.keyCache), maxKeyCacheEntries)
		assert.NotNil(t, c.cachedKey("https://fresh.example/key"))
	})
}

func TestFlushKeyCache(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	addr := Address("billing.invopop.com")
	url := addr.KeyURL(key.ID())

	pubData, err := json.Marshal(key.Public())
	require.NoError(t, err)
	f := &countingFetcher{
		inner: &mapFetcher{data: map[string][]byte{url: pubData}},
		n:     map[string]int{},
	}
	c := NewClient(WithFetcher(f))

	_, err = c.FetchKey(ctx, addr, key.ID())
	require.NoError(t, err)
	c.FlushKeyCache()
	_, err = c.FetchKey(ctx, addr, key.ID())
	require.NoError(t, err)
	assert.Equal(t, 2, f.n[url], "flush forces a refetch")
}

func TestHTTPFetcherRefusesRedirects(t *testing.T) {
	var followed bool
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		followed = true
	}))
	defer target.Close()
	for _, code := range []int{http.StatusFound, http.StatusTemporaryRedirect, http.StatusPermanentRedirect} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, target.URL+r.URL.Path, code)
		}))
		f := newHTTPFetcher(true)

		_, err := f.Fetch(context.Background(), srv.URL+"/.well-known/gobl/keys/k", nil)
		require.ErrorIs(t, err, ErrFetchFailed)
		assert.ErrorContains(t, err, "redirects are not followed")

		err = f.Post(context.Background(), srv.URL+"/.well-known/gobl/inbox", []byte("{}"))
		require.ErrorIs(t, err, ErrFetchFailed)
		assert.ErrorContains(t, err, "redirects are not followed")
		srv.Close()
	}
	assert.False(t, followed)
}

func TestHTTPFetcherWrapsTransportErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	f := newHTTPFetcher(true)

	_, err := f.Fetch(ctx, srv.URL+"/.well-known/gobl/keys/k", nil)
	require.ErrorIs(t, err, ErrUnavailable)
	assert.ErrorIs(t, err, context.Canceled)

	err = f.Post(ctx, srv.URL+"/.well-known/gobl/inbox", []byte("{}"))
	require.ErrorIs(t, err, ErrUnavailable)
	assert.ErrorIs(t, err, context.Canceled)
}
