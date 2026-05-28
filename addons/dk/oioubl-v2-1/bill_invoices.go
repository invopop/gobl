package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
)

// billInvoiceRules returns the OIOUBL 2.1 invoice rule set.
//
// Rules are translated from the OIOUBL Schematron v1.17.1; each assertion
// cites its OIOUBL rule ID (e.g. F-LIB401) so failures map back to the
// canonical spec. ID_Error_List.csv in the schematron release lists all
// 1715 rule IDs across the 18 doctypes.
//
// Rules are added incrementally — this skeleton compiles and registers
// the empty set so downstream wiring (gobl.ubl Context, apps/unimaze)
// can target the addon while specific rules land in follow-up commits.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice))
}
