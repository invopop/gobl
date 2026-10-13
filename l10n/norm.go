package l10n

import (
	"strings"

	"github.com/invopop/gobl/norm"
)

// init registers the intrinsic normalizers for country codes. They run on
// every country code field in a document before the normalizers of the
// struct that holds it, so a regime guard such as tax.IdentityIn("DE")
// sees the uppercased code.
func init() {
	norm.Register(
		norm.For(func(c *ISOCountryCode) { *c = ISOCountryCode(normalizeCountryCode(string(*c))) }),
		norm.For(func(c *TaxCountryCode) { *c = TaxCountryCode(normalizeCountryCode(string(*c))) }),
	)
}

func normalizeCountryCode(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
