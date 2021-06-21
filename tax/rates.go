package tax

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
)

// Rates contains a list of taxes, usually applied to an
// invoice line or item.
type Rates []*Rate

// Rate shows the type of tax and percentage that should be applied and
// represented in the tax totals.
type Rate struct {
	Category Code           `json:"cat" jsonschema:"title=Category Code"`
	Code     Code           `json:"code" jsonschema:"title=Code"`
	Value    num.Percentage `json:"value" jsonschema:"title=Value"`
	//Base     num.Amount     `json:"base" jsonschema:"title=Base,description=Base value to which taxes are added"`
	//Value    num.Amount     `json:"value" jsonschema:"title=Value,description=The amount of tax applied"`
	//Retained bool           `json:"retained,omitempty" jsonschema:"title=Retained,description=True when this tax is retained by the client."`
}

// Validate ensures the Rate contains all the details required.
func (r *Rate) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Category, validation.Required),
		validation.Field(&r.Code, validation.Required),
		validation.Field(&r.Value, validation.Required),
		// validation.Field(&r.Percent, validation.Required),
		// validation.Field(&r.Base, validation.Required),
	)
}
