package gobl

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Taxes contains a list of taxes, usually applied to an
// invoice line or item.
type Taxes []Tax

// Tax shows the type of tax, rate, and base that should be applied and
// represented in the tax totals.
type Tax struct {
	Category tax.Category   `json:"cat" jsonschema:"title=Category"`
	Code     tax.Code       `json:"code" jsonschema:"title=Code"`
	Base     num.Amount     `json:"base" jsonschema:"title=Base"`
	Rate     num.Percentage `json:"rate" jsonschema:"title=Rate"`
	Retained bool           `json:"retained,omitempty" jsonschema:"title=Retained,description=True when this tax is retained by the client."`
}

// TaxCodeTotal contains a sum of all the taxes in the document with
// a matching code. Rate is redundant, but serves as a copy of the
// data used for calculations.
type TaxCodeTotal struct {
	Code  tax.Code       `json:"code"`
	Base  num.Amount     `json:"base"`
	Rate  num.Percentage `json:"rate"`
	Value num.Amount     `json:"value"`
}

// TaxCategoryTotal contains the calculation of all the taxes
// with a matching category.
type TaxCategoryTotal struct {
	Category tax.Category   `json:"category"`
	Codes    []TaxCodeTotal `json:"codes"`
	Base     num.Amount     `json:"base"`
	Value    num.Amount     `json:"value"`
	Retained bool           `json:"retained,omitempty"`
}

// TaxTotal contains the final calculations of all the tax categories
// and sub-codes into a grand total. The sum of the tax total is
// most likely to be the final amount payable.
type TaxTotal struct {
	Categories []TaxCategoryTotal `json:"categories"`
	Base       num.Amount         `json:"base"`
	Value      num.Amount         `json:"value"`
	Sum        num.Amount         `json:"sum" jsonschema:"title=Sum,description=Total of Base plus Value"`
}
