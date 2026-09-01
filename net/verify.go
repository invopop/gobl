package net

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
)

// VerifyParty returns the address declared by a party envelope's gobl endpoint
// after verifying a self-signature from that address. It does not verify
// authority endorsements.
func (c *Client) VerifyParty(ctx context.Context, env *gobl.Envelope) (Address, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return "", fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	party, ok := env.Extract().(*org.Party)
	if !ok {
		return "", ErrPartyMissing
	}
	ep := party.Endpoint(Scheme)
	if ep == nil {
		return "", fmt.Errorf("%w: party declares no gobl: endpoint", ErrPartyMissing)
	}
	subject, err := ParseAddress(ep.URI.Opaque())
	if err != nil {
		return "", fmt.Errorf("%w: party endpoint %q is not a valid address: %v", ErrVerifyFailed, ep.URI, err)
	}
	ok, err = c.subjectSignatureFor(ctx, env, subject, func(*head.SigningPayload) bool {
		return true
	})
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("%w: no valid self-signature by %q", ErrVerifyFailed, subject)
	}
	return subject, nil
}

// VerifyDelivery returns the sole issuer of valid signatures addressed to
// self. It fails if no issuer or multiple issuers bind the delivery.
func (c *Client) VerifyDelivery(ctx context.Context, env *gobl.Envelope, self Address) (Address, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return "", fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	want, err := ParseAddress(string(self))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	if len(env.Signatures) > maxEnvelopeSignatures {
		return "", fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	var sender Address
	var unavailable error
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		aud, err := ParseAddress(p.Aud)
		if err != nil || aud != want {
			continue
		}
		iss, err := ParseAddress(p.Iss)
		if err != nil {
			continue
		}
		pub, err := c.FetchKey(ctx, iss, sig.KeyID())
		if err != nil {
			if errors.Is(err, ErrUnavailable) {
				unavailable = err
			}
			continue
		}
		if err := env.VerifySignature(sig, pub); err != nil {
			continue
		}
		if sender != "" && sender != iss {
			return "", fmt.Errorf("%w: ambiguous delivery: bound to %q by both %q and %q",
				ErrVerifyFailed, want, sender, iss)
		}
		sender = iss
	}
	if sender == "" {
		if unavailable != nil {
			return "", unavailable
		}
		return "", fmt.Errorf("%w: no valid signature bound to %q", ErrVerifyFailed, want)
	}
	return sender, nil
}

// subjectSignatureFor reports whether the envelope carries a valid
// signature by subject whose payload satisfies match. The signature
// count is capped (each candidate can cost a key fetch); a transient
// key-fetch failure returns ErrUnavailable when nothing else matched.
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
