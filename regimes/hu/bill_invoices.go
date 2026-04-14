package hu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules defines the Hungarian invoice validation rules.
// Supplier adószám is required per Section 169 of Act CXXVII of 2007.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
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
