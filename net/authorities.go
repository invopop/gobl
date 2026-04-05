package net

import "github.com/invopop/gobl/dsig"

// Authorities contains the built-in trusted authority public keys used
// as the first check when verifying KeySet signatures. Initially empty;
// populate via RegisterAuthority or WithAuthorities on Client.
var Authorities []*dsig.PublicKey

// RegisterAuthority adds a public key to the global set of trusted
// authority keys.
func RegisterAuthority(key *dsig.PublicKey) {
	Authorities = append(Authorities, key)
}
