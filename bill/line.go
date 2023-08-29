package bill

import (
	"context"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Line is a single row in an invoice.
type Line struct {
	// Unique identifier for this line
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Line number inside the parent (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Number of items
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Details about what is being sold
	Item *org.Item `json:"item" jsonschema:"title=Item"`
	// Result of quantity multiplied by the item's price (calculated)
	Sum num.Amount `json:"sum" jsonschema:"title=Sum" jsonschema_extras:"calculated=true"`
	// Discounts applied to this line
	Discounts []*LineDiscount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges applied to this line
	Charges []*LineCharge `json:"charges,omitempty" jsonschema:"title=Charges"`
	// Map of taxes to be applied and used in the invoice totals
	Taxes tax.Set `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Total line amount after applying discounts to the sum (calculated).
	Total num.Amount `json:"total" jsonschema:"title=Total"  jsonschema_extras:"calculated=true"`
	// Set of specific notes for this line that may be required for
	// clarification.
	Notes []*cbc.Note `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// GetTaxes responds with the array of tax rates applied to this line.
func (l *Line) GetTaxes() tax.Set {
	return l.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (l *Line) GetTotal() num.Amount {
	return l.Total
}

// ValidateWithContext ensures the line contains everything required using
// the provided context that should include the regime.
func (l *Line) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, l,
		validation.Field(&l.UUID),
		validation.Field(&l.Index, validation.Required),
		validation.Field(&l.Quantity, validation.Required),
		validation.Field(&l.Item, validation.Required),
		validation.Field(&l.Sum, validation.Required),
		validation.Field(&l.Discounts),
		validation.Field(&l.Charges),
		validation.Field(&l.Taxes),
		validation.Field(&l.Total, validation.Required),
		validation.Field(&l.Notes),
	)
}

// calculate figures out the totals according to quantity and discounts.
func (l *Line) calculate(r *tax.Regime, zero num.Amount) error {
	if l.Item == nil {
		return nil
	}

	if err := r.CalculateObject(l); err != nil {
		return err
	}
	if err := r.CalculateObject(l.Item); err != nil {
		return validation.Errors{"item": err}
	}

	// Ensure the Price precision is set correctly according to the currency
	l.Item.Price = l.Item.Price.MatchPrecision(zero)

	// Calculate the line sum and total
	l.Sum = l.Item.Price.Multiply(l.Quantity)
	l.Total = l.Sum

	for _, d := range l.Discounts {
		if d.Percent != nil && !d.Percent.IsZero() {
			d.Amount = d.Percent.Of(l.Sum) // always override
		}
		d.Amount = d.Amount.MatchPrecision(zero)
		l.Total = l.Total.Subtract(d.Amount)
	}

	for _, c := range l.Charges {
		if c.Percent != nil && !c.Percent.IsZero() {
			c.Amount = c.Percent.Of(l.Sum) // always override
		}
		c.Amount = c.Amount.MatchPrecision(zero)
		l.Total = l.Total.Add(c.Amount)
	}
	return nil
}

func (l *Line) removeIncludedTaxes(cat cbc.Code, accuracy uint32) *Line {
	rate := l.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return l
	}

	l2 := *l
	l2i := *l.Item

	l2.Sum = l2.Sum.Upscale(accuracy).Remove(*rate.Percent)
	l2i.Price = l2.Sum.Divide(l2.Quantity)
	l2.Total = l2.Total.Upscale(accuracy).Remove(*rate.Percent)

	if len(l2.Discounts) > 0 {
		rows := make([]*LineDiscount, len(l2.Discounts))
		for i, v := range l.Discounts {
			d := *v
			d.Amount = d.Amount.Upscale(accuracy).Remove(*rate.Percent)
			rows[i] = &d
		}
		l2.Discounts = rows
	}

	if len(l2.Charges) > 0 {
		rows := make([]*LineCharge, len(l2.Charges))
		for i, v := range l.Charges {
			d := *v
			d.Amount = d.Amount.Upscale(accuracy).Remove(*rate.Percent)
			rows[i] = &d
		}
		l2.Charges = rows
	}

	l2.Item = &l2i
	return &l2
}

func calculateLines(r *tax.Regime, zero num.Amount, lines []*Line) error {
	for i, l := range lines {
		l.Index = i + 1
		if err := l.calculate(r, zero); err != nil {
			return validation.Errors{strconv.Itoa(i): err}
		}
	}
	return nil
}

func calculateLineSum(zero num.Amount, lines []*Line) num.Amount {
	sum := zero
	for _, l := range lines {
		sum = sum.MatchPrecision(l.Total)
		sum = sum.Add(l.Total)
	}
	return sum
}
