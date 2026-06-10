package currency

import (
	"github.com/invopop/gobl/num"
)

// An Amount represents a monetary value in a specific currency.
type Amount struct {
	// Label allows for additional information to be added to the
	// currency Amount that may be useful.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Code defines the currency for this amount.
	Currency Code `json:"currency" jsonschema:"title=Currency"`
	// Value is the amount in the currency.
	Value num.Amount `json:"value" jsonschema:"title=Value"`
}
