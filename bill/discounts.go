package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
)

// LineDiscount represents an amount deducted from the line, and will be
// applied before taxes.
// TODO: use UNTDID 5189 code list
type LineDiscount struct {
	// Percentage if fixed amount not applied
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Fixed discount amount to apply (calculated if percent present).
	Amount num.Amount `json:"amount" jsonschema:"title=Value" jsonschema_extras:"calculated=true"`
	// Reason code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the discount was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
}

// Validate checks the line discount's fields.
func (ld *LineDiscount) Validate() error {
	return validation.ValidateStruct(ld,
		validation.Field(&ld.Percent),
		validation.Field(&ld.Amount, validation.Required),
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
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Line number inside the list of discounts (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Reference or ID for this Discount
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Base represents the value used as a base for percent calculations.
	// If not already provided, we'll take the invoices sum.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the invoice's Sum.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Amount to apply (calculated if percent present).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// List of taxes to apply to the discount
	Taxes tax.Set `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Code for the reason this discount applied
	Code string `json:"code,omitempty" jsonschema:"title=Reason Code"`
	// Text description as to why the discount was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Additional semi-structured information.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks the discount's fields.
func (m *Discount) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.UUID),
		validation.Field(&m.Base),
		validation.Field(&m.Percent),
		validation.Field(&m.Amount, validation.Required),
		validation.Field(&m.Taxes),
		validation.Field(&m.Meta),
	)
}

// GetTaxes responds with the array of tax rates applied to this line.
func (m *Discount) GetTaxes() tax.Set {
	return m.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
// We return a negative value so that discounts will be applied correctly.
func (m *Discount) GetTotal() num.Amount {
	return m.Amount.Invert()
}

func (m *Discount) removeIncludedTaxes(cat org.Code, accuracy uint32) *Discount {
	rate := m.Taxes.Get(cat)
	if rate == nil {
		return m
	}
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Remove(rate.Percent)
	return &m2
}
