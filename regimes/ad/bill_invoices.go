package ad

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules defines Andorran invoice-level requirements.
// Suppliers must have an NRT tax ID code to issue IGI invoices.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", "supplier must have an NRT tax ID code",
					is.Func("has NRT tax ID", func(value any) bool {
						party, _ := value.(*org.Party)
						return party != nil && party.TaxID != nil && party.TaxID.Code != ""
					}),
				),
			),
		),
	)
}
