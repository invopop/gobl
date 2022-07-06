package bill

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// TaxScheme allows for defining a specific or special scheme that applies to the
// billing document. Schemes are defined as needed for each region.
type TaxScheme string

// Tax defines a summary of the taxes which may be applied to an invoice.
type Tax struct {
	// Category of the tax already included in the line item prices, especially
	// useful for B2C retailers with customers who prefer final prices inclusive of
	// tax.
	PricesInclude org.Code `json:"prices_include,omitempty" jsonschema:"title=Prices Include"`

	// Special tax schemes that apply to this invoice according to local requirements.
	Schemes tax.SchemeKeys `json:"schemes,omitempty" jsonschema:"title=Schemes"`

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}
