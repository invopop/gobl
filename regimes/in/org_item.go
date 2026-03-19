package in

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.When(
			is.HasContext(tax.RegimeIn(CountryCode)),
			rules.Field("identities",
				rules.Assert("01", "all items must have an HSN identity code",
					org.IdentitiesTypeIn(IdentityTypeHSN),
				),
			),
		),
	)

}
