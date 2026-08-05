package be

import (
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
			is.InContext(tax.RegimeIn(l10n.BE.Tax())),
			rules.Field("supplier",
				rules.When(
					is.Not(org.PartyHasIdentityTypeIn(IdentityTypeBCE)),
					rules.Field("tax_id",
						rules.Assert("01", "supplier tax ID required for Belgian regime", is.Present),
						rules.Field("code",
							rules.Assert("02", "supplier tax ID code required for Belgian regime", is.Present),
						),
					),
				),
				rules.When(
					is.Not(org.PartyHasTaxIDCode()),
					rules.Field("identities",
						rules.Assert("03", "supplier identities must include BCE type",
							org.IdentitiesTypeIn(IdentityTypeBCE)),
					),
				),
			),
		),
	)
}
