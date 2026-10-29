package sk

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules checks the supplier's tax ID, which Slovak VAT Act
// 222/2004 § 74(1) requires on a tax invoice.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", "invoice SK supplier must have a tax ID code (DIČ or IČ DPH)",
					org.PartyHasTaxIDCode(),
				),
			),
		),
	)
}
