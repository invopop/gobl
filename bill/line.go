package bill

import (
	"context"
	"fmt"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Line is a single row in an invoice.
type Line struct {
	uuid.Identify
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

	// internal amount provided with greater precision
	total num.Amount
}

// GetTaxes responds with the array of tax rates applied to this line.
func (l *Line) GetTaxes() tax.Set {
	return l.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (l *Line) GetTotal() num.Amount {
	return l.total
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

// calculate figures out the totals according to quantity and discounts
// always using the currency of the item if different from the parents
// currency.
func (l *Line) calculate(r *tax.Regime, cur currency.Code) error {
	if l.Item == nil {
		return nil
	}
	if err := r.CalculateObject(l); err != nil {
		return err
	}

	// Ensure Item looks good
	if err := l.Item.Calculate(); err != nil { // Normalizes
		return validation.Errors{"item": err}
	}
	if err := r.CalculateObject(l.Item); err != nil {
		return validation.Errors{"item": err}
	}

	// Perform currency manipulation to ensure item's price is
	// in local currency.
	zero := cur.Def().Zero()
	if l.Item.Currency != currency.CodeEmpty {
		zero = l.Item.Currency.Def().Zero()
	}
	l.Item.Price = l.Item.Price.MatchPrecision(zero)

	// Increase price accuracy for calculations
	price := l.Item.Price
	price = price.RescaleUp(zero.Exp() + 2)

	// Calculate the line sum and total
	l.Sum = price.Multiply(l.Quantity)
	l.total = l.Sum

	for _, d := range l.Discounts {
		if d.Percent != nil && !d.Percent.IsZero() {
			d.Amount = d.Percent.Of(l.Sum) // always override
		}
		d.Amount = d.Amount.MatchPrecision(zero)
		l.total = l.total.Subtract(d.Amount)
		d.Amount = d.Amount.Rescale(l.Item.Price.Exp())
	}

	for _, c := range l.Charges {
		if c.Percent != nil && !c.Percent.IsZero() {
			c.Amount = c.Percent.Of(l.Sum) // always override
		}
		c.Amount = c.Amount.MatchPrecision(zero)
		l.total = l.total.Add(c.Amount)
		c.Amount = c.Amount.Rescale(l.Item.Price.Exp())
	}

	// Rescale the final sum and total
	l.Sum = l.Sum.Rescale(l.Item.Price.Exp())

	// Perform currency conversion on the total

	l.Total = l.total.Rescale(l.Item.Price.Exp())

	return nil
}

func (l *Line) removeIncludedTaxes(cat cbc.Code) *Line {
	accuracy := defaultTaxRemovalAccuracy
	rate := l.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return l
	}

	l2 := *l
	l2i := *l.Item

	l2i.Price = l.Item.Price.Upscale(accuracy).Remove(*rate.Percent)
	// assume sum and total will be calculated automatically

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

func calculateLines(r *tax.Regime, lines []*Line, cur currency.Code) error {
	for i, l := range lines {
		l.Index = i + 1
		if err := l.calculate(r, cur); err != nil {
			return validation.Errors{strconv.Itoa(i): err}
		}
	}
	return nil
}

func calculateLineSum(lines []*Line, cur currency.Code, rates []*currency.ExchangeRate) (num.Amount, error) {
	sum := cur.Def().Zero()
	for i, l := range lines {
		total := l.total
		lc := l.Item.Currency
		if lc != currency.CodeEmpty {
			np := currency.Exchange(rates, lc, cur, total)
			if np == nil {
				err := validation.Errors{
					strconv.Itoa(i): fmt.Errorf("no exchange rate found from '%v' into '%v'", lc, cur),
				}
				return sum, err
			}
			total = *np
		}
		sum = sum.MatchPrecision(total)
		sum = sum.Add(total)
	}
	return sum, nil
}
