package bill

import (
	"github.com/invopop/gobl/cbc"
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

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ContainsScheme returns true if the tax contains the given scheme.
func (t *Tax) ContainsTag(key cbc.Key) bool {
	for _, s := range t.Tags {
		if s == key {
			return true
		}
	}
	return false
}

// Validate ensures the tax details look valid.
func (t *Tax) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.PricesInclude),
		validation.Field(&t.Tags),
		validation.Field(&t.Meta),
	)
}
