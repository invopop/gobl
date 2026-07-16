package net

import (
	"context"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
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
// published key, and a candidate whose signed exp claim has passed
// is rejected: expired endorsements are not evidence.
//
// Returns nil on the first authority signature that verifies. If no
// signature is from a known authority, returns ErrUnknownAuthority.
// If a verified authority signature has expired, returns
// ErrSignatureExpired. If all candidates fail crypto verification,
// returns ErrVerifyFailed wrapping the last error.
//
// Callers that want to accept self-signed (no-authority) envelopes
// should skip this call rather than ignore its error. Callers with
// a minimum scope requirement use VerifyAuthorityWithScope.
func (c *Client) VerifyAuthority(ctx context.Context, env *gobl.Envelope) error {
	return c.verifyAuthority(ctx, env, "")
}

// VerifyAuthorityWithScope performs the VerifyAuthority check and
// additionally requires the countersignature's signed scope claim to
// satisfy minScope: the claim must be the same key or rank higher
// (registered < verified); custom scope keys only satisfy an exact
// match. A verified authority signature falling short of the minimum
// returns ErrScopeInsufficient.
func (c *Client) VerifyAuthorityWithScope(ctx context.Context, env *gobl.Envelope, minScope cbc.Key) error {
	return c.verifyAuthority(ctx, env, minScope)
}

func (c *Client) verifyAuthority(ctx context.Context, env *gobl.Envelope, minScope cbc.Key) error {
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

	// claimErr records a candidate that verified cryptographically but
	// failed a claim check (expired or insufficient scope); it takes
	// precedence over crypto failures from other candidates.
	var lastErr, claimErr error
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
		// Claims are only trusted after the signature verifies. Per
		// RFC 7519, the current time must be strictly before exp.
		if p.ExpiresAt > 0 && time.Now().UTC().Unix() >= p.ExpiresAt {
			claimErr = fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
				time.Unix(p.ExpiresAt, 0).UTC().Format(time.RFC3339))
			continue
		}
		if !scopeSatisfies(p.Scope, minScope) {
			claimErr = fmt.Errorf("%w: %q does not meet minimum %q", ErrScopeInsufficient, p.Scope, minScope)
			continue
		}
		return nil
	}
	if claimErr != nil {
		return claimErr
	}
	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrVerifyFailed, lastErr)
	}
	return ErrUnknownAuthority
}

// scopeRank orders the baseline scope keys: none < registered <
// verified. Unrecognised keys rank 0.
func scopeRank(scope cbc.Key) int {
	switch scope {
	case head.ScopeRegistered:
		return 1
	case head.ScopeVerified:
		return 2
	}
	return 0
}

// scopeSatisfies reports whether a signed scope claim meets the
// required minimum: any claim satisfies an empty minimum, an exact
// match always satisfies, and baseline keys satisfy lower-ranked
// baseline minimums.
func scopeSatisfies(got, minScope cbc.Key) bool {
	if minScope == "" || got == minScope {
		return true
	}
	r := scopeRank(minScope)
	return r > 0 && scopeRank(got) >= r
}
