package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
)

// Outlays holds an array of Outlay objects used inside a billing document.
type Outlays []*Outlay

// Outlay represents a reimbursable expense that was paid for by the supplier and invoiced separately
// by the third party directly to the customer.
// Most suppliers will want to include the expenses of their providers as part of their
// own operational costs. However, outlays are common in countries like Spain where it is typical
// for an accountant or lawyer to pay for notary fees, but forward the invoice to the
// customer.
type Outlay struct {
	// Unique identity for this outlay.
	UUID string `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Outlay number index inside the invoice for ordering.
	Index int `json:"i" jsonschema:"title=Index"`
	// When was the outlay made.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// Invoice number or other reference detail used to identify the outlay.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Series of the outlay invoice.
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Details on what the outlay was.
	Description string `json:"desc" jsonschema:"title=Description"`
	// Who was the supplier of the outlay
	Supplier *org.Party `json:"supplier,omitempty" jsonschema:"title=Supplier"`
	// Amount paid by the supplier.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// Validate ensures the outlay contains everything required.
func (o *Outlay) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Index, validation.Required),
		validation.Field(&o.Description, validation.Required),
		validation.Field(&o.Amount, validation.Required),
	)
}
