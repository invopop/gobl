package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Identity is used to define a code for a specific context.
type Identity struct {
	// Unique identity for this identity object.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Optional label useful for non-standard identities to give a bit more context.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The type of Code being represented and usually specific for
	// a particular context, country, or tax regime.
	Type cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`
	// The actual value of the identity code.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
}

// Validate ensures the identity looks valid.
func (i *Identity) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the identity looks valid inside the provided context.
func (i *Identity) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.Label),
		validation.Field(&i.Type),
		validation.Field(&i.Code, validation.Required),
	)
}
