package net

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

// Authorities is the hardcoded set of GOBL Net Addresses considered
// trusted registration authorities. A /who response is only
// considered endorsed if it is signed by at least one of these.
//
// lookup.gobl.org is the network's default authority; operators may
// add their own with RegisterAuthority or WithAuthorities.
var Authorities = []Address{
	"lookup.gobl.org",
}

// SandboxAuthorities is the default trust list for sandbox
// deployments: lookup.sandbox.gobl.org runs the same registration
// service as the live authority but endorses test identities through
// relaxed verification providers. The live and sandbox lists are
// disjoint by construction — endorsements from a sandbox authority
// MUST never be accepted in live contexts, and vice versa. Clients
// opt in with WithSandbox.
var SandboxAuthorities = []Address{
	"lookup.sandbox.gobl.org",
}

// maxEnvelopeSignatures bounds how many signatures VerifyAuthority is
// willing to examine: each candidate can cost a network key fetch, so
// an unbounded envelope would be a fetch-amplification gadget.
// Legitimate envelopes carry a handful of signatures.
const maxEnvelopeSignatures = 32

// RegisterAuthority adds an address to the global set of trusted
// authority addresses.
func RegisterAuthority(addr Address) {
	Authorities = append(Authorities, addr)
}

// Endorsement describes a successful authority check on an envelope:
// the trusted authority that countersigned it and, when that
// countersignature names a verifier whose own countersignature also
// verifies, the verifier's address.
type Endorsement struct {
	// Authority is the trusted address whose countersignature
	// endorsed the envelope: the subject is a registered identity.
	Authority Address
	// Verifier is the address the authority named as having performed
	// identity verification (KYC/KYB) of the subject, confirmed by
	// the verifier's own countersignature on the same envelope. Empty
	// when the identity is registered but not verified, including
	// when a named verifier's countersignature is absent, invalid, or
	// expired — verification degrades, registration stands.
	Verifier Address
}

// Verified reports whether the endorsement includes confirmed
// identity verification.
func (e *Endorsement) Verified() bool {
	return e != nil && e.Verifier != ""
}

// VerifyAuthority checks that the envelope carries at least one
// signature whose signed `iss` resolves to an address in the
// client's known authorities (the package-level Authorities slice
// plus anything added via WithAuthorities). Each candidate signature
// is cryptographically verified against the authority's own
// published key, and a candidate whose signed exp claim has passed
// is rejected: expired endorsements are not evidence.
//
// When the authority's countersignature carries a `verifier` claim,
// the named address's own countersignature on the same envelope is
// looked up and verified the same way; if it holds, the returned
// endorsement carries the verifier's address. A verifier that names
// itself needs no second signature — its own countersignature serves
// as both attestations. A named verifier whose countersignature is
// missing, invalid, or expired degrades the endorsement to
// registered rather than failing it: the registration stands on its
// own. Callers that require verification check Endorsement.Verified.
//
// Every candidate authority signature is considered: an endorsement
// with a confirmed verifier is preferred over a registered-only one
// regardless of signature order. If no signature is from a known
// authority, returns ErrUnknownAuthority. If a verified authority
// signature has expired, returns ErrSignatureExpired. If all
// candidates fail crypto verification, returns ErrVerifyFailed
// wrapping the last error.
//
// Callers that want to accept self-signed (no-authority) envelopes
// should skip this call rather than ignore its error.
func (c *Client) VerifyAuthority(ctx context.Context, env *gobl.Envelope) (*Endorsement, error) {
	// A malformed envelope may carry signatures without a header;
	// reject rather than let header verification panic.
	if env == nil || env.Head == nil || len(env.Signatures) == 0 {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if len(c.authorities) == 0 {
		return nil, fmt.Errorf("%w: no authorities registered on this client", ErrUnknownAuthority)
	}
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	// Both sides of the lookup are canonicalized so U-Label or
	// trailing-dot forms — configured or signed — compare equal.
	auths := make(map[Address]bool, len(c.authorities))
	for _, a := range c.authorities {
		if canon, err := ParseAddress(string(a)); err == nil {
			auths[canon] = true
		}
	}

	// claimErr records a candidate that verified cryptographically but
	// carried an expired exp claim; it takes precedence over crypto
	// failures from other candidates. registered holds the first
	// valid registered-only endorsement while the remaining
	// signatures are searched for one with a confirmed verifier.
	var lastErr, claimErr error
	var registered *Endorsement
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		issuer, err := ParseAddress(p.Iss)
		if err != nil {
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
		// RFC 7519, the current time must be strictly before exp. Zero
		// means the claim is absent (the JSON zero value); any other
		// value — including a pre-1970 NumericDate — is enforced.
		if expired(p.ExpiresAt) {
			claimErr = fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
				time.Unix(p.ExpiresAt, 0).UTC().Format(time.RFC3339))
			continue
		}
		verifier, err := c.confirmVerifier(ctx, env, issuer, p)
		if err != nil {
			// The verifier's key could not be checked (transient):
			// surface the condition rather than silently degrading a
			// genuinely verified endorsement to registered.
			return nil, err
		}
		end := &Endorsement{Authority: issuer, Verifier: verifier}
		if end.Verified() {
			return end, nil
		}
		if registered == nil {
			registered = end
		}
	}
	if registered != nil {
		return registered, nil
	}
	if claimErr != nil {
		return nil, claimErr
	}
	if lastErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, lastErr)
	}
	return nil, ErrUnknownAuthority
}

