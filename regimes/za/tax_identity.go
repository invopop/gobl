package za

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeRegexp matches South African VAT registration numbers: always 10
// digits, always starting with the digit 4 (e.g. 4480152117).
//
// Unlike most countries covered by GOBL, SARS has never published a check
// digit algorithm for this number, so format is the only thing that can be
// validated offline here. Authoritative validation requires a live lookup
// against SARS's VAT Vendor Search service (https://secure.sarsefiling.co.za/vatvendorsearch.aspx),
// which is a runtime/network concern outside the scope of this library.
//
// Reference: https://lookuptax.com/docs/how-to-verify/vat-verification-south-africa
var taxCodeRegexp = regexp.MustCompile(`^4\d{9}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid South African VAT identity code",
					is.MatchesRegexp(taxCodeRegexp),
				),
			),
		),
	)
}
