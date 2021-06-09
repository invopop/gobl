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
	Category Code           `json:"cat" jsonschema:"title=Category Code"`
	Code     Code           `json:"code" jsonschema:"title=Code"`
	Base     num.Amount     `json:"base" jsonschema:"title=Base,description=Base value to which taxes are added"`
	Percent  num.Percentage `json:"percent" jsonschema:"title=Percentage"`
	Value    num.Amount     `json:"value" jsonschema:"title=Value,description=The amount of tax applied"`
	Retained bool           `json:"retained,omitempty" jsonschema:"title=Retained,description=True when this tax is retained by the client."`
}
