package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// Tag describes a tax tag that can be used to identify additional
// requirements in an electronic invoice.
type Tag struct {
	// Key used to identify the tag
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Name of this scheme.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Human details describing what this scheme is used for.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`
	// List of schemes that this tag can appear under.
	Schemes []cbc.Key `json:"schemes,omitempty" jsonschema:"title=Schemes"`
	// Additional local
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the tax tag looks valid.
func (t *Tag) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Key, validation.Required),
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.Description),
		validation.Field(&t.Meta),
	)
}
