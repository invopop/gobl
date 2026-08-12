package net

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

// Authorities is the default list of trusted registration authorities.
var Authorities = []Address{
	"lookup.gobl.org",
}

// SandboxAuthorities is the trust list selected by WithSandbox. Sandbox
// endorsements must not be accepted in live contexts.
var SandboxAuthorities = []Address{
	"lookup.sandbox.gobl.org",
}

// maxEnvelopeSignatures bounds how many signatures VerifyAuthority is
// willing to examine: each candidate can cost a network key fetch, so
// an unbounded envelope would be a fetch-amplification gadget.
// Legitimate envelopes carry a handful of signatures.
const maxEnvelopeSignatures = 32

// RegisterAuthority appends addr to the default authority list.
func RegisterAuthority(addr Address) {
	Authorities = append(Authorities, addr)
}

// Endorsement identifies the trusted authority behind an envelope and any
// confirmed, independent identity verifier it named.
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

// VerifyAuthority verifies countersignatures from the client's trusted
// authorities. It returns the strongest endorsement found, preferring one with
// confirmed verifier evidence. The verifier must differ from the authority;
// self-referential, missing, or invalid verifier evidence leaves a
// registration-only endorsement valid. It returns ErrUnknownAuthority when no
// trusted candidate exists, ErrSignatureExpired for expired authority evidence,
// ErrUnavailable for transient key lookup failures, or ErrVerifyFailed when
// candidates fail cryptographic verification.
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

	// claimErr records a candidate that verified but carries an
	// expired exp; unavailableErr records a transient key-fetch
	// failure and outranks it. registered holds the first valid
	// registered-only endorsement while the search continues for a
	// confirmed verifier.
	var lastErr, claimErr, unavailableErr error
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
			if errors.Is(err, ErrUnavailable) {
				unavailableErr = err
			} else {
				lastErr = err
			}
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
	if unavailableErr != nil {
		return nil, unavailableErr
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
// names a valid address, distinct from the authority, whose countersignature
// verifies. An absent, self-referential, malformed, or invalid claim returns
// "", degrading the endorsement to registered. A transient key-fetch failure
// returns ErrUnavailable instead of silently downgrading the endorsement.
func (c *Client) confirmVerifier(ctx context.Context, env *gobl.Envelope, authority Address, p *head.SigningPayload) (Address, error) {
	verifier, err := ParseAddress(p.Verifier)
	if err != nil {
		return "", nil //nolint:nilerr // malformed claim degrades to registered
	}
	// Registration and identity verification must be performed by distinct
	// participants. A self-referential claim provides no verifier evidence.
	if verifier == authority {
		return "", nil
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
