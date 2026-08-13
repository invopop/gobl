package net

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
)

// Who fetches and verifies addr's party envelope. The subject must equal addr
// and have an audience-free self-signature. WithIdentity authenticates the
// request. HTTP 204 and 202 responses return ErrNoContent and ErrPending.
func (c *Client) Who(ctx context.Context, addr Address) (*gobl.Envelope, error) {
	// Canonicalize so well-known URLs and the issuer comparison use
	// the ASCII form regardless of how the address was written.
	addr, err := ParseAddress(string(addr))
	if err != nil {
		return nil, err
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
	subject, err := c.VerifyParty(ctx, env)
	if err != nil {
		return nil, err
	}
	if subject != addr {
		return nil, fmt.Errorf("%w: who response is the identity of %q, not %q", ErrVerifyFailed, subject, addr)
	}
	// A who response is a public document: it must carry at least
	// one audience-free self-signature. Audience-bound self-
	// signatures are ignored; an envelope with only caller-bound
	// signatures is not a public identity.
	ok, err := c.subjectSignatureFor(ctx, env, addr, func(p *head.SigningPayload) bool {
		return p.Aud == ""
	})
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: who response carries no audience-free self-signature", ErrVerifyFailed)
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
