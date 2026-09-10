package pe

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
			// Every Peruvian invoice is issued under the supplier's RUC.
			// PartyHasTaxIDCode also rejects a tax identity present with
			// an empty code, which a nil-check alone would let through.
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier tax ID code required for Peruvian regime",
					org.PartyHasTaxIDCode(),
				),
				// Calculate normally derives the regime from this country; the
				// explicit check also protects documents that declare PE directly.
				rules.Field("tax_id",
					rules.Field("country",
						rules.Assert("02", "invoice supplier tax ID country must be PE for Peruvian regime",
							is.In(CountryCode),
						),
					),
				),
			),
		),
	)
}
