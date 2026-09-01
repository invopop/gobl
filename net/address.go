package net

import (
	"fmt"
	"strings"

	"golang.org/x/net/idna"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules/is"
)

// addressIDNA canonicalizes and validates names using the IDNA lookup profile.
var addressIDNA = idna.New(
	idna.MapForLookup(),
	idna.Transitional(false),
	idna.StrictDomainName(true),
)

// Address is a GOBL Net participant's fully qualified domain name.
type Address string

// ParseAddress validates an FQDN and returns its canonical ASCII form.
// It accepts internationalized names, trims whitespace and a trailing dot,
// and rejects schemes, ports, paths, and single-label names.
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

// URI returns the address as a gobl URI for party endpoints and routing fields.
// Signed iss, aud, and verifier claims use String instead.
func (a Address) URI() cbc.URI {
	return cbc.URI(Scheme + ":" + string(a))
}

// JWKSURL returns the conventional bulk JWK Set URL for the address.
// GOBL Net verification uses KeyURL; signatures do not carry jku.
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
