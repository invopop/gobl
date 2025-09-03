package tax

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Combo represents the tax combination of a category code and rate key. The percent
// and retained attributes will be determined automatically from the Rate key if set
// during calculation.
type Combo struct {
	// Country code override when issuing with taxes applied from different countries.
	Country l10n.TaxCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Key helps determine the tax situation within the category.
	Key cbc.Key `json:"key,omitempty"`
	// Rate within a category and for a given key to apply.
	Rate cbc.Key `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Percent defines the percentage set manually or determined from the
	// key. A nil percent implies that this tax combo is either exempt or not-subject.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent" jsonschema_extras:"calculated=true"`
	// Some countries require an additional surcharge (may be determined if key present).
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge" jsonschema_extras:"calculated=true"`
	// Local codes that apply for a given rate or percentage that need to be identified and validated.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Copied from the category definition, implies this tax combo is retained
	retained bool `json:"-"`
	// Copied from the category definition, implies this tax combo is informative
	informative bool `json:"-"`
}

// ValidateWithContext ensures the Combo has the correct details.
func (c *Combo) ValidateWithContext(ctx context.Context) error {
	// First perform combo validation with the regime from the context,
	// or the country override.

	var r *RegimeDef
	if c.Country.Empty() {
		r = RegimeDefFromContext(ctx)
	} else {
		r = RegimeDefFor(c.Country.Code())
	}
	return ValidateStructWithContext(ctx, c,
		validation.Field(&c.Category,
			validation.Required,
			r.InCategories(),
		),
		validation.Field(&c.Key,
			r.InCategoryKeys(c.Category),
		),
		validation.Field(&c.Rate,
			r.InCategoryRates(c.Category, c.Key),
		),
		validation.Field(&c.Percent,
			r.RequiresPercent(c.Category, c.Key),
		),
		validation.Field(&c.Surcharge, validation.When(
			c.Percent == nil,
			validation.Nil.Error("required with percent"),
		)),
		validation.Field(&c.Ext),
	)
}

// Normalize tries to normalize the data inside the tax combo.
func (c *Combo) Normalize(normalizers Normalizers) {
	if c == nil {
		return
	}

	switch c.Category {
	case CategoryVAT:
		switch c.Rate {
		case KeyZero:
			c.Key = KeyZero
			c.Rate = cbc.KeyEmpty
			if c.Percent == nil {
				c.Percent = num.NewPercentage(0, 2)
			}
		case KeyExempt:
			c.Key = KeyExempt
			c.Rate = cbc.KeyEmpty
		case KeyExempt.With("reverse-charge"):
			c.Key = KeyReverseCharge
			c.Rate = cbc.KeyEmpty
			c.Percent = nil
		case KeyExempt.With("export"):
			c.Key = KeyExport
			c.Rate = cbc.KeyEmpty
		case KeyExempt.With("eea"), KeyExempt.With("export").With("eea"):
			c.Key = KeyIntraCommunity
			c.Rate = cbc.KeyEmpty
		default:
			// Make no further assumptions about the key, but try to replace standard
			// rate with general.
			if c.Rate == KeyStandard {
				c.Rate = RateGeneral
			} else if found, ok := strings.CutPrefix(c.Rate.String(), "standard+"); ok {
				c.Rate = cbc.Key(RateGeneral.String() + "+" + found)
			}
		}

		if c.Key == cbc.KeyEmpty {
			// Special case for zero percent which has no additional rates
			if c.Percent != nil && c.Percent.IsZero() {
				c.Key = KeyZero
			}
		}
	}

	c.Ext = CleanExtensions(c.Ext)
	normalizers.Each(c)
}

func (c *Combo) calculate(country l10n.TaxCountryCode, date cal.Date) error {
	if c.Country == country {
		c.Country = ""
	} else if c.Country != "" {
		country = c.Country
	}

	r := RegimeDefFor(country.Code())
	cd := r.CategoryDef(c.Category) // may provide global category
	if r != nil && cd == nil {
		return ErrInvalidCategory.WithMessage("'%s' not defined in regime", c.Category.String())
	} else if cd == nil {
		return nil // no category, nothing to do
	}

	c.retained = cd.Retained
	c.informative = cd.Informative

	// If there are keys defined for the category, but the combo does not
	// have a key, then we will use the standard key.
	if len(cd.Keys) > 0 && c.Key == cbc.KeyEmpty {
		c.Key = KeyStandard
	}

	// If there is an associated key definition that does not expect percentages,
	// clear them to avoid unexpected output.
	if kd := cd.KeyDef(c.Key); kd != nil && kd.NoPercent {
		c.Percent = nil
		c.Surcharge = nil
		return nil
	}

	return c.prepareRate(cd, date)
}

// prepare updates the Combo object's Percent and Retained properties using the base totals
// as a source of additional data for making decisions.
func (c *Combo) prepareRate(cd *CategoryDef, date cal.Date) error {
	// If there is no rate for the combo, there isn't much else we can prepare.
	if c.Rate == cbc.KeyEmpty {
		return nil
	}

	rate := cd.RateDef(c.Key, c.Rate)
	if rate == nil {
		if c.Key == cbc.KeyEmpty {
			return ErrInvalid.WithMessage("'%s' rate not defined in category '%s'", c.Rate.String(), c.Category.String())
		}
		return ErrInvalid.WithMessage("'%s' rate not defined for key '%s' in category '%s'", c.Rate.String(), c.Key.String(), c.Category.String())
	}

	// if there are no rate values, don't attempt to prepare anything else.
	if len(rate.Values) == 0 {
		return nil
	}

	value := rate.Value(date, c.Ext)
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

// UnmarshalJSON is a migration helper that will prepare the Combo's
// key from either the old tags or rate fields.
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
	if c.Rate == cbc.KeyEmpty && len(aux.Tags) > 0 {
		c.Rate = aux.Tags[0]
	}
	return nil
}

// JSONSchemaExtend will extend the JSON schema for the Combo object with
// global tax category keys.
func (c Combo) JSONSchemaExtend(s *jsonschema.Schema) {
	for _, cd := range globalCategories {
		s.AnyOf = append(s.AnyOf, c.jsonSchemaBuildCategory(cd))
	}
}

func (Combo) jsonSchemaBuildCategory(cd *CategoryDef) *jsonschema.Schema {
	// Build the JSON schema for the category definition.
	s := new(jsonschema.Schema)
	s.If = &jsonschema.Schema{
		Properties: jsonschema.NewProperties(),
	}
	s.If.Properties.Set("cat", &jsonschema.Schema{
		Const: cd.Code.String(),
	})
	oneOf := make([]*jsonschema.Schema, len(cd.Keys))
	for i, kd := range cd.Keys {
		oneOf[i] = &jsonschema.Schema{
			Const: kd.Key.String(),
			Title: kd.Name.String(),
		}
	}
	// Add the key definitions to the schema
	s.Then = &jsonschema.Schema{
		Properties: jsonschema.NewProperties(),
	}
	s.Then.Properties.Set("key", &jsonschema.Schema{
		OneOf: oneOf,
	})
	return s
}
