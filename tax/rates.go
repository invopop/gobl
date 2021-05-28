package tax

import (
	"github.com/invopop/gobl/num"
)

// Rates contains a list of taxes, usually applied to an
// invoice line or item.
type Rates []Rate

// Rate shows the type of tax, rate, and base that should be applied and
// represented in the tax totals.
type Rate struct {
	CategoryCode Code           `json:"cat_code" jsonschema:"title=Category Code"`
	RateCode     Code           `json:"rate_code" jsonschema:"title=Rate Code"`
	Base         num.Amount     `json:"base" jsonschema:"title=Base"`
	Percent      num.Percentage `json:"percent" jsonschema:"title=Percentage"`
	Retained     bool           `json:"retained,omitempty" jsonschema:"title=Retained,description=True when this tax is retained by the client."`
}

// RateTotal contains a sum of all the tax rates in the document with
// a matching category and definition. Def field is redundant, but serves as a copy of the
// data used for calculations.
type RateTotal struct {
	Code    Code           `json:"code"`
	Base    num.Amount     `json:"base"`
	Percent num.Percentage `json:"percent"`
	Value   num.Amount     `json:"value"`
}

// CategoryTotal contains the calculation of all the taxes
// with a matching category.
type CategoryTotal struct {
	Code     Code        `json:"category"`
	Rates    []RateTotal `json:"rates"`
	Base     num.Amount  `json:"base"`
	Value    num.Amount  `json:"value"`
	Retained bool        `json:"retained,omitempty"`
}

// Totals contains the final calculations of all the tax categories
// and sub-codes into a grand total. The sum of the tax total is
// most likely to be the final amount payable.
type Totals struct {
	Categories []CategoryTotal `json:"categories"`
	Base       num.Amount      `json:"base"`
	Value      num.Amount      `json:"value"`
	Sum        num.Amount      `json:"sum" jsonschema:"title=Sum,description=Total of Base plus Value"`
}
