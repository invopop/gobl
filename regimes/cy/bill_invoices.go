package cy

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
			rules.Field("supplier",
				rules.Field("tax_id",
					rules.Assert("01", "supplier tax ID is required", is.Present),
					rules.Field("code",
						rules.Assert("02", "supplier tax ID code is required", is.Present),
					),
				),
				rules.Field("addresses",
					rules.Assert("03", "supplier address is required", is.Present),
				),
			),
		),
	)
}
