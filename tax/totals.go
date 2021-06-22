package tax

import (
	"fmt"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
)

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and definition.
type RateTotal struct {
	Code    Code           `json:"code" jsonschema:"title=Code"`
	Base    num.Amount     `json:"base" jsonschema:"title=Base"`
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent"`
	Value   num.Amount     `json:"value" jsonschema:"title=Value"`
	sum     num.Amount     `json:"-"` // used for internal calculations when tax included
}

// CategoryTotal groups together all rates inside a given category.
type CategoryTotal struct {
	Code     Code         `json:"code" jsonschema:"title=Code"`
	Retained bool         `json:"retained,omitempty" jsonschema:"title=Retained"`
	Rates    []*RateTotal `json:"rates" jsonschema:"title=Rates"`
	Base     num.Amount   `json:"base" jsonschema:"title=Base"`
	Value    num.Amount   `json:"value" jsonschema:"title=Value"`
}

// Total contains a set of Category Totals which in turn
// contain all the accumulated taxes contained in the document. The resulting
// `sum` is that value that should be added to the payable total.
type Total struct {
	Categories []*CategoryTotal `json:"categories,omitempty" jsonschema:"title=Categories"`
	Sum        num.Amount       `json:"sum" jsonschema:"title=Sum,description=Total value of all the taxes to be added or retained."`
}

// TaxableLine defines what we expect from a line in order to subsequently calculate
// the taxes that need to be added or retained.
type TaxableLine interface {
	GetTaxRates() Rates
	GetTotal() num.Amount
}

// NewTotal initiates a new total instance.
func NewTotal(zero num.Amount) *Total {
	t := new(Total)
	t.Categories = make([]*CategoryTotal, 0)
	t.Sum = zero
	return t
}

// NewCategoryTotal prepares a category total calculation.
func NewCategoryTotal(code Code, retained bool, zero num.Amount) *CategoryTotal {
	ct := new(CategoryTotal)
	ct.Code = code
	ct.Rates = make([]*RateTotal, 0)
	ct.Base = zero
	ct.Value = zero
	ct.Retained = retained
	return ct
}

func NewRateTotal(code Code, percent num.Percentage, zero num.Amount) *RateTotal {
	rt := new(RateTotal)
	rt.Code = code
	rt.Percent = percent
	rt.Base = zero
	rt.Value = zero
	rt.sum = zero
	return rt
}

// calculate goes through each rate defined inside the category
func (ct *CategoryTotal) calculate(zero num.Amount) {
	ct.Base = zero
	ct.Value = zero
	for _, rt := range ct.Rates {
		ct.Base = ct.Base.Add(rt.Base)
		ct.Value = ct.Value.Add(rt.Value)
	}
}

// Calculate figures out the total taxes for the set of `TaxableLine`s provided.
func (t *Total) Calculate(reg *Region, lines []TaxableLine, taxIncluded bool, date org.Date, zero num.Amount) error {
	// NOTE: This method looks more complex than it could be as we're providing
	// additional logic that will deal with situations whereby taxes are included
	// in line prices potentially alongside retained taxes. Retained taxes cannot be
	// included under the "price includes tax" umbrella, so we need to use the base
	// calculated after removing regular VAT or Sales taxes.
	// A typical use case for this is in Spain whereby regular VAT needs to be applied
	// alongside IRPF (income tax) which is retained by the client.
	// As a general rule, invoice taxes must always be calculated at the last possible
	// moment to avoid accumulating rounding errors

	// Group all the line tax combinations together
	groups := groupLines(lines, zero)

	// For each group, figure out the categories and then bases, first for the regular
	// taxes, then for the retained taxes so that if taxes are included in prices, we
	// have a base price that retained taxes can be applied to.
	for _, g := range groups {
		base := zero
		// determine base from non-retained taxes
		for _, r := range g.rates {
			ct, rt, err := t.categoryAndRateTotals(reg, r, date, zero)
			if err != nil {
				return err
			}
			if ct.Retained {
				continue
			}
			if taxIncluded {
				rt.sum = rt.sum.Add(g.sum)
				rt.Value = rt.Percent.From(rt.sum)
				rt.Base = rt.sum.Subtract(rt.Value)
			} else {
				rt.Base = rt.Base.Add(g.sum)
				rt.Value = rt.Percent.Of(rt.Base)
			}
			base = base.Add(rt.Base)
		}
		// use base for retained taxes
		for _, r := range g.rates {
			ct, rt, err := t.categoryAndRateTotals(reg, r, date, zero)
			if err != nil {
				return err
			}
			if !ct.Retained {
				continue
			}
			rt.Base = rt.Base.Add(base)
			rt.Value = rt.Percent.Of(rt.Base)
		}
	}

	t.Sum = zero
	for _, ct := range t.Categories {
		ct.calculate(zero)
		if ct.Retained {
			t.Sum = t.Sum.Subtract(ct.Value)
		} else if !taxIncluded {
			t.Sum = t.Sum.Add(ct.Value)
		}
	}

	return nil
}

func (t *Total) categoryAndRateTotals(reg *Region, r *Rate, date org.Date, zero num.Amount) (*CategoryTotal, *RateTotal, error) {
	cat, ok := reg.Category(r.Category)
	if !ok {
		return nil, nil, fmt.Errorf("failed to find category, invalid code: %v", r.Category)
	}
	def, ok := cat.Def(r.Code)
	if !ok {
		return nil, nil, fmt.Errorf("failed to find rate definition, invalid code: %v", r.Code)
	}
	val, ok := def.On(date)
	if !ok {
		return nil, nil, fmt.Errorf("tax rate cannot be provided for date")
	}

	var catTotal *CategoryTotal
	for _, ct := range t.Categories {
		if ct.Code == r.Category {
			catTotal = ct
			break
		}
	}
	if catTotal == nil {
		catTotal = NewCategoryTotal(r.Category, cat.Retained, zero)
		t.Categories = append(t.Categories, catTotal)
	}

	// Prepare the Rate
	var rateTotal *RateTotal
	for _, rt := range catTotal.Rates {
		if rt.Code == r.Code {
			rateTotal = rt
			break
		}
	}
	if rateTotal == nil {
		rateTotal = NewRateTotal(r.Code, val.Percent, zero)
		catTotal.Rates = append(catTotal.Rates, rateTotal)
	}

	return catTotal, rateTotal, nil
}

type lineGroups []*lineGroup

type lineGroup struct {
	rates Rates
	sum   num.Amount
}

func (lgs lineGroups) find(trs Rates) *lineGroup {
	// find an existing tax combo
	for _, lg := range lgs {
		if lg.rates.Equals(trs) {
			return lg
		}
	}
	return nil
}

func groupLines(lines []TaxableLine, zero num.Amount) lineGroups {
	// group all the rates together to form a set of rate
	// combinations and totals
	lgs := make(lineGroups, 0)
	for _, l := range lines {
		lg := lgs.find(l.GetTaxRates())
		if lg == nil {
			lg = new(lineGroup)
			lg.rates = l.GetTaxRates()
			lg.sum = zero
			lgs = append(lgs, lg)
		}
		lg.sum = lg.sum.Add(l.GetTotal())
	}
	return lgs
}
