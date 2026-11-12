package net

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
)

// VerifyParty returns the subject of a party envelope, the address in its
// header's from field, after verifying a self-signature by that address.
// The party's endpoints are not consulted. It does not verify authority
// endorsements.
func (c *Client) VerifyParty(ctx context.Context, env *gobl.Envelope) (Address, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return "", fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if _, ok := env.Extract().(*org.Party); !ok {
		return "", ErrPartyMissing
	}
	subject, _, err := routing(env)
	if err != nil {
		return "", err
	}
	if subject == "" {
		return "", fmt.Errorf("%w: envelope has no %s: from address", ErrVerifyFailed, Scheme)
	}
	ok, err := c.subjectSignatureFor(ctx, env, subject, func(*head.SigningPayload) bool {
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

// routing returns the canonical addresses in the envelope header's from and
// to fields, each "" when empty. A value that is not a gobl: address means
// the envelope is not meant for GOBL Net, and is an error. The fields are
// unsigned: callers must confirm them against the signatures.
func routing(env *gobl.Envelope) (from, to Address, err error) {
	if from, err = routingAddress("from", env.Head.From); err != nil {
		return "", "", err
	}
	if to, err = routingAddress("to", env.Head.To); err != nil {
		return "", "", err
	}
	return from, to, nil
}

func routingAddress(field string, uri cbc.URI) (Address, error) {
	if uri == "" {
		return "", nil
	}
	if uri.Scheme() != Scheme {
		return "", fmt.Errorf("%w: %s %q is not a %s: address", ErrVerifyFailed, field, uri, Scheme)
	}
	addr, err := ParseAddress(uri.Opaque())
	if err != nil {
		return "", fmt.Errorf("%w: %s %q: %v", ErrVerifyFailed, field, uri, err)
	}
	return addr, nil
}

// VerifyDelivery returns the sole issuer of valid signatures addressed to
// self. It fails if no issuer or multiple issuers bind the delivery, or if
// the header's from or to field is set to anything but the sender and self.
func (c *Client) VerifyDelivery(ctx context.Context, env *gobl.Envelope, self Address) (Address, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return "", fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	want, err := ParseAddress(string(self))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	from, to, err := routing(env)
	if err != nil {
		return "", err
	}
	if to != "" && to != want {
		return "", fmt.Errorf("%w: to %q is not this inbox %q", ErrVerifyFailed, to, want)
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
	if from != "" && from != sender {
		return "", fmt.Errorf("%w: from %q is not the sender %q", ErrVerifyFailed, from, sender)
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
