package bill

import (
	"context"
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// TaxScheme allows for defining a specific or special scheme that applies to the
// billing document. Schemes are defined as needed for each region.
type TaxScheme string

// Tax defines a summary of the taxes which may be applied to an invoice.
type Tax struct {
	// Category of the tax already included in the line item prices, especially
	// useful for B2C retailers with customers who prefer final prices inclusive of
	// tax.
	PricesInclude cbc.Code `json:"prices_include,omitempty" jsonschema:"title=Prices Include"`

	// Special tax tags that apply to this invoice according to local requirements.
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Calculator defines the rule to use when calculating the taxes.
	// Currently supported options: `line`, or `total` (default).
	Calculator cbc.Key `json:"calculator,omitempty" jsonschema:"title=Calculator"`

	// Rounding defines the rounding rule to use when calculating the invoice sums.
	Rounding cbc.Key `json:"rounding,omitempty" jsonschema:"title=Rounding"`

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ContainsTag returns true if the tax contains the given tag.
func (t *Tax) ContainsTag(key cbc.Key) bool {
	if t == nil {
		return false
	}
	return key.In(t.Tags...)
}

// ValidateWithContext ensures the tax details look valid.
func (t *Tax) ValidateWithContext(ctx context.Context) error {
	r, _ := ctx.Value(tax.KeyRegime).(*tax.Regime)
	if r == nil {
		return errors.New("tax regime not found in context")
	}
	return validation.ValidateStructWithContext(ctx, t,
		validation.Field(&t.PricesInclude),
		validation.Field(&t.Tags, validation.Each(r.InTags())),
		validation.Field(&t.Meta),
	)
}

// JSONSchemaExtend extends the schema with additional property details
func (Tax) JSONSchemaExtend(schema *jsonschema.Schema) {
	props := schema.Properties
	if val, ok := props.Get("calculator"); ok {
		its := val.(*jsonschema.Schema)
		its.OneOf = make([]*jsonschema.Schema, len(tax.TotalCalculatorDefs))
		for i, v := range tax.TotalCalculatorDefs {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       v.Key.String(),
				Title:       v.Name.String(i18n.EN),
				Description: v.Desc.String(i18n.EN),
			}
		}
	}
	if val, ok := props.Get("rounding"); ok {
		its := val.(*jsonschema.Schema)
		its.OneOf = make([]*jsonschema.Schema, len(tax.TotalRoundingDefs))
		for i, v := range tax.TotalRoundingDefs {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       v.Key.String(),
				Title:       v.Name.String(i18n.EN),
				Description: v.Desc.String(i18n.EN),
			}
		}
	}
}
