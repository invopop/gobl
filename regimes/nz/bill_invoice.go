package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			tax.RegimeIn(CountryCode),
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier in New Zealand must have a GST tax ID code",
					org.PartyHasTaxIDCode(),
				),
			),
		),
	)
}
