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
	authorities []*dsig.PublicKey
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithFetcher sets a custom Fetcher implementation.
func WithFetcher(f Fetcher) ClientOption {
	return func(c *Client) {
		c.fetcher = f
	}
}

// WithAuthorities adds trusted authority public keys to the client,
// supplementing the built-in Authorities.
func WithAuthorities(keys ...*dsig.PublicKey) ClientOption {
	return func(c *Client) {
		c.authorities = append(c.authorities, keys...)
	}
}

// NewClient creates a new GOBL Net client.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		fetcher:     NewHTTPFetcher(),
		authorities: append([]*dsig.PublicKey{}, Authorities...),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// FetchKeySet retrieves the KeySet from the well-known URL derived from
// the given Address. The response is compatible with both standard JWKS
// and GOBL's extended KeySet format.
func (c *Client) FetchKeySet(ctx context.Context, addr Address) (*KeySet, error) {
	if err := addr.Validate(); err != nil {
		return nil, err
	}
	url := addr.JWKSURL()
	data, err := c.fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	ks := new(KeySet)
	if err := json.Unmarshal(data, ks); err != nil {
		return nil, fmt.Errorf("%w: invalid KeySet response: %v", ErrFetchFailed, err)
	}
	return ks, nil
}

// FetchPublicKey retrieves the KeySet for the given address and finds
// the key matching the provided kid. Returns a dsig.PublicKey ready
// for verification.
func (c *Client) FetchPublicKey(ctx context.Context, addr Address, kid string) (*dsig.PublicKey, error) {
	ks, err := c.FetchKeySet(ctx, addr)
	if err != nil {
		return nil, err
	}
	return ks.Key(kid)
}
