package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/dsig"
)

// KeySet extends a standard JSON Web Key Set (RFC 7517) with an optional
// digest and signatures. This allows a third party (e.g., a GOBL Net authority)
// to vouch for the authenticity of the keys by signing the digest.
//
// Standard JWKS consumers see {"keys":[...]} and ignore the extra fields.
// GOBL-aware consumers can additionally verify endorsement signatures.
type KeySet struct {
	// Standard JWKS keys array.
	Keys []jose.JSONWebKey `json:"keys"`
	// Digest of the canonical JSON of the keys array.
	Digest *dsig.Digest `json:"dig,omitempty"`
	// Signatures over the digest, typically from a trusted authority.
	Signatures []*dsig.Signature `json:"sigs,omitempty"`
}

// NewKeySet creates a new KeySet from the provided keys and calculates
// the digest.
func NewKeySet(keys ...jose.JSONWebKey) (*KeySet, error) {
	ks := &KeySet{
		Keys: keys,
	}
	if err := ks.Calculate(); err != nil {
		return nil, err
	}
	return ks, nil
}

// Calculate refreshes the digest from the current keys.
func (ks *KeySet) Calculate() error {
	d, err := ks.computeDigest()
	if err != nil {
		return err
	}
	ks.Digest = d
	return nil
}

// computeDigest generates a SHA256 digest of the canonical JSON
// representation of the keys array.
func (ks *KeySet) computeDigest() (*dsig.Digest, error) {
	cd, err := c14n.MarshalJSON(ks.Keys)
	if err != nil {
		return nil, fmt.Errorf("net: canonical JSON error: %w", err)
	}
	return dsig.NewSHA256Digest(cd), nil
}

// Sign calculates the digest and signs it with the provided private key.
// The signature is appended to the Signatures array. Optional signer
// options (e.g., dsig.WithGN) can be provided.
func (ks *KeySet) Sign(key *dsig.PrivateKey, opts ...dsig.SignerOption) error {
	if err := ks.Calculate(); err != nil {
		return err
	}
	sig, err := dsig.NewSignature(key, ks.Digest, opts...)
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}
	ks.Signatures = append(ks.Signatures, sig)
	return nil
}

// Signed returns true if the KeySet has signatures.
func (ks *KeySet) Signed() bool {
	return len(ks.Signatures) > 0
}

// Verify checks that all signatures are valid against at least one of the
// provided public keys, and that the signed digest matches the current
// keys. If no keys are provided, only the digest consistency is checked.
func (ks *KeySet) Verify(keys ...*dsig.PublicKey) error {
	if len(ks.Signatures) == 0 {
		return errors.New("net: no signatures to verify")
	}

	// Recompute digest to ensure keys haven't been tampered with
	currentDigest, err := ks.computeDigest()
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}

	var msgs []string
	for i, sig := range ks.Signatures {
		if err := ks.verifySignature(sig, currentDigest, keys...); err != nil {
			msgs = append(msgs, "sigs["+strconv.Itoa(i)+"]: "+err.Error())
		}
	}
	if len(msgs) > 0 {
		return fmt.Errorf("net: %s", strings.Join(msgs, "; "))
	}
	return nil
}

func (ks *KeySet) verifySignature(sig *dsig.Signature, currentDigest *dsig.Digest, keys ...*dsig.PublicKey) error {
	if len(keys) == 0 {
		// No keys provided, only check digest consistency
		d := new(dsig.Digest)
		if err := sig.UnsafePayload(d); err != nil {
			return errors.New("invalid signature payload")
		}
		return currentDigest.Equals(d)
	}
	for _, k := range keys {
		d := new(dsig.Digest)
		if err := sig.VerifyPayload(k, d); err != nil {
			continue
		}
		return currentDigest.Equals(d)
	}
	return errors.New("no key match found")
}

// Key finds a key by its key ID and returns it as a dsig.PublicKey.
func (ks *KeySet) Key(kid string) (*dsig.PublicKey, error) {
	for _, k := range ks.Keys {
		if k.KeyID == kid {
			pub := k.Public()
			return dsig.NewPublicKey(pub)
		}
	}
	return nil, fmt.Errorf("%w: kid=%q", ErrKeyNotFound, kid)
}

// MarshalJSON provides custom JSON marshaling to ensure the keys array
// is always present (never null).
func (ks *KeySet) MarshalJSON() ([]byte, error) {
	type Alias KeySet
	a := &struct {
		*Alias
	}{
		Alias: (*Alias)(ks),
	}
	if a.Keys == nil {
		a.Keys = []jose.JSONWebKey{}
	}
	return json.Marshal(a)
}
