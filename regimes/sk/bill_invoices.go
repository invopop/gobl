package sk

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

/*
	A Slovak invoice carries the supplier's DIČ, or the IČ DPH of a VAT payer.

	Source: https://www.slov-lex.sk/ezbierky/pravne-predpisy/SK/ZZ/2004/222/
	(Act 222/2004 Z. z. on VAT, § 74(1)(a))
*/

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				// PartyHasTaxIDCode also rejects a tax identity with an empty code
				rules.Assert("01", "invoice supplier tax ID code (DIČ or IČ DPH) required for Slovak regime",
					org.PartyHasTaxIDCode(),
				),
				// Calculate derives the regime from this country, the check covers documents declaring SK directly
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
