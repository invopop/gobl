package pe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			// Every Peruvian invoice is issued under the supplier's RUC.
			rules.Field("supplier",
				rules.Field("tax_id",
					rules.Assert("01", "invoice supplier tax ID required for Peruvian regime", is.Present),
				),
			),
		),
	)
}
