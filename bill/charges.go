package bill

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// LineCharge represents an amount added to the line, and will be
// applied before taxes.
// TODO: use UNTDID 7161 code list
type LineCharge struct {
	// Percentage if fixed amount not applied
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Fixed or resulting charge amount to apply (calculated if percent present).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// Reference code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the charge was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
}

// Validate checks the line charge's fields.
func (lc *LineCharge) Validate() error {
	return validation.ValidateStruct(lc,
		validation.Field(&lc.Percent),
		validation.Field(&lc.Amount, validation.Required),
	)
}

// Charge represents a surchange applied to the complete document
// independent from the individual lines.
type Charge struct {
	// Unique identifying for the discount entry
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Key for grouping or identifying charges for tax purposes.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Line number inside the list of discounts (calculated).
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Code to used to refer to the this charge
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Base represents the value used as a base for percent calculations.
	// If not already provided, we'll take the invoices sum before
	// discounts.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the invoice's Sum
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Amount to apply (calculated if percent present)
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// List of taxes to apply to the charge
	Taxes tax.Set `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Code for why was this charge applied?
	Code string `json:"code,omitempty" jsonschema:"title=Reason Code"`
	// Text description as to why the charge was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Additional semi-structured information.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ValidateWithContext checks the discount's fields.
func (m *Charge) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, m,
		validation.Field(&m.UUID),
		validation.Field(&m.Base),
		validation.Field(&m.Percent),
		validation.Field(&m.Amount, validation.Required),
		validation.Field(&m.Taxes),
		validation.Field(&m.Meta),
	)
}

// GetTaxes responds with the array of tax rates applied to this line.
func (m *Charge) GetTaxes() tax.Set {
	return m.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (m *Charge) GetTotal() num.Amount {
	return m.Amount
}

func (m *Charge) removeIncludedTaxes(cat cbc.Code, accuracy uint32) *Charge {
	rate := m.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return m
	}
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Remove(*rate.Percent)
	return &m2
}
