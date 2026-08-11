package net

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

// VerifyEnvelope performs remote verification of a signed GOBL envelope.
// It reads the subject's GOBL Net identity (iss) from the first
// signature's signed payload, fetches that address's public keys, and
// verifies that signature. When expectedAud is non-empty, at least one
// valid signature by the subject must carry that signed audience —
// searched across all signatures, since the subject appends one
// audience-bound signature per delivery hop and their order is not
// significant. The verified subject address is returned. Signatures by
// other parties (e.g. authority countersignatures) are not checked
// here; use VerifyAuthority for those.
func (c *Client) VerifyEnvelope(ctx context.Context, env *gobl.Envelope, expectedAud Address) (Address, error) {
	// A malformed envelope may carry signatures without a header;
	// reject rather than let header verification panic.
	if env == nil || env.Head == nil || !env.Signed() {
		return "", fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}

	sig := env.Signatures[0]
	p, err := head.SignedPayload(sig)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	if p.Iss == "" {
		return "", fmt.Errorf("%w: signature has no iss", ErrVerifyFailed)
	}
	// Canonicalize the issuer so key-fetch URLs and comparisons use
	// the ASCII form regardless of how the iss was written. Anything
	// that is not a bare FQDN — including URI forms — is rejected.
	issuer, err := ParseAddress(p.Iss)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	kid := sig.KeyID()
	if kid == "" {
		return "", fmt.Errorf("%w: signature has no key ID", ErrVerifyFailed)
	}

	pubKey, err := c.FetchKey(ctx, issuer, kid)
	if err != nil {
		if errors.Is(err, ErrUnavailable) {
			return "", err
		}
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	// VerifySignature enforces the key's validity window against the
	// signed `iat` via head.Header.Verify.
	if err := env.VerifySignature(sig, pubKey); err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	if expectedAud != "" {
		want, err := ParseAddress(string(expectedAud))
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
		}
		ok, err := c.subjectSignatureFor(ctx, env, issuer, func(sp *head.SigningPayload) bool {
			aud, aerr := ParseAddress(sp.Aud)
			return aerr == nil && aud == want
		})
		if err != nil {
			return "", err
		}
		if !ok {
			return "", fmt.Errorf("%w: no valid signature by %q for audience %q", ErrVerifyFailed, issuer, want)
		}
	}

	return issuer, nil
}

// subjectSignatureFor reports whether the envelope carries at least
// one cryptographically valid signature by subject whose signed
// payload satisfies match. The signature count is bounded the same
// way as VerifyAuthority (each candidate can cost a key fetch); a
// transient key-fetch failure is returned as ErrUnavailable when no
// other candidate matched.
func (c *Client) subjectSignatureFor(ctx context.Context, env *gobl.Envelope, subject Address, match func(*head.SigningPayload) bool) (bool, error) {
	if len(env.Signatures) > maxEnvelopeSignatures {
		return false, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	var unavailable error
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		iss, err := ParseAddress(p.Iss)
		if err != nil || iss != subject {
			continue
		}
		if !match(p) {
			continue
		}
		pub, err := c.FetchKey(ctx, subject, sig.KeyID())
		if err != nil {
			if errors.Is(err, ErrUnavailable) {
				unavailable = err
			}
			continue
		}
		if err := env.VerifySignature(sig, pub); err != nil {
			continue
		}
		return true, nil
	}
	if unavailable != nil {
		return false, unavailable
	}
	return false, nil
}
