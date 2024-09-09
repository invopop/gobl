package bill

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
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

	// Additional extensions that are applied to the invoice as a whole as opposed to specific
	// sections.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

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
	return validation.ValidateStructWithContext(ctx, t,
		validation.Field(&t.PricesInclude),
		validation.Field(&t.Tags, validation.Each(r.InTags())),
		validation.Field(&t.Ext),
		validation.Field(&t.Meta),
	)
}
