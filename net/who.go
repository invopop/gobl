package net

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
)

// Who fetches the public identity for the given address with a GET on
// its well-known who endpoint and verifies it. The response envelope's
// first signature must be the subject's self-signature: its signed
// `iss` is resolved to a published key and must name the fetched
// address. The document must be an org.Party. Authority
// countersignatures, if any, are preserved on the returned envelope
// for VerifyAuthority.
//
// The request carries a bearer request token minted from the client's
// identity (WithIdentity); conforming servers reject requests without
// one.
//
// A 204 response returns ErrNoContent: the address exists but does not
// publish identity details (a receive-only account). A 202 response
// returns ErrPending: the request was recorded and the owner may
// deliver its party envelope to the requester's inbox later.
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
	issuer, err := c.VerifyEnvelope(ctx, env, "")
	if err != nil {
		return nil, err
	}
	if issuer != addr {
		return nil, fmt.Errorf("%w: who issuer %q does not match address %q", ErrVerifyFailed, issuer, addr)
	}
	// A who response is a public document: a caller-bound (aud-carrying)
	// envelope is not a conforming identity and must not be treated as
	// one.
	p, err := head.SignedPayload(env.Signatures[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	if p.Aud != "" {
		return nil, fmt.Errorf("%w: who response must not be audience-bound (aud %q)", ErrVerifyFailed, p.Aud)
	}
	if _, ok := env.Extract().(*org.Party); !ok {
		return nil, ErrPartyMissing
	}
	return env, nil
}

// VerifySender confirms that the given address is approved to send
// documents: its who identity must verify (see Who) and carry a
// countersignature from one of the client's trusted authorities.
// When requireVerified is true the endorsement must additionally
// carry a confirmed verifier (see VerifyAuthority), else
// ErrNotVerified. Returns the sender's endorsed party on success.
// Receiving inboxes call this with the verified issuer of an
// incoming envelope before accepting it.
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
