package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier in Singapore must have a GST tax ID code or a UEN identity",
					is.AnyOf(
						org.PartyHasTaxIDCode(),
						org.PartyHasIdentityTypeIn(IdentityTypeUEN),
					),
				),
			),
		),
	)
}
