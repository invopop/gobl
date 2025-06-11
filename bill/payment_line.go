package bill

import (
	"context"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// PaymentLine defines the details of a line item in a payment document.
type PaymentLine struct {
	uuid.Identify

	// Line number within the parent document (automatically calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`

	// Reference to the document being paid
	Document *org.DocumentRef `json:"document,omitempty" jsonschema:"title=Document"`

	// When making multiple payments for a single document, this specifies the
	// installment number for this payment line.
	Installment int `json:"installment,omitempty" jsonschema:"title=Installment"`

	// Amount already paid in previous installments, which may be required
	// by some tax regimes or specific use cases.
	Advances *num.Amount `json:"advances,omitempty" jsonschema:"title=Advances"`

	// Amount of the total payment allocated to the referenced document.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`

	// Additional notes specific to this line item for clarification purposes
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// ValidateWithContext ensures that the fields contained in the PaymentLine look correct.
func (pl *PaymentLine) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, pl,
		validation.Field(&pl.Document),
		validation.Field(&pl.Installment, validation.Min(1), validation.Max(999)),
		validation.Field(&pl.Advances, num.Max(pl.Amount)),
		validation.Field(&pl.Amount),
		validation.Field(&pl.Notes),
	)
}
