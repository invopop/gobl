package in

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
)

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Field("identities",
			rules.Assert("01", "all items must have an HSN identity code",
				org.IdentitiesTypeIn(IdentityTypeHSN),
			),
		),
	)

}
