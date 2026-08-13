package fr

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
			is.InContext(tax.RegimeIn(l10n.FR.Tax())),
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier must have a tax ID code or a SIREN/SIRET identity",
					is.AnyOf(
						org.PartyHasTaxIDCode(),
						org.PartyHasIdentityTypeIn(IdentityTypeSIREN, IdentityTypeSIRET),
					),
				),
			),
		),
	)
}
