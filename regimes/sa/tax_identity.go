package sa

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// vatIDRegexp matches a Saudi VAT Identification Number: 15 digits starting and ending with 3.
var vatIDRegexp = regexp.MustCompile(`^3[0-9]{13}3$`)

// normalizeTaxIdentity strips whitespace, separators and country prefix from a
// Saudi VAT Identification Number.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.Assert("01", "invoice tax id code must be 15 digits long, and start and end with 3",
					is.MatchesRegexp(vatIDRegexp)),
			),
		),
	)
}
