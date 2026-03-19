package pl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.HasContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Field("tax_id",
					rules.Field("code",
						rules.Assert("01", "supplier tax ID code required",
							is.Present,
						),
					),
				),
			),
		),
	)
}
