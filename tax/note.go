package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Note represents a tax-related note, typically used for exemption reasons
// or other tax-specific explanations that need to be associated with a
// particular tax category.
type Note struct {
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat,omitempty" jsonschema:"title=Category"`
	// Key usually identifies the tax rate key this note applies to (e.g. "exempt",
	// "reverse-charge"), but may also be used for other identifiers depending on context.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Text contains the exemption reason or explanation.
	Text string `json:"text" jsonschema:"title=Text"`
	// Extensions for additional structured data.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// SameAs returns true if the two notes refer to the same tax situation,
// based on Category and Key.
func (n *Note) SameAs(n2 *Note) bool {
	if n == nil || n2 == nil {
		return false
	}
	return n.Category == n2.Category && n.Key == n2.Key
}

// Normalize cleans up the note's extensions.
func (n *Note) Normalize(normalizers Normalizers) {
	if n == nil {
		return
	}
	n.Ext = CleanExtensions(n.Ext)
	normalizers.Each(n)
}

func noteRules() *rules.Set {
	return rules.For(new(Note),
		rules.Field("text",
			rules.Assert("01", "tax note text is required", is.Present),
		),
	)
}
