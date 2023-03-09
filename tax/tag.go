package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// TagDef describes a tax tag that can be used to identify additional
// requirements in an electronic invoice.
type TagDef struct {
	// Key used to identify the tag
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Name of this tag.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Additional local
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the tax tag looks valid.
func (td *TagDef) Validate() error {
	return validation.ValidateStruct(td,
		validation.Field(&td.Key, validation.Required),
		validation.Field(&td.Name),
		validation.Field(&td.Meta),
	)
}
