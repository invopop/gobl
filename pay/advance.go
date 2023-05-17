package pay

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
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
	// Key is the payment means used to make the advance.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code of the payment means if not defined in the standard list of keys.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// ID or reference for the advance.
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// If this "advance" payment has come from a public grant or subsidy, set this to true.
	Grant bool `json:"grant,omitempty" jsonschema:"title=Grant"`
	// Details about the advance.
	Description string `json:"desc" jsonschema:"title=Description"`
	// How much as a percentage of the total with tax was paid
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// How much was paid.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// If different from the parent document's base currency.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
}

// Validate checks the advance looks okay
func (a *Advance) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Amount, validation.Required),
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

// JSONSchemaExtend adds the method key definitions to the schema.
func (Advance) JSONSchemaExtend(schema *jsonschema.Schema) {
	val, _ := schema.Properties.Get("key")
	prop, ok := val.(*jsonschema.Schema)
	if ok {
		prop.OneOf = make([]*jsonschema.Schema, len(MeansKeyDefinitions))
		for i, v := range MeansKeyDefinitions {
			prop.OneOf[i] = &jsonschema.Schema{
				Const:       v.Key,
				Title:       v.Title,
				Description: v.Description,
			}
		}
	}
}
