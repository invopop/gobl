package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// billInvoiceRules returns the OIOUBL 2.1 rule set for bill.Invoice
// (covers both invoices and credit notes).
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("01", "supplier tax ID is required (F-INV034)", is.Present),
			),
			rules.Field("inboxes",
				rules.Assert("02", "supplier inboxes are required (F-INV031)", is.Present),
			),
		),
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("03", "customer inboxes are required (F-INV044)", is.Present),
			),
			// F-INV046 requires exactly one Contact in OIOUBL output;
			// gobl.ubl picks one Person at emit time, so the addon asserts presence only.
			rules.Field("people",
				rules.Assert("04", "customer people are required (F-INV046)", is.Present),
			),
		),
		rules.Field("ordering",
			rules.Field("code",
				rules.Assert("05", "ordering code is required when ordering is set (F-INV024)", is.Present),
			),
		),
	)
}
