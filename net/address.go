package net

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/rules/is"
)

const (
	// WellKnownPath is the base path for GOBL Net well-known URLs.
	WellKnownPath = "/.well-known"
	// JWKSPath is the full well-known path for JWKS discovery.
	JWKSPath = WellKnownPath + "/jwks.json"
)

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

// JWKSURL returns the deterministic JWKS discovery URL for this address.
func (a Address) JWKSURL() string {
	return "https://" + string(a) + JWKSPath
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
