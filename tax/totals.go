package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
)

// CategoryTotal groups together all rates inside a given category.
type CategoryTotal struct {
	Code      cbc.Code     `json:"code" jsonschema:"title=Code"`
	Retained  bool         `json:"retained,omitempty" jsonschema:"title=Retained"`
	Rates     []*RateTotal `json:"rates" jsonschema:"title=Rates"`
	Amount    num.Amount   `json:"amount" jsonschema:"title=Amount"`
	Surcharge *num.Amount  `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`
}

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and rate. The Key is optional as we may be using
// the percentage to group rates.
type RateTotal struct {
	// Optional rate key is required when grouping.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// If the rate is defined with extensions, they'll be used to group by also.
	Ext cbc.CodeMap `json:"ext,omitempty" jsonschema:"title=Ext"`
	// Base amount that the percentage is applied to.
	Base num.Amount `json:"base" jsonschema:"title=Base"`
	// Percentage of the rate, which may be nil for exempt rates.
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

// PreciseSum is used internally to provide a more precise sum that maintains
// the accuracy provided by the original line totals.
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
	ct.Retained = c.category.Retained
	return ct
}

// newRateTotal returns a rate total.
func newRateTotal(c *Combo, zero num.Amount) *RateTotal {
	rt := new(RateTotal)
	rt.Key = c.Rate // may be empty!
	rt.Ext = c.Ext  // may be empty!
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
