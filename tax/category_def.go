package tax

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// CategoryDef contains the definition of a general type of tax inside a region.
type CategoryDef struct {
	// Code to be used in documents
	Code cbc.Code `json:"code" jsonschema:"title=Code"`

	// Short name of the category to be used instead of code in output
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Human name for the code to use for titles
	Title i18n.String `json:"title,omitempty" jsonschema:"title=Title"`

	// Useful description of the category.
	Description *i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Retained when true implies that the tax amount will be retained
	// by the buyer on behalf of the supplier, and thus subtracted from
	// the invoice taxable base total. Typically used for taxes related to
	// income.
	Retained bool `json:"retained,omitempty" jsonschema:"title=Retained"`

	// Informative when true implies that the tax amount will be calculated
	// and reported but will not affect the invoice totals. Typically used
	// for taxes that are embedded in the base amount or don't impact the
	// final payable amount.
	Informative bool `json:"informative,omitempty" jsonschema:"title=Informative"`

	// Specific tax definitions inside this category.
	Keys []*KeyDef `json:"keys,omitempty" jsonschema:"title=Keys"`

	// Rates defines the set of rates that can be used with this category.
	Rates []*RateDef `json:"rates,omitempty" jsonschema:"title=Rates"`

	// Extensions defines a list of extension keys that may be used or required
	// as an alternative or alongside choosing a rate for the tax category.
	// Every key must be defined in the Regime's extensions table.
	Extensions []cbc.Key `json:"extensions,omitempty" jsonschema:"title=Extensions"`

	// Map defines a set of regime specific code mappings.
	Map cbc.CodeMap `json:"map,omitempty" jsonschema:"title=Map"`

	// List of sources for the information contained in this category.
	Sources []*cbc.Source `json:"sources,omitempty" jsonschema:"title=Sources"`

	// Extension key-value pairs that will be copied to the tax combo if this
	// category is used.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Meta contains additional information about the category that is relevant
	// for local frequently used formats.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ValidateWithContext ensures the Category's contents are correct.
func (c *CategoryDef) ValidateWithContext(ctx context.Context) error {
	r := RegimeDefFromContext(ctx)
	err := validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Code, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Title, validation.Required),
		validation.Field(&c.Description),
		validation.Field(&c.Sources),
		validation.Field(&c.Keys),
		validation.Field(&c.Rates),
		validation.Field(&c.Extensions,
			validation.Each(cbc.InKeyDefs(r.Extensions)),
		),
		validation.Field(&c.Map),
		validation.Field(&c.Retained, validation.When(c.Informative,
			validation.In(false).Error("cannot be true when informative is true"),
		)),
	)
	return err
}

// KeyDef provides the key definition for the category, if it exists.
func (c *CategoryDef) KeyDef(key cbc.Key) *KeyDef {
	if c == nil {
		return nil
	}
	for _, k := range c.Keys {
		if k.Key == key {
			return k
		}
	}
	return nil
}

// RateDef provides the rate definition for the category, if it exists,
// for a given key and rate. The key may be empty, in which case the rate
// must be defined without any keys.
func (c *CategoryDef) RateDef(key, rate cbc.Key) *RateDef {
	if c == nil {
		return nil
	}
	// First try to find exact match
	var scnd *RateDef
	for _, r := range c.Rates {
		if !r.HasKey(key) {
			continue
		}
		if r.Rate == rate {
			return r
		}
		// Use the second match as a fallback
		if rate.HasPrefix(r.Rate) && scnd == nil {
			scnd = r
		}
	}
	return scnd
}
