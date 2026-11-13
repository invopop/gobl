package net

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
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
	addr, err := ParseAddress(uri.Opaque())
	if err != nil {
		return "", fmt.Errorf("%w: %s %q: %v", ErrVerifyFailed, field, uri, err)
	}
	return addr, nil
}

// VerifyDelivery checks that env is a GOBL Net envelope for self: its
// header's from and to fields must hold gobl: addresses, to must be self,
// and a valid signature must carry from as its iss and to as its aud. The
// from address is then the sender.
func (c *Client) VerifyDelivery(ctx context.Context, env *gobl.Envelope, self Address) error {
	if env == nil || env.Head == nil || !env.Signed() {
		return fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
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
	if len(env.Signatures) > maxEnvelopeSignatures {
		return fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	var unavailable error
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
		return nil
	}
	if unavailable != nil {
		return unavailable
	}
	return fmt.Errorf("%w: no valid signature from %q to %q", ErrVerifyFailed, from, to)
}
