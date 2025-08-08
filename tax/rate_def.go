package tax

import (
	"context"
	"errors"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// RateDef defines a single rate inside a category
type RateDef struct {
	// Rate defines the key for which this rate applies.
	Rate cbc.Key `json:"rate" jsonschema:"title=Rate"`

	// Keys identifies the set of tax keys defined in the category that this
	// rate can be used with.
	Keys []cbc.Key `json:"keys,omitempty" jsonschema:"title=Keys"`

	// Human name of the rate
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Useful description of the rate.
	Description i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Values contains a list of Value objects that contain the
	// current and historical percentage values for the rate and
	// additional filters.
	// Order is important, newer values should come before
	// older values.
	Values []*RateValueDef `json:"values,omitempty" jsonschema:"title=Values"`

	// Meta contains additional information about the rate that is relevant
	// for local frequently used implementations.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// RateValueDef contains a percentage rate or fixed amount for a given date range.
// Fiscal policy changes mean that rates are not static so we need to
// be able to apply the correct rate for a given period.
type RateValueDef struct {
	// Only apply this rate if one of the tags is present in the invoice.
	// Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tags"`
	// Ext map of keys that can be used to filter to determine if the rate applies.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Date from which this value should be applied.
	Since *cal.Date `json:"since,omitempty" jsonschema:"title=Since"`
	// Percent rate that should be applied
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent"`
	// An additional surcharge to apply.
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`
	// When true, this value should no longer be used.
	Disabled bool `json:"disabled,omitempty" jsonschema:"title=Disabled"`
}

// ValidateWithContext checks that our tax definition is valid. This is only really
// meant to be used when testing new regional tax definitions.
func (r *RateDef) ValidateWithContext(ctx context.Context) error {
	err := validation.ValidateStructWithContext(ctx, r,
		validation.Field(&r.Rate, validation.Required),
		validation.Field(&r.Keys),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Values,
			validation.By(checkRateValuesOrder),
		),
		validation.Field(&r.Meta),
	)
	return err
}

// Validate ensures the tax rate contains all the required fields.
func (rv *RateValueDef) Validate() error {
	return validation.ValidateStruct(rv,
		validation.Field(&rv.Percent, validation.Required),
	)
}

// HasKey returns true if we make a match for the provided key, or if
// rate has no keys and the provided key is empty.
func (r *RateDef) HasKey(key cbc.Key) bool {
	if r == nil {
		return false
	}
	if len(r.Keys) == 0 && key.IsEmpty() {
		return true
	}
	for _, k := range r.Keys {
		if k == key {
			return true
		}
	}
	return false
}

// Value determines the tax rate value for the provided date and zone, if applicable.
func (r *RateDef) Value(date cal.Date, ext Extensions) *RateValueDef {
	for _, rv := range r.Values {
		if len(rv.Ext) > 0 {
			if !ext.Contains(rv.Ext) {
				continue
			}
		}
		if rv.Since == nil || !rv.Since.IsValid() || rv.Since.Before(date.Date) {
			return rv
		}
	}
	return nil
}

func checkRateValuesOrder(list interface{}) error {
	values, ok := list.([]*RateValueDef)
	if !ok {
		return errors.New("must be a tax rate value array")
	}
	var date *cal.Date
	// loop through and check order of Since value
	for i := range values {
		v := values[i]
		if len(v.Ext) > 0 {
			// TODO: check extensions order also
			// Not too important at the moment.
			continue
		}
		if date != nil && date.IsValid() {
			if v.Since.IsValid() && !v.Since.Before(date.Date) {
				return errors.New("invalid date order")
			}
		}
		date = v.Since
	}
	return nil
}
