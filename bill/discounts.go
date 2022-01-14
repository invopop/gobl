package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// LineDiscount represents an amount deducted from the line, and will be
// applied before taxes.
// TODO: use UNTDID 5189 code list
type LineDiscount struct {
	// Percentage rate if fixed amount not applied
	Rate *num.Percentage `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Fixed discount amount to apply
	Value num.Amount `json:"value" jsonschema:"title=Value"`
	// Reason code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the discount was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
}

// Validate checks the line discount's fields.
func (ld *LineDiscount) Validate() error {
	return validation.ValidateStruct(ld,
		validation.Field(&ld.Value, validation.Required),
	)
}

// Discounts represents an array of discounts.
type Discounts []*Discount

// Discount represents an allowance applied to the complete document
// independent from the individual lines. These have more in common with
// Invoice Lines than anything else, as each discount must have the
// correct taxes defined.
type Discount struct {
	// Unique identifying for the discount entry
	UUID string `json:"uuid,omitempty" jsonschema:"title=UUID`
	// Line number inside the list of discounts
	Index int `json:"i" jsonschema:"title=Index"`
	// Percentage rate to apply to the invoice's Sum
	Rate *num.Percentage `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Amount to apply
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// List of taxes to apply to the discount
	Taxes tax.Rates `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Why was this discount applied?
	Code string `json:"code,omitempty" jsonschema:"title=Reason Code"`
	// Text description as to why the discount was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
}

// Validate checks the discount's fields.
func (d *Discount) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Amount, validation.Required),
	)
}

// GetTaxRates responds with the array of tax rates applied to this line.
func (l *Discount) GetTaxRates() tax.Rates {
	return l.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (l *Discount) GetTotal() num.Amount {
	return l.Amount
}
