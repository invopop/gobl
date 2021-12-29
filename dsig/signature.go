package dsig

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/jsonschema"
	"github.com/square/go-jose/v3"
)

// Signature represents a stored JSON Web Signature and provides helper
// methods to be able to extract and verify contents.
type Signature struct {
	jws *jose.JSONWebSignature
}

// NewSignature instantiates a new Signature object by signing the provided
// data using the private key. The signature will use the same algorithm as
// defined by the key.
func NewSignature(key *PrivateKey, data interface{}) (*Signature, error) {
	if err := key.Validate(); err != nil {
		return nil, ErrKeyInvalid
	}

	alg, err := key.signatureAlgorithm()
	if err != nil {
		return nil, fmt.Errorf("dsig: %w", err)
	}
	sk := jose.SigningKey{
		Algorithm: alg,
		Key:       key.jwk,
	}
	signer, err := jose.NewSigner(sk, nil)
	if err != nil {
		return nil, fmt.Errorf("dsig: %w", err)
	}

	// get a JSON string of the payload
	p, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("dsig: %w", err)
	}

	s := new(Signature)
	s.jws, err = signer.Sign(p)
	if err != nil {
		return nil, fmt.Errorf("dsig: %w", err)
	}
	// correct issue in copying Key ID header
	s.jws.Signatures[0].Header.KeyID = key.ID()

	return s, nil
}

// ParseSignature converts raw signature data into an object that
// can be used to extract and validate.
func ParseSignature(data string) (*Signature, error) {
	s := new(Signature)
	err := s.parse(data)
	return s, err
}

func (s *Signature) parse(data string) error {
	o, err := jose.ParseSigned(data)
	if err != nil {
		return fmt.Errorf("dsig: %w", err)
	}
	s.jws = o
	return nil
}

// KeyID extracts the ID used to generate the signature from the
// headers.
func (s *Signature) KeyID() string {
	if s.jws == nil || len(s.jws.Signatures) == 0 {
		return ""
	}
	return s.jws.Signatures[0].Header.KeyID
}

// String provides the compact form signature.
func (s *Signature) String() string {
	if s.jws == nil {
		return ""
	}
	d, err := s.jws.CompactSerialize()
	if err != nil {
		return ""
	}
	return d
}

// VerifyPayload verifies that the provided key was indeed used to
// sign the original payload and will parse the data ready to use.
func (s *Signature) VerifyPayload(key *PublicKey, payload interface{}) error {
	data, err := s.jws.Verify(key.jwk)
	if err != nil {
		// at the risk of hiding useful errors, provide our own
		return ErrKeyMismatch
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return fmt.Errorf("dsig verify: %w", err)
	}

	return nil
}

// UnsafePayload will extract the payload data into the provided
// object but will not perform any signature checking. Only
// recommended for specific use cases when the original key is
// not available or has already been confirmed elsewhere.
func (s *Signature) UnsafePayload(payload interface{}) error {
	data := s.jws.UnsafePayloadWithoutVerification()
	if err := json.Unmarshal(data, payload); err != nil {
		return fmt.Errorf("dsig unsafe payload: %w", err)
	}
	return nil
}

// JSONWebSignature provides underlying JOSE object.
func (s *Signature) JSONWebSignature() *jose.JSONWebSignature {
	return s.jws
}

// MarshalJSON provides the compact string signature ready to be
// using as a JSON string.
func (s *Signature) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(s.String())
	if err != nil {
		return nil, fmt.Errorf("dsig: %w", err)
	}
	return data, nil
}

// UnmarshalJSON parses the compact signature string.
func (s *Signature) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("dsig: %w", err)
	}
	if len(str) == 0 {
		return nil
	}
	return s.parse(str)
}

// JSONSchemaType returns the json schema type.
func (Signature) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:        "string",
		Title:       "Signature",
		Description: "JSON Web Signature in compact form.",
	}
}
