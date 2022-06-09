package tax

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
)

// CategoryTotal groups together all rates inside a given category.
type CategoryTotal struct {
	Code     Code         `json:"code" jsonschema:"title=Code"`
	Retained bool         `json:"retained,omitempty" jsonschema:"title=Retained"`
	Rates    []*RateTotal `json:"rates" jsonschema:"title=Rates"`
	Base     num.Amount   `json:"base" jsonschema:"title=Base"`
	Amount   num.Amount   `json:"amount" jsonschema:"title=Amount"`
}

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and rate. The Key is optional as we may be using
// the percentage to group rates.
type RateTotal struct {
	Key     Key            `json:"key,omitempty" jsonschema:"title=Key"`
	Base    num.Amount     `json:"base" jsonschema:"title=Base"`
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
}

// TaxableLine defines what we expect from a line in order to subsequently calculate
// the taxes that need to be added or retained.
type TaxableLine interface {
	GetTaxes() Set
	GetTotal() num.Amount
}

// NewTotal initiates a new total instance.
func NewTotal(zero num.Amount) *Total {
	t := new(Total)
	t.Categories = make([]*CategoryTotal, 0)
	t.Sum = zero
	return t
}

// newCategoryTotal prepares a category total calculation.
func newCategoryTotal(c *Combo, zero num.Amount) *CategoryTotal {
	ct := new(CategoryTotal)
	ct.Code = c.Category
	ct.Rates = make([]*RateTotal, 0)
	ct.Base = zero
	ct.Amount = zero
	ct.Retained = c.category.Retained
	return ct
}

// newRateTotal returns a rate total.
func newRateTotal(c *Combo, zero num.Amount) *RateTotal {
	rt := new(RateTotal)
	rt.Key = c.Rate // may be empty!
	rt.Percent = c.Percent
	rt.Base = zero
	rt.Amount = zero
	return rt
}

// Category provides the category total for the matching code.
func (t *Total) Category(code Code) *CategoryTotal {
	for _, ct := range t.Categories {
		if ct.Code == code {
			return ct
		}
	}
	return nil
}

// Rate grabs the matching rate from the category total, or nil.
func (ct *CategoryTotal) Rate(key Key) *RateTotal {
	for _, rt := range ct.Rates {
		if rt.Key == key {
			return rt
		}
	}
	return nil
}

// Calculate figures out the total taxes for the set of `TaxableLine`s provided.
func (t *Total) Calculate(reg *Region, lines []TaxableLine, taxIncluded Code, date cal.Date, zero num.Amount) error {
	if reg == nil {
		return ErrMissingRegion
	}

	// get a simplified list of lines we can manipulate if needed
	taxLines := mapTaxLines(lines)

	// First, prepare all tax combos with the region and date details
	for _, tl := range taxLines {
		for _, c := range tl.taxes {
			if err := c.prepare(reg, date); err != nil {
				return err
			}
		}
	}

	// If prices include a tax, perform a pre-loop to update all the line prices with
	// the price minus the defined tax. To help reduce the risk of rounding errors,
	// we'll add an extra couple of 0s.
	if !taxIncluded.IsEmpty() {
		for _, tl := range taxLines {
			if rate := tl.rateForCategory(taxIncluded); rate != "" {
				c := tl.taxes.Get(taxIncluded)
				if c.category.Retained {
					return ErrInvalidPricesInclude.WithMessage("cannot include retained category '%s'", taxIncluded.String())
				}

				// update the price scale, add two 0s, this will be removed later.
				tl.price = tl.price.Rescale(tl.price.Exp() + 2)
				tl.price = tl.price.Subtract(c.Percent.From(tl.price))
			}
		}
	}

	// Go through each line and add the price to the base of each tax
	for _, tl := range taxLines {
		for _, c := range tl.taxes {
			rt := t.rateTotalFor(c, zero)
			rt.Base = rt.Base.MatchPrecision(tl.price)
			rt.Base = rt.Base.Add(tl.price)
		}
	}

	// Now go through each category to apply the percentage and calculate the final sums
	t.Sum = zero
	for _, ct := range t.Categories {
		ct.calculate(zero)
		if ct.Retained {
			t.Sum = t.Sum.Subtract(ct.Amount)
		} else {
			t.Sum = t.Sum.Add(ct.Amount)
		}
	}

	return nil
}

// calculate goes through each rate defined inside the category, ensures
// the amounts are correct, and adds each to the category base.
func (ct *CategoryTotal) calculate(zero num.Amount) {
	ct.Base = zero
	ct.Amount = zero
	for _, rt := range ct.Rates {
		rt.Amount = rt.Percent.Of(rt.Base).Rescale(zero.Exp())
		rt.Base = rt.Base.Rescale(zero.Exp())
		ct.Base = ct.Base.Add(rt.Base)
		ct.Amount = ct.Amount.Add(rt.Amount)
	}
}

// rateTotalFor either finds of creates total objects for the category and rate.
// May error if we detect and incorrect combination.
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
		if rt.Percent.Equals(c.Percent) {
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

// taxLine is used to replace
type taxLine struct {
	price num.Amount
	taxes Set
}

func (tl *taxLine) rateForCategory(code Code) Key {
	return tl.taxes.Rate(code)
}

func mapTaxLines(lines []TaxableLine) []*taxLine {
	tls := make([]*taxLine, len(lines))
	for i, v := range lines {
		tls[i] = &taxLine{
			price: v.GetTotal(),
			taxes: v.GetTaxes(),
		}
	}
	return tls
}
