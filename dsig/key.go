package dsig

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/invopop/gobl/cal"
)

const defaultKeyUse = "sig"

// The crypto/elliptic package doesn't provide constants for this.
const (
	curveAlgorithmP256 = "P-256"
)

// PrivateKey makes it easy to deal with private keys used to sign data
// and created signatures.
// These should obviously be kept secure and be used to generate the public
// keys.
type PrivateKey struct {
	jwk *jose.JSONWebKey
}

// PublicKey is generated from the private key and can be shared freely
// as it cannot be used to create signatures.
//
// In addition to the standard JWK members carried in the embedded JWK,
// a PublicKey may declare a validity window via ValidFrom and
// ValidUntil. Both are optional; when set they bound the signing time
// any signature produced by this key is allowed to carry (see Allows
// and the signed `ts` in head.SigningPayload). The fields serialise as
// the RFC 7517 §4 extension members `valid_from` / `valid_until`,
// which JOSE consumers that do not recognise them MUST ignore.
type PublicKey struct {
	jwk *jose.JSONWebKey

	// ValidFrom is the earliest time at which this key may sign. nil
	// means no lower bound.
	ValidFrom *cal.Timestamp
	// ValidUntil is the latest time at which this key may sign. nil
	// means no upper bound; a value in the past indicates a retired
	// key whose historical signatures still verify.
	ValidUntil *cal.Timestamp
}

// NewPublicKey creates a PublicKey from a jose.JSONWebKey.
// The key must be a valid public key.
func NewPublicKey(jwk jose.JSONWebKey) (*PublicKey, error) {
	pk := &PublicKey{jwk: &jwk}
	if err := pk.Validate(); err != nil {
		return nil, err
	}
	return pk, nil
}

// NewES256Key provides a new ECDSA 256 bit private key and assigns it
// an ID.
func NewES256Key() *PrivateKey {
	pubCurve := elliptic.P256()
	pk, _ := ecdsa.GenerateKey(pubCurve, rand.Reader)
	return newKey(pk, string(jose.ES256))
}

func newKey(pk interface{}, alg string) *PrivateKey {
	k := new(PrivateKey)
	k.jwk = new(jose.JSONWebKey)
	k.jwk.Key = pk
	k.jwk.Algorithm = alg
	k.jwk.Use = defaultKeyUse
	k.jwk.KeyID = uuid.Must(uuid.NewRandom()).String()
	return k
}

// ID provides the private key's UUID
func (k *PrivateKey) ID() string {
	return k.jwk.KeyID
}

// ID provides the public key's UUID
func (k *PublicKey) ID() string {
	return k.jwk.KeyID
}

// signatureAlgorithm attempts to determine the key's algorithm based on the
// key fields. This is a bit more reliable than depending on the
// optional `alg` property. Algorithm names provided match those
// required for signatures. Anything not defined here will not be supported
// for the time being.
func (k *PrivateKey) signatureAlgorithm() (jose.SignatureAlgorithm, error) {
	if pk, ok := k.jwk.Key.(*ecdsa.PrivateKey); ok {
		switch pk.Params().Name {
		case curveAlgorithmP256:
			return jose.ES256, nil
		}
	}
	return "", errors.New("unrecognized key signature algorithm")
}

// Validate let's us know if the private key was generated or parsed correctly.
func (k *PrivateKey) Validate() error {
	if k.jwk == nil {
		return errors.New("key not set")
	}
	if k.ID() == "" {
		return errors.New("id required")
	}
	if !k.jwk.Valid() {
		return errors.New("jose key is invalid")
	}
	if k.jwk.IsPublic() {
		return errors.New("private key only contains public part")
	}
	return nil
}

// Validate let's us know if the public key was parsed correctly.
func (k *PublicKey) Validate() error {
	if k.jwk == nil {
		return errors.New("key not set")
	}
	if k.ID() == "" {
		return errors.New("id required")
	}
	if !k.jwk.Valid() {
		return errors.New("jose key is invalid")
	}
	if !k.jwk.IsPublic() {
		return errors.New("public key is private")
	}
	return nil
}

