package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn("ES")),
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier in Spain is required", is.Present),
				rules.Field("tax_id",
					rules.Assert("02", "invoice supplier tax ID in Spain is required", is.Present),
					rules.Field("code",
						rules.Assert("03", "invoice supplier tax ID code in Spain is required", is.Present),
					),
				),
			),
		),
	)
}
