package net

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
)

// routing returns the canonical addresses in the envelope header's from and
// to fields. Both are required, and a value that is not a gobl: address
// means the envelope is not meant for GOBL Net. The fields are unsigned:
// callers must confirm them against the signatures.
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
		return "", fmt.Errorf("%w: envelope has no %s address", ErrVerifyFailed, field)
	}
	if uri.Scheme() != Scheme {
		return "", fmt.Errorf("%w: %s %q is not a %s: address", ErrVerifyFailed, field, uri, Scheme)
	}
	if strings.ContainsAny(string(uri), "?#") {
		return "", fmt.Errorf("%w: %s %q has a query or fragment", ErrVerifyFailed, field, uri)
	}
	addr, err := ParseAddress(uri.Opaque())
	if err != nil {
		return "", fmt.Errorf("%w: %s %q: %v", ErrVerifyFailed, field, uri, err)
	}
	return addr, nil
}

// VerifyDelivery checks that env is a GOBL Net envelope for self: its
// document must match its digest, its header's from and to fields must
// hold gobl: addresses, to must be self, and a valid signature must carry
// from as its iss and to as its aud, with an iat no later than now, no exp
// that has passed, and a key that is still valid. The from address is
// then the sender.
func (c *Client) VerifyDelivery(ctx context.Context, env *gobl.Envelope, self Address) error {
	if env == nil || env.Head == nil || !env.Signed() {
		return fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if err := checkDigest(env); err != nil {
		return err
	}
	want, err := ParseAddress(string(self))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	from, to, err := routing(env)
	if err != nil {
		return err
	}
	if to != want {
		return fmt.Errorf("%w: to %q is not %q", ErrVerifyFailed, to, want)
	}
	return c.confirmed(ctx, env, from, to)
}

// confirmed fails unless env carries a valid signature by from with to as
// its audience, an iat no later than now, no exp that has passed, and a
// key that is still valid.
func (c *Client) confirmed(ctx context.Context, env *gobl.Envelope, from, to Address) error {
	if len(env.Signatures) > maxEnvelopeSignatures {
		return fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	var unavailable, reason error
	candidates := 0
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		if iss, err := ParseAddress(p.Iss); err != nil || iss != from {
			continue
		}
		if aud, err := ParseAddress(p.Aud); err != nil || aud != to {
			continue
		}
		if err := checkIssued(p); err != nil {
			reason = err
			continue
		}
		if expired(p.ExpiresAt) {
			reason = fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
				time.Unix(p.ExpiresAt, 0).UTC().Format(time.RFC3339))
			continue
		}
		if candidates++; candidates > maxIssuerCandidates {
			break
		}
		pub, err := c.FetchKey(ctx, from, sig.KeyID())
		if err != nil {
			if errors.Is(err, ErrUnavailable) {
				unavailable = err
			} else {
				reason = err
			}
			continue
		}
		if err := env.VerifySignature(sig, pub); err != nil {
			reason = err
			continue
		}
		// The signature's iat is within the key's validity; the key must
		// also not have been retired since.
		if err := stillValid(pub); err != nil {
			reason = err
			continue
		}
		return nil
	}
	if unavailable != nil {
		return unavailable
	}
	if reason != nil {
		return fmt.Errorf("%w: no valid signature from %q to %q: %w", ErrVerifyFailed, from, to, reason)
	}
	return fmt.Errorf("%w: no valid signature from %q to %q", ErrVerifyFailed, from, to)
}

// stillValid fails when pub's valid_until has passed.
func stillValid(pub *dsig.PublicKey) error {
	if pub.ValidUntil != nil && time.Now().After(pub.ValidUntil.Time) {
		return fmt.Errorf("%w: key %q was retired at %s", ErrVerifyFailed, pub.ID(), pub.ValidUntil)
	}
	return nil
}

// maxIssuerCandidates bounds how many of an envelope's signatures by one
// issuer are checked against its keys: each can cost a key fetch from
// that issuer.
const maxIssuerCandidates = 4

// checkDigest fails unless env has a document matching the digest in
// its header, which is what signatures cover.
func checkDigest(env *gobl.Envelope) error {
	if env.Document == nil || env.Document.IsEmpty() {
		return fmt.Errorf("%w: envelope has no document", ErrVerifyFailed)
	}
	if env.Head.Digest == nil {
		return fmt.Errorf("%w: envelope has no digest", ErrVerifyFailed)
	}
	if d, err := env.Digest(); err != nil || env.Head.Digest.Equals(d) != nil {
		return fmt.Errorf("%w: envelope digest does not match its document", ErrVerifyFailed)
	}
	return nil
}

// checkIssued reports whether a signature carries an iat no later than
// now, allowing clock skew.
func checkIssued(p *head.SigningPayload) error {
	if p.IssuedAt == 0 {
		return fmt.Errorf("%w: signature by %q has no iat", ErrVerifyFailed, p.Iss)
	}
	if time.Unix(p.IssuedAt, 0).After(time.Now().Add(tokenClockSkew)) {
		return fmt.Errorf("%w: signature by %q is issued in the future", ErrVerifyFailed, p.Iss)
	}
	return nil
}
