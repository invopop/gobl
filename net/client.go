package net

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	stdnet "net"
	"net/http"
	"sync"
	"time"

	"github.com/invopop/gobl/dsig"
)

const (
	defaultTimeout = 10 * time.Second
	maxBodySize    = 1 << 20 // 1MB
)

// Key cache limits. Published keys are immutable per kid, so fetched
// keys are safe to reuse for a short TTL; the size cap stops hostile
// iss/kid values growing the cache without bound.
const (
	defaultKeyCacheTTL = 5 * time.Minute
	maxKeyCacheEntries = 1024
)

// dialTimeout is the per-attempt timeout for the SSRF-safe dialer.
const dialTimeout = 5 * time.Second

// safeDialContext is the DialContext used by the default HTTPFetcher.
// It resolves the target host and refuses to connect when any of the
// resolved IPs is a loopback, private, link-local, or unspecified
// address — the standard SSRF defense for a client that dials hosts
// derived from signed payloads (a `gobl:` `iss` URI). Tests and local
// development should inject a custom Fetcher rather than relaxing
// this default.
func safeDialContext(ctx context.Context, network, addr string) (stdnet.Conn, error) {
	host, port, err := stdnet.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ips, err := stdnet.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	for _, ip := range ips {
		if !isPublicIP(ip) {
			return nil, fmt.Errorf("%w: refusing to dial non-public address %s (%s)", ErrFetchFailed, host, ip)
		}
	}
	d := &stdnet.Dialer{Timeout: dialTimeout}
	return d.DialContext(ctx, network, stdnet.JoinHostPort(host, port))
}

// isPublicIP reports whether ip is a routable, non-special address.
// A loopback, private (RFC 1918 / RFC 6598), link-local, multicast,
// unspecified, or interface-local-multicast IP is rejected.
func isPublicIP(ip stdnet.IP) bool {
	if ip == nil {
		return false
	}
	switch {
	case ip.IsLoopback(),
		ip.IsPrivate(),
		ip.IsLinkLocalUnicast(),
		ip.IsLinkLocalMulticast(),
		ip.IsInterfaceLocalMulticast(),
		ip.IsUnspecified(),
		ip.IsMulticast():
		return false
	}
	return true
}

// Fetcher is the transport used by a Client: Fetch retrieves
// well-known resources and Post delivers envelopes to an inbox, as
// used by Client.Send. The header carries any Authorization request
// token the Client has minted for the request; implementations must
// forward it. Implementations that never deliver (verification-only
// transports, test fixtures) may implement Post as a plain error
// return.
type Fetcher interface {
	Fetch(ctx context.Context, url string, header http.Header) ([]byte, error)
	Post(ctx context.Context, url string, body []byte, header http.Header) error
}

// HTTPFetcher implements Fetcher using net/http.
type HTTPFetcher struct {
	Client *http.Client
}

// NewHTTPFetcher creates an HTTPFetcher with sensible defaults.
// The fetcher's transport rejects any dial whose resolved IP is
// loopback, private, link-local, multicast, or unspecified, to
// prevent SSRF attacks via a signed `iss` URI. There is no public
// escape hatch for the SSRF guard; in-process test fixtures should
// inject their own Fetcher via WithFetcher.
func NewHTTPFetcher() *HTTPFetcher {
	return newHTTPFetcher(false)
}

// newHTTPFetcher is the internal constructor. allowLoopback bypasses
// the SSRF guard so unit tests can talk to httptest servers bound to
// 127.0.0.1. Not exported.
func newHTTPFetcher(allowLoopback bool) *HTTPFetcher {
	transport := &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   dialTimeout,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if !allowLoopback {
		transport.DialContext = safeDialContext
	}
	return &HTTPFetcher{
		Client: &http.Client{
			Timeout:   defaultTimeout,
			Transport: transport,
		},
	}
}

