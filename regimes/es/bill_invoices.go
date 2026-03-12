package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			tax.RegimeIn("ES"),
			rules.Field("supplier",
				rules.Assert("01", "supplier is required", rules.Required),
				rules.Field("tax_id",
					rules.Assert("02", "supplier tax ID is required", rules.Required),
					rules.Field("code",
						rules.Assert("03", "supplier tax ID must have a code", rules.Required),
					),
				),
			),
		),
	)
}
