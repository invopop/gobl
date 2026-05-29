package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// billInvoiceRules returns the OIOUBL 2.1 rule set for bill.Invoice
// (covers both invoices and credit notes).
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("code",
			rules.Assert("05", "invoice code is required (F-INV009)", is.Present),
		),
		rules.Field("supplier",
			rules.Field("inboxes",
				rules.Assert("01", "supplier inboxes are required (F-INV031)", is.Present),
			),
		),
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("02", "customer inboxes are required (F-INV044)", is.Present),
			),
			// F-INV046 requires exactly one Contact in OIOUBL output;
			// gobl.ubl picks one Person at emit time, so the addon asserts presence only.
			rules.Field("people",
				rules.Assert("03", "customer people are required (F-INV046)", is.Present),
			),
		),
		rules.When(is.Func("any line has order ref", anyLineHasOrderRef),
			rules.Field("ordering",
				rules.Assert("07", "ordering is required when any line has an order reference (F-INV142)", is.Present),
			),
		),
		rules.Field("ordering",
			rules.Field("code",
				rules.Assert("04", "ordering code is required when ordering is set (F-INV024)", is.Present),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Field("quantity",
					rules.Assert("06", "line quantity must not be zero (F-INV147)", is.Func("non-zero amount", quantityNonZero)),
				),
			),
		),
	)
}

func anyLineHasOrderRef(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	for _, line := range inv.Lines {
		if !line.Order.IsEmpty() {
			return true
		}
	}
	return false
}

func quantityNonZero(val any) bool {
	switch a := val.(type) {
	case num.Amount:
		return !a.IsZero()
	case *num.Amount:
		return a == nil || !a.IsZero()
	}
	return true
}
