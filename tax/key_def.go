package tax

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// KeyDef defines a key that can be used inside a tax category. Rates may also be
// defined in addition to the key, and reference them for filtering purposes.
type KeyDef struct {
	// Key identifies this rate within the system
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`

	// Human name of the rate set
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Useful description of the rate.
	Description i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// NoPercent when true implies that the rate when used in a tax Combo should
	// not define a percent value.
	NoPercent bool `json:"no_percent,omitempty" jsonschema:"title=No Percent"`
}

// ValidateWithContext checks that the rate set is valid.
func (r *KeyDef) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, r,
		validation.Field(&r.Key),
		validation.Field(&r.Name),
	)
}
