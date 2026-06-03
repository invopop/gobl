package net

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules/is"
)

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

// Address represents a GOBL Net address, which is a fully qualified
// domain name (FQDN) used for key discovery and network identification.
type Address string

// ParseAddress validates and returns an Address from a string.
// The input must be a valid FQDN (no scheme, no path, no port).
func ParseAddress(fqdn string) (Address, error) {
	fqdn = strings.TrimSpace(fqdn)
	fqdn = strings.ToLower(fqdn)
	if fqdn == "" {
		return "", ErrAddressEmpty
	}
	// Strip trailing dot if present (DNS canonical form)
	fqdn = strings.TrimSuffix(fqdn, ".")
	if !is.DNSName.Check(fqdn) {
		return "", fmt.Errorf("%w: %q", ErrAddressInvalid, fqdn)
	}
	// Must have at least two labels (e.g., "example.com")
	if strings.Count(fqdn, ".") < 1 {
		return "", fmt.Errorf("%w: must be a fully qualified domain name", ErrAddressInvalid)
	}
	return Address(fqdn), nil
}

// String returns the FQDN string.
func (a Address) String() string {
	return string(a)
}

// URI returns the address as a gobl: scheme cbc.URI, e.g.
// "gobl:acme.example.com", suitable for a signature's iss/aud. The
// scheme labels the identity as a GOBL Net address rather than a
// generic HTTPS service.
func (a Address) URI() cbc.URI {
	return cbc.URI(Scheme + ":" + string(a))
}

// JWKSURL returns the deterministic JWK Set discovery URL for this
// address. The matching JOSE `jku` header on a signature points here
// so generic JWT verifiers can fetch the public keys automatically.
func (a Address) JWKSURL() string {
	return "https://" + string(a) + JWKSPath
}

// KeyURL returns the deterministic discovery URL for a single public
// key (by kid) published by this address.
func (a Address) KeyURL(kid string) string {
	return "https://" + string(a) + KeyPath(kid)
}

// WhoURL returns the deterministic identity (who) URL for this address.
func (a Address) WhoURL() string {
	return "https://" + string(a) + WhoPath
}

// InboxURL returns the deterministic inbox URL for this address.
func (a Address) InboxURL() string {
	return "https://" + string(a) + InboxPath
}

// Topic reverses the FQDN labels to produce a notification topic string.
// For example, "billing.invopop.com" becomes "com.invopop.billing".
func (a Address) Topic() string {
	parts := strings.Split(string(a), ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}

// Validate checks that the address is a valid FQDN.
func (a Address) Validate() error {
	_, err := ParseAddress(string(a))
	return err
}
