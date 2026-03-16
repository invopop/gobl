package tax

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/validation"
)

// Note represents a tax-related note, typically used for exemption reasons
// or other tax-specific explanations that need to be associated with a
// particular tax category.
type Note struct {
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Country code override when issuing with taxes applied from different countries.
	Country l10n.TaxCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// Key identifies the tax situation this note applies to (e.g. "exempt", "reverse-charge").
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Text contains the exemption reason or explanation.
	Text string `json:"text" jsonschema:"title=Text"`
	// Extensions for additional structured data.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize cleans up the note's extensions.
func (n *Note) Normalize(normalizers Normalizers) {
	if n == nil {
		return
	}
	n.Ext = CleanExtensions(n.Ext)
	normalizers.Each(n)
}

// ValidateWithContext ensures the note is valid.
func (n *Note) ValidateWithContext(ctx context.Context) error {
	var r *RegimeDef
	if n.Country.Empty() {
		r = RegimeDefFromContext(ctx)
	} else {
		r = RegimeDefFor(n.Country.Code())
	}
	return ValidateStructWithContext(ctx, n,
		validation.Field(&n.Category,
			r.InCategories(),
		),
		validation.Field(&n.Country),
		validation.Field(&n.Key,
			r.InCategoryKeys(n.Category),
		),
		validation.Field(&n.Text, validation.Required),
		validation.Field(&n.Ext),
	)
}
