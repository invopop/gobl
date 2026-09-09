package de

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
			is.InContext(tax.RegimeIn(l10n.DE.Tax())),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf("invoice DE supplier must have either tax ID code or identity with '%s' key", IdentityKeyTaxNumber),
					is.AnyOf(
						org.PartyHasTaxIDCode(),
						org.PartyHasIdentityKeyIn(IdentityKeyTaxNumber),
					),
				),
			),
		),
	)
}
