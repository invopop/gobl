package tax

import (
	"errors"

	"cloud.google.com/go/civil"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/invopop/gobl/num"
)

// RegionID defines a unique code used to identify each region.
type RegionID string

// Category defines a grouping of taxes whereby only one
// definition inside a tax category can be applied to a given
// invoice line.
type Category string

// Code defines a simple code used to identify a tax
// within a category.
type Code string

// Defs is used to define a list of tax definitions.
type Defs []Def

// ByCode finds the matching tax definition by it's code, or
// returns nil if it is not recognised.
func (tds Defs) ByCode(code Code) *Def {
	for i := range tds {
		td := &tds[i]
		if td.Code == code {
			return td
		}
	}
	return nil
}

// Def defines a tax combination of code and rate.
type Def struct {
	Name     string   `json:"name" jsonschema:"title=Name"`
	Desc     string   `json:"desc,omitempty" jsonschema:"title=Description"`
	Category Category `json:"category"`

	// Code represents an internal unique reference for the tax definition.
	// Once set, the codes must not be changed.
	Code Code `json:"code" jsonschema:"title=Code"`

	// Rates contains a list of TaxRate objects that contain the
	// current and historical tax rates for the current category
	// and code. Order is important, newer rates should come before
	// older rates.
	Rates []Rate `json:"rates" jsonschema:"title=Rate"`

	// Retained when true implies that the tax amount will be retained
	// by the buyer on behalf of the supplier, and thus subtracted from
	// the invoice total to be paid.
	Retained bool `json:"retained,omitempty"`
}

// Validate checks that our tax definition is valid. This is only really
// meant to be used when testing new regional tax definitions.
func (td Def) Validate() error {
	err := validation.ValidateStruct(&td,
		validation.Field(&td.Category, validation.Required),
		validation.Field(&td.Code, validation.Required),
		validation.Field(&td.Name, validation.Required),
		validation.Field(&td.Rates, validation.Required, validation.By(checkTaxRateOrder)),
	)
	return err
}

// Rate contains a percentage tax rate for a given date range.
// Fiscal policy changes meen that rates are not fixed so we need to
// be able to apply the correct rate for a given period.
type Rate struct {
	Since    civil.Date     `json:"since,omitempty"`
	Value    num.Percentage `json:"value"`
	Disabled bool           `json:"disabled,omitempty"`
}

// Validate ensures the tax rate contains all the required fields.
func (tr Rate) Validate() error {
	return validation.ValidateStruct(&tr,
		validation.Field(&tr.Value, validation.Required),
	)
}

func checkTaxRateOrder(value interface{}) error {
	trs, ok := value.([]Rate)
	if !ok {
		return errors.New("must be a tax rate array")
	}
	var d civil.Date
	// loop through and check order of Since value
	for i := range trs {
		r := &trs[i]
		if d.IsValid() {
			if r.Since.IsValid() && !r.Since.Before(d) {
				return errors.New("invalid tax rate since date order")
			}
		}
		d = r.Since
	}
	return nil
}

// RateOn determines the tax rate for the provided date.
func (td *Def) RateOn(date civil.Date) num.Amount {

	return num.Amount{}
}