// Fetch retrieves the body from the given URL, forwarding the given
// headers (e.g. an Authorization request token) on the request.
func (f *HTTPFetcher) Fetch(ctx context.Context, url string, header http.Header) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	copyHeader(req.Header, header)
	req.Header.Set("Accept", "application/json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", transportErr(err), err)
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("%w: %s", ErrNoContent, url)
	}
	if resp.StatusCode == http.StatusAccepted {
		return nil, fmt.Errorf("%w: %s", ErrPending, url)
	}
	if retryable(resp.StatusCode) {
		return nil, fmt.Errorf("%w: HTTP %d from %s", ErrUnavailable, resp.StatusCode, url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: HTTP %d from %s", ErrFetchFailed, resp.StatusCode, url)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	return body, nil
}

// retryable reports whether a response status signals a transient
// condition worth retrying: 429 or any 5xx.
func retryable(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500
}

// transportErr classifies a round-trip error: the SSRF guard's
// refusal to dial is permanent (ErrFetchFailed); anything else —
// timeouts, resets, DNS failures — is transient (ErrUnavailable).
func transportErr(err error) error {
	if errors.Is(err, ErrFetchFailed) {
		return ErrFetchFailed
	}
	return ErrUnavailable
}

// Post delivers a JSON body to the given URL, forwarding the given
// headers on the request. A 202 response is success. Any other 4xx —
// except 429 — reports ErrInboxRejected; 429, 5xx, and transport
// failures report the retryable ErrUnavailable; anything else
// reports ErrFetchFailed.
func (f *HTTPFetcher) Post(ctx context.Context, url string, body []byte, header http.Header) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	copyHeader(req.Header, header)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", transportErr(err), err)
	}
	defer resp.Body.Close() // nolint:errcheck

	switch {
	case resp.StatusCode == http.StatusAccepted:
		return nil
	case resp.StatusCode >= 400 && resp.StatusCode < 500 &&
		resp.StatusCode != http.StatusTooManyRequests:
		return fmt.Errorf("%w: HTTP %d from %s", ErrInboxRejected, resp.StatusCode, url)
	case retryable(resp.StatusCode):
		return fmt.Errorf("%w: HTTP %d from %s", ErrUnavailable, resp.StatusCode, url)
	default:
		return fmt.Errorf("%w: HTTP %d from %s", ErrFetchFailed, resp.StatusCode, url)
	}
}

