package org

import (
	"context"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Telephone describes what is expected for a telephone number.
type Telephone struct {
	uuid.Identify
	// Identifier for this number.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Free-text string that represents the telephone number.
	Number string `json:"num" jsonschema:"title=Number"`
}

// Validate checks the telephone objects number to ensure it looks correct.
func (t *Telephone) Validate() error {
	return t.ValidateWithContext(context.Background())
}

// Normalize will try to remove any unnecessary whitespace from the telephone number.
func (t *Telephone) Normalize() {
	if t == nil {
		return
	}
	uuid.Normalize(&t.UUID)
	t.Label = cbc.NormalizeString(t.Label)
	t.Number = strings.TrimSpace(t.Number)
}

// ValidateWithContext checks the telephone objects number to ensure it looks correct inside the provided context.
func (t *Telephone) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, t,
		validation.Field(&t.UUID),
		validation.Field(&t.Number, validation.Required),
	)
}
