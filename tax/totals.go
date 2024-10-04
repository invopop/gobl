package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// CategoryTotal groups together all rates inside a given category.
type CategoryTotal struct {
	Code      cbc.Code     `json:"code" jsonschema:"title=Code"`
	Retained  bool         `json:"retained,omitempty" jsonschema:"title=Retained"`
	Rates     []*RateTotal `json:"rates" jsonschema:"title=Rates"`
	Amount    num.Amount   `json:"amount" jsonschema:"title=Amount"`
	Surcharge *num.Amount  `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`

	amount num.Amount // internal amount with greater accuracy
}

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and rate. The Key is optional as we may be using
// the percentage to group rates.
type RateTotal struct {
	// Optional rate key is required when grouping.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Country code override when issuing with taxes applied from different countries,
	// it'd be very strange to mix rates from different countries, but in theory
	// this would be possible.
	Country l10n.TaxCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
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
	// Total value of all the taxes applied.
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`

	// Precise sum in the background, in case needed for calculations
	sum num.Amount
}

// PreciseAmount contains the intermediary amount generated from the calculator
// with the original precision. This is useful when a Category Total needs
// to be used for further calculations, such as when an invoice includes taxes.
func (ct *CategoryTotal) PreciseAmount() num.Amount {
	if !ct.amount.IsZero() {
		return ct.amount
	}
	return ct.Amount
}

// PreciseSum contains an intermediary sum generated from the calculator
// with the original precision. If no calculations were made on the totals,
// such as when loading, the original sum will be provided instead.
func (t *Total) PreciseSum() num.Amount {
	if !t.sum.IsZero() {
		return t.sum
	}
	return t.Sum
}

// newCategoryTotal prepares a category total calculation.
func newCategoryTotal(c *Combo, zero num.Amount) *CategoryTotal {
	ct := new(CategoryTotal)
	ct.Code = c.Category
	ct.Rates = make([]*RateTotal, 0)
	ct.Amount = zero
	ct.amount = zero
	ct.Retained = c.retained
	return ct
}

// newRateTotal returns a rate total.
func newRateTotal(c *Combo, zero num.Amount) *RateTotal {
	rt := new(RateTotal)
	rt.Key = c.Rate        // may be empty!
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
func (t *Total) Negate() *Total {
	nt := t.Clone()
	for _, ct := range nt.Categories {
		ct.Amount = ct.Amount.Negate()
		ct.amount = ct.amount.Negate()
		for _, rt := range ct.Rates {
			rt.Base = rt.Base.Negate()
			rt.Amount = rt.Amount.Negate()
		}
	}
	nt.Sum = t.Sum.Negate()
	nt.sum = t.sum.Negate()
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
		nt.Categories[i].amount = ct.amount
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
	nt.sum = t.sum
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
			catTotal.amount = ct.amount
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
	nt.sum = nt.sum.Add(t2.sum)

	return nt
}
