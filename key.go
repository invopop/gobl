package gobl

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
	jose "gopkg.in/square/go-jose.v2"
)

// Key encapsulates JSON Web Key handling to make it easier for us
// to build a parse keys expected for use with GoBL.
type Key struct {
	jwk *jose.JSONWebKey
}

// NewECDSAKey generates and instantiates a new private key using
// the eliptic curve algorithm.
func NewECDSAKey(id string) (*Key, error) {
	k := new(Key)
	curve := elliptic.P256()
	pkey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating ecdsa key: %w", err)
	}

	k.jwk = new(jose.JSONWebKey)
	k.jwk.Key = pkey
	k.jwk.Use = "sig"

	if id == "" {
		id = uuid.Must(uuid.NewRandom()).String()
	}
	k.jwk.KeyID = id

	return k, nil
}

// ID return's the Key's ID
func (k *Key) ID() string {
	return k.jwk.KeyID
}

// IsPublic returns true if this Key only contains the public
// component.
func (k *Key) IsPublic() bool {
	return k.jwk.IsPublic()
}

// Public returns the public version of this private key.
// If this is already a public key, we'll just return this
// object.
func (k *Key) Public() *Key {
	if k.IsPublic() {
		return k
	}
	pk := new(Key)
	jwkp := k.jwk.Public()
	pk.jwk = &jwkp
	return pk
}

// Valid checks that the key contains the correct contents.
func (k *Key) Valid() bool {
	if k.jwk == nil {
		return false
	}
	if k.jwk.KeyID == "" {
		return false
	}
	return k.jwk.Valid()
}

// Thumbprint returns a hexidecimal SHA256 string of the key's
// thumbprint, useful for comparing keys.
func (k *Key) Thumbprint() (string, error) {
	d, err := k.jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("key thumbprint: %w", err)
	}
	return fmt.Sprintf("%x", d), nil
}

// JSONWebKey provides access to the raw underlying web key implementation.
// If you are using this, there is a chance you're doing something wrong.
func (k *Key) JSONWebKey() *jose.JSONWebKey {
	return k.jwk
}

// MarshalJSON takes this key and converts it into a JSON version.
func (k *Key) MarshalJSON() ([]byte, error) {
	if k.jwk == nil {
		return nil, nil
	}
	return k.jwk.MarshalJSON()
}

// UnmarshalJSON reads the provided raw data and builds our key.
func (k *Key) UnmarshalJSON(data []byte) error {
	if k.jwk == nil {
		k.jwk = new(jose.JSONWebKey)
	}
	return k.jwk.UnmarshalJSON(data)
}
