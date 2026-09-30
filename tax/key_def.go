package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// KeyDef defines a key that can be used inside a tax category. Rates may also be
// defined in addition to the key, and reference them for filtering purposes.
type KeyDef struct {
	// Unique identifier for the key within the category.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`

	// Human readable name of the key.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Useful description of the key.
	Description i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// When true, a tax combo using this key should not define a percent value.
	NoPercent bool `json:"no_percent,omitempty" jsonschema:"title=No Percent"`
}
