package bill

import (
	"context"
	"encoding/json"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Tax defines a summary of the taxes which may be applied to an invoice.
type Tax struct {
	// Category of the tax already included in the line item prices, especially
	// useful for B2C retailers with customers who prefer final prices inclusive of
	// tax.
	PricesInclude cbc.Code `json:"prices_include,omitempty" jsonschema:"title=Prices Include"`

	// Additional extensions that are applied to the invoice as a whole as opposed to specific
	// sections.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	tags []cbc.Key
}

// Normalize performs normalization on the tax and embedded objects using the
// provided list of normalizers.
func (t *Tax) Normalize(normalizers tax.Normalizers) {
	if t == nil {
		return
	}
	t.Ext = tax.CleanExtensions(t.Ext)
	normalizers.Each(t)
}

// ValidateWithContext ensures the tax details look valid.
func (t *Tax) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, t,
		validation.Field(&t.PricesInclude),
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
