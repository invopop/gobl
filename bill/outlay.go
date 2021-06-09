package bill

import "github.com/invopop/gobl/num"

// Outlays holds an array of Outlay objects used inside a billing document.
type Outlays []*Outlay

// Outlay represents a expense that was paid for by the supplier and invoiced separately
// by the third party directly to the customer.
// Most suppliers will want to include the expenses of their providers as part of their
// own operational costs. However, outlays are common in countries like Spain, for example
// when an accountant or lawyer will pay for notary fees, but forward the invoice to the
// customer.
type Outlay struct {
	UUID        string     `json:"uuid,omitempty"`
	Index       int        `json:"index" jsonschema:"title=Index,description=Line number inside the invoice, starting from 0."`
	Ref         string     `json:"ref,omitempty" jsonschema:"title=Reference,description=A code, invoice number, or other reference detail used to identify the outlay."`
	Description string     `json:"desc" jsonschema:"title=Description,description=Details on what the outlay was."`
	Paid        num.Amount `json:"paid" jsonschema:"title=Paid,description=Amount paid by the supplier."`
}