// confirmVerifier resolves the verifier claim on a verified authority
// countersignature. It returns the verifier's address when the claim
// names a valid bare address whose own countersignature on the
// envelope verifies (or names the authority itself, whose signature
// has already been checked), and "" when the claim is absent,
// malformed, or its countersignature is definitively invalid — an
// unconfirmed verifier degrades to registered, it does not
// invalidate the endorsement. A transient failure to fetch the
// verifier's key (ErrUnavailable) is returned as an error instead:
// "could not check" must not silently downgrade a verified identity.
func (c *Client) confirmVerifier(ctx context.Context, env *gobl.Envelope, authority Address, p *head.SigningPayload) (Address, error) {
	verifier, err := ParseAddress(p.Verifier)
	if err != nil {
		return "", nil //nolint:nilerr // malformed claim degrades to registered
	}
	// An authority that performs verification itself names itself;
	// the countersignature just checked serves as both attestations.
	if verifier == authority {
		return verifier, nil
	}
	if err := c.verifySignatureBy(ctx, env, verifier); err != nil {
		if errors.Is(err, ErrUnavailable) {
			return "", err
		}
		return "", nil //nolint:nilerr // invalid or absent signature degrades
	}
	return verifier, nil
}

// verifySignatureBy reports whether the envelope carries a valid,
// unexpired countersignature whose signed iss resolves to addr.
func (c *Client) verifySignatureBy(ctx context.Context, env *gobl.Envelope, addr Address) error {
	var lastErr error
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		issuer, err := ParseAddress(p.Iss)
		if err != nil || issuer != addr {
			continue
		}
		pub, err := c.FetchKey(ctx, issuer, sig.KeyID())
		if err != nil {
			lastErr = err
			continue
		}
		if err := env.Head.Verify(sig, pub); err != nil {
			lastErr = err
			continue
		}
		if expired(p.ExpiresAt) {
			lastErr = fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
				time.Unix(p.ExpiresAt, 0).UTC().Format(time.RFC3339))
			continue
		}
		return nil
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("%w: no signature from %q", ErrVerifyFailed, addr)
}

// expired reports whether a signed exp claim (zero when absent) has
// passed.
func expired(exp int64) bool {
	return exp != 0 && time.Now().UTC().Unix() >= exp
}
