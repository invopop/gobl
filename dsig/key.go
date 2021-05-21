package dsig

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/square/go-jose/v3"
)

const defaultKeyUse = "sig"

// The crypto/elliptic package doesn't provide constants for this.
const (
	curveAlgorithmP256 = "P-256"
)

// Key wraps around the underlying JSON Web Key to simplify and
// standardise the experience.
type Key struct {
	jwk *jose.JSONWebKey
}

// NewES256Key provides a new ECDSA 256 bit private key and assign it
// an ID.
func NewES256Key() *Key {
	pubCurve := elliptic.P256()
	pk, _ := ecdsa.GenerateKey(pubCurve, rand.Reader)
	return newKey(pk, string(jose.ES256))
}

func newKey(pk interface{}, alg string) *Key {
	k := new(Key)
	k.jwk = new(jose.JSONWebKey)
	k.jwk.Key = pk
	k.jwk.Algorithm = alg
	k.jwk.Use = defaultKeyUse
	k.jwk.KeyID = uuid.Must(uuid.NewRandom()).String()
	return k
}

// ID provides the key's UUID
func (k *Key) ID() string {
	return k.jwk.KeyID
}

// signatureAlgorithm attempts to determine the key's algorithm based on the
// key fields. This is a bit more reliable than depending on the
// optional `alg` property. Algorithm names provided match those
// required for signatures. Anything not defined here will not be supported
// for the time being.
func (k *Key) signatureAlgorithm() (jose.SignatureAlgorithm, error) {
	if pk, ok := k.jwk.Key.(*ecdsa.PrivateKey); ok {
		switch pk.Params().Name {
		case curveAlgorithmP256:
			return jose.ES256, nil
		}
	}
	return "", errors.New("unrecognized key signature algorithm")
}

// IsPublic returns true if this key only contains the public part.
func (k *Key) IsPublic() bool {
	return k.jwk.IsPublic()
}

// Valid let's us know if the key was generated correctly.
func (k *Key) Valid() bool {
	if k.jwk == nil {
		return false
	}
	if k.ID() == "" {
		return false
	}
	return k.jwk.Valid()
}

// Public provides the public counterpart of a private key. If this
// key is already public, nil is provided.
func (k *Key) Public() *Key {
	if k.IsPublic() {
		return nil
	}
	pk := new(Key)
	jwk := k.jwk.Public()
	pk.jwk = &jwk
	return pk
}

// Thumbprint returns teh SHA256 hex string of this key's thumbprint.
// Extremely useful for quickly checking that two keys are the same.
func (k *Key) Thumbprint() string {
	d, err := k.jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", d)
}

// MarshalJSON provides the JSON version of the key.
func (k *Key) MarshalJSON() ([]byte, error) {
	if !k.Valid() {
		return []byte{}, errors.New("cannot marshal invalid key")
	}
	return k.jwk.MarshalJSON()
}

// UnmarshalJSON parses the provided key. If parsing the key is
// successful but the key is still invalid, like for example if
// it doesn't contain an ID, it will be rejected.
func (k *Key) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if k.jwk == nil {
		k.jwk = new(jose.JSONWebKey)
	}
	if err := k.jwk.UnmarshalJSON(data); err != nil {
		return err
	}
	if !k.Valid() {
		return errors.New("invalid key")
	}
	return nil
}
