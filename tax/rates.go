package tax

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Rates contains a list of taxes, usually applied to an
// invoice line or item.
type Rates []*Rate

// Rate references the tax category and rate code that should be applied to
// this line item when calculating the final taxes.
type Rate struct {
	Category Code `json:"cat" jsonschema:"title=Category Code,description=From the available options for the region."`
	Code     Code `json:"code" jsonschema:"title=Code,description=As defined for the region and category."`
}

// Validate ensures the Rate contains all the details required.
func (r *Rate) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Category, validation.Required),
		validation.Field(&r.Code, validation.Required),
	)
}

// Validate ensures the rates array looks correct.
func (rs Rates) Validate() error {
	combos := make(map[Code]Code)
	for i, r := range rs {
		if _, ok := combos[r.Category]; ok {
			return fmt.Errorf("%d: category %v can only be defined once per line", i, r.Category)
		}
		if err := r.Validate(); err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
		combos[r.Category] = r.Code
	}
	return nil
}

// Equals returns true if the array of rates match, regardless of order.
func (rs Rates) Equals(rs2 Rates) bool {
	for _, a := range rs {
		match := false
		for _, b := range rs2 {
			if a.Category == b.Category && a.Code == b.Code {
				match = true
			}
		}
		if !match {
			// implies the code defined in the base, was not present in the second
			// array.
			return false
		}
	}
	return true
}
