package net

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/org"
)

// Who fetches addr's party envelope and verifies it as a delivery from
// addr to the client's identity (see VerifyDelivery), which WithIdentity
// sets and the request token presents. HTTP 204 and 202 responses return
// ErrNoContent and ErrPending.
func (c *Client) Who(ctx context.Context, addr Address) (*gobl.Envelope, error) {
	// Canonicalize so well-known URLs and the from comparison use the
	// ASCII form regardless of how the address was written.
	addr, err := ParseAddress(string(addr))
	if err != nil {
		return nil, err
	}
	if c.identity == nil {
		return nil, fmt.Errorf("%w: who requires a client identity", ErrFetchFailed)
	}
	header, err := c.authHeader(addr)
	if err != nil {
		return nil, err
	}
	data, err := c.fetcher.Fetch(ctx, addr.WhoURL(), header)
	if err != nil {
		return nil, err
	}
	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err != nil {
		return nil, fmt.Errorf("%w: invalid who envelope: %v", ErrFetchFailed, err)
	}
	if err := c.VerifyDelivery(ctx, env, c.identity.addr); err != nil {
		return nil, err
	}
	if from, _ := ParseAddress(env.Head.From.Opaque()); from != addr {
		return nil, fmt.Errorf("%w: who response is from %q, not %q", ErrVerifyFailed, from, addr)
	}
	if _, ok := env.Extract().(*org.Party); !ok {
		return nil, ErrPartyMissing
	}
	return env, nil
}

// VerifySender returns addr's party after verifying its identity and a trusted
// authority endorsement. If requireVerified is true, the endorsement must also
// include confirmed verifier evidence.
func (c *Client) VerifySender(ctx context.Context, addr Address, requireVerified bool) (*org.Party, error) {
	env, err := c.Who(ctx, addr)
	if err != nil {
		return nil, err
	}
	end, err := c.VerifyAuthority(ctx, env)
	if err != nil {
		return nil, err
	}
	if requireVerified && !end.Verified() {
		return nil, fmt.Errorf("%w: %s", ErrNotVerified, addr)
	}
	return env.Extract().(*org.Party), nil
}