// copyHeader merges src into dst without sharing the underlying
// slices.
func copyHeader(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

// Client provides GOBL Net operations including KeySet fetching
// and remote verification.
type Client struct {
	fetcher     Fetcher
	authorities []Address
	identity    *identity

	keyTTL   time.Duration
	keyMu    sync.Mutex
	keyCache map[string]keyCacheEntry
}

// keyCacheEntry pairs a fetched public key with its expiry.
type keyCacheEntry struct {
	key *dsig.PublicKey
	exp time.Time
}

// identity holds the address and signing key a Client uses to mint
// request tokens for outbound who and inbox requests.
type identity struct {
	addr Address
	key  *dsig.PrivateKey
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithFetcher sets a custom Fetcher implementation.
func WithFetcher(f Fetcher) ClientOption {
	return func(c *Client) {
		c.fetcher = f
	}
}

// WithAuthorities sets the client's trusted registration authorities,
// replacing the default list (net.Authorities, which carries the
// network's default authority). Operators who want to supplement the
// default rather than replace it include it explicitly, or grow the
// global default with RegisterAuthority. Replacement semantics let
// closed deployments exclude the default authority entirely.
func WithAuthorities(addrs ...Address) ClientOption {
	return func(c *Client) {
		c.authorities = append([]Address{}, addrs...)
	}
}

// WithKeyCacheTTL overrides how long fetched public keys are reused
// before being fetched again. Zero disables the cache; the default is
// five minutes.
func WithKeyCacheTTL(d time.Duration) ClientOption {
	return func(c *Client) {
		c.keyTTL = d
	}
}

// WithSandbox switches the client's trusted registration authorities
// to the sandbox list (net.SandboxAuthorities), replacing the live
// default. Sandbox authorities endorse test identities with relaxed
// verification, so the two environments never trust each other's
// endorsements. Like WithAuthorities this replaces the trust list —
// when both options are given, the last one wins.
func WithSandbox() ClientOption {
	return func(c *Client) {
		c.authorities = append([]Address{}, SandboxAuthorities...)
	}
}

// WithIdentity sets the client's own address and private key, used to
// mint the request tokens that who and inbox requests carry in their
// Authorization header. A client without an identity sends those
// requests bare and conforming servers reject them.
func WithIdentity(addr Address, key *dsig.PrivateKey) ClientOption {
	return func(c *Client) {
		c.identity = &identity{addr: addr, key: key}
	}
}

// authHeader returns the headers for a request to aud, containing a
// freshly minted bearer token when the client has an identity.
func (c *Client) authHeader(aud Address) (http.Header, error) {
	h := http.Header{}
	if c.identity == nil {
		return h, nil
	}
	token, err := NewToken(c.identity.key, c.identity.addr, aud, 0)
	if err != nil {
		return nil, err
	}
	h.Set("Authorization", "Bearer "+token)
	return h, nil
}

// NewClient creates a new GOBL Net client.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		fetcher:     NewHTTPFetcher(),
		authorities: append([]Address{}, Authorities...),
		keyTTL:      defaultKeyCacheTTL,
		keyCache:    map[string]keyCacheEntry{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// FetchKey retrieves a single public key (with its optional validity
// window) from the well-known per-key URL derived from the given
// address and kid. The response body is a JWK (RFC 7517) possibly
// augmented with the `valid_from` / `valid_until` extension members
// understood by dsig.PublicKey.
func (c *Client) FetchKey(ctx context.Context, addr Address, kid string) (*dsig.PublicKey, error) {
	// Canonicalize so the per-key URL uses the ASCII form regardless
	// of how the address was written.
	addr, err := ParseAddress(string(addr))
	if err != nil {
		return nil, err
	}
	if kid == "" {
		return nil, fmt.Errorf("%w: kid is required", ErrFetchFailed)
	}
	url := addr.KeyURL(kid)
	if pk := c.cachedKey(url); pk != nil {
		return pk, nil
	}
	// Key endpoints are open: no request token is attached.
	data, err := c.fetcher.Fetch(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	pk := new(dsig.PublicKey)
	if err := json.Unmarshal(data, pk); err != nil {
		return nil, fmt.Errorf("%w: invalid JWK response: %v", ErrFetchFailed, err)
	}
	if pk.ID() != kid {
		return nil, fmt.Errorf("%w: kid mismatch (got %q, want %q)", ErrFetchFailed, pk.ID(), kid)
	}
	c.storeKey(url, pk)
	return pk, nil
}

// FlushKeyCache empties the client's key cache so subsequent
// verifications refetch published keys — observing a key's removal
// from the issuer's endpoint immediately, e.g. after a reported
// compromise, instead of after the cache TTL.
func (c *Client) FlushKeyCache() {
	c.keyMu.Lock()
	defer c.keyMu.Unlock()
	clear(c.keyCache)
}

// cachedKey returns the unexpired cached key for the given URL, or nil.
func (c *Client) cachedKey(url string) *dsig.PublicKey {
	if c.keyTTL <= 0 {
		return nil
	}
	c.keyMu.Lock()
	defer c.keyMu.Unlock()
	e, ok := c.keyCache[url]
	if !ok {
		return nil
	}
	if time.Now().After(e.exp) {
		delete(c.keyCache, url)
		return nil
	}
	return e.key
}

// storeKey caches a verified-shape key for the client's TTL. When the
// cache is full, expired entries are purged first; if it is still full
// an arbitrary entry is evicted so the cache never exceeds its cap.
func (c *Client) storeKey(url string, pk *dsig.PublicKey) {
	if c.keyTTL <= 0 {
		return
	}
	c.keyMu.Lock()
	defer c.keyMu.Unlock()
	now := time.Now()
	if len(c.keyCache) >= maxKeyCacheEntries {
		for k, e := range c.keyCache {
			if now.After(e.exp) {
				delete(c.keyCache, k)
			}
		}
		for k := range c.keyCache {
			if len(c.keyCache) < maxKeyCacheEntries {
				break
			}
			delete(c.keyCache, k)
		}
	}
	c.keyCache[url] = keyCacheEntry{key: pk, exp: now.Add(c.keyTTL)}
}

// FetchPublicKey is an alias for FetchKey retained for clarity at
// call sites that only need the verification primitive.
func (c *Client) FetchPublicKey(ctx context.Context, addr Address, kid string) (*dsig.PublicKey, error) {
	return c.FetchKey(ctx, addr, kid)
}
