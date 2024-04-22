package org

import (
	"context"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Telephone describes what is expected for a telephone number.
type Telephone struct {
	uuid.Identify
	// Identifier for this number.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The number to be dialed in ITU E.164 international format.
	Number string `json:"num" jsonschema:"title=Number"`
}

// Validate checks the telephone objects number to ensure it looks correct.
func (t *Telephone) Validate() error {
	return t.ValidateWithContext(context.Background())
}

// ValidateWithContext checks the telephone objects number to ensure it looks correct inside the provided context.
func (t *Telephone) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, t,
		validation.Field(&t.UUID),
		validation.Field(&t.Number, validation.Required, is.E164),
	)
}
