package se

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
			is.InContext(tax.RegimeIn(l10n.SE.Tax())),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf("invoice SE supplier must have either tax ID code or identity with %s, %s, or %s type", IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr),
					is.AnyOf(
						org.PartyHasTaxIDCode(),
						org.PartyHasIdentityTypeIn(IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr),
					),
				),
			),
		),
	)
}
