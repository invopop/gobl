package bill

import (
	"context"
	"encoding/json"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Tax defines a summary of the taxes which may be applied to an invoice.
type Tax struct {
	// Category of the tax already included in the line item prices, especially
	// useful for B2C retailers with customers who prefer final prices inclusive of
	// tax.
	PricesInclude cbc.Code `json:"prices_include,omitempty" jsonschema:"title=Prices Include"`

	// Rounding model used to perform tax calculations on the invoice. This
	// will be configured automatically based on the tax regime, or
	// `sum-then-round` by default, but you can override here if needed.
	// Use with caution, as some conversion tools may make assumptions about
	// the rounding model used.
	Rounding cbc.Key `json:"rounding,omitempty" jsonschema:"title=Rounding Model"`

	// Additional extensions that are applied to the invoice as a whole as opposed to specific
	// sections.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	tags []cbc.Key
}

// MergeExtensions makes it easier to add extensions to the tax object
// by automatically handling nil data, and replying a new updated instance.
func (t *Tax) MergeExtensions(ext tax.Extensions) *Tax {
	if len(ext) == 0 {
		return t
	}
	if t == nil {
		t = new(Tax)
	}
	t.Ext = t.Ext.Merge(ext)
	return t
}

// GetExt is a convenience method to retrieve an extension value while
// providing nil checks on the tax object.
func (t *Tax) GetExt(key cbc.Key) cbc.Code {
	if t == nil {
		return cbc.CodeEmpty
	}
	return t.Ext.Get(key)
}

// HasExt is a convenience method to check for an extension value while
// providing nil checks on the tax object.
func (t *Tax) HasExt(key cbc.Key) bool {
	if t == nil {
		return false
	}
	return t.Ext.Has(key)
}

// Normalize performs normalization on the tax and embedded objects using the
// provided list of normalizers.
func (t *Tax) Normalize(normalizers tax.Normalizers) {
	if t == nil {
		return
	}
	// migration for old rounding rules
	switch t.Rounding {
	case "sum-then-round":
		t.Rounding = tax.RoundingRulePrecise
	case "round-then-sum":
		t.Rounding = tax.RoundingRuleCurrency
	}
	t.Ext = tax.CleanExtensions(t.Ext)
	normalizers.Each(t)
}

// ValidateWithContext ensures the tax details look valid.
func (t *Tax) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, t,
		validation.Field(&t.PricesInclude),
		validation.Field(&t.Rounding,
			cbc.InKeyDefs(tax.RoundingRules),
		),
		validation.Field(&t.Ext),
		validation.Field(&t.Meta),
	)
}

// UnmarshalJSON helps migrate the desc field to description.
func (t *Tax) UnmarshalJSON(data []byte) error {
	type Alias Tax
	aux := struct {
		Tags []cbc.Key `json:"tags,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.tags = aux.Tags
	return nil
}

// JSONSchemaExtend is used to add the additional options to the JSON schema.
func (t Tax) JSONSchemaExtend(schema *jsonschema.Schema) {
	if p, ok := schema.Properties.Get("rounding"); ok {
		p.OneOf = make([]*jsonschema.Schema, len(tax.RoundingRules))
		for i, r := range tax.RoundingRules {
			p.OneOf[i] = &jsonschema.Schema{
				Const:       r.Key.String(),
				Title:       r.Name.String(),
				Description: r.Desc.String(),
			}
		}
	}
}
