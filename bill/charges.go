package bill

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
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
	uuid.Identify
	// Key for grouping or identifying charges for tax purposes.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Line number inside the list of charges (calculated).
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Code to used to refer to the this charge
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Base represents the value used as a base for percent calculations instead
	// of the invoice's sum of lines.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the Base or Invoice Sum
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

	// internal amount for calculations
	amount num.Amount
}

// ValidateWithContext checks the charge's fields.
func (m *Charge) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, m,
		validation.Field(&m.UUID),
		validation.Field(&m.Base),
		validation.Field(&m.Percent,
			validation.When(
				m.Base != nil,
				validation.Required,
			),
		),
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

func (m *Charge) removeIncludedTaxes(cat cbc.Code) *Charge {
	accuracy := defaultTaxRemovalAccuracy
	rate := m.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return m
	}
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Remove(*rate.Percent)
	return &m2
}

func (m *Charge) convertInto(ex *currency.ExchangeRate) *Charge {
	accuracy := defaultCurrencyConversionAccuracy
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Multiply(ex.Amount)
	return &m2
}

func calculateCharges(lines []*Charge, sum, zero num.Amount) {
	// COPIED FROM discount.go
	if len(lines) == 0 {
		return
	}
	for i, l := range lines {
		l.Index = i + 1
		if l.Percent != nil && !l.Percent.IsZero() {
			base := sum
			exp := zero.Exp()
			if l.Base != nil {
				base = l.Base.RescaleUp(exp)
				exp = base.Exp()
			}
			l.Amount = l.Percent.Of(base)
			l.amount = l.Amount
			l.Amount = l.Amount.Rescale(exp)
		} else {
			l.Amount = l.Amount.MatchPrecision(zero)
			l.amount = l.Amount
		}
	}
}

func calculateChargeSum(charges []*Charge, zero num.Amount) *num.Amount {
	if len(charges) == 0 {
		return nil
	}
	total := zero
	for _, l := range charges {
		total = total.MatchPrecision(l.amount)
		total = total.Add(l.amount)
	}
	return &total
}
