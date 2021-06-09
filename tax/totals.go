package tax

import (
	"errors"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
)

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and definition.
type RateTotal struct {
	Code    Code           `json:"code"`
	Base    num.Amount     `json:"base"`
	Percent num.Percentage `json:"percent"`
	Value   num.Amount     `json:"value"`
}

// CategoryTotal groups together a
type CategoryTotal struct {
	Code     Code         `json:"code"`
	Retained bool         `json:"retained,omitempty"`
	Rates    []*RateTotal `json:"rates"`
	Base     num.Amount   `json:"base"`
	Value    num.Amount   `json:"value"`
}

// Total contains a set of Category Totals which in turn
// contain all the accumulated taxes contained in the document.
type Total struct {
	sync.Mutex
	Categories []*CategoryTotal `json:"categories,omitempty"`
	Sum        num.Amount       `json:"sum" jsonschema:"title=Sum,description=Total value of all the taxes to be added or retained."`
}

// Validate ensures the Rate contains all the details required.
func (r *Rate) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Category, validation.Required),
		validation.Field(&r.Code, validation.Required),
		validation.Field(&r.Base, validation.Required),
		validation.Field(&r.Percent, validation.Required),
		validation.Field(&r.Value, validation.Required),
	)
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
	return rt
}

// AddRate makes it easier to add a new rate to the totals. It'll automatically
// handle splitting by category. A zero value is required so we know what to base
// calculations on.
func (t *Total) AddRate(r Rate, zero num.Amount) error {
	// Just in case we use this in multiple requests
	t.Lock()
	defer t.Unlock()

	// Prepare the category
	var cat *CategoryTotal
	for _, ct := range t.Categories {
		if ct.Code == r.Category {
			cat = ct
			if cat.Retained != r.Retained {
				return errors.New("category retained value does not match previous values")
			}
			break
		}
	}
	if cat == nil {
		cat = NewCategoryTotal(r.Category, r.Retained, zero)
		t.Categories = append(t.Categories, cat)
	}

	// Prepare the Rate
	var rate *RateTotal
	for _, rt := range cat.Rates {
		if rt.Code == r.Code {
			rate = rt
			if !rt.Percent.Equals(r.Percent) {
				return errors.New("rate percent does not match previous values")
			}
			break
		}
	}
	if rate == nil {
		rate = NewRateTotal(r.Code, r.Percent, zero)
		cat.Rates = append(cat.Rates, rate)
	}

	// Add the rate to the totals
	rate.Base = rate.Base.Add(r.Base)
	rate.Value = rate.Value.Add(r.Value)

	// Let's recalculate again
	cat.Calculate(zero)
	t.Calculate(zero)

	return nil
}

// Calculate goes through each rate defined inside the category
func (ct *CategoryTotal) Calculate(zero num.Amount) {
	ct.Base = zero
	ct.Value = zero
	for _, rt := range ct.Rates {
		ct.Base = ct.Base.Add(rt.Base)
		ct.Value = ct.Value.Add(rt.Value)
	}
}

// Calculate figures out how much total tax needs to be added or taken
// away from the resulting document. The resulting sum should be added to
// the invoice totals to reflect the final payment amount.
func (t *Total) Calculate(zero num.Amount) {
	t.Sum = zero
	for _, ct := range t.Categories {
		if ct.Retained {
			t.Sum = t.Sum.Subtract(ct.Value)
		} else {
			t.Sum = t.Sum.Add(ct.Value)
		}
	}
}
