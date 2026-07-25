package net

import (
	"context"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

// VerifyEnvelope performs remote verification of a signed GOBL envelope.
// It reads the signer's GOBL Net identity (iss) from the first
// signature's signed payload, fetches that address's public keys, and
// verifies that signature. When expectedAud is non-empty, the signature's
// signed audience (aud) must equal it. The verified issuer address is
// returned. Additional signatures (e.g. authority countersignatures) are
// not checked here; use VerifyAuthority for those.
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
		got, err := ParseAddress(p.Aud)
		if err != nil || got != want {
			return "", fmt.Errorf("%w: audience mismatch (got %q, want %q)", ErrVerifyFailed, p.Aud, want)
		}
	}

	return issuer, nil
}
