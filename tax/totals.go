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
	Amount  num.Amount     `json:"amount" jsonschema:"title=Amount"`
}

// CategoryTotal groups together all rates inside a given category.
type CategoryTotal struct {
	Code     Code         `json:"code" jsonschema:"title=Code"`
	Retained bool         `json:"retained,omitempty" jsonschema:"title=Retained"`
	Rates    []*RateTotal `json:"rates" jsonschema:"title=Rates"`
	Base     num.Amount   `json:"base" jsonschema:"title=Base"`
	Amount   num.Amount   `json:"amount" jsonschema:"title=Amount"`
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
	ct.Amount = zero
	ct.Retained = retained
	return ct
}

// NewRateTotal returns a rate total.
func NewRateTotal(code Code, percent num.Percentage, zero num.Amount) *RateTotal {
	rt := new(RateTotal)
	rt.Code = code
	rt.Percent = percent
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
func (ct *CategoryTotal) Rate(code Code) *RateTotal {
	for _, rt := range ct.Rates {
		if rt.Code == code {
			return rt
		}
	}
	return nil
}

// Calculate figures out the total taxes for the set of `TaxableLine`s provided.
func (t *Total) Calculate(reg *Region, lines []TaxableLine, taxIncluded Code, date org.Date, zero num.Amount) error {
	// NOTE: This method looks more complex than it could be as we're providing
	// additional logic that will deal with situations whereby a tax is included
	// in line prices potentially with other taxes.
	//
	// A typical use case for this is in Spain whereby regular VAT needs to be applied
	// alongside IRPF (income tax) which is retained by the client.
	//
	// Tax surcharges (another very rare addition) are also not included in prices that
	// include tax.
	//
	// As a general rule, invoice taxes must always be calculated at the last possible
	// moment to avoid accumulating rounding errors.

	// get a simplified list of lines we can manipulate if needed
	taxLines := mapTaxLines(lines)

	// If prices include a tax, perform a pre-loop to update all the line prices with
	// the price minus the defined tax. To help reduce the risk of rounding errors,
	// we'll add an extra couple of 0s.
	if !taxIncluded.IsEmpty() {
		for _, tl := range taxLines {
			if rate := tl.rateForCategory(taxIncluded); rate != nil {
				c, err := reg.comboOn(rate, date)
				if err != nil {
					return err
				}
				if c.category.Retained {
					return fmt.Errorf("cannot include retained tax category '%v' in price", taxIncluded)
				}

				// update the price scale, add two 0s, this will be removed later.
				tl.price = tl.price.Rescale(tl.price.Exp() + 2)
				tl.price = tl.price.Subtract(c.value.Percent.From(tl.price))
			}
		}
	}

	// Go through each line and add the price to the base of each tax
	for _, tl := range taxLines {
		for _, r := range tl.rates {
			rt, err := t.rateTotalFor(reg, r, date, zero)
			if err != nil {
				return err
			}

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
func (t *Total) rateTotalFor(reg *Region, r *Rate, date org.Date, zero num.Amount) (*RateTotal, error) {
	c, err := reg.comboOn(r, date)
	if err != nil {
		return nil, err
	}

	var catTotal *CategoryTotal
	for _, ct := range t.Categories {
		if ct.Code == r.Category {
			catTotal = ct
			break
		}
	}
	if catTotal == nil {
		catTotal = NewCategoryTotal(r.Category, c.category.Retained, zero)
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
		rateTotal = NewRateTotal(r.Code, c.value.Percent, zero)
		catTotal.Rates = append(catTotal.Rates, rateTotal)
	}

	return rateTotal, nil
}

// taxLine is used to replace
type taxLine struct {
	price num.Amount
	rates Rates
}

func (tl *taxLine) rateForCategory(code Code) *Rate {
	for _, r := range tl.rates {
		if r.Category == code {
			return r
		}
	}
	return nil
}

func mapTaxLines(lines []TaxableLine) []*taxLine {
	tls := make([]*taxLine, len(lines))
	for i, v := range lines {
		tls[i] = &taxLine{
			price: v.GetTotal(),
			rates: v.GetTaxRates(),
		}
	}
	return tls
}
