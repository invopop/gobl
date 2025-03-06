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
		sum = sum.MatchPrecision(l.total)
		sum = sum.Add(l.total)
	}
	return sum
}

// calculate figures out the totals according to quantity and discounts.
func calculateLine(l *Line, cur currency.Code, rates []*currency.ExchangeRate, rr cbc.Key) error {
	if l.Item == nil { // implies invalid, so just skip
		return nil
	}
	zero := cur.Def().Zero()
	price := zero

	if len(l.Substituted) > 0 {
		// Calculate the substituted line items, which have no consequence on the
		// final calculations, but still need some kind of normalization.
		for i, sl := range l.Substituted {
			sl.Index = i + 1
			if err := calculateSubLine(sl, cur, rates); err != nil {
				return validation.Errors{
					"substituted": validation.Errors{strconv.Itoa(i): err},
				}
			}
		}
	}

	if len(l.Breakdown) > 0 {
		// Use the breakdown to calculate and replace the item price
		// while also maintaining precision.
		np := zero
		hasPrice := false
		for i, sl := range l.Breakdown {
			sl.Index = i + 1
			if err := calculateSubLine(sl, cur, rates); err != nil {
				return validation.Errors{
					"breakdown": validation.Errors{strconv.Itoa(i): err},
				}
			}
			if sl.Total != nil {
				hasPrice = true
				np = np.MatchPrecision(*sl.Total).Add(*sl.Total)
				price = price.MatchPrecision(sl.total).Add(sl.total)
			}
		}
		if hasPrice {
			l.Item.Currency = cur
			l.Item.Price = &np
			l.Item.AltPrices = nil
		}
	} else if l.Item.Price != nil {
		// Perform currency manipulation to ensure item's price is
		// in the document's currency.
		if err := calculateLineItemPrice(l.Item, cur, rates); err != nil {
			return validation.Errors{
				"item": err,
			}
		}
		// Increase price accuracy for calculations
		exp := zero.Exp()
		if rr == tax.RoundingRuleSumThenRound {
			exp += 2
		}
		price = l.Item.Price.RescaleUp(exp)
	}

	if l.Item.Price == nil {
		l.Item.AltPrices = nil
		l.Sum = nil
		l.Total = nil
		return nil
	}

	// Calculate the line sum and total
	sum := price.Multiply(l.Quantity)
	total := tax.ApplyRoundingRule(rr, cur, sum)
	total = calculateLineDiscounts(l.Discounts, sum, total, cur)
	total = calculateLineCharges(l.Charges, l.Quantity, sum, total, cur)

	// Rescale the final sum to match item's price
	sum = sum.Rescale(l.Item.Price.Exp())
	l.Sum = &sum
	l.total = tax.ApplyRoundingRule(rr, cur, total)
	total = l.total.Rescale(zero.Exp())
	l.Total = &total

	return nil
}

// calculateSubline figures out the totals according to quantity and discounts.
// We don't apply rounding rules here, as the objective is to have
// maximum precision to determine the final line item price.
func calculateSubLine(sl *SubLine, cur currency.Code, rates []*currency.ExchangeRate) error {
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

	// Increase price accuracy for calculations
	zero := cur.Def().Zero()
	price := sl.Item.Price.RescaleUp(zero.Exp() + 2)

	// Calculate the line sum and total
	sum := price.Multiply(sl.Quantity)
	total := sum
	total = calculateLineDiscounts(sl.Discounts, sum, total, cur)
	total = calculateLineCharges(sl.Charges, sl.Quantity, sum, total, cur)

	// Rescale the final sum and total
	sl.total = total
	total = total.Rescale(sl.Item.Price.Exp())
	sum = sum.Rescale(sl.Item.Price.Exp())
	sl.Sum = &sum
	sl.Total = &total

	return nil
}

func calculateLineDiscounts(discounts []*LineDiscount, sum, total num.Amount, cur currency.Code) num.Amount {
	cd := cur.Def()
	for _, d := range discounts {
		if d.Percent != nil && !d.Percent.IsZero() {
			base := sum
			if d.Base != nil {
				base = cd.Rescale(*d.Base)
				d.Base = &base
			}
			d.Amount = d.Percent.Of(base) // always override
		}
		total = total.Subtract(d.Amount)
		// As per EN16931 specs, discount amounts have same number of
		// decimal places as the currency.
		d.Amount = cd.Rescale(d.Amount)
	}
	return total
}

func calculateLineCharges(charges []*LineCharge, quantity, sum, total num.Amount, cur currency.Code) num.Amount {
	cd := cur.Def()
	for _, c := range charges {
		if c.Percent != nil && !c.Percent.IsZero() {
			base := sum
			if c.Base != nil {
				base = cd.Rescale(*c.Base)
				c.Base = &base
			}
			c.Amount = c.Percent.Of(sum) // always override
		}
		// Charges also support setting a rate and quantity
		if c.Rate != nil {
			q := quantity
			if c.Quantity != nil {
				q = *c.Quantity
			}
			c.Amount = c.Rate.Multiply(q)
		}

		total = total.Add(c.Amount)
		// As per EN16931 specs, charge amounts have same number of
		// decimal places as the currency.
		c.Amount = cd.Rescale(c.Amount)
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
