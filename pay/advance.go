package pay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
)

// Advances contains an array of advance objects.
type Advances []*Advance

// Advance represents a single payment that has been made already, such
// as a deposit on an intent to purchase, or as credit from a previous
// invoice which was later corrected or cancelled.
type Advance struct {
	// Unique identifier for this advance.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// When the advance was made.
	Date *org.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// ID or reference for the advance.
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// If this "advance" payment has come from a public grant or subsidy, set this to true.
	Grant bool `json:"grant,omitempty" jsonschema:"title=Grant"`
	// Details about the advance.
	Description string `json:"desc" jsonschema:"title=Description"`
	// How much as a percentage of the total with tax was paid
	Rate *num.Percentage `json:"rate,omitempty" jsonschema:"title=Rate"`
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

// Calculate will update the amount using the rate of the provided
// total, if defined.
func (a *Advance) Calculate(totalWithTax num.Amount) {
	if a.Rate != nil {
		a.Amount = a.Rate.Of(totalWithTax)
	}
}
