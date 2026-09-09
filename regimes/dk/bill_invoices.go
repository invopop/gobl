package dk

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(l10n.DK.Tax())),
			rules.Field("supplier",
				rules.When(
					is.Not(org.PartyHasTaxIDCode()),
					rules.Field("identities",
						rules.Assert("01", fmt.Sprintf(
							"invoice DK supplier without a tax ID code requires an identity with '%s' or '%s' type",
							IdentityTypeCVR, IdentityTypeCPR,
						),
							org.IdentitiesTypeIn(IdentityTypeCVR, IdentityTypeCPR),
						),
					),
				),
			),
		),
	)
}
