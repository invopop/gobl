package bs

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Bahamas TIN is 9 digits
var tinRegexp = regexp.MustCompile(`^[0-9]{9}$`)

// normalizeTaxIdentity gets rid of occassional hyphens
// in TINs.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.Assert("01", "tax id code must be a valid TIN",
					is.MatchesRegexp(tinRegexp)),
			),
		),
	)
}
