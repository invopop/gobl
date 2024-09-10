package tax

import (
	"context"
	"encoding/json"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// Combo represents the tax combination of a category code and rate key. The percent
// and retained attributes will be determined automatically from the Rate key if set
// during calculation.
type Combo struct {
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Country code override when issuing with taxes applied from different countries.
	Country l10n.TaxCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// Rate within a category to apply.
	Rate cbc.Key `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Percent defines the percentage set manually or determined from the rate
	// key (calculated if rate present). A nil percent implies that this tax combo
	// is **exempt** from tax.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent" jsonschema_extras:"calculated=true"`
	// Some countries require an additional surcharge (calculated if rate present).
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge" jsonschema_extras:"calculated=true"`
	// Local codes that apply for a given rate or percentage that need to be identified and validated.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Copied from the category definition, implies this tax combo is retained
	retained bool `json:"-"`
}

// ValidateWithContext ensures the Combo has the correct details.
func (c *Combo) ValidateWithContext(ctx context.Context) error {
	// First perform combo validation with the regime from the context,
	// or the country override.
	var r *Regime
	if c.Country.Empty() {
		r, _ = ctx.Value(KeyRegime).(*Regime)
	} else {
		r = RegimeFor(c.Country.Code())
	}
	err := validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Category,
			validation.Required,
			r.InCategories(),
		),
		validation.Field(&c.Rate,
			r.InCategoryRates(c.Category),
		),
		validation.Field(&c.Percent),
		validation.Field(&c.Surcharge, validation.When(
			c.Percent == nil,
			validation.Nil.Error("required with percent"),
		)),
	)
	if err != nil {
		return err
	}

	// Always use base regime to check extensions
	r, _ = ctx.Value(KeyRegime).(*Regime)
	cat := r.Category(c.Category)
	exts := make([]cbc.Key, 0)
	if cat != nil {
		exts = append(exts, cat.Extensions...)
	}
	err = validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Ext,
			ExtensionsHas(exts...),
			validation.When(
				len(exts) == 0,
				validation.Empty,
				validation.Skip,
			),
		),
	)
	if err != nil {
		return err
	}
	if cat != nil && cat.Validation != nil {
		if err := cat.Validation(c); err != nil {
			return err
		}
	}
	return r.ValidateObject(c)
}

// NormalizeCombo tries to normalize the data inside the tax combo.
func NormalizeCombo(c *Combo) *Combo {
	if c == nil {
		return nil
	}
	c.Ext = NormalizeExtensions(c.Ext)
	return c
}

func (c *Combo) calculate(country l10n.TaxCountryCode, tags []cbc.Key, date cal.Date) error {
	if c.Country == country {
		c.Country = ""
	} else if c.Country != "" {
		country = c.Country
	}

	r := RegimeFor(country.Code())
	if r == nil {
		// if the tax regime is not yet defined, don't try to perform
		// any extra calculations.
		return nil
	}

	return c.calculateForRegime(r, tags, date)
}

func (c *Combo) calculateForRegime(r *Regime, tags []cbc.Key, date cal.Date) error {
	category := r.Category(c.Category)
	if category == nil {
		return ErrInvalidCategory.WithMessage("'%s' not defined in regime", c.Category.String())
	}
	c.retained = category.Retained

	if err := c.prepareRate(category, tags, date); err != nil {
		return err
	}

	// Run the regime's calculations and normalisations, but only if this is not
	// a country-specific combo.
	if c.Country == "" {
		return r.CalculateObject(c)
	}

	return nil
}

// prepare updates the Combo object's Percent and Retained properties using the base totals
// as a source of additional data for making decisions.
func (c *Combo) prepareRate(category *Category, tags []cbc.Key, date cal.Date) error {
	// If there is no rate for the combo, there isn't much else we can do.
	if c.Rate == cbc.KeyEmpty {
		return nil
	}

	rate := category.Rate(c.Rate)
	if rate == nil {
		return ErrInvalidRate.WithMessage("'%s' rate not defined in category '%s'", c.Rate.String(), c.Category.String())
	}

	// Copy over the predefined extensions from the rate to the combo.
	if c.Country == "" && len(rate.Ext) > 0 {
		if c.Ext == nil {
			c.Ext = make(Extensions)
		}
		for k, v := range rate.Ext {
			c.Ext[k] = v
		}
	}

	if rate.Exempt {
		c.Percent = nil
		c.Surcharge = nil
		return nil
	}

	// if there are no rate values, don't attempt to prepare anything else.
	if len(rate.Values) == 0 {
		return nil
	}

	value := rate.Value(date, tags, c.Ext)
	if value == nil {
		return ErrInvalidDate.WithMessage("rate value unavailable for '%s' in '%s' on '%s'", c.Rate.String(), c.Category.String(), date.String())
	}

	p := value.Percent // copy
	c.Percent = &p

	if value.Surcharge != nil {
		s := *value.Surcharge // copy
		c.Surcharge = &s
	} else {
		c.Surcharge = nil
	}

	return nil
}

// UnmarshalJSON is a temporary migration helper that will move the
// first of the "tags" array used in earlier versions of GOBL into
// the rate field.
func (c *Combo) UnmarshalJSON(data []byte) error {
	type Alias Combo
	aux := struct {
		*Alias
		Tags []cbc.Key `json:"tags"`
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Tags) > 0 && c.Rate == cbc.KeyEmpty {
		c.Rate = aux.Tags[0]
	}
	return nil
}
