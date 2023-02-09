package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cbc"
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

	// Special tax schemes that apply to this invoice according to local requirements.
	Schemes []cbc.Key `json:"schemes,omitempty" jsonschema:"title=Schemes"`

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ContainsScheme returns true if the tax contains the given scheme.
func (t *Tax) ContainsScheme(key cbc.Key) bool {
	for _, s := range t.Schemes {
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
		validation.Field(&t.Schemes),
		validation.Field(&t.Meta),
	)
}
