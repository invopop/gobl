package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
)

// Outlays holds an array of Outlay objects used inside a billing document.
type Outlays []*Outlay

// Outlay represents a expense that was paid for by the supplier and invoiced separately
// by the third party directly to the customer.
// Most suppliers will want to include the expenses of their providers as part of their
// own operational costs. However, outlays are common in countries like Spain, for example
// when an accountant or lawyer will pay for notary fees, but forward the invoice to the
// customer.
type Outlay struct {
	// Unique identity for this outlay.
	UUID string `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Outlay number index inside the invoice for ordering.
	Index int `json:"i" jsonschema:"title=Index"`
	// A code, invoice number, or other reference detail used to identify the outlay.
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Details on what the outlay was.
	Description string `json:"desc" jsonschema:"title=Description"`
	// Amount paid by the supplier.
	Paid num.Amount `json:"paid" jsonschema:"title=Paid"`
}

// Validate ensures the outlay contains everything required.
func (o *Outlay) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Index, validation.Required),
		validation.Field(&o.Description, validation.Required),
		validation.Field(&o.Paid, validation.Required),
	)
}
