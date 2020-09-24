package gobl

import (
	jose "gopkg.in/square/go-jose.v2"
)

// Signature provides a JSON Web Signature implementation that
// will always serialize and parse the signature into a
// compact form.
type Signature struct {
	sig *jose.JSONWebSignature
}

// NewSignature helps build a new signature of the data
func NewSignature(k *Key, data interface{}) *Signature {
	s := new(Signature)
	return s
}
