package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// buyerIdentityThreshold is the total price, including GST, from which a tax
// invoice must also show the buyer's identity or ABN.
var buyerIdentityThreshold = num.MakeAmount(1000, 0)

// billInvoiceRules applies the tax invoice requirements to standard invoices
// that charge GST: the seller's ABN, and the buyer for sales of A$1,000 or more.
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
			rules.When(
				is.AllOf(
					bill.InvoiceTypeIn(bill.InvoiceTypeStandard),
					is.Func("charges GST on A$1,000 or more", invoiceMeetsBuyerThreshold),
				),
				rules.Field("customer",
					rules.Assert("02", "invoice customer is required when GST is charged on a total of A$1,000 or more",
						is.Present,
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

func invoiceMeetsBuyerThreshold(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || !invoiceChargesGST(inv) || inv.Currency != currency.AUD {
		return false
	}
	return inv.Totals.TotalWithTax.Compare(buyerIdentityThreshold) >= 0
}
