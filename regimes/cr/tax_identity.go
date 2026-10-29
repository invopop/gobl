package cr

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Costa Rica issues four tax identity numbers, none carrying a check digit
// ("Anexos y Estructuras v4.4", campo "Número de cédula"):
//
//   - Cédula física: 9 digits, no leading zero.
//   - Cédula jurídica or NITE: 10 characters, alphanumeric for legal entities.
//   - DIMEX: 11 or 12 digits, no leading zero.
//
// Numbers are stored unpadded; the zero-padded form belongs to the clave
// numérica. "Extranjero No Domiciliado" and "No Contribuyente" are excluded, as
// the annex exempts both from registration with the tax administration.
var taxCodeRegexp = regexp.MustCompile(`^([1-9]\d{8}|[0-9A-Z]{10}|[1-9]\d{10,11})$`)

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax id code must be 9, 11 or 12 digits without a leading zero, or 10 alphanumeric characters",
					is.MatchesRegexp(taxCodeRegexp),
				),
			),
		),
	)
}
