package net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdnet "net"
	"net/http"
	"time"

	"github.com/invopop/gobl/dsig"
)

const (
	defaultTimeout = 10 * time.Second
	maxBodySize    = 1 << 20 // 1MB
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
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
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

// Fetcher defines the interface for fetching data from a URL.
type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
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

// Fetch retrieves the body from the given URL.
func (f *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: HTTP %d from %s", ErrFetchFailed, resp.StatusCode, url)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	return body, nil
}

// Client provides GOBL Net operations including KeySet fetching
// and remote verification.
type Client struct {
	fetcher     Fetcher
	authorities []Address
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithFetcher sets a custom Fetcher implementation.
func WithFetcher(f Fetcher) ClientOption {
	return func(c *Client) {
		c.fetcher = f
	}
}

// WithAuthorities adds trusted authority GOBL Net Addresses to the
// client, supplementing the built-in Authorities.
func WithAuthorities(addrs ...Address) ClientOption {
	return func(c *Client) {
		c.authorities = append(c.authorities, addrs...)
	}
}

// NewClient creates a new GOBL Net client.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		fetcher:     NewHTTPFetcher(),
		authorities: append([]Address{}, Authorities...),
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
	if err := addr.Validate(); err != nil {
		return nil, err
	}
	if kid == "" {
		return nil, fmt.Errorf("%w: kid is required", ErrFetchFailed)
	}
	data, err := c.fetcher.Fetch(ctx, addr.KeyURL(kid))
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
	return pk, nil
}

// FetchPublicKey is an alias for FetchKey retained for clarity at
// call sites that only need the verification primitive.
func (c *Client) FetchPublicKey(ctx context.Context, addr Address, kid string) (*dsig.PublicKey, error) {
	return c.FetchKey(ctx, addr, kid)
}
