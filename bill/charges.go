package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// LineCharge represents an amount added to the line, and will be
// applied before taxes.
// TODO: use UNTDID 7161 code list
type LineCharge struct {
	// Percentage rate if fixed amount not applied
	Rate *num.Percentage `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Fixed or resulting charge amount to apply
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// Reference code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the charge was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
}

// Validate checks the line charge's fields.
func (lc *LineCharge) Validate() error {
	return validation.ValidateStruct(lc,
		validation.Field(&lc.Amount, validation.Required),
	)
}

// Charges represents an array of charge objects
type Charges []*Charge

// Charge represents a surchange applied to the complete document
// independent from the individual lines.
type Charge struct {
	// Unique identifying for the discount entry
	UUID string `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Line number inside the list of discounts
	Index int `json:"i" jsonschema:"title=Index"`
	// Code to used to refer to the this charge
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Base represents the value used as a base for rate calculations.
	// If not already provided, we'll take the invoices sum before
	// discounts.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage rate to apply to the invoice's Sum
	Rate *num.Percentage `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Amount to apply
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// List of taxes to apply to the charge
	Taxes tax.Rates `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Code for why was this charge applied?
	Code string `json:"code,omitempty" jsonschema:"title=Reason Code"`
	// Text description as to why the charge was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Additional semi-structured information.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks the discount's fields.
func (m *Charge) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Amount, validation.Required),
	)
}

// GetTaxRates responds with the array of tax rates applied to this line.
func (m *Charge) GetTaxRates() tax.Rates {
	return m.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (m *Charge) GetTotal() num.Amount {
	return m.Amount
}
