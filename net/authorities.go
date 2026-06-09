package net

import (
	"context"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

// Authorities is the hardcoded set of GOBL Net Addresses considered
// trusted KYC vendors. A /who response is only considered endorsed
// if it is signed by at least one of these.
//
// The list is intentionally empty in this release; entries are added
// here as vendors come online.
var Authorities = []Address{}

// RegisterAuthority adds an address to the global set of trusted
// authority addresses.
func RegisterAuthority(addr Address) {
	Authorities = append(Authorities, addr)
}

// VerifyAuthority checks that the envelope carries at least one
// signature whose signed `iss` resolves to an address in the
// client's known authorities (the package-level Authorities slice
// plus anything added via WithAuthorities). Each candidate signature
// is cryptographically verified against the authority's own
// published key.
//
// Returns nil on the first authority signature that verifies. If no
// signature is from a known authority, returns ErrUnknownAuthority.
// If all candidates fail crypto verification, returns
// ErrVerifyFailed wrapping the last error.
//
// Callers that want to accept self-signed (no-authority) envelopes
// should skip this call rather than ignore its error.
func (c *Client) VerifyAuthority(ctx context.Context, env *gobl.Envelope) error {
	if env == nil || len(env.Signatures) == 0 {
		return fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if len(c.authorities) == 0 {
		return fmt.Errorf("%w: no authorities registered on this client", ErrUnknownAuthority)
	}
	auths := make(map[Address]bool, len(c.authorities))
	for _, a := range c.authorities {
		auths[a] = true
	}

	var lastErr error
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		if p.Iss.Scheme() != Scheme {
			continue
		}
		issuer := Address(p.Iss.Opaque())
		if err := issuer.Validate(); err != nil {
			continue
		}
		if !auths[issuer] {
			continue
		}
		// Candidate from a known authority — verify it crypto-wise.
		pub, err := c.FetchKey(ctx, issuer, sig.KeyID())
		if err != nil {
			lastErr = err
			continue
		}
		if err := env.Head.Verify(sig, pub); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrVerifyFailed, lastErr)
	}
	return ErrUnknownAuthority
}
