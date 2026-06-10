package net

import (
	"context"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
)

// VerifyEnvelope performs remote verification of a signed GOBL envelope.
// It reads the signer's GOBL Net identity (iss) from the first
// signature's signed payload, fetches that address's public keys, and
// verifies the signature. When expectedAud is non-empty, the signature's
// signed audience (aud) must equal it. The verified issuer address is
// returned.
func (c *Client) VerifyEnvelope(ctx context.Context, env *gobl.Envelope, expectedAud cbc.URI) (Address, error) {
	if !env.Signed() {
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
	if p.Iss.Scheme() != Scheme {
		return "", fmt.Errorf("%w: iss %q is not a gobl address", ErrVerifyFailed, p.Iss)
	}
	issuer := Address(p.Iss.Opaque())
	if err := issuer.Validate(); err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	kid := sig.KeyID()
	if kid == "" {
		return "", fmt.Errorf("%w: signature has no key ID", ErrVerifyFailed)
	}

	pubKey, err := c.FetchKey(ctx, issuer, kid)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	// env.Verify enforces the key's validity window against the signed
	// `ts` via head.Header.Verify.
	if err := env.Verify(pubKey); err != nil {
		return "", fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	if expectedAud != "" && p.Aud != expectedAud {
		return "", fmt.Errorf("%w: audience mismatch (got %q, want %q)", ErrVerifyFailed, p.Aud, expectedAud)
	}

	return issuer, nil
}
