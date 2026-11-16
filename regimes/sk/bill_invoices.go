package sk

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
			// A Slovak invoice carries the supplier's DIČ or IČ DPH
			// (Act 222/2004 Z. z. § 74(1)(a)). PartyHasTaxIDCode also
			// rejects a tax identity present with an empty code.
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier tax ID code (DIČ or IČ DPH) required for Slovak regime",
					org.PartyHasTaxIDCode(),
				),
				// Calculate normally derives the regime from this country; the
				// explicit check also protects documents that declare SK directly.
				rules.Field("tax_id",
					rules.Field("country",
						rules.Assert("02", "invoice supplier tax ID country must be SK for Slovak regime",
							is.In(CountryCode),
						),
					),
				),
			),
		),
	)
}
