package tax

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// Note represents a tax-related note, typically used for exemption reasons
// or other tax-specific explanations that need to be associated with a
// particular tax category.
type Note struct {
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
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
	rd := RegimeDefFromContext(ctx)
	return ValidateStructWithContext(ctx, n,
		validation.Field(&n.Category,
			rd.InCategories(),
		),
		validation.Field(&n.Key,
			rd.InCategoryKeys(n.Category),
		),
		validation.Field(&n.Text, validation.Required),
		validation.Field(&n.Ext),
	)
}
