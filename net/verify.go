package net

import (
	"context"
	"fmt"

	"github.com/invopop/gobl"
)

// KeySetVerification contains the result of verifying a KeySet's
// endorsement signatures.
type KeySetVerification struct {
	// Signed indicates whether the KeySet had any signatures.
	Signed bool
	// Authority is the GOBL Net address of the signer, extracted from the
	// gn header. Empty if the KeySet was verified against a pinned authority.
	Authority Address
	// Pinned is true if the KeySet was verified against a built-in
	// or client-configured authority key.
	Pinned bool
}

// VerifyEnvelope performs remote verification of a signed GOBL envelope.
// It extracts the "gn" header from the first signature to derive the JWKS URL,
// fetches the key set, finds the signing key by kid, and verifies the
// signature and header contents.
//
// If addr is provided (non-empty), it overrides the gn header in the signature.
func (c *Client) VerifyEnvelope(ctx context.Context, env *gobl.Envelope, addr Address) error {
	if !env.Signed() {
		return fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}

	sig := env.Signatures[0]

	// Determine address: explicit parameter or from gn header
	if addr == "" {
		gn := sig.GN()
		if gn == "" {
			return ErrNoGNHeader
		}
		var err error
		addr, err = ParseAddress(gn)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrVerifyFailed, err)
		}
	}

	kid := sig.KeyID()
	if kid == "" {
		return fmt.Errorf("%w: signature has no key ID", ErrVerifyFailed)
	}

	pubKey, err := c.FetchPublicKey(ctx, addr, kid)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	if err := env.Verify(pubKey); err != nil {
		return fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	return nil
}

// VerifyKeySet checks the endorsement signatures on a KeySet.
//
// The verification flow:
//  1. If the KeySet is not signed, returns {Signed: false} (not an error).
//  2. Tries to verify against pinned authority keys (built-in + client-configured).
//  3. If no pinned key matches, extracts the "gn" header from the signature,
//     fetches that authority's KeySet, and verifies (max 1 hop, no further recursion).
func (c *Client) VerifyKeySet(ctx context.Context, ks *KeySet) (*KeySetVerification, error) {
	if !ks.Signed() {
		return &KeySetVerification{Signed: false}, nil
	}

	// Try pinned authority keys first
	if len(c.authorities) > 0 {
		if err := ks.Verify(c.authorities...); err == nil {
			return &KeySetVerification{
				Signed: true,
				Pinned: true,
			}, nil
		}
	}

	// Extract authority address from the first signature's gn header
	sig := ks.Signatures[0]
	gn := sig.GN()
	if gn == "" {
		return nil, fmt.Errorf("%w: KeySet signature has no gn header and no pinned authority matched", ErrVerifyFailed)
	}

	authorityAddr, err := ParseAddress(gn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	kid := sig.KeyID()
	if kid == "" {
		return nil, fmt.Errorf("%w: KeySet signature has no key ID", ErrVerifyFailed)
	}

	// Fetch the authority's KeySet (1 hop, no further recursion)
	authorityKS, err := c.FetchKeySet(ctx, authorityAddr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	authorityKey, err := authorityKS.Key(kid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	if err := ks.Verify(authorityKey); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}

	return &KeySetVerification{
		Signed:    true,
		Authority: authorityAddr,
		Pinned:    false,
	}, nil
}
