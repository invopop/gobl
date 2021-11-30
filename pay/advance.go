package pay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
)

// Advance represents a single payment that has been made already, such
// as a deposit on an intent to purchase, or as credit from a previous
// invoice which was later corrected or cancelled.
type Advance struct {
	UUID        *uuid.UUID    `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identifier for this advance."`
	Date        *org.Date     `json:"date,omitempty" jsonschema:"title=Date,description=When the advance was made."`
	Code        string        `json:"code,omitempty" jsonschema:"title=Code,description=Reference for the advance."`
	Description string        `json:"desc" jsonschema:"title=Description,description=Details about the advance."`
	Amount      num.Amount    `json:"amount" jsonschema:"title=Amount,description=How much was paid."`
	Currency    currency.Code `json:"currency,omitempty" jsonschema:"title=Currency,description=If different from the parent document's base currency."`
}

// Validate checks the advance looks okay
func (a *Advance) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Amount, validation.Required),
		validation.Field(&a.Description, validation.Required),
	)
}
