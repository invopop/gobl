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

// Discount represents an allowance applied to the complete document
// independent from the individual lines. These have more in common with
// Invoice Lines than anything else, as each discount must have the
// correct taxes defined.
type Discount struct {
	uuid.Identify
	// Line number inside the list of discounts (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Reference or ID for this Discount
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Base represents the value used as a base for percent calculations instead
	// of the invoice's sum of lines.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the base or invoice's sum.
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
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// internal amount for calculations
	amount num.Amount
}

// ValidateWithContext checks the discount's fields.
func (m *Discount) ValidateWithContext(ctx context.Context) error {
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
func (m *Discount) GetTaxes() tax.Set {
	return m.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
// We return a negative value so that discounts will be applied correctly.
func (m *Discount) GetTotal() num.Amount {
	return m.Amount.Invert()
}

func (m *Discount) removeIncludedTaxes(cat cbc.Code) *Discount {
	accuracy := defaultTaxRemovalAccuracy
	rate := m.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return m
	}
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Remove(*rate.Percent)
	return &m2
}

func (m *Discount) convertInto(ex *currency.ExchangeRate) *Discount {
	accuracy := defaultCurrencyConversionAccuracy
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Multiply(ex.Amount)
	return &m2
}

func calculateDiscounts(lines []*Discount, sum, zero num.Amount) {
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

func calculateDiscountSum(discounts []*Discount, zero num.Amount) *num.Amount {
	if len(discounts) == 0 {
		return nil
	}
	total := zero
	for _, l := range discounts {
		total = total.MatchPrecision(l.amount)
		total = total.Add(l.amount)
	}
	return &total
}
