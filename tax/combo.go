package tax

import (
	"encoding/json"
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
)

// Combo represents the tax combination of a category code and rate key. The percent
// and retained attributes will be determined automatically from the Rate key if set
// during calculation.
type Combo struct {
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Country code override when issuing with taxes applied from different countries.
	Country l10n.TaxCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
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
	Ext Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`

	// Copied from the category definition, implies this tax combo is retained
	retained bool `json:"-"`
	// Copied from the category definition, implies this tax combo is informative
	informative bool `json:"-"`
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
			// This can cause problems with backwards compatibility as the "exempt"
			// rate was used too widely. Addons will need to try and account for this.
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

		switch c.Key {
		case cbc.KeyEmpty:
			// Special case for zero percent which has no additional rates
			if c.Percent != nil && c.Percent.IsZero() {
				c.Key = KeyZero
			}
		case KeyZero:
			if c.Percent == nil {
				zp := num.PercentageZero
				c.Percent = &zp
			}
		}
	}

	c.Ext = c.Ext.Clean()
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

func comboRules() *rules.Set {
	return rules.For(new(Combo),
		rules.Assert("01", "tax category not valid for regime",
			is.FuncContext("category in regime", comboCategoryValid),
		),
		rules.Assert("02", "tax combo key not valid for category in regime",
			is.FuncContext("key in category", comboKeyValid),
		),
		rules.Assert("03", "tax combo rate not valid for key in regime",
			is.FuncContext("rate in key", comboRateValid),
		),
		rules.Assert("04", "tax combo percent required or invalid for key in regime",
			is.FuncContext("percent valid for key", comboPercentValid),
		),
		rules.Assert("05", "tax combo surcharge requires percent in regime",
			is.Func("surcharge requires percent", comboSurchargeValid),
		),
		rules.Assert("06", "tax combo extension key not defined in regime",
			is.Func("extension keys defined", comboExtensionsValid),
		),
	)
}

// regimeDefFromContext returns the RegimeDef from the validation context.
func regimeDefFromContext(ctx rules.Context) *RegimeDef {
	if r, ok := ctx.Value(regimeContextKey).(Regime); ok {
		return r.RegimeDef()
	}
	return nil
}

// regimeDefForCombo returns the RegimeDef for the combo, using the combo's
// Country override when set, otherwise falling back to the context regime.
func regimeDefForCombo(ctx rules.Context, combo *Combo) *RegimeDef {
	if combo != nil && !combo.Country.Empty() {
		return RegimeDefFor(combo.Country.Code())
	}
	return regimeDefFromContext(ctx)
}

func comboCategoryValid(ctx rules.Context, val any) bool {
	combo, ok := val.(*Combo)
	if !ok {
		return true
	}
	rd := regimeDefForCombo(ctx, combo)
	if rd == nil {
		return true
	}
	return rd.CategoryDef(combo.Category) != nil
}

func comboKeyValid(ctx rules.Context, val any) bool {
	combo, ok := val.(*Combo)
	if !ok || combo.Key == cbc.KeyEmpty {
		return true
	}
	rd := regimeDefForCombo(ctx, combo)
	if rd == nil {
		return true
	}
	cd := rd.CategoryDef(combo.Category)
	if cd == nil {
		return true
	}
	return cd.KeyDef(combo.Key) != nil
}

func comboRateValid(ctx rules.Context, val any) bool {
	combo, ok := val.(*Combo)
	if !ok || combo.Rate == cbc.KeyEmpty {
		return true
	}
	rd := regimeDefForCombo(ctx, combo)
	if rd == nil {
		return true
	}
	cd := rd.CategoryDef(combo.Category)
	if cd == nil {
		return true
	}
	return cd.RateDef(combo.Key, combo.Rate) != nil
}

func comboPercentValid(ctx rules.Context, val any) bool {
	combo, ok := val.(*Combo)
	if !ok {
		return true
	}
	rd := regimeDefForCombo(ctx, combo)
	if rd == nil {
		return true
	}
	cd := rd.CategoryDef(combo.Category)
	if cd == nil {
		return true
	}
	// No key: percent is always required.
	if combo.Key == cbc.KeyEmpty {
		return combo.Percent != nil
	}
	kd := cd.KeyDef(combo.Key)
	if kd == nil {
		return true // unknown key, skip percent check
	}
	if kd.NoPercent {
		return combo.Percent == nil
	}
	return combo.Percent != nil
}

func comboSurchargeValid(val any) bool {
	combo, ok := val.(*Combo)
	if !ok {
		return true
	}
	if combo.Surcharge != nil && combo.Percent == nil {
		return false
	}
	return true
}

func comboExtensionsValid(val any) bool {
	combo, ok := val.(*Combo)
	if !ok {
		return true
	}
	for _, k := range combo.Ext.Keys() {
		if ExtensionForKey(k) == nil {
			return false
		}
	}
	return true
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
