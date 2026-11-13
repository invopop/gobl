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
// header's from and to fields must hold gobl: addresses, to must be self,
// and a valid signature must carry from as its iss and to as its aud, with
// an iat no later than now and a key that is still valid. The from address
// is then the sender.
func (c *Client) VerifyDelivery(ctx context.Context, env *gobl.Envelope, self Address) error {
	_, err := c.DeliveryKeys(ctx, env, self)
	return err
}

// DeliveryKeys makes the checks of VerifyDelivery and returns the keys of
// every valid signature confirming the envelope's from and to. A sender
// rotating its key signs with both the old and the new one.
func (c *Client) DeliveryKeys(ctx context.Context, env *gobl.Envelope, self Address) ([]*dsig.PublicKey, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	want, err := ParseAddress(string(self))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	from, to, err := routing(env)
	if err != nil {
		return nil, err
	}
	if to != want {
		return nil, fmt.Errorf("%w: to %q is not %q", ErrVerifyFailed, to, want)
	}
	return c.pairKeys(ctx, env, from, to)
}

// pairKeys returns the distinct keys of the valid signatures by from with
// to as their audience, each with an iat no later than now and a key that
// is still valid.
func (c *Client) pairKeys(ctx context.Context, env *gobl.Envelope, from, to Address) ([]*dsig.PublicKey, error) {
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	var keys []*dsig.PublicKey
	var unavailable error
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
		if checkIssued(p) != nil {
			continue
		}
		if candidates++; candidates > maxIssuerCandidates {
			break
		}
		pub, err := c.FetchKey(ctx, from, sig.KeyID())
		if err != nil {
			if errors.Is(err, ErrUnavailable) {
				unavailable = err
			}
			continue
		}
		if err := env.VerifySignature(sig, pub); err != nil {
			continue
		}
		if pub.Allows(time.Now()) != nil {
			continue
		}
		if !hasKey(keys, pub.JKT()) {
			keys = append(keys, pub)
		}
	}
	if len(keys) > 0 {
		return keys, nil
	}
	if unavailable != nil {
		return nil, unavailable
	}
	return nil, fmt.Errorf("%w: no valid signature from %q to %q", ErrVerifyFailed, from, to)
}

// hasKey reports whether keys includes the key with thumbprint jkt.
func hasKey(keys []*dsig.PublicKey, jkt string) bool {
	for _, k := range keys {
		if k.JKT() == jkt {
			return true
		}
	}
	return false
}

// maxIssuerCandidates bounds how many of an envelope's signatures by one
// issuer are checked against its keys: each can cost a key fetch from
// that issuer.
const maxIssuerCandidates = 4

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
