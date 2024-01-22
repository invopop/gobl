package pay

import (
	"context"
	"encoding/json"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Advance represents a single payment that has been made already, such
// as a deposit on an intent to purchase, or as credit from a previous
// invoice which was later corrected or cancelled.
type Advance struct {
	// Unique identifier for this advance.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// When the advance was made.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// The payment means used to make the advance.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// ID or reference for the advance.
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// If this "advance" payment has come from a public grant or subsidy, set this to true.
	Grant bool `json:"grant,omitempty" jsonschema:"title=Grant"`
	// Details about the advance.
	Description string `json:"description" jsonschema:"title=Description"`
	// How much as a percentage of the total with tax was paid
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// How much was paid.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// If different from the parent document's base currency.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
}

// Validate checks the advance looks okay
func (a *Advance) Validate() error {
	return a.ValidateWithContext(context.Background())
}

// ValidateWithContext checks the advance looks okay inside the context.
func (a *Advance) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, a,
		validation.Field(&a.Amount, validation.Required),
		validation.Field(&a.Key, HasValidMeansKey),
		validation.Field(&a.Description, validation.Required),
	)
}

// CalculateFrom will update the amount using the rate of the provided
// total, if defined.
func (a *Advance) CalculateFrom(totalWithTax num.Amount) {
	if a.Percent != nil {
		a.Amount = a.Percent.Of(totalWithTax)
	}
}

// JSONUnmarshal helps migrate the desc field to description.
func (a *Advance) JSONUnmarshal(data []byte) error {
	type Alias Advance
	aux := &struct {
		Desc string `json:"desc,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Desc != "" {
		a.Description = aux.Desc
	}
	return nil
}

// JSONSchemaExtend extends the JSONSchema for the Instructions type.
func (Advance) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithMeansKey(schema, "key")
}
