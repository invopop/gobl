package bill

import (
	"fmt"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	linePrecisionExtra uint32 = 2
)

func calculateLines(lines []*Line, cur currency.Code, rates []*currency.ExchangeRate, rr cbc.Key) error {
	for i, l := range lines {
		l.Index = i + 1
		if err := calculateLine(l, cur, rates, rr); err != nil {
			return validation.Errors{strconv.Itoa(i): err}
		}
	}
	return nil
}

func calculateLineSum(lines []*Line, cur currency.Code) num.Amount {
	sum := cur.Def().Zero()
	for _, l := range lines {
		if l.Total != nil {
			sum = sum.MatchPrecision(*l.Total)
			sum = sum.Add(*l.Total)
		}
	}
	return sum
}

// calculate figures out the totals according to quantity and discounts.
func calculateLine(l *Line, cur currency.Code, rates []*currency.ExchangeRate, rr cbc.Key) error {
	if l.Item == nil { // implies invalid, so just skip
		return nil
	}
	zero := cur.Def().Zero()

	if len(l.Substituted) > 0 {
		// Calculate the substituted line items, which have no consequence on the
		// final calculations, but still need some kind of normalization.
		for i, sl := range l.Substituted {
			sl.Index = i + 1
			if err := calculateSubLine(sl, cur, rates, rr); err != nil {
				return validation.Errors{
					"substituted": validation.Errors{strconv.Itoa(i): err},
				}
			}
		}
	}

	// Use the breakdown to calculate and replace the item price
	// while also maintaining precision.
	if len(l.Breakdown) > 0 {
		np := zero
		hasPrice := false
		for i, sl := range l.Breakdown {
			sl.Index = i + 1
			if err := calculateSubLine(sl, cur, rates, rr); err != nil {
				return validation.Errors{
					"breakdown": validation.Errors{strconv.Itoa(i): err},
				}
			}
			if sl.Total != nil {
				hasPrice = true
				np = np.MatchPrecision(*sl.Total).Add(*sl.Total)
			}
		}
		if hasPrice {
			np = np.Rescale(determineSubLinePrecision(l.Breakdown))
			l.Item.Currency = cur
			l.Item.Price = &np
			l.Item.AltPrices = nil
		}
	}

	if l.Item.Price == nil {
		l.Item.AltPrices = nil
		l.Sum = nil
		l.Total = nil
		return nil
	}

	// Perform currency manipulation to ensure item's price is
	// in the document's currency.
	if err := calculateLineItemPrice(l.Item, cur, rates); err != nil {
		return validation.Errors{
			"item": err,
		}
	}
	// Increase price accuracy for calculations
	exp := zero.Exp()
	if rr == tax.RoundingRulePrecise {
		exp += linePrecisionExtra
	}
	price := l.Item.Price.RescaleUp(exp)

	// Calculate the line sum and total
	sum := price.Multiply(l.Quantity)
	sum = tax.ApplyRoundingRule(rr, cur, sum)
	total := sum
	total = calculateLineDiscounts(l.Discounts, sum, total, cur, rr)
	total = calculateLineCharges(l.Charges, l.Quantity, sum, total, cur, rr)

	// Assume the updated sum and total
	l.Sum = &sum
	l.Total = &total

	return nil
}

// calculateSubline figures out the totals according to quantity and discounts.
// We don't apply rounding rules here, as the objective is to have
// maximum precision to determine the final line item price.
func calculateSubLine(sl *SubLine, cur currency.Code, rates []*currency.ExchangeRate, rr cbc.Key) error {
	if sl.Item == nil {
		return nil
	}

	if sl.Item.Price == nil {
		sl.Sum = nil
		sl.Total = nil
		return nil
	}

	// Perform currency manipulation to ensure item's price is
	// in the document's currency.
	if err := calculateLineItemPrice(sl.Item, cur, rates); err != nil {
		return err
	}

	// Increase price accuracy for calculations depending on rounding rule
	zero := cur.Def().Zero()
	price := *sl.Item.Price
	if rr == tax.RoundingRulePrecise {
		price = price.RescaleUp(zero.Exp() + linePrecisionExtra)
	}

	// Calculate the line sum and total
	sum := price.Multiply(sl.Quantity)
	sum = tax.ApplyRoundingRule(rr, cur, sum)
	total := sum
	total = calculateLineDiscounts(sl.Discounts, sum, total, cur, rr)
	total = calculateLineCharges(sl.Charges, sl.Quantity, sum, total, cur, rr)

	// Rescale the final sum and total
	sl.Sum = &sum
	sl.Total = &total

	return nil
}

