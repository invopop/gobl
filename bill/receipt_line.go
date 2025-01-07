package bill

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// ReceiptLine defines the details of a line required in an invoice.
type ReceiptLine struct {
	uuid.Identify
	// Line number inside the parent (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`

	// Direction for the flow of money, either debit (+) or credit (-) from the perspective
	// of the supplier's asset account. Debiting an asset account increases its value, implying
	// incoming payment. Crediting an asset account decreases its value, implying outgoing payment.
	Type cbc.Key `json:"type" jsonschema:"title=Type" jsonschema_extras:"calculated=true"`

	// The document reference related to the payment.
	Document *org.DocumentRef `json:"document" jsonschema:"title=Document"`

	// Tax total breakdown from the original document, if required and available.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax"`

	// Total amount to be paid for this line.
	Total num.Amount `json:"total" jsonschema:"title=Total"`
}

// ValidateWithContext ensures that the fields contained in the ReceiptLine look correct.
func (rl *ReceiptLine) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, rl,
		validation.Field(&rl.Type, validation.Required, isValidReceiptLineType),
		validation.Field(&rl.Document, validation.Required),
		validation.Field(&rl.Tax),
		validation.Field(&rl.Total, validation.Required),
	)
}

// JSONSchemaExtend extends the schema with additional property details
func (rl ReceiptLine) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(ReceiptLineTypes))
		for i, kd := range InvoiceTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
}
