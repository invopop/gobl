package currency

import (
	"github.com/invopop/gobl/num"
)

// An Amount represents a monetary value in a specific currency.
type Amount struct {
	// Additional information about the amount that may be useful.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Currency of the amount.
	Currency Code `json:"currency" jsonschema:"title=Currency"`
	// Amount in the currency.
	Value num.Amount `json:"value" jsonschema:"title=Value"`
}