func calculateLineDiscounts(discounts []*LineDiscount, sum, total num.Amount, cur currency.Code, rr cbc.Key) num.Amount {
	cd := cur.Def()
	for _, d := range discounts {
		if d.Percent != nil && !d.Percent.IsZero() {
			base := sum
			if d.Base != nil {
				b := d.Base.RescaleUp(cd.Subunits)
				d.Base = &b
				base = d.Base.RescaleUp(cd.Subunits + linePrecisionExtra)
				base = tax.ApplyRoundingRule(rr, cur, base)
			}
			d.Amount = d.Percent.Of(base) // always override
		}
		d.Amount = cd.RescaleUp(d.Amount)
		total = total.Subtract(d.Amount)
	}
	return total
}

func calculateLineCharges(charges []*LineCharge, quantity, sum, total num.Amount, cur currency.Code, rr cbc.Key) num.Amount {
	cd := cur.Def()
	for _, c := range charges {
		if c.Percent != nil && !c.Percent.IsZero() {
			base := sum
			if c.Base != nil {
				b := c.Base.RescaleUp(cd.Subunits)
				c.Base = &b
				base = c.Base.RescaleUp(cd.Subunits + linePrecisionExtra)
				base = tax.ApplyRoundingRule(rr, cur, base)
			}
			c.Amount = c.Percent.Of(base) // always override
		}
		// Charges also support setting a rate and quantity
		if c.Rate != nil {
			q := quantity
			if c.Quantity != nil {
				q = *c.Quantity
			}
			c.Amount = c.Rate.Multiply(q)
		}
		c.Amount = cd.RescaleUp(c.Amount)
		total = total.Add(c.Amount)
	}
	return total
}

// calculateItemPrice will attempt to perform any currency conversion process on
// the line item's data so that the currency always matches that of the
// document.
func calculateLineItemPrice(item *org.Item, cur currency.Code, rates []*currency.ExchangeRate) error {
	icur := item.Currency
	if icur == currency.CodeEmpty {
		icur = cur
	}
	price := item.Price.MatchPrecision(icur.Def().Zero())
	if item.Currency == currency.CodeEmpty || item.Currency == cur {
		item.Price = &price
		return nil
	}

	// Grab a copy of the base price
	nap := &currency.Amount{
		Currency: item.Currency,
		Value:    price,
	}

	// First check the alt prices
	for _, ap := range item.AltPrices {
		if ap.Currency == cur {
			item.Currency = ap.Currency
			price = ap.Value.MatchPrecision(ap.Currency.Def().Zero())
			item.Price = &price
			item.AltPrices = []*currency.Amount{nap}
			return nil
		}
	}

	// Try to perform a currency exchange
	np := currency.Convert(rates, item.Currency, cur, price)
	if np == nil {
		return fmt.Errorf("no exchange rate found from '%v' to '%v'", item.Currency, cur)
	}
	item.Price = np
	item.Currency = cur
	item.AltPrices = []*currency.Amount{nap}
	return nil
}

// determineSubLinePrecision will iterate through the provided sublines to try and
// determine the precision to use for the final price.
func determineSubLinePrecision(sls []*SubLine) uint32 {
	e := uint32(0)
	for _, sl := range sls {
		if sl.Item == nil || sl.Item.Price == nil {
			continue
		}
		x := sl.Item.Price.Exp()
		if x > e {
			e = x
		}
	}
	return e
}

// roundLines is a convenience function to round all the lines in a document.
func roundLines(lines []*Line, cur currency.Code) {
	for _, l := range lines {
		l.round(cur)
	}
}

// round performs a rounding operation on key numbers of the line
// so that no number exceeds the item price's or currency's precision.
// This method is only useful for precision rounding, as currency
// round will automatically apply the correct rounding rules.
func (l *Line) round(cur currency.Code) {
	e := cur.Def().Subunits
	if l.Item.Price != nil {
		e = l.Item.Price.Exp()
	}
	if l.Sum != nil {
		sum := l.Sum.RescaleDown(e)
		l.Sum = &sum
	}
	if l.Total != nil {
		e := l.Item.Price.Exp()
		total := l.Total.RescaleDown(e)
		l.Total = &total
	}
	for _, d := range l.Discounts {
		d.round(e)
	}
	for _, c := range l.Charges {
		c.round(e)
	}
	for _, sl := range l.Breakdown {
		sl.round(e)
	}
	for _, sl := range l.Substituted {
		sl.round(e)
	}
}

func (d *LineDiscount) round(e uint32) {
	d.Amount = d.Amount.RescaleDown(e)
}

func (c *LineCharge) round(e uint32) {
	c.Amount = c.Amount.RescaleDown(e)
}

// round performs a rounding operation on the sub-line's totals
// so that everything is aligned with the currency's precision.
func (sl *SubLine) round(e uint32) {
	if sl.Sum != nil {
		// Ensure sum precision is aligned with price
		sum := sl.Sum.RescaleDown(e)
		sl.Sum = &sum
	}
	if sl.Total != nil {
		total := sl.Total.RescaleDown(e)
		sl.Total = &total
	}
}
