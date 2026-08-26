package mt

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
			rules.Field("supplier",
				// Cap. 406 art. 50(5) and Twelfth Schedule item 3(c): a tax invoice
				// must show the supplier's VAT identification number.
				rules.Assert("01", "invoice supplier in Malta must have a VAT tax ID code",
					org.PartyHasTaxIDCode(),
				),
			),
		),
	)
}