// Public provides the public counterpart of a private key, ready to be used
// to be persisted in a key store and verify signatures.
func (k *PrivateKey) Public() *PublicKey {
	pk := new(PublicKey)
	jwk := k.jwk.Public()
	pk.jwk = &jwk
	return pk
}

// Sign is a helper method that will generate a signature using the
// private key.
func (k *PrivateKey) Sign(data interface{}) (*Signature, error) {
	return NewSignature(k, data)
}

// Verify is a wrapper around the signature's VerifyPayload method for
// the sake of convenience.
func (k *PublicKey) Verify(sig *Signature, payload interface{}) error {
	return sig.VerifyPayload(k, payload)
}

// Thumbprint returns the SHA256 hex string of the private key's thumbprint.
// Extremely useful for quickly checking that two keys, either public or private,
// are the same.
func (k *PrivateKey) Thumbprint() string {
	return keyThumbprint(k.jwk)
}

// Thumbprint returns the SHA256 hex string of the public key's thumbprint.
// Extremely useful for quickly checking that two keys are the same.
func (k *PublicKey) Thumbprint() string {
	return keyThumbprint(k.jwk)
}

func keyThumbprint(jwk *jose.JSONWebKey) string {
	d, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", d)
}

// MarshalJSON provides the JSON version of the key.
func (k *PrivateKey) MarshalJSON() ([]byte, error) {
	return k.jwk.MarshalJSON()
}

// MarshalJSON emits the standard JWK fields, with `valid_from` and
// `valid_until` flattened into the same JSON object when set.
func (k *PublicKey) MarshalJSON() ([]byte, error) {
	b, err := k.jwk.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if k.ValidFrom == nil && k.ValidUntil == nil {
		return b, nil
	}
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("dsig: marshal public key: %w", err)
	}
	if k.ValidFrom != nil {
		v, err := json.Marshal(k.ValidFrom)
		if err != nil {
			return nil, err
		}
		m["valid_from"] = v
	}
	if k.ValidUntil != nil {
		v, err := json.Marshal(k.ValidUntil)
		if err != nil {
			return nil, err
		}
		m["valid_until"] = v
	}
	return json.Marshal(m)
}

// Allows reports whether the given signing time falls within this
// key's declared validity window. A nil ts (signature without a
// timestamp) skips the check; absent bounds on the key skip their
// respective half of the check.
func (k *PublicKey) Allows(ts *cal.Timestamp) error {
	if ts == nil {
		return nil
	}
	if k.ValidFrom != nil && ts.Time.Before(k.ValidFrom.Time) {
		return fmt.Errorf("dsig: signing time %s is before key's valid_from %s", ts, k.ValidFrom)
	}
	if k.ValidUntil != nil && ts.Time.After(k.ValidUntil.Time) {
		return fmt.Errorf("dsig: signing time %s is after key's valid_until %s", ts, k.ValidUntil)
	}
	return nil
}

// UnmarshalJSON parses the JSON private key data. You should perform
// validation on the key to ensure it was provided correctly.
func (k *PrivateKey) UnmarshalJSON(data []byte) error {
	if k.jwk == nil {
		k.jwk = new(jose.JSONWebKey)
	}
	return k.jwk.UnmarshalJSON(data)
}

// UnmarshalJSON parses the JSON public key data, including the
// optional `valid_from` / `valid_until` extension members. You should
// perform validation on the key to ensure it was provided correctly.
func (k *PublicKey) UnmarshalJSON(data []byte) error {
	if k.jwk == nil {
		k.jwk = new(jose.JSONWebKey)
	}
	if err := k.jwk.UnmarshalJSON(data); err != nil {
		return err
	}
	aux := struct {
		ValidFrom  *cal.Timestamp `json:"valid_from"`
		ValidUntil *cal.Timestamp `json:"valid_until"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	k.ValidFrom = aux.ValidFrom
	k.ValidUntil = aux.ValidUntil
	return nil
}
