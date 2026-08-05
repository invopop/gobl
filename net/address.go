package net

import (
	"fmt"
	"strings"

	"golang.org/x/net/idna"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules/is"
)

// addressIDNA is the Lookup profile (the strict IDN form used for DNS
// lookups). ToASCII converts U-Labels to A-Labels, lowercases ASCII
// labels, and rejects labels that don't satisfy the IDNA2008 lookup
// rules — exactly the canonicalization GOBL Net needs.
var addressIDNA = idna.New(
	idna.MapForLookup(),
	idna.Transitional(false),
	idna.StrictDomainName(true),
)

// Address represents a GOBL Net address, which is a fully qualified
// domain name (FQDN) used for key discovery and network identification.
type Address string

// ParseAddress validates and returns an Address from a string.
// The input must be a valid FQDN (no scheme, no path, no port).
//
// Internationalised domain names (IDN) in U-Label form are accepted
// and normalised to their ASCII (A-Label / Punycode) representation,
// so the canonical form on the wire and in `iss`/`aud` is always
// ASCII. `München.DE` and `xn--mnchen-3ya.de` parse to the same
// Address.
func ParseAddress(fqdn string) (Address, error) {
	fqdn = strings.TrimSpace(fqdn)
	if fqdn == "" {
		return "", ErrAddressEmpty
	}
	// Strip trailing dot if present (DNS canonical form)
	fqdn = strings.TrimSuffix(fqdn, ".")
	// idna.ToASCII lowercases ASCII labels, converts any U-Labels to
	// A-Labels, and rejects labels that don't satisfy the IDNA2008
	// lookup rules.
	ascii, err := addressIDNA.ToASCII(fqdn)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAddressInvalid, err)
	}
	if !is.DNSName.Check(ascii) {
		return "", fmt.Errorf("%w: %q", ErrAddressInvalid, ascii)
	}
	// Must have at least two labels (e.g., "example.com")
	if strings.Count(ascii, ".") < 1 {
		return "", fmt.Errorf("%w: must be a fully qualified domain name", ErrAddressInvalid)
	}
	return Address(ascii), nil
}

// String returns the FQDN string.
func (a Address) String() string {
	return string(a)
}

// URI returns the address as a gobl: scheme cbc.URI, e.g.
// "gobl:acme.example.com", for use where multiple schemes coexist —
// org.Endpoint lists and the envelope header's unsigned from/to
// routing fields. Signed iss/aud/verifier claims carry the bare
// address (Address.String) instead: within the protocol they can
// only be GOBL Net addresses, so no scheme is needed.
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

// Validate checks that the address is a valid FQDN.
func (a Address) Validate() error {
	_, err := ParseAddress(string(a))
	return err
}
