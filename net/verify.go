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
// document must match its digest, its header's from and to fields must hold gobl: addresses, to must be self,
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
	cs, err := c.DeliveryConfirmations(ctx, env, self)
	if err != nil {
		return nil, err
	}
	return keysOf(cs), nil
}

// MaxRequestAge is how long after its signing time a confirming
// signature counts towards a registration request (README.md §7.3): the
// maximum age of a request token (§6).
const MaxRequestAge = maxTokenTTL

// ConfirmingSignature is a valid signature confirming an envelope's from
// and to.
type ConfirmingSignature struct {
	// Key is the signer's published key.
	Key *dsig.PublicKey
	// IssuedAt is the signature's iat.
	IssuedAt time.Time
	// Successor is the thumbprint in the signature's cnf claim, if any:
	// in a key change, the new key the current key names (README.md
	// §7.3).
	Successor string
}

// DeliveryConfirmations makes the checks of VerifyDelivery and returns
// the valid signatures confirming the envelope's from and to, one per
// key: its most recent. Authorities use their signing times and cnf
// claims to process registration requests (README.md §7.3).
func (c *Client) DeliveryConfirmations(ctx context.Context, env *gobl.Envelope, self Address) ([]*ConfirmingSignature, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if err := checkDigest(env); err != nil {
		return nil, err
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
	return c.pairConfirmations(ctx, env, from, to)
}

// pairKeys returns the distinct keys of the valid signatures by from with
// to as their audience.
func (c *Client) pairKeys(ctx context.Context, env *gobl.Envelope, from, to Address) ([]*dsig.PublicKey, error) {
	cs, err := c.pairConfirmations(ctx, env, from, to)
	if err != nil {
		return nil, err
	}
	return keysOf(cs), nil
}

// keysOf returns the keys of cs.
func keysOf(cs []*ConfirmingSignature) []*dsig.PublicKey {
	keys := make([]*dsig.PublicKey, len(cs))
	for i, c := range cs {
		keys[i] = c.Key
	}
	return keys
}

// pairConfirmations returns the valid signatures by from with to as their
// audience, each with an iat no later than now and a key that is still
// valid, keeping the most recent per key.
func (c *Client) pairConfirmations(ctx context.Context, env *gobl.Envelope, from, to Address) ([]*ConfirmingSignature, error) {
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	var out []*ConfirmingSignature
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
		cs := &ConfirmingSignature{Key: pub, IssuedAt: time.Unix(p.IssuedAt, 0).UTC()}
		if p.Confirmation != nil {
			cs.Successor = p.Confirmation.JKT
		}
		out = addConfirmation(out, cs)
	}
	if len(out) > 0 {
		return out, nil
	}
	if unavailable != nil {
		return nil, unavailable
	}
	if reason != nil {
		return nil, fmt.Errorf("%w: no valid signature from %q to %q: %v", ErrVerifyFailed, from, to, reason)
	}
	return nil, fmt.Errorf("%w: no valid signature from %q to %q", ErrVerifyFailed, from, to)
}

// addConfirmation adds cs to out, keeping one signature per key: one
// naming a successor over one that does not, unless that one is more
// than MaxRequestAge newer, and otherwise the more recent.
func addConfirmation(out []*ConfirmingSignature, cs *ConfirmingSignature) []*ConfirmingSignature {
	for i, o := range out {
		if o.Key.JKT() == cs.Key.JKT() {
			if preferConfirmation(cs, o) {
				out[i] = cs
			}
			return out
		}
	}
	return append(out, cs)
}

// preferConfirmation reports whether a should be kept over b, two
// signatures by the same key.
func preferConfirmation(a, b *ConfirmingSignature) bool {
	switch {
	case a.Successor != "" && b.Successor == "":
		return b.IssuedAt.Sub(a.IssuedAt) <= MaxRequestAge
	case a.Successor == "" && b.Successor != "":
		return a.IssuedAt.Sub(b.IssuedAt) > MaxRequestAge
	case a.IssuedAt.Equal(b.IssuedAt):
		return a.Successor > b.Successor // independent of signature order
	}
	return a.IssuedAt.After(b.IssuedAt)
}

// stillValid fails when pub's valid_until has passed.
func stillValid(pub *dsig.PublicKey) error {
	if pub.ValidUntil != nil && time.Now().After(pub.ValidUntil.Time) {
		return fmt.Errorf("%w: key %q was retired at %s", ErrVerifyFailed, pub.ID(), pub.ValidUntil)
	}
	return nil
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
