package sa

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Saudi VAT Identification Number is a 15-digit number.
var vatIDPattern = `^3[0-9]{13}3$`

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("SA"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Saudi tax identity code",
					is.Matches(vatIDPattern),
				),
			),
		),
	)
}
