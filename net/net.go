// Package net provides models and utilities for GOBL Net, a decentralized identity and communication protocol.
//
// EXPERIMENTAL: GOBL Net is under active development. Its API and wire
// protocol may change without notice and carry no stability guarantee yet.
package net

const (
	// Scheme is the URI scheme used by GOBL Net addresses, e.g.
	// "gobl:acme.example.com".
	Scheme = "gobl"

	// WellKnownPath is the base path for GOBL Net well-known URLs.
	WellKnownPath = "/.well-known/gobl"
	// KeysPath is the base of the per-key endpoint; the full path for a
	// single key is KeysPath + "/" + kid.
	KeysPath = WellKnownPath + "/keys"
	// WhoPath is the well-known path serving the signed Party envelope.
	WhoPath = WellKnownPath + "/who"
	// InboxPath is the well-known path accepting envelope deliveries.
	InboxPath = WellKnownPath + "/inbox"
	// JWKSPath is the bulk JWK Set endpoint published at the root
	// well-known directory so generic JWT tooling (jwt.io, OIDC-style
	// verifiers) can resolve `jku` and verify signatures without
	// out-of-band key exchange.
	JWKSPath = "/.well-known/jwks.json"
)

// KeyPath returns the well-known path serving a single public key by
// its key ID. Use this to construct lookup URLs.
func KeyPath(kid string) string {
	return KeysPath + "/" + kid
}
