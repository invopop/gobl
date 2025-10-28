package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// CategoryTotal groups together all rates inside a given category.
type CategoryTotal struct {
	Code        cbc.Code     `json:"code" jsonschema:"title=Code"`
	Retained    bool         `json:"retained,omitempty" jsonschema:"title=Retained"`
	Informative bool         `json:"informative,omitempty" jsonschema:"title=Informative"`
	Rates       []*RateTotal `json:"rates" jsonschema:"title=Rates"`
	Amount      num.Amount   `json:"amount" jsonschema:"title=Amount"`
	Surcharge   *num.Amount  `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`
}

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and rate. The Key is optional as we may be using
// the percentage to group rates.
type RateTotal struct {
	// Country code override when issuing with taxes applied from different countries,
	// it'd be very strange to mix rates from different countries, but in theory
	// this would be possible.
	Country l10n.TaxCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// Tax key if supported by the category.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// If the rate is defined with extensions, they'll be used to group by also.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Base amount that the percentage is applied to.
	Base num.Amount `json:"base" jsonschema:"title=Base"`
	// Percentage of the rate. Will be nil when taxes are **exempt**.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Surcharge applied to the rate.
	Surcharge *RateTotalSurcharge `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`
	// Total amount of rate, excluding surcharges
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// RateTotalSurcharge reflects the sum surcharges inside the rate.
type RateTotalSurcharge struct {
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent"`
	Amount  num.Amount     `json:"amount" jsonschema:"title=Amount"`
}

// Total contains a set of Category Totals which in turn
// contain all the accumulated taxes contained in the document. The resulting
// `sum` is that value that should be added to the payable total.
type Total struct {
	// Grouping of all the taxes by their category
	Categories []*CategoryTotal `json:"categories,omitempty" jsonschema:"title=Categories"`
	// Total value of all non-retained or indirect taxes.
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`
	// Sum of retained or withheld tax amounts
	Retained *num.Amount `json:"retained,omitempty" jsonschema:"title=Retained"`
}

// newCategoryTotal prepares a category total calculation.
func newCategoryTotal(c *Combo, zero num.Amount) *CategoryTotal {
	ct := new(CategoryTotal)
	ct.Code = c.Category
	ct.Rates = make([]*RateTotal, 0)
	ct.Amount = zero
	ct.Retained = c.retained
	ct.Informative = c.informative
	return ct
}

// newRateTotal returns a rate total.
func newRateTotal(c *Combo, zero num.Amount) *RateTotal {
	rt := new(RateTotal)
	rt.Key = c.Key         // may be empty!
	rt.Country = c.Country // usually empty
	rt.Ext = c.Ext         // may be empty!
	if c.Percent != nil {
		pc := *c.Percent
		rt.Percent = &pc
	}
	rt.Base = zero
	rt.Amount = zero
	if c.Surcharge != nil {
		rt.Surcharge = &RateTotalSurcharge{
			Percent: *c.Surcharge,
			Amount:  zero,
		}
	}
	return rt
}

// Category provides the category total for the matching code.
func (t *Total) Category(code cbc.Code) *CategoryTotal {
	if t == nil {
		return nil
	}
	for _, ct := range t.Categories {
		if ct.Code == code {
			return ct
		}
	}
	return nil
}

func (rt *RateTotal) matches(c *Combo) bool {
	if !rt.Ext.Equals(c.Ext) {
		// Extensions, if set, should always match
		return false
	}
	// assume country will always be empty if same as the base
	if rt.Country != c.Country {
		return false
	}
	if rt.Percent == nil || c.Percent == nil {
		return rt.Percent == nil && c.Percent == nil
	}
	if rt.Surcharge != nil || c.Surcharge != nil {
		if rt.Surcharge == nil || c.Surcharge == nil {
			return false
		}
		if !rt.Surcharge.Percent.Equals(*c.Surcharge) {
			return false
		}
	}
	return rt.Percent.Equals(*c.Percent)
}

// rateTotalFor either finds or creates total objects for the category and rate.
func (t *Total) rateTotalFor(c *Combo, zero num.Amount) *RateTotal {
	var catTotal *CategoryTotal
	for _, ct := range t.Categories {
		if ct.Code == c.Category {
			catTotal = ct
			break
		}
	}
	if catTotal == nil {
		catTotal = newCategoryTotal(c, zero)
		t.Categories = append(t.Categories, catTotal)
	}

	// Prepare the Rate, match using percent value
	var rateTotal *RateTotal
	for _, rt := range catTotal.Rates {
		if rt.matches(c) {
			rateTotal = rt
			break
		}
	}
	if rateTotal == nil {
		rateTotal = newRateTotal(c, zero)
		catTotal.Rates = append(catTotal.Rates, rateTotal)
	}

	return rateTotal
}

// Negate provides a new total with all the values inverted (positive to negative and vice versa).
// Will return nil if total is nil.
func (t *Total) Negate() *Total {
	if t == nil {
		return nil
	}
	nt := t.Clone()
	for _, ct := range nt.Categories {
		ct.Amount = ct.Amount.Negate()
		for _, rt := range ct.Rates {
			rt.Base = rt.Base.Negate()
			rt.Amount = rt.Amount.Negate()
		}
	}
	nt.Sum = t.Sum.Negate()
	return nt
}

// Matches returns true if the two totals have the same percent, surchage,
// and extensions, but potentially different amounts. Keys are ignored during
// comparisons as only the output is relevant.
func (rt *RateTotal) Matches(rt2 *RateTotal) bool {
	if rt.Ext.Equals(rt2.Ext) {
		if rt.Country == rt2.Country {
			if rt.Percent == nil && rt2.Percent == nil {
				return true
			}
			if rt.Percent != nil && rt2.Percent != nil {
				if rt.Percent.Equals(*rt2.Percent) {
					if rt.Surcharge == nil && rt2.Surcharge == nil {
						return true
					}
					if rt.Surcharge != nil && rt2.Surcharge != nil {
						if rt.Surcharge.Percent.Equals(rt2.Surcharge.Percent) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// Clone creates a new total with the same values as the original, but in an
// independent object.
func (t *Total) Clone() *Total {
	if t == nil {
		return nil
	}
	nt := new(Total)
	nt.Categories = make([]*CategoryTotal, len(t.Categories))
	for i, ct := range t.Categories {
		nt.Categories[i] = new(CategoryTotal)
		nt.Categories[i].Code = ct.Code
		nt.Categories[i].Retained = ct.Retained
		nt.Categories[i].Amount = ct.Amount
		nt.Categories[i].Surcharge = ct.Surcharge
		nt.Categories[i].Rates = make([]*RateTotal, len(ct.Rates))
		for j, rt := range ct.Rates {
			nt.Categories[i].Rates[j] = new(RateTotal)
			nt.Categories[i].Rates[j].Key = rt.Key
			nt.Categories[i].Rates[j].Country = rt.Country
			nt.Categories[i].Rates[j].Ext = rt.Ext
			nt.Categories[i].Rates[j].Base = rt.Base
			nt.Categories[i].Rates[j].Percent = rt.Percent
			nt.Categories[i].Rates[j].Amount = rt.Amount
			if rt.Surcharge != nil {
				nt.Categories[i].Rates[j].Surcharge = &RateTotalSurcharge{
					Percent: rt.Surcharge.Percent,
					Amount:  rt.Surcharge.Amount,
				}
			}
		}
	}
	nt.Sum = t.Sum
	nt.Retained = t.Retained
	return nt
}

// Merge will combine two totals objects into a new one, summing up the values
// of the categories and rates. The original totals will not be modified.
// The totals may contain zero amounts if the amounts in the second total are negative.
func (t *Total) Merge(t2 *Total) *Total {
	// Create a new total with the same categories
	nt := t.Clone()

	// Now merge the second total
	for _, ct := range t2.Categories {
		// Find the category in the nt total
		var catTotal *CategoryTotal
		for _, mct := range nt.Categories {
			if mct.Code == ct.Code {
				catTotal = mct
				break
			}
		}
		if catTotal == nil {
			catTotal = new(CategoryTotal)
			catTotal.Code = ct.Code
			catTotal.Retained = ct.Retained
			catTotal.Amount = ct.Amount
			catTotal.Surcharge = ct.Surcharge
			catTotal.Rates = append(catTotal.Rates, ct.Rates...)
			nt.Categories = append(nt.Categories, catTotal)
		} else {
			catTotal.Amount = catTotal.Amount.Add(ct.Amount)
			if ct.Surcharge != nil && catTotal.Surcharge != nil {
				ns := catTotal.Surcharge.Add(*ct.Surcharge)
				catTotal.Surcharge = &ns
			} else {
				catTotal.Surcharge = ct.Surcharge
			}
			// Merge the rates
			for _, rt := range ct.Rates {
				// Find the rate in the nt category
				var rateTotal *RateTotal
				for _, mrt := range catTotal.Rates {
					// match against the values, not the key
					if mrt.Matches(rt) {
						rateTotal = mrt
						break
					}
				}
				if rateTotal == nil {
					rateTotal = new(RateTotal)
					rateTotal.Key = rt.Key
					rateTotal.Country = rt.Country
					rateTotal.Ext = rt.Ext
					rateTotal.Base = rt.Base
					rateTotal.Percent = rt.Percent
					if rt.Surcharge != nil {
						rateTotal.Surcharge = &RateTotalSurcharge{
							Percent: rt.Surcharge.Percent,
							Amount:  rt.Surcharge.Amount,
						}
					}
					rateTotal.Amount = rt.Amount
					catTotal.Rates = append(catTotal.Rates, rateTotal)
				} else {
					// Merge the amounts
					rateTotal.Base = rateTotal.Base.Add(rt.Base)
					rateTotal.Amount = rateTotal.Amount.Add(rt.Amount)
					if rt.Surcharge != nil {
						rateTotal.Surcharge.Amount = rateTotal.Surcharge.Amount.Add(rt.Surcharge.Amount)
					}
				}
			}
		}
	}

	// Merge the sum
	nt.Sum = nt.Sum.Add(t2.Sum)

	return nt
}

// Scale will recalculate the taxable bases and amounts by multiplying them
// by the provided factor. This is used in scenarios where the total represents
// a percentage or portion of the original amount.
func (t *Total) Scale(factor num.Amount, cur currency.Code, rr cbc.Key) {
	if factor.IsZero() || t == nil {
		return
	}
	for _, ct := range t.Categories {
		for _, rt := range ct.Rates {
			rt.Base = rt.Base.Multiply(factor)
		}
	}
	t.Calculate(cur, rr)
}

// Exchange will recalculate the total with the same values, but in a different
// currency. This is used in scenarios where tax totals are embedded into
// sub-documents and need to be combined.
// If either the rate ot total is nil, nothing will be done. The rounding rule
// will be used for re-calculations.
func (t *Total) Exchange(rate *currency.ExchangeRate, rr cbc.Key) {
	if rate == nil || t == nil {
		return
	}
	for _, ct := range t.Categories {
		for _, rt := range ct.Rates {
			rt.Base = rate.Convert(rt.Base)
		}
	}
	t.Calculate(rate.To, rr)
}

// Calculate will go through all the categories and rates to calculate the final
// sum of the taxes. The rounding rule will be applied to the final sums.
func (t *Total) Calculate(cur currency.Code, rr cbc.Key) {
	if t == nil {
		return
	}
	zero := cur.Def().Zero()
	t.calculateFinalSum(zero, rr)
}

func (t *Total) calculateFinalSum(zero num.Amount, rr cbc.Key) {
	// Now go through each category to apply the percentage and calculate the final sums
	t.Sum = zero
	for _, ct := range t.Categories {
		t.calculateBaseCategoryTotal(ct, zero, rr)

		if ct.Informative {
			// Informative taxes don't affect Sum or Retained
			continue
		}
		if ct.Retained {
			if t.Retained == nil {
				t.Retained = &zero
			}
			tr := *t.Retained
			tr = matchRoundingPrecision(rr, tr, ct.Amount)
			tr = tr.Add(ct.Amount)
			if ct.Surcharge != nil {
				tr = tr.Add(*ct.Surcharge)
			}
			t.Retained = &tr
		} else {
			t.Sum = matchRoundingPrecision(rr, t.Sum, ct.Amount)
			t.Sum = t.Sum.Add(ct.Amount)
			if ct.Surcharge != nil {
				t.Sum = t.Sum.Add(*ct.Surcharge)
			}
		}
	}
}

func (t *Total) calculateBaseCategoryTotal(ct *CategoryTotal, zero num.Amount, rr cbc.Key) {
	ct.Amount = zero
	for _, rt := range ct.Rates {
		if rt.Percent == nil {
			rt.Amount = zero
			continue // exempt, nothing else to do
		}
		base := rt.Base
		rt.Amount = rt.Percent.Of(rt.Base)
		ct.Amount = matchRoundingPrecision(rr, ct.Amount, rt.Amount)
		ct.Amount = ct.Amount.Add(rt.Amount)
		if rt.Surcharge != nil {
			rt.Surcharge.Amount = rt.Surcharge.Percent.Of(base)
			if ct.Surcharge == nil {
				ct.Surcharge = &zero
			}
			a := rt.Surcharge.Amount
			x := *ct.Surcharge
			x = matchRoundingPrecision(rr, x, a)
			x = x.Add(a)
			ct.Surcharge = &x
		}
	}
}

// matchPrecision will decide what precision to maintain on the amount based on
// the rounding rule.
func matchRoundingPrecision(rr cbc.Key, a, b num.Amount) num.Amount {
	switch rr {
	case RoundingRuleCurrency:
		return a // maintain original precision
	}
	return a.MatchPrecision(b)
}

// Round will go through all the values generated and round them to the currency's
// preferred precision.
func (t *Total) Round(zero num.Amount) {
	for _, ct := range t.Categories {
		for _, rt := range ct.Rates {
			rt.Amount = rt.Amount.Rescale(zero.Exp())
			rt.Base = rt.Base.Rescale(zero.Exp())
			if rt.Surcharge != nil {
				rt.Surcharge.Amount = rt.Surcharge.Amount.Rescale(zero.Exp())
			}
		}
		ct.Amount = ct.Amount.Rescale(zero.Exp())
		if ct.Surcharge != nil {
			*ct.Surcharge = ct.Surcharge.Rescale(zero.Exp())
		}
	}
	t.Sum = t.Sum.Rescale(zero.Exp())
	if t.Retained != nil {
		*t.Retained = t.Retained.Rescale(zero.Exp())
	}
}
