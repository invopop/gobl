package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
)

// billInvoiceRules returns the OIOUBL 2.1 rule set for bill.Invoice
// (covers both invoices and credit notes).
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice))
}
