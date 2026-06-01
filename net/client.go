package net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/invopop/gobl/dsig"
)

const (
	defaultTimeout = 10 * time.Second
	maxBodySize    = 1 << 20 // 1MB
)

// Fetcher defines the interface for fetching data from a URL.
type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// HTTPFetcher implements Fetcher using net/http.
type HTTPFetcher struct {
	Client *http.Client
}

// NewHTTPFetcher creates an HTTPFetcher with sensible defaults.
func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		Client: &http.Client{
			Timeout: defaultTimeout,
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
