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

// calculate figures out the totals according to quantity and discounts.
func (l *Line) calculate(r *tax.Regime, cur currency.Code, rates []*currency.ExchangeRate) error {
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
	// in the document's currency.
	if err := l.normalizeItemPrice(cur, rates); err != nil {
		return err
	}

	// Increase price accuracy for calculations
	zero := cur.Def().Zero()
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

// normalizeItemPrice will attempt to perform any currency conversion process on
// the line item's data so that the currency always matches that of the
// document.
func (l *Line) normalizeItemPrice(cur currency.Code, rates []*currency.ExchangeRate) error {
	item := l.Item
	icur := item.Currency
	if icur == currency.CodeEmpty {
		icur = cur
	}
	item.Price = item.Price.MatchPrecision(icur.Def().Zero())
	if item.Currency == currency.CodeEmpty || item.Currency == cur {
		return nil
	}

	// Grab a copy of the base price
	nap := &currency.Amount{
		Currency: item.Currency,
		Value:    item.Price,
	}

	// First check the alt prices
	for _, ap := range item.AltPrices {
		if ap.Currency == cur {
			item.Currency = ap.Currency
			item.Price = ap.Value.MatchPrecision(ap.Currency.Def().Zero())
			item.AltPrices = []*currency.Amount{nap}
			return nil
		}
	}

	// Try to perform a currency exchange
	np := currency.Convert(rates, item.Currency, cur, item.Price)
	if np == nil {
		return fmt.Errorf("no exchange rate found from '%v' to '%v'", item.Currency, cur)
	}
	item.Price = *np
	item.Currency = cur
	item.AltPrices = []*currency.Amount{nap}
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

	l2i.AltPrices = nil // empty alternative prices
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

func calculateLines(r *tax.Regime, lines []*Line, cur currency.Code, rates []*currency.ExchangeRate) error {
	for i, l := range lines {
		l.Index = i + 1
		if err := l.calculate(r, cur, rates); err != nil {
			return validation.Errors{strconv.Itoa(i): err}
		}
	}
	return nil
}

func calculateLineSum(lines []*Line, cur currency.Code) num.Amount {
	sum := cur.Def().Zero()
	for _, l := range lines {
		sum = sum.MatchPrecision(l.total)
		sum = sum.Add(l.total)
	}
	return sum
}
