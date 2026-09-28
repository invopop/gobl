package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules applies the tax invoice requirement to standard invoices
// that charge GST: the seller's ABN.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				is.AllOf(
					bill.InvoiceTypeIn(bill.InvoiceTypeStandard),
					is.Func("charges GST", invoiceChargesGST),
				),
				rules.Field("supplier",
					rules.Assert("01", "invoice supplier must have an ABN when GST is charged",
						org.PartyHasTaxIDCode(),
					),
				),
			),
		),
	)
}

func invoiceChargesGST(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Totals == nil {
		return false
	}
	ct := inv.Totals.Taxes.Category(tax.CategoryGST)
	return ct != nil && !ct.Amount.IsZero()
}
